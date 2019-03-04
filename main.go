package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type output struct {
	Status        bool     `json:"status"`
	Summary       []entry  `json:"summary"`
	AddedFiles    []string `json:"added_files"`
	RemovedFiles  []string `json:"removed_files"`
	ModifiedFiles []string `json:"modified_files"`
}

type entry struct {
	Name     string `json:"name"`
	Severity int    `json:"severity"`
	Added    int    `json:"added"`
	Removed  int    `json:"removed"`
	Modified int    `json:"modified"`
}

var (
	flagFile   = flag.String("file", "", "Report file to read. If this is not set, stdin will be used")
	flagPretty = flag.Bool("pretty", false, "Pretty print the output")
)

func init() {
	flag.Parse()
}

func main() {
	var input io.Reader

	if *flagFile != "" {
		var err error
		input, err = os.Open(*flagFile)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		input = bufio.NewReader(os.Stdin)
	}

	scanner := bufio.NewScanner(input)

	result := output{
		Status:        true,
		Summary:       []entry{},
		AddedFiles:    []string{},
		RemovedFiles:  []string{},
		ModifiedFiles: []string{},
	}

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "Added:"):
			appendFiles(scanner, &result.AddedFiles)
			break
		case strings.HasPrefix(line, "Removed:"):
			appendFiles(scanner, &result.RemovedFiles)
			break
		case strings.HasPrefix(line, "Modified:"):
			appendFiles(scanner, &result.ModifiedFiles)
			break
		default:
			processMiscLine(line, scanner, &result)
			break
		}
	}

	if len(result.AddedFiles) > 0 ||
		len(result.AddedFiles) > 0 ||
		len(result.ModifiedFiles) > 0 ||
		len(result.Summary) > 0 {
		result.Status = false
	}

	encoder := json.NewEncoder(os.Stdout)
	if *flagPretty {
		encoder.SetIndent("", "  ")
	}
	encoder.Encode(result)
}

func appendFiles(scanner *bufio.Scanner, slice *[]string) {
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		*slice = append(*slice, strings.Trim(line, "\""))
	}
}

func processMiscLine(line string, scanner *bufio.Scanner, result *output) {
	rgx := regexp.MustCompile(`^(\*)?\s+((?:\S+\s)*(?:\S+))\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)`)
	fields := rgx.FindStringSubmatch(line)
	if len(fields) < 6 {
		return
	}
	if fields[1] == "*" {
		severity, err := strconv.Atoi(fields[3])
		added, err := strconv.Atoi(fields[4])
		removed, err := strconv.Atoi(fields[5])
		modified, err := strconv.Atoi(fields[6])
		if err != nil {
			log.Fatal(err)
		}
		result.Summary = append(result.Summary, entry{
			Name:     fields[2],
			Severity: severity,
			Added:    added,
			Removed:  removed,
			Modified: modified,
		})
	}
}
