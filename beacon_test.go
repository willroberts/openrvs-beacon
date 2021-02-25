package beacon

import (
	"testing"
	"time"
)

var (
	// TODO: Use fixtures.
	testHost    = "64.225.54.237"
	testPort    = 6776 + 1000 // Beacon port, not game server port!
	testTimeout = 5 * time.Second
)

func TestGetBeacon(t *testing.T) {
	_, err := GetServerReport(testHost, testPort, testTimeout)
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseBeacon(t *testing.T) {
	reportBytes, err := GetServerReport(testHost, testPort, testTimeout)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ParseServerReport(testHost, reportBytes)
	if err != nil {
		t.Fatal(err)
	}
}
