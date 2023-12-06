package junit

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"testfmt/internal/outfmt"
	"testfmt/internal/result"
)

func init() {
	outfmt.Register("junit", &Formatter{})
}

type Formatter struct{}

func (f *Formatter) Format(dst io.Writer, result result.Result) error {
	ret := &Testsuites{}

	var totalDuration time.Duration
	var totalFailures int

	for _, r := range result.Tests {
		totalDuration += r.Duration
		testSuite := TestSuite{
			Name:     r.Name,
			Failures: "0",
			TestCases: []TestCase{
				{
					Classname: r.SourceFile,
					File:      r.SourceFile,
					Name:      r.Name,
					Time:      fmt.Sprintf("%f", r.Duration.Seconds()),
				},
			},
		}

		if r.Failed {
			totalFailures++
			testSuite.Failures = "1"
			testSuite.TestCases[0].Failure = &Failure{
				Message: "Test failed",
				Type:    "Failure",
				Text:    "\n" + r.Output + "\n",
			}
		}

		testSuite.Time = fmt.Sprintf("%f", r.Duration.Seconds())

		testSuite.Tests = fmt.Sprintf("%d", len(testSuite.TestCases))
		ret.TestsuiteEntries = append(ret.TestsuiteEntries, testSuite)
	}

	ret.Tests = fmt.Sprintf("%d", len(result.Tests))
	ret.Time = fmt.Sprintf("%f", totalDuration.Seconds())
	ret.Failures = fmt.Sprintf("%d", totalFailures)

	encoder := xml.NewEncoder(dst)
	encoder.Indent("", "  ")
	return encoder.Encode(ret)
}
