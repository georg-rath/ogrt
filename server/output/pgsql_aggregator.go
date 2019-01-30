package output

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/georg-rath/ogrt/protocol"

	"github.com/golang/protobuf/jsonpb"
	"github.com/lib/pq"
	"github.com/rcrowley/go-metrics"
)

/**
Aggregate ProcessInfo messages into job summaries, stored in a PostgreSQL database.

Each process sends a ProcessStart message on startup and a ProcessEnd message on termination.
The ProcessEnd message might not get sent in case a process terminates without calling the finishing
handler, which happens on SIGKILL or rare corner cases. ProcessStart messages should always arrive, but might
get lost in transit.
Incoming Process and ResourceInfos are stored ProcessTuple, which is kept in memory until it is complete (containing
both messages). After completion the ProcessTuple is added to a JobData struct, which is cached.

On completion of ProcessTuple:
1. Check if JobData is cached
2. If not pull it in.
3. Add Tuple to JobData
-> there is a flushing goroutine, which flushes the cache to the database every n seconds, we need to mark the JobData dirty on update and load
-> there is a purging goroutine, which purges the cache of aged out items every n seconds

*/
type PgSqlAggregatorOutput struct {
	Emitter

	completionFn func(n int)
}

var jobCache *JobDataCache
var tupleCache *TupleCache

var db *sql.DB
var loadJobDataSql string = `SELECT "ID", "JobID", "User", "FirstCommand", "LastCommand", "MaxRSS", "MaxRSSCmdline", "RuUtime", "RuStime", "RuMinflt", "RuMajflt", "RuInblock", "RuOublock", "RuNvcsw", "RuNivcsw", "Modules", "SharedObjects", "Hosts" FROM "Jobs" WHERE "JobID" = $1`
var loadJobDataStmt *sql.Stmt

var insertJobDataSql string = `
	INSERT INTO "Jobs"
		("ID", "JobID", "User", "FirstCommand", "LastCommand", "MaxRSS", "MaxRSSCmdline", "RuUtime", "RuStime", "RuMinflt", "RuMajflt", "RuInblock", "RuOublock", "RuNvcsw", "RuNivcsw", "SharedObjects", "Modules", "Hosts")
		VALUES
		(DEFAULT, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
  RETURNING "ID";`
var updateJobDataSql string = `
	UPDATE "Jobs"
	SET
		"JobID"        =  $2,
		"User"         =  $3,
		"FirstCommand" =  $4,
		"LastCommand"  =  $5,
		"MaxRSS"       =  $6,
		"MaxRSSCmdline"=  $7,
		"RuUtime"      =  $8,
		"RuStime"      =  $9,
		"RuMinflt"     = $10,
		"RuMajflt"     = $11,
		"RuInblock"    = $12,
		"RuOublock"    = $13,
		"RuNvcsw"      = $14,
		"RuNivcsw"     = $15,
		"SharedObjects"= $16,
		"Modules"      = $17,
		"Hosts"        = $18
	WHERE "ID" = $1;`
var updateJobDataStmt *sql.Stmt
var insertJobDataStmt *sql.Stmt

type JobData struct {
	RowId         int64
	JobId         string
	User          string
	Hosts         []string
	FirstCommand  time.Time
	LastCommand   time.Time
	MaxRSS        int64
	MaxRSSCmdline string
	RuUtime       int64
	RuStime       int64
	RuMinflt      int64
	RuMajflt      int64
	RuInblock     int64
	RuOublock     int64
	RuNvcsw       int64
	RuNivcsw      int64
	Modules       []*msg.Module
	SharedObjects []*msg.SharedObject

	dirty bool
	mu    *sync.Mutex
}

