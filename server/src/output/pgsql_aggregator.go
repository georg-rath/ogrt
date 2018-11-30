package output

import (
	"database/sql"
	"encoding/hex"
	"log"
	"strings"
	"time"

	"github.com/georg-rath/ogrt/src/protocol"

	"github.com/golang/protobuf/jsonpb"
	_ "github.com/lib/pq"
)

/**
Aggregate ProcessInfo messages into job summaries, stored in a PostgreSQL database.

Each process sends a ProcessInfo message on startup and a ProcessResourceInfo message on termination.
The ProcessResourceInfo message might not get sent in case a process terminates without calling the finishing
handler, which happens on SIGKILL or rare corner cases. ProcessInfo messages should always arrive, but might
get lost in transit.
Incoming Process and ResourceInfos are stored ProcessTuple, which is kept in memory until it is complete (containing
both messages). After completion the ProcessTuple is added to a JobData struct, which is cached.

On completion of ProcessTuple:
1. Check if JobData is cached
2. If not pull it in. TODO: think about if this needs to be protected by a lock
3. Add Tuple to JobData
-> there is a flushing goroutine, which flushes out the cache to the database every n seconds, we need to mark the JobData dirty on update and load

*/
type PgSqlAggregatorOutput struct {
	OGWriter
}

var jobCache map[string]*JobData
var uuidCache map[string]*ProcessTuple

var db *sql.DB
var loadJobDataSql string = `SELECT "ID", "JobID", "User", "FirstCommand", "LastCommand", "MaxRSS", "MaxRSSCmdline", "RuUtime", "RuStime", "RuMinflt", "RuMajflt", "RuInblock", "RuOublock", "RuNvcsw", "RuNivcsw", "Modules", "SharedObjects" FROM "Jobs" WHERE "JobID" = $1`
var loadJobDataStmt *sql.Stmt

var insertJobDataSql string = `
	INSERT INTO "Jobs"
		("ID", "JobID", "User", "FirstCommand", "LastCommand", "MaxRSS", "MaxRSSCmdline", "RuUtime", "RuStime", "RuMinflt", "RuMajflt", "RuInblock", "RuOublock", "RuNvcsw", "RuNivcsw", "SharedObjects", "Modules")
		VALUES
		(DEFAULT, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);`
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
		"Modules"      = $17
	WHERE "ID" = $1;`
var updateJobDataStmt *sql.Stmt
var insertJobDataStmt *sql.Stmt

type JobData struct {
	RowId         int64
	JobId         string
	User          string
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
	Modules       []*OGRT.Module
	SharedObjects []*OGRT.SharedObject

	dirty bool
}

func LoadJobData(pi *OGRT.ProcessInfo) (jd *JobData) {
	rows, err := loadJobDataStmt.Query(pi.JobId)
	if err != nil {
		log.Fatal("failed loading job data")
	}

	// here we can get multiple entries, in case jobs rolled over or some other kind of funkiness going on.
	// we take the job that fits based on start time, as this is likely to be "the right thing".
	var modules, sharedObjects string
	for rows.Next() {
		rjd := &JobData{}
		if err := rows.Scan(&rjd.RowId, &rjd.JobId, &rjd.User, &rjd.FirstCommand, &rjd.LastCommand, &rjd.MaxRSS, &rjd.MaxRSSCmdline, &rjd.RuUtime, &rjd.RuStime, &rjd.RuMinflt, &rjd.RuMajflt, &rjd.RuInblock, &rjd.RuOublock, &rjd.RuNvcsw, &rjd.RuNivcsw, &modules, &sharedObjects); err != nil {
			log.Fatal(err)
		}
		if jd == nil {
			jd = rjd
			continue
		}
		// get the one closest to the current one
		if convertTime(pi.Time).Sub(rjd.FirstCommand) < convertTime(pi.Time).Sub(jd.FirstCommand) {
			jd = rjd
		}
	}

	if jd == nil {
		jd = &JobData{
			JobId:         pi.JobId,
			User:          pi.Username,
			Modules:       pi.LoadedModules,
			SharedObjects: pi.SharedObjects,
			dirty:         true,
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
			pb := &OGRT.Module{}
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
			pb := &OGRT.SharedObject{}
			if err := jsonpb.UnmarshalString(sharedObjects[begin:end], pb); err != nil {
				log.Fatal("unmarshal of json to pb (shared objects)", err)
			}
			jd.SharedObjects = append(jd.SharedObjects, pb)
			begin = end + 1
		}
	}

	return
}

func convertTime(raw int64) time.Time {
	return time.Unix(0, raw*int64(time.Millisecond))
}

func (jd *JobData) AddTuple(pt *ProcessTuple) {
	if !pt.Complete() {
		return
	}

	jd.RuUtime += pt.ProcessResourceInfo.RuUtime
	jd.RuStime += pt.ProcessResourceInfo.RuStime
	jd.RuMinflt += pt.ProcessResourceInfo.RuMinflt
	jd.RuMajflt += pt.ProcessResourceInfo.RuMajflt
	jd.RuInblock += pt.ProcessResourceInfo.RuInblock
	jd.RuOublock += pt.ProcessResourceInfo.RuOublock
	jd.RuNvcsw += pt.ProcessResourceInfo.RuNvcsw
	jd.RuNivcsw = jd.RuNivcsw + pt.ProcessResourceInfo.RuNivcsw

	if pt.ProcessResourceInfo.RuMaxrss > jd.MaxRSS {
		jd.MaxRSS = pt.ProcessResourceInfo.RuMaxrss
		jd.MaxRSSCmdline = pt.ProcessInfo.Cmdline
	}

	t := convertTime(pt.ProcessInfo.Time)
	if jd.FirstCommand.IsZero() || t.Before(jd.FirstCommand) {
		jd.FirstCommand = t
	}
	if jd.LastCommand.IsZero() || t.After(jd.LastCommand) {
		jd.LastCommand = t
	}

nextModule:
	for _, newModule := range pt.ProcessInfo.LoadedModules {
		for _, oldModule := range jd.Modules {
			if newModule.Name == oldModule.Name {
				continue nextModule
			}
		}
		jd.Modules = append(jd.Modules, newModule)
	}

nextSo:
	for _, newSo := range pt.ProcessInfo.SharedObjects {
		for _, oldSo := range jd.SharedObjects {
			if newSo.Path == oldSo.Path && newSo.Signature == oldSo.Signature {
				continue nextSo
			}
		}
		jd.SharedObjects = append(jd.SharedObjects, newSo)
	}
	return
}

func (jd *JobData) Save() {
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
		if _, err := insertJobDataStmt.Exec(jd.JobId, jd.User, jd.FirstCommand, jd.LastCommand, jd.MaxRSS, jd.MaxRSSCmdline, jd.RuUtime, jd.RuStime, jd.RuMinflt, jd.RuMajflt, jd.RuInblock, jd.RuOublock, jd.RuNvcsw, jd.RuNivcsw, soBuilder.String(), moduleBuilder.String()); err != nil {
			log.Println("failed inserting job data", err)
		}
		return
	}
	if _, err := updateJobDataStmt.Exec(jd.RowId, jd.JobId, jd.User, jd.FirstCommand, jd.LastCommand, jd.MaxRSS, jd.MaxRSSCmdline, jd.RuUtime, jd.RuStime, jd.RuMinflt, jd.RuMajflt, jd.RuInblock, jd.RuOublock, jd.RuNvcsw, jd.RuNivcsw, soBuilder.String(), moduleBuilder.String()); err != nil {
		log.Println("failed inserting job data", err)
	}
}

type ProcessTuple struct {
	ProcessInfo         *OGRT.ProcessInfo
	ProcessResourceInfo *OGRT.ProcessResourceInfo
	LastUpdate          time.Time
}

func (pt *ProcessTuple) Complete() bool {
	if pt.ProcessInfo != nil && pt.ProcessResourceInfo != nil {
		return true
	}
	return false
}

func (fw *PgSqlAggregatorOutput) Open(params string) {
	var err error
	db, err = sql.Open("postgres", params)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	createTable := `CREATE TABLE IF NOT EXISTS "Jobs" (
		"ID" serial NOT NULL,
		"JobID" text NOT NULL,
		"User" text NOT NULL,
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

	jobCache = make(map[string]*JobData)
	uuidCache = make(map[string]*ProcessTuple)
}

