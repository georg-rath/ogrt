package main

import "github.com/gin-gonic/gin"
import "github.com/rcrowley/go-metrics"
import "time"
import "log"
import "strings"

func StartWebAPI() {
	r := gin.Default()
	r.GET("/metrics", func(c *gin.Context) {
		results := outputMetrics(metrics.DefaultRegistry, log.Lmicroseconds)
		c.JSON(200, gin.H{
			"outputs": results,
		})
	})
	r.Run()
}

func outputMetrics(r metrics.Registry, scale time.Duration) map[string]map[string]float64 {
	//du := float64(scale)
	//duSuffix := scale.String()[1:]
	results := make(map[string]map[string]float64)

	r.Each(func(name string, i interface{}) {
		if strings.HasPrefix(name, "output_") {
			switch metric := i.(type) {
			case metrics.Timer:
				results[name] = make(map[string]float64)
				t := metric.Snapshot()
				//ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})

				results[name]["count"] = float64(t.Count())
				results[name]["rate1m"] = t.Rate1()
				results[name]["rate5m"] = t.Rate5()
				results[name]["rate15m"] = t.Rate15()
				results[name]["ratemean"] = t.RateMean()
				//log.Printf("  min:         %12.2f%s\n", float64(t.Min())/du, duSuffix)
				//log.Printf("  max:         %12.2f%s\n", float64(t.Max())/du, duSuffix)
				//log.Printf("  mean:        %12.2f%s\n", t.Mean()/du, duSuffix)
				//log.Printf("  stddev:      %12.2f%s\n", t.StdDev()/du, duSuffix)
				//log.Printf("  median:      %12.2f%s\n", ps[0]/du, duSuffix)
				//log.Printf("  75%%:         %12.2f%s\n", ps[1]/du, duSuffix)
				//log.Printf("  95%%:         %12.2f%s\n", ps[2]/du, duSuffix)
				//log.Printf("  99%%:         %12.2f%s\n", ps[3]/du, duSuffix)
				//log.Printf("  99.9%%:       %12.2f%s\n", ps[4]/du, duSuffix)
			}
		}
	})
	return results
}