func LoadJobData(pi *msg.ProcessStart) (jd *JobData) {
	log.Println("load job data for", pi.JobId)
	rows, err := loadJobDataStmt.Query(pi.JobId)
	if err != nil {
		log.Fatal("failed loading job data")
	}

	// here we can get multiple entries, in case jobs rolled over or some other kind of funkiness going on.
	// we take the job that fits based on start time, as this is likely to be "the right thing".
	var modules, sharedObjects string
	for rows.Next() {
		rjd := &JobData{
			dirty: false,
			mu:    &sync.Mutex{},
		}
		if err := rows.Scan(&rjd.RowId, &rjd.JobId, &rjd.User, &rjd.FirstCommand, &rjd.LastCommand, &rjd.MaxRSS, &rjd.MaxRSSCmdline, &rjd.RuUtime, &rjd.RuStime, &rjd.RuMinflt, &rjd.RuMajflt, &rjd.RuInblock, &rjd.RuOublock, &rjd.RuNvcsw, &rjd.RuNivcsw, &modules, &sharedObjects, pq.Array(&rjd.Hosts)); err != nil {
			log.Fatal(err)
		}
		if jd == nil {
			jd = rjd
			continue
		}
		// get the one closest to the current one
		if convertFromMsToTime(pi.Time).Sub(rjd.FirstCommand) < convertFromMsToTime(pi.Time).Sub(jd.FirstCommand) {
			jd = rjd
		}
	}

	if jd == nil {
		jd = &JobData{
			JobId:         pi.JobId,
			User:          pi.Username,
			Hosts:         []string{pi.Hostname},
			Modules:       pi.LoadedModules,
			SharedObjects: pi.SharedObjects,
			dirty:         true,
			mu:            &sync.Mutex{},
		}
	}

	// this is awkward but fast unmarshalling, can't go from array of JSON to protobuf directly
	// we assume we get a valid json array of decodeable protobufs.
	// note to self: when doing servers, a boundary between protocol and internal data format might pay off...
	if modules != "" && modules != "[]" {
		begin, end := 0, 0
		for {
			begin = begin + strings.Index(modules[begin:], "{")
			newEnd := strings.Index(modules[begin:], "}")
			if newEnd == -1 {
				break
			}
			end = end + newEnd + (begin - end) + 1
			pb := &msg.Module{}
			if err := jsonpb.UnmarshalString(modules[begin:end], pb); err != nil {
				log.Fatal("unmarshal of json to pb (modules)", err)
			}
			jd.Modules = append(jd.Modules, pb)
			begin = end + 1
		}
	}
	if sharedObjects != "" && sharedObjects != "[]" {
		begin, end := 0, 0
		for {
			begin = begin + strings.Index(sharedObjects[begin:], "{")
			newEnd := strings.Index(sharedObjects[begin:], "}")
			if newEnd == -1 {
				break
			}
			end = end + newEnd + (begin - end) + 1
			pb := &msg.SharedObject{}
			if err := jsonpb.UnmarshalString(sharedObjects[begin:end], pb); err != nil {
				log.Fatal("unmarshal of json to pb (shared objects)", err)
			}
			jd.SharedObjects = append(jd.SharedObjects, pb)
			begin = end + 1
		}
	}

	return
}

func (jd *JobData) AddTuple(pt *ProcessTuple) {
	jd.mu.Lock()
	defer jd.mu.Unlock()
	if !pt.Complete() {
		return
	}

	jd.RuUtime += pt.ProcessEnd.RuUtime
	jd.RuStime += pt.ProcessEnd.RuStime
	jd.RuMinflt += pt.ProcessEnd.RuMinflt
	jd.RuMajflt += pt.ProcessEnd.RuMajflt
	jd.RuInblock += pt.ProcessEnd.RuInblock
	jd.RuOublock += pt.ProcessEnd.RuOublock
	jd.RuNvcsw += pt.ProcessEnd.RuNvcsw
	jd.RuNivcsw += pt.ProcessEnd.RuNivcsw

	// merge hosts
	addHost := true
	for _, host := range jd.Hosts {
		if host == pt.ProcessStart.Hostname {
			addHost = false
			break
		}
	}
	if addHost {
		jd.Hosts = append(jd.Hosts, pt.ProcessStart.Hostname)
		sort.Strings(jd.Hosts)
	}

	if pt.ProcessEnd.RuMaxrss > jd.MaxRSS {
		jd.MaxRSS = pt.ProcessEnd.RuMaxrss
		jd.MaxRSSCmdline = pt.ProcessStart.Cmdline
	}

	t := convertFromMsToTime(pt.ProcessStart.Time)
	if jd.FirstCommand.IsZero() || t.Before(jd.FirstCommand) {
		jd.FirstCommand = t
	}
	if jd.LastCommand.IsZero() || t.After(jd.LastCommand) {
		jd.LastCommand = t
	}

nextModule:
	for _, newModule := range pt.ProcessStart.LoadedModules {
		for _, oldModule := range jd.Modules {
			if newModule.Name == oldModule.Name {
				continue nextModule
			}
		}
		jd.Modules = append(jd.Modules, newModule)
	}

nextSo:
	for _, newSo := range pt.ProcessStart.SharedObjects {
		for _, oldSo := range jd.SharedObjects {
			if newSo.Path == oldSo.Path && newSo.Signature == oldSo.Signature {
				continue nextSo
			}
		}
		jd.SharedObjects = append(jd.SharedObjects, newSo)
	}

	jd.dirty = true
	return
}

