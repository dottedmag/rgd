package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	goroutineHeaderRx = regexp.MustCompile(`^goroutine (\d+) (.*) (\[.*\]):$`)
	homebrewRx        = regexp.MustCompile(`^\t/opt/homebrew/Cellar/go/(.*)/libexec/src(.*)$`)
	funcRx            = regexp.MustCompile(`^([a-zA-Z0-9-_.]+/)?.*\(.*\)$`)
)

func formatDump(homeDir string, out io.Writer, in io.Reader) error {
	modCacheDir := homeDir + "/go/pkg/mod"

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := strings.TrimSuffix(scanner.Text(), "\n")
		switch {
		case strings.HasPrefix(line, "goroutine "):
			m := goroutineHeaderRx.FindStringSubmatch(line)
			if m == nil {
				return fmt.Errorf("failed to parse %q as a goroutine header", line)
			}
			fmt.Fprintf(out, "%-4s %s %s:\n", m[1], m[3], m[2])
		case strings.HasPrefix(line, "\t"):
			if strings.HasPrefix(line, "\t"+modCacheDir) {
				fmt.Fprintln(out, "    +CACHE+"+strings.TrimPrefix(line, "\t"+modCacheDir))
				continue
			}
			if strings.HasPrefix(line, "\t"+homeDir) {
				fmt.Fprintln(out, "    <HOME>"+strings.TrimPrefix(line, "\t"+homeDir))
				continue
			}
			if m := homebrewRx.FindStringSubmatch(line); m != nil {
				fmt.Fprintln(out, "    -GO "+m[1]+"-"+m[2])
				continue
			}
			fmt.Fprintln(out, "    "+strings.TrimPrefix(line, "\t"))
		case strings.HasPrefix(line, "created by "):
			fmt.Fprintln(out, "  ."+line)
		default:
			if m := funcRx.FindStringSubmatch(line); m != nil {
				// not a stdlib, first component is DNS name, or it's a main package
				if strings.HasPrefix(line, "main.") || strings.Contains(m[1], ".") {
					fmt.Fprintln(out, "  "+line)
				} else {
					fmt.Fprintln(out, "  -"+line)
				}
				continue
			}
			fmt.Fprintln(out, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: rgd <file>|-\nMakes Go goroutine dump more human-readable\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
	}

	inputFile := flag.Arg(0)
	var input io.Reader
	if inputFile == "-" {
		input = os.Stdin
	} else {
		fh, err := os.Open(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open %q: %v\n", inputFile, err)
			os.Exit(1)
		}
		input = fh
	}

	homeDir := os.Getenv("HOME")

	if err := formatDump(homeDir, os.Stdout, input); err != nil {
		fmt.Fprintf(os.Stderr, "failed to format goroutines dump: %v\n", err)
		os.Exit(1)
	}
}
