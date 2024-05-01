package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/dcreager/dcreager.net/generate"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

const domain = "dcreager.net"
const outputDir = ".html"

func run() error {
	start := time.Now()
	overallCount := 0

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	count, err := generate.ProcessSourceDir(domain, path.Join(home, "notes"), outputDir)
	if err != nil {
		return err
	}
	overallCount += count

	count, err = generate.ProcessSourceDir(domain, "overrides", outputDir)
	if err != nil {
		return err
	}
	overallCount += count

	duration := time.Since(start)
	fmt.Printf("Processed %d files in %.2f seconds.\n", overallCount, duration.Seconds())
	return nil
}
