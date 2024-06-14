package main

import (
	"flag"
	"fmt"
	"github.com/bojand/ghz/runner"
	"os"
	"time"

	"github.com/bojand/ghz/printer"
	loadTestRunner "github.com/kyma-project/networking-dev-tools/grpc-loadtest/internal/runner"
)

var (
	service     string
	address     string
	insecure    bool
	requests    uint
	concurrency uint
	duration    int
	rps         uint
)

func init() {
	flag.StringVar(&service, "service", "HelloService.SayHello", "The service and method to call")
	flag.StringVar(&address, "address", "localhost:50051", "The address of the service host and port")
	flag.BoolVar(&insecure, "insecure", true, "Use an insecure connection")
	flag.UintVar(&requests, "number", 10000, "The number of requests to send")
	flag.IntVar(&duration, "duration", 0, "The duration in seconds to send requests. If different than 0, requests will be ignored")
	flag.UintVar(&concurrency, "concurrency", 1, "The number of requests to run concurrently")
	flag.UintVar(&rps, "rps", 0, "The target requests per second")

	flag.Parse()
}

func main() {
	var report *runner.Report
	if duration == 0 {
		fmt.Printf("Running test of %s for %d requests with %d concurrency\n", service, requests, concurrency)
		r, err := loadTestRunner.RunTestForRequestNumber(service, address, insecure, requests, concurrency, rps)
		if err != nil {
			panic(err)
		}
		report = r
	} else {
		fmt.Printf("Running test of %s for %d seconds with %d concurrency\n", service, duration, concurrency)
		r, err := loadTestRunner.RunTestForDuration(service, address, insecure, time.Duration(duration)*time.Second, concurrency, rps)
		if err != nil {
			panic(err)
		}
		report = r
	}

	printerStdout := printer.ReportPrinter{
		Out:    os.Stdout,
		Report: report,
	}

	err := printerStdout.Print("summary")
	if err != nil {
		panic(err)
	}
	file, err := os.OpenFile("/tmp/report.html", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}

	printerHtml := printer.ReportPrinter{
		Out:    file,
		Report: report,
	}

	err = printerHtml.Print("html")
	if err != nil {
		panic(err)
	}
}
