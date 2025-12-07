package writer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gathercode/pkg/gather"
)

func WriteAggregatedMarkdown(entries []gather.Entry, outPath string) error {
	sort.SliceStable(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].DisplayPath) < strings.ToLower(entries[j].DisplayPath)
	})
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	sep := "---------------"
	for i, e := range entries {
		if _, err := outFile.WriteString("// " + e.DisplayPath + "\n"); err != nil {
			return err
		}
		src, err := os.Open(e.AbsPath)
		if err != nil {
			_, _ = outFile.WriteString(fmt.Sprintf("/* ERROR opening file: %v */\n", err))
		} else {
			_, err = io.Copy(outFile, src)
			src.Close()
			if err != nil {
				return err
			}
			if _, err := outFile.WriteString("\n"); err != nil {
				return err
			}
		}
		if _, err := outFile.WriteString(sep + "\n"); err != nil {
			return err
		}
		if i < len(entries)-1 {
			if _, err := outFile.WriteString("\n"); err != nil {
				return err
			}
		}
	}
	absOut, _ := filepath.Abs(outPath)
	fmt.Printf("wrote %d entries to %s\n", len(entries), absOut)
	return nil
}