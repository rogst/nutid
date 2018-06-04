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

// Config holds cmdline parameters etc
type Config struct {
	Unix     uint64
	Add      time.Duration
	NoColors bool
}

// RegisterFlags loads cmdline params into config
func (c *Config) RegisterFlags(f *flag.FlagSet) {
	f.Uint64Var(&c.Unix, "unix", 0, "Convert Unix timestamp")
	f.DurationVar(&c.Add, "add", 0, "Add time duration")
	f.BoolVar(&c.NoColors, "no-colors", false, "Disable coloring of time")
}

func main() {
	var config Config
	config.RegisterFlags(flag.CommandLine)
	flag.Parse()

	now := time.Now()
	if config.Unix > 0 {
		now = time.Unix(int64(config.Unix), 0)
	}

	if config.Add != 0 {
		now = now.Add(config.Add)
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

	_, utcOffset := now.Zone()
	if fi.Mode()&os.ModeNamedPipe != 0 {
		// We got data over Stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Println(string(re.ReplaceAllFunc(scanner.Bytes(), func(b []byte) []byte {
				t := strings.Split(string(b), ":")
				h, _ := strconv.Atoi(t[0])
				h = h - utcOffset
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
			fmt.Printf("Local: %s\n", now.Local())
			fmt.Printf("  UTC: %s\n", now.UTC())
			fmt.Printf(" Unix: %d\n", now.Unix())
		} else {
			fmt.Printf("\033[1;37mLocal: %s\033[0m\n", now.Local())
			fmt.Printf("\033[1;32m  UTC: %s\033[0m\n", now.UTC())
			fmt.Printf("\033[1;36m Unix: %d\033[0m\n", now.Unix())
		}
	}
}