func (jd *JobData) Save() {
	jd.mu.Lock()
	defer jd.mu.Unlock()

	// marshal shared objects to JSON
	// needs a bit of a workaround, cause we can't directly marshal an array of protobufs
	m := &jsonpb.Marshaler{}
	var soBuilder strings.Builder
	lastElement := len(jd.SharedObjects) - 1
	soBuilder.WriteString("[")
	for i, so := range jd.SharedObjects {
		if s, err := m.MarshalToString(so); err == nil {
			soBuilder.WriteString(s)
			if i != lastElement {
				soBuilder.WriteString(",")
			}
		}
	}
	soBuilder.WriteString("]")

	// marshal modules to JSON - see above
	var moduleBuilder strings.Builder
	lastElement = len(jd.Modules) - 1
	moduleBuilder.WriteString("[")
	for i, module := range jd.Modules {
		if s, err := m.MarshalToString(module); err == nil {
			moduleBuilder.WriteString(s)
			if i != lastElement {
				moduleBuilder.WriteString(",")
			}
		}
	}
	moduleBuilder.WriteString("]")

	if jd.RowId == 0 {
		if err := insertJobDataStmt.QueryRow(jd.JobId, jd.User, jd.FirstCommand, jd.LastCommand, jd.MaxRSS, jd.MaxRSSCmdline, jd.RuUtime, jd.RuStime, jd.RuMinflt, jd.RuMajflt, jd.RuInblock, jd.RuOublock, jd.RuNvcsw, jd.RuNivcsw, soBuilder.String(), moduleBuilder.String(), pq.Array(jd.Hosts)).Scan(&jd.RowId); err != nil {
			log.Println("failed inserting job data", err)
			return
		}
		jd.dirty = false
		return
	}
	if _, err := updateJobDataStmt.Exec(jd.RowId, jd.JobId, jd.User, jd.FirstCommand, jd.LastCommand, jd.MaxRSS, jd.MaxRSSCmdline, jd.RuUtime, jd.RuStime, jd.RuMinflt, jd.RuMajflt, jd.RuInblock, jd.RuOublock, jd.RuNvcsw, jd.RuNivcsw, soBuilder.String(), moduleBuilder.String(), pq.Array(jd.Hosts)); err != nil {
		log.Println("failed updating job data", err)
		return
	}
	jd.dirty = false
}

func (jd *JobData) Dirty() (dirty bool) {
	jd.mu.Lock()
	defer jd.mu.Unlock()
	return jd.dirty
}

type ProcessTuple struct {
	ProcessStart *msg.ProcessStart
	ProcessEnd   *msg.ProcessEnd
	LastUpdate   time.Time
}

func (pt *ProcessTuple) Complete() bool {
	if pt.ProcessStart != nil && pt.ProcessEnd != nil {
		return true
	}
	return false
}

