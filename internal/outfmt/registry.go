package outfmt

import (
	"errors"
	"io"

	"testfmt/internal/result"
)

type Formatter interface {
	Format(dst io.Writer, result result.Result) error
}

type FormatterFunc func(dst io.Writer, result result.Result) error

var (
	registry = map[string]Formatter{}
)

func Register(name string, formatter Formatter) {
	registry[name] = formatter
}

func Get(name string) (Formatter, error) {
	formatter, found := registry[name]
	if !found {
		return nil, errors.New("formatter not found")
	}
	return formatter, nil
}

func ListFormatters() []string {
	var names []string
	for name := range registry {
		names = append(names, name)
	}
	return names
}
