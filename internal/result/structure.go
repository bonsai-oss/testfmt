package result

import (
	"time"
)

type Result struct {
	Tests []Test
}

type Test struct {
	Name       string
	SourceFile string
	Output     string
	Duration   time.Duration
	Failed     bool
}