func (fw *PgSqlAggregatorOutput) PersistJobStart(job_start *OGRT.JobStart) {
}

func (fw *PgSqlAggregatorOutput) PersistJobEnd(job_end *OGRT.JobEnd) {
}

func (fw *PgSqlAggregatorOutput) PersistProcessResourceInfo(pri *OGRT.ProcessResourceInfo) {
	uuid := hex.EncodeToString(pri.Uuid)

	if _, found := uuidCache[uuid]; !found {
		pt := &ProcessTuple{
			ProcessResourceInfo: pri,
			LastUpdate:          time.Now(),
		}

		uuidCache[uuid] = pt
		return
	}

	// processinfo arrived already
	pt := uuidCache[uuid]
	pt.ProcessResourceInfo = pri
	//TODO: until we got the flusher right, the following statement will crash
	jd := jobCache[pt.ProcessInfo.JobId]
	jd.AddTuple(pt)
	delete(uuidCache, uuid)

	//TODO: move this to flusher goroutine
	jd.Save()
	delete(jobCache, pt.ProcessInfo.JobId)
}

func (fw *PgSqlAggregatorOutput) PersistProcessInfo(pi *OGRT.ProcessInfo) {
	uuid := hex.EncodeToString(pi.Uuid)
	if _, found := uuidCache[uuid]; !found {
		pt := &ProcessTuple{
			ProcessInfo: pi,
			LastUpdate:  time.Now(),
		}

		if _, found := jobCache[pi.JobId]; !found {
			jd := LoadJobData(pi)
			jobCache[pi.JobId] = jd
		}
		uuidCache[uuid] = pt
		return
	}

	// resourceinfo arrived already
	pt := uuidCache[uuid]
	pt.ProcessInfo = pi
	if _, found := jobCache[pi.JobId]; !found {
		jd := LoadJobData(pi)
		jobCache[pi.JobId] = jd
	}
	jobCache[pi.JobId].AddTuple(pt)
	delete(uuidCache, uuid)
}

func (fw *PgSqlAggregatorOutput) Close() {
	db.Close()
}
