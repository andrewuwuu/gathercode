package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gathercode/pkg/gather"
	"gathercode/pkg/writer"
)

func main() {
	var inputsArg string
	var output string
	var extsArg string
	var includeHidden bool
	var branch string
	var parallel int
	var timeoutStr string

	flag.StringVar(&inputsArg, "inputs", "", "comma-separated list of local paths or repo URLs")
	flag.StringVar(&output, "output", "aggregated.md", "output markdown filename")
	flag.StringVar(&extsArg, "ext", ".go,.sql", "comma-separated extensions to include")
	flag.BoolVar(&includeHidden, "include-hidden", false, "include hidden files and directories")
	flag.StringVar(&branch, "branch", "", "branch to use when fetching repos")
	flag.IntVar(&parallel, "parallel", 0, "number of concurrent workers (0 = default 4)")
	flag.StringVar(&timeoutStr, "timeout", "10m", "operation timeout (e.g. 30s, 2m, 10m)")
	flag.Parse()

	inputs := []string{}
	if inputsArg != "" {
		for _, p := range strings.Split(inputsArg, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				inputs = append(inputs, p)
			}
		}
	}
	inputs = append(inputs, flag.Args()...)

	if len(inputs) == 0 {
		fmt.Fprintln(os.Stderr, "error: provide inputs via --inputs or positional args")
		flag.Usage()
		os.Exit(2)
	}

	exts := []string{}
	for _, e := range strings.Split(extsArg, ",") {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}
		if !strings.HasPrefix(e, ".") {
			e = "." + e
		}
		exts = append(exts, strings.ToLower(e))
	}

	options := gather.Options{
		IncludeHidden: includeHidden,
		Branch:        branch,
	}

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid timeout: %v\n", err)
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var entries []gather.Entry
	if parallel > 0 {
		entries, err = gather.CollectConcurrent(ctx, inputs, exts, options, parallel)
	} else {
		entries, err = gather.Collect(inputs, exts, options)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error collecting files: %v\n", err)
		os.Exit(1)
	}
	if len(entries) == 0 {
		fmt.Fprintln(os.Stderr, "no files found matching extensions")
		os.Exit(2)
	}

	outDir := filepath.Dir(output)
	if outDir != "." && outDir != "" {
		_ = os.MkdirAll(outDir, 0o755)
	}

	if err := writer.WriteAggregatedMarkdown(entries, output); err != nil {
		fmt.Fprintf(os.Stderr, "error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Written %d files to %s\n", len(entries), output)
	cancel()
}
