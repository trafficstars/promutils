package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/prometheus/common/model"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
)

const (
	dateTimeLayout = "2006-01-02 15:04:05"
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	urlString := flag.String("url", "http://127.0.0.1:9090/", `URL to the prometheus HTTP interface`)
	startString := flag.String("start", "2019-06-21 13:00:00", `start date-time, format: YYYY-MM-DD HH:MM:SS`)
	endString := flag.String("end", "", `end date-time, format: YYYY-MM-DD HH:MM:SS`)
	step := flag.Duration("step", 0, `step`)
	flag.Parse()

	prometheusClient, err := api.NewClient(api.Config{
		Address:      *urlString,
	})

	panicIf(err)
	prometheus := v1.NewAPI(prometheusClient)

	start, err := time.ParseInLocation(dateTimeLayout, *startString, time.UTC)
	panicIf(err)

	var end time.Time
	if *endString != "" {
		var err error
		end, err = time.ParseInLocation(dateTimeLayout, *endString, time.UTC)
		panicIf(err)
	} else {
		end = time.Now()
	}


	if *step == 0 {
		*step = end.Sub(start)
	}


	value, warnings, err := prometheus.QueryRange(context.Background(), strings.Join(flag.Args(), " "), v1.Range{
		Start: start,
		End:   end,
		Step: *step,
	})
	if err != nil {
		fmt.Println(err)
	}

	switch v := value.(type) {
	case model.Matrix:
		for _, stream := range v {
			fmt.Println("stream, metric", stream.Metric)
			for _, row := range stream.Values {
				fmt.Println("row", row.Timestamp, row.Value)
			}
		}
		fmt.Println(warnings, err)
	default:
		fmt.Println(value, warnings, err, value.Type())
	}

}
