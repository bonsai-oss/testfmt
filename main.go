package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/fatih/color"

	"testfmt/internal/outfmt"
	_ "testfmt/internal/outfmt/junit"
	"testfmt/internal/result"
)

type prefixWriter struct {
	prefix string
	w      io.Writer
}

func (pw *prefixWriter) Write(p []byte) (n int, err error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	for scanner.Scan() {
		line := scanner.Text()
		_, err := pw.w.Write([]byte(pw.prefix + line + "\n"))
		if err != nil {
			return 0, err
		}
	}
	return len(p), scanner.Err()
}

func listExecutables(directory string) ([]string, error) {
	var executables []string
	walkError := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() && info.Mode().Perm()&0100 != 0 {
			executables = append(executables, path)
		}
		return nil
	})
	if walkError != nil {
		return []string{}, walkError
	}
	return executables, nil
}

func init() {
	color.NoColor = false
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Lmsgprefix)
}

func main() {
	globalStart := time.Now()
	globalIsFailed := false
	defer func() {
		if globalIsFailed {
			os.Exit(1)
		}
	}()

	defer func() {
		color.New(color.Bold, color.FgGreen).Println(fmt.Sprintf("Total time: %s", time.Since(globalStart)))
	}()

	formatters := outfmt.ListFormatters()
	formatters = append(formatters, "none")
	app := kingpin.New("resultfmt", "Test code formatter")
	app.HelpFlag.Short('h')
	dir := app.Flag("dir", "Directory to search").Short('d').Default("tests").ExistingDir()
	formatterName := app.Flag("format", "Output format").Short('f').Default("none").Enum(formatters...)
	outputFilePath := app.Flag("output", "Output file path; only used on --format != none").Short('o').Default("output.xml").String()

	kingpin.MustParse(app.Parse(os.Args[1:]))

	executables, findExecutablesError := listExecutables(*dir)
	if findExecutablesError != nil || len(executables) == 0 {
		kingpin.Fatalf("Error finding executables: %s", findExecutablesError)
	}

	results := result.Result{}

	for _, executable := range executables {
		start := time.Now()
		isFailed := false
		color.New(color.Bold, color.FgHiWhite).Println(executable)
		output := bytes.NewBuffer(nil)
		prefixStderr := &prefixWriter{"├─ ", os.Stderr}
		cmd := exec.Command(executable)
		cmd.Stdout = io.MultiWriter(output, prefixStderr)
		cmd.Stderr = io.MultiWriter(output, prefixStderr)
		runExecutableError := cmd.Run()
		if runExecutableError != nil {
			isFailed = true
			globalIsFailed = true
			color.New(color.Bold, color.FgRed).Println("├─ ", runExecutableError)
		}

		// flush the output
		prefixStderr.Write(nil)

		log.Print("└─ ", color.New(color.Bold, color.FgHiWhite).Sprintf("Runtime: %s", time.Since(start)))
		results.Tests = append(results.Tests, result.Test{
			Name:       path.Base(executable),
			SourceFile: executable,
			Output:     output.String(),
			Duration:   time.Since(start),
			Failed:     isFailed,
		})
	}

	color.New(color.Bold, color.FgHiWhite).Println("-----------\nResults:")
	for _, r := range results.Tests {
		if r.Failed {
			color.New(color.Bold, color.FgRed).Print("✘\t")
		} else {
			color.New(color.Bold, color.FgGreen).Print("✔\t")
		}
		color.New(color.FgHiWhite).Println(r.SourceFile)
	}

	if *formatterName == "none" {
		return
	}

	outputFile, createOutputFileError := os.Create(*outputFilePath)
	if createOutputFileError != nil {
		kingpin.Fatalf("Error creating output file: %s", createOutputFileError)
	}
	defer outputFile.Close()

	formatter, getFormatterError := outfmt.Get(*formatterName)
	if getFormatterError != nil {
		kingpin.Fatalf("Error getting formatter: %s", getFormatterError)
	}
	formatError := formatter.Format(outputFile, results)
	if formatError != nil {
		kingpin.Fatalf("Error formatting: %s", formatError)
	}
}
