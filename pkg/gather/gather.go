package gather

import "context"

type Entry struct {
	DisplayPath string
	AbsPath     string
}

type Options struct {
	IncludeHidden bool
	Branch        string
}

func Collect(inputs []string, exts []string, opts Options) ([]Entry, error) {
	ctx := context.Background()
	entries, err := CollectConcurrent(ctx, inputs, exts, opts, 4)
	if err != nil {
		return nil, err
	}
	return entries, nil
}