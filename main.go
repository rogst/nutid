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
)

// TODO, capture stdin and translate timestamps to localtime

type Config struct {
	UTCOffset int
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	f.IntVar(&c.UTCOffset, "utc-offset", 0, "UTC Offset in hours")
}

func main() {
	var config Config
	config.RegisterFlags(flag.CommandLine)
	flag.Parse()

	fi, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	re, err := regexp.Compile("([0-5]{1}[0-9]{1}:[0-5]{1}[0-9]{1}:[0-5]{1}[0-9]{1})")
	if err != nil {
		log.Fatal(err)
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

				return []byte(fmt.Sprintf("%d:%s:%s", h, t[1], t[2]))
			})))
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}
