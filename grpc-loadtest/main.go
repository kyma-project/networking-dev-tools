package main

import (
	"flag"
	"fmt"
	"github.com/bojand/ghz/runner"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
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

	authority string

	htmlReportPath string

	sftpHost string
	sshUser  string
	sshPass  string
)

func init() {
	flag.StringVar(&service, "service", "HelloService.SayHello", "The service and method to call")
	flag.StringVar(&address, "address", "localhost:50051", "The address of the service host and port")
	flag.BoolVar(&insecure, "insecure", true, "Use an insecure connection")
	flag.UintVar(&requests, "requests", 10000, "The number of requests to send")
	flag.IntVar(&duration, "duration", 0, "The duration in seconds to send requests. If different than 0, requests will be ignored")
	flag.UintVar(&concurrency, "concurrency", 1, "The number of requests to run concurrently")
	flag.UintVar(&rps, "rps", 0, "The target requests per second")

	flag.StringVar(&authority, "authority", "", "The authority pseudo-header to use in the request")

	flag.StringVar(&htmlReportPath, "htmlReportPath", "/tmp/report.html", "The path to the html report")

	flag.StringVar(&sftpHost, "sftpHost", "", "The address of the sftp host and port. If not provided, the results will not be uploaded to the sftp server.")
	flag.StringVar(&sshUser, "sshUser", "goat", "The username for the sftp server")
	flag.StringVar(&sshPass, "sshPass", "load", "The password for the sftp server")

	flag.Parse()
}

func main() {
	var report *runner.Report
	if duration == 0 {
		fmt.Printf("Running test of %s for %d requests with %d concurrency\n", service, requests, concurrency)
		r, err := loadTestRunner.RunTestForRequestNumber(service, authority, address, insecure, requests, concurrency, rps)
		if err != nil {
			panic(err)
		}
		report = r
	} else {
		fmt.Printf("Running test of %s for %d seconds with %d concurrency\n", service, duration, concurrency)
		r, err := loadTestRunner.RunTestForDuration(service, authority, address, insecure, time.Duration(duration)*time.Second, concurrency, rps)
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

	var printerHtml printer.ReportPrinter
	if sftpHost != "" {
		sshClient, err := DialSSH(sftpHost)
		if err != nil {
			panic(err)
		}

		sftpClient, err := sftp.NewClient(sshClient)
		if err != nil {
			panic(err)
		}

		fileHandle, err := sftpClient.OpenFile(htmlReportPath, os.O_CREATE|os.O_WRONLY)
		if err != nil {
			panic(err)
		}

		printerHtml = printer.ReportPrinter{
			Out:    fileHandle,
			Report: report,
		}
	} else {
		file, err := os.OpenFile(htmlReportPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err.Error())
		}

		printerHtml = printer.ReportPrinter{
			Out:    file,
			Report: report,
		}
	}

	err = printerHtml.Print("html")
	if err != nil {
		panic(err)
	}
}

func DialSSH(sftpHost string) (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshPass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", sftpHost, sshConfig)
	if err != nil {
		return nil, err
	}

	return sshClient, nil
}
