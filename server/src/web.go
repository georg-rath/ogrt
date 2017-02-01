package main

import "github.com/gin-gonic/gin"
import "github.com/rcrowley/go-metrics"
import "time"
import "strings"

func StartWebAPI() {
	r := gin.Default()
	r.GET("/metrics", func(c *gin.Context) {
		outputs := outputMetrics(metrics.DefaultRegistry, time.Millisecond)
		inputs := inputMetrics(metrics.DefaultRegistry, time.Millisecond)
		c.JSON(200, gin.H{
			"outputs": outputs,
			"inputs":  inputs,
		})
	})
	r.Static("/web", "./web")
	r.Run()
}

func inputMetrics(r metrics.Registry, scale time.Duration) map[string]map[string]float64 {
	du := float64(scale)
	results := make(map[string]map[string]float64)

	r.Each(func(name string, i interface{}) {
		if strings.HasPrefix(name, "input_") {
			switch metric := i.(type) {
			case metrics.Timer:
				cleanName := strings.SplitAfter(name, "_")[1]
				results[cleanName] = make(map[string]float64)
				t := metric.Snapshot()
				ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})

				results[cleanName]["count"] = float64(t.Count())
				results[cleanName]["rate1m"] = t.Rate1()
				results[cleanName]["rate5m"] = t.Rate5()
				results[cleanName]["rate15m"] = t.Rate15()
				results[cleanName]["ratemean"] = t.RateMean()
				results[cleanName]["duration_mean"] = t.Mean() / du
				results[cleanName]["duration_min"] = float64(t.Min()) / du
				results[cleanName]["duration_max"] = float64(t.Max()) / du
				results[cleanName]["duration_stddev"] = t.StdDev() / du
				results[cleanName]["duration_median"] = ps[0] / du
				results[cleanName]["duration_p75"] = ps[1] / du
				results[cleanName]["duration_p95"] = ps[2] / du
				results[cleanName]["duration_p99"] = ps[3] / du
				results[cleanName]["duration_p999"] = ps[4] / du
			}
		}
	})
	return results

}

func outputMetrics(r metrics.Registry, scale time.Duration) map[string]map[string]float64 {
	du := float64(scale)
	//duSuffix := scale.String()[1:]
	results := make(map[string]map[string]float64)

	r.Each(func(name string, i interface{}) {
		if strings.HasPrefix(name, "output_") {
			switch metric := i.(type) {
			case metrics.Timer:
				cleanName := strings.SplitAfter(name, "_")[1]
				results[cleanName] = make(map[string]float64)
				t := metric.Snapshot()
				ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})

				results[cleanName]["count"] = float64(t.Count())
				results[cleanName]["rate1m"] = t.Rate1()
				results[cleanName]["rate5m"] = t.Rate5()
				results[cleanName]["rate15m"] = t.Rate15()
				results[cleanName]["ratemean"] = t.RateMean()
				results[cleanName]["duration_mean"] = t.Mean() / du
				results[cleanName]["duration_min"] = float64(t.Min()) / du
				results[cleanName]["duration_max"] = float64(t.Max()) / du
				results[cleanName]["duration_stddev"] = t.StdDev() / du
				results[cleanName]["duration_median"] = ps[0] / du
				results[cleanName]["duration_p75"] = ps[1] / du
				results[cleanName]["duration_p95"] = ps[2] / du
				results[cleanName]["duration_p99"] = ps[3] / du
				results[cleanName]["duration_p999"] = ps[4] / du
			}
		}
	})
	return results
}
