package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TODO, capture stdin and translate timestamps to localtime

type Config struct {
	UTCOffset int
	NoColors  bool
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	f.IntVar(&c.UTCOffset, "utc-offset", 0, "UTC Offset in hours (default 0 = auto-detect)")
	f.BoolVar(&c.NoColors, "no-colors", false, "Disable coloring of time")
}

func main() {
	var config Config
	config.RegisterFlags(flag.CommandLine)
	flag.Parse()

	now := time.Now()
	if config.UTCOffset == 0 {
		_, offset := now.Zone()
		config.UTCOffset = offset / 3600
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	re, err := regexp.Compile("([0-5]{1}[0-9]{1}:[0-5]{1}[0-9]{1}:[0-5]{1}[0-9]{1})")
	if err != nil {
		log.Fatal(err)
	}

	outputTmpl := "\033[1;32m%d:%s:%s\033[0m"
	if config.NoColors {
		outputTmpl = "%d:%s:%s"
	}

	if fi.Mode()&os.ModeNamedPipe != 0 {
		// We got data over Stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			//fmt.Println(scanner.Text())
			fmt.Println(string(re.ReplaceAllFunc(scanner.Bytes(), func(b []byte) []byte {
				t := strings.Split(string(b), ":")
				h, _ := strconv.Atoi(t[0])
				h = h - config.UTCOffset
				if h < 0 {
					h = h + 24
				}

				return []byte(fmt.Sprintf(outputTmpl, h, t[1], t[2]))
			})))
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	} else {
		if config.NoColors {
			fmt.Printf("Local: %s\n", now)
			fmt.Printf("  UTC: %s\n", now.UTC())
		} else {
			fmt.Printf("\033[1;37mLocal: %s\033[0m\n", now)
			fmt.Printf("\033[1;32m  UTC: %s\033[0m\n", now.UTC())
		}
	}
}