func (fw *PgSqlAggregatorOutput) Open(completionFn func(n int), config map[string]interface{}) {
	host := config["Host"].(string)
	port := config["Port"].(int64)
	user := config["User"].(string)
	password := config["Password"].(string)
	database := config["Database"].(string)
	var flushInterval time.Duration
	if flushIntervalRaw, ok := config["FlushInterval"].(int64); !ok {
		flushInterval = 10 * time.Second
	} else {
		flushInterval = time.Duration(flushIntervalRaw) * time.Second
	}
	var purgeInterval time.Duration
	if purgeIntervalRaw, ok := config["PurgeInterval"].(int64); !ok {
		purgeInterval = 3 * time.Minute
	} else {
		purgeInterval = time.Duration(purgeIntervalRaw) * time.Second
	}
	var purgeAge time.Duration
	if purgeAgeRaw, ok := config["PurgeAge"].(int64); !ok {
		purgeAge = 3 * time.Minute
	} else {
		purgeAge = time.Duration(purgeAgeRaw) * time.Second
	}

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, database))
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	createTable := `CREATE TABLE IF NOT EXISTS "Jobs" (
		"ID" serial NOT NULL,
		"JobID" text NULL,
		"User" text NOT NULL,
		"Hosts" text[] NOT NULL,
		"FirstCommand" timestamp NOT NULL,
		"LastCommand" timestamp NOT NULL,
		"MaxRSS" bigint NOT NULL,
		"MaxRSSCmdline" text NOT NULL,
		"RuUtime" bigint NOT NULL,
		"RuStime" integer NOT NULL,
		"RuMinflt" bigint NOT NULL,
		"RuMajflt" bigint NOT NULL,
		"RuInblock" bigint NOT NULL,
		"RuOublock" bigint NOT NULL,
		"RuNvcsw" bigint NOT NULL,
		"RuNivcsw" bigint NOT NULL,
		"Modules" jsonb NOT NULL,
		"SharedObjects" jsonb NOT NULL
	);`

	if _, err := db.Exec(createTable); err != nil {
		log.Fatal("failed creating table", err)
	}
	if loadJobDataStmt, err = db.Prepare(loadJobDataSql); err != nil {
		log.Fatal("failed to prepare query", err)
	}
	if updateJobDataStmt, err = db.Prepare(updateJobDataSql); err != nil {
		log.Fatal("failed to prepare query", err)
	}
	if insertJobDataStmt, err = db.Prepare(insertJobDataSql); err != nil {
		log.Fatal("failed to prepare query", err)
	}

	jobCache = NewJobDataCache(flushInterval, purgeInterval, purgeAge)
	tupleCache = NewTupleCache()

	fw.completionFn = completionFn
}

func (fw *PgSqlAggregatorOutput) EmitProcessStart(pi *msg.ProcessStart) {
	uuid := hex.EncodeToString(pi.Uuid)
	pt := &ProcessTuple{
		ProcessStart: pi,
		LastUpdate:   time.Now(),
	}
	if tuple, loaded := tupleCache.LoadOrStore(uuid, pt); !loaded {
		jobCache.GetOrLoad(pi)
	} else {
		// resourceinfo arrived already
		jobCache.GetOrLoad(pi).AddTuple(tuple)
		tupleCache.Delete(uuid)
	}

	fw.completionFn(-1)
}

func (fw *PgSqlAggregatorOutput) EmitProcessEnd(pri *msg.ProcessEnd) {
	uuid := hex.EncodeToString(pri.Uuid)
	pt := &ProcessTuple{
		ProcessEnd: pri,
		LastUpdate: time.Now(),
	}
	if pt, loaded := tupleCache.LoadOrStore(uuid, pt); loaded {
		// ProcessStart was already in cache
		pt.ProcessEnd = pri
		if pt.ProcessStart == nil {
			//TODO: this is an edge case
			//TODO: in case of multiple end messages this loses valuable info
			log.Println("No Start for UUID", uuid)
			return
		}
		jobCache.GetOrLoad(pt.ProcessStart).AddTuple(pt)
		tupleCache.Delete(uuid)
	}

	fw.completionFn(-1)
}

func (fw *PgSqlAggregatorOutput) Close() {
	jobCache.Close()
	db.Close()
}

