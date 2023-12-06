package junit

import (
	"encoding/xml"
)

type Testsuites struct {
	XMLName          xml.Name    `xml:"testsuites"`
	Text             string      `xml:",chardata"`
	ID               string      `xml:"id,attr,omitempty"`
	Name             string      `xml:"name,attr,omitempty"`
	Tests            string      `xml:"tests,attr"`
	Failures         string      `xml:"failures,attr"`
	Time             string      `xml:"time,attr"`
	TestsuiteEntries []TestSuite `xml:"testsuite"`
}

type TestSuite struct {
	Text      string     `xml:",chardata"`
	ID        string     `xml:"id,attr,omitempty"`
	Name      string     `xml:"name,attr"`
	Tests     string     `xml:"tests,attr"`
	Failures  string     `xml:"failures,attr"`
	Time      string     `xml:"time,attr"`
	TestCases []TestCase `xml:"testcase"`
}

type TestCase struct {
	Classname string   `xml:"classname,attr,omitempty"`
	Text      string   `xml:",chardata"`
	ID        string   `xml:"id,attr,omitempty"`
	Name      string   `xml:"name,attr"`
	Time      string   `xml:"time,attr"`
	File      string   `xml:"file,attr,omitempty"`
	Failure   *Failure `xml:"failure,omitempty"`
}

type Failure struct {
	Text    string `xml:",chardata"`
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
}
