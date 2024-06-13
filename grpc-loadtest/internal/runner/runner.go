package runner

import (
	"github.com/bojand/ghz/runner"

	"time"
)

func RunTestForDuration(service string, address string, insecure bool, duration time.Duration, concurrency uint) (*runner.Report, error) {
	run, err := runner.Run(
		service,
		address,
		runner.WithProtoFile("hello.proto", []string{}),
		runner.WithData(map[string]interface{}{
			"name": "world",
		}),
		runner.WithInsecure(insecure),
		runner.WithRunDuration(duration),
		runner.WithConcurrency(concurrency),
	)

	return run, err
}

func RunTestForRequestNumber(service string, address string, insecure bool, requests uint, concurrency uint) (*runner.Report, error) {
	run, err := runner.Run(
		service,
		address,
		runner.WithProtoFile("hello.proto", []string{}),
		runner.WithData(map[string]interface{}{
			"name": "world",
		}),
		runner.WithInsecure(insecure),
		runner.WithTotalRequests(requests),
		runner.WithConcurrency(concurrency),
	)

	return run, err
}