// JobDataCache caches aggregated job data.
// Assumption is that this cache is read-mostly, which is why
// the map that stores the data is protected with a RW lock, with
// each element in the map having a separate mutex. This should make
// it so that we incur locking penalty only in the insert case.
type JobDataCache struct {
	cache map[string]*JobData
	mu    *sync.RWMutex

	stop    chan struct{}
	stopped chan struct{}
}

func NewJobDataCache(flushInterval, purgeInterval, purgeAge time.Duration) (jdc *JobDataCache) {
	jdc = &JobDataCache{
		cache: make(map[string]*JobData),
		mu:    &sync.RWMutex{},

		stop:    make(chan struct{}),
		stopped: make(chan struct{}),
	}

	go func() {
		flushTicker := time.NewTicker(flushInterval)
		defer flushTicker.Stop()
		purgeTicker := time.NewTicker(purgeInterval)
		defer purgeTicker.Stop()

		flushFn := func() {
			// log.Printf("%d elements in tupleCache\n", tupleCache.Size())
			// log.Printf("%d elements in jobCache\n", jobCache.Size())
			jdc.mu.RLock()
			flushed := 0
			for _, v := range jdc.cache {
				if v.Dirty() {
					v.Save()
					flushed++
				}
			}
			jdc.mu.RUnlock()
			log.Printf("flushed %d dirty elements from jobCache\n", flushed)
		}
		for {
			select {
			case <-jdc.stop:
				// log.Println("flushing aggregator cache...")
				flushFn()
				jdc.stopped <- struct{}{}
				break
			case <-flushTicker.C:
				flushFn()
			case <-purgeTicker.C:
				jdc.mu.Lock()
				purged := 0
				now := time.Now()
				for _, v := range jdc.cache {
					if !v.Dirty() {
						age := now.Sub(v.LastCommand)
						if age > purgeAge {
							delete(jdc.cache, v.JobId)
							purged++
						}
					}
				}
				log.Printf("purged %d elements from jobCache\n", purged)
				jdc.mu.Unlock()
			}
		}
	}()

	return
}

func (jdc *JobDataCache) GetOrLoad(pi *msg.ProcessStart) (jd *JobData) {
	jdc.mu.RLock()
	value, found := jdc.cache[pi.JobId]
	// we did not find the entry, lets fetch it from the db
	if !found {
		// drop read lock and get write lock
		jdc.mu.RUnlock()
		jdc.mu.Lock()
		defer jdc.mu.Unlock()
		// check again, there might have been a concurrent Get
		if value, found := jdc.cache[pi.JobId]; found {
			return value
		}
		// there wasn't - fetch from db
		jd = LoadJobData(pi)
		jdc.cache[pi.JobId] = jd
		return jd
	}
	jdc.mu.RUnlock()
	return value
}

func (jdc *JobDataCache) Size() uint {
	return uint(len(jdc.cache))
}

func (jdc *JobDataCache) Close() {
	jdc.stop <- struct{}{}
	<-jdc.stopped
}

// TupleCache caches incomplete tuples. Tuples get evicted once complete.
// We expect high insert rate, which is why this cache was implemented using
// a sync.Map, as we expect it to perform better than a lock protected map
// for this use case.
type TupleCache struct {
	cache *sync.Map
	size  uint64
}

func NewTupleCache() (tc *TupleCache) {
	tc = &TupleCache{
		cache: &sync.Map{},
	}

	metrics.NewRegisteredFunctionalGauge("tuplecache_size", metrics.DefaultRegistry, func() int64 {
		return int64(tc.Size())
	})

	return
}

func (tc *TupleCache) LoadOrStore(key string, value *ProcessTuple) (actual *ProcessTuple, loaded bool) {
	actualUntyped, loaded := tc.cache.LoadOrStore(key, value)
	actual = actualUntyped.(*ProcessTuple)
	if !loaded {
		atomic.AddUint64(&tc.size, 1)
	}
	return
}

func (tc *TupleCache) Delete(key string) {
	tc.cache.Delete(key)
	atomic.AddUint64(&tc.size, ^uint64(0))
}

func (tc *TupleCache) Size() uint64 {
	return atomic.LoadUint64(&tc.size)
}
