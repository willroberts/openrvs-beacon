package beacon

import (
	"testing"
)

var (
	testHost        = "64.225.54.237" // openrvs.org
	testReportBytes = []byte("")
)

func TestParseBeacon(t *testing.T) {
	if _, err := ParseServerReport(testHost, testReportBytes); err != nil {
		t.Fatal(err)
	}
}
