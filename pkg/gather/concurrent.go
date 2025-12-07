package gather

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gathercode/pkg/github"
	"gathercode/pkg/utils"
)

func CollectConcurrent(ctx context.Context, inputs []string, exts []string, opts Options, maxWorkers int) ([]Entry, error) {
	if maxWorkers <= 0 {
		maxWorkers = 4
	}
	extSet := make(map[string]bool)
	for _, e := range exts {
		extSet[strings.ToLower(e)] = true
	}
	type result struct {
		entries []Entry
		err     error
	}
	inCh := make(chan string)
	outCh := make(chan result)
	var wg sync.WaitGroup
	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case inp, ok := <-inCh:
				if !ok {
					return
				}
				var collected []Entry
				if utils.PathExists(inp) {
					fi, err := os.Stat(inp)
					if err != nil {
						outCh <- result{nil, nil}
						continue
					}
					if fi.IsDir() {
						ents, err := utils.GatherFilesFromRoot(inp, extSet, opts.IncludeHidden)
						if err == nil {
							for _, e := range ents {
								collected = append(collected, Entry{DisplayPath: e.DisplayPath, AbsPath: e.AbsPath})
							}
						}
						outCh <- result{collected, nil}
						continue
					}
					ext := strings.ToLower(filepath.Ext(fi.Name()))
					if extSet[ext] {
						display := filepath.ToSlash(filepath.Base(inp))
						abs, _ := filepath.Abs(inp)
						collected = append(collected, Entry{DisplayPath: display, AbsPath: abs})
					}
					outCh <- result{collected, nil}
					continue
				}
				if github.IsRepoURL(inp) {
					repoRoot, repoFolder, err := github.FetchRepo(inp, opts.Branch)
					if err != nil {
						outCh <- result{nil, nil}
						continue
					}
					ents, err := utils.GatherFilesFromRoot(repoRoot, extSet, opts.IncludeHidden)
					if err == nil {
						for _, e := range ents {
							parts := strings.SplitN(e.DisplayPath, "/", 2)
							var rest string
							if len(parts) == 2 {
								rest = parts[1]
							} else {
								rest = parts[0]
							}
							display := filepath.ToSlash(filepath.Join(repoFolder, rest))
							collected = append(collected, Entry{DisplayPath: display, AbsPath: e.AbsPath})
						}
					}
					_ = os.RemoveAll(repoRoot)
					outCh <- result{collected, nil}
					continue
				}
				outCh <- result{nil, nil}
			}
		}
	}
	wg.Add(maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go worker()
	}
	go func() {
		defer close(inCh)
		for _, inp := range inputs {
			select {
			case <-ctx.Done():
				return
			case inCh <- inp:
			}
		}
	}()
	go func() {
		wg.Wait()
		close(outCh)
	}()
	var final []Entry
	for {
		select {
		case <-ctx.Done():
			return final, ctx.Err()
		case res, ok := <-outCh:
			if !ok {
				return final, nil
			}
			if res.err != nil {
				continue
			}
			if len(res.entries) > 0 {
				final = append(final, res.entries...)
			}
		}
	}
}
