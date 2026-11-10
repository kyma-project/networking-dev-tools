package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bojand/ghz/runner"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"

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

	sftpHost       string
	sshUser        string
	sshPass        string
	sshKeyPath     string
	knownHostsPath string
)

func defaultKnownHostsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".ssh", "known_hosts")
}

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
	flag.StringVar(&sshPass, "sshPass", "load", "The password for the sftp server. Ignored when sshKeyPath is set.")
	flag.StringVar(&sshKeyPath, "sshKeyPath", "", "Path to a private key for the sftp server. Takes precedence over sshPass when set.")
	flag.StringVar(&knownHostsPath, "knownHostsPath", defaultKnownHostsPath(), "Path to the known_hosts file used to verify the sftp server host key")

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
	if knownHostsPath == "" {
		return nil, fmt.Errorf("knownHostsPath is empty: cannot verify host key, refusing to connect")
	}

	hostKeyCallback, err := knownhosts.New(knownHostsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load known_hosts from %q: %w", knownHostsPath, err)
	}

	authMethod, err := sshAuthMethod()
	if err != nil {
		return nil, err
	}

	sshConfig := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: hostKeyCallback,
	}

	sshClient, err := ssh.Dial("tcp", sftpHost, sshConfig)
	if err != nil {
		return nil, err
	}

	return sshClient, nil
}

func sshAuthMethod() (ssh.AuthMethod, error) {
	if sshKeyPath != "" {
		key, err := os.ReadFile(sshKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key %q: %w", sshKeyPath, err)
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key %q: %w", sshKeyPath, err)
		}
		return ssh.PublicKeys(signer), nil
	}
	return ssh.Password(sshPass), nil
}
