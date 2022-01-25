//go:build integration

package beacon

import (
	"testing"
	"time"
)

var (
	testIntegHost   = "64.225.54.237" // openrvs.org
	testInvalidHost = "127.0.0.1"
	testPort        = 7776 // Classic Maps Terrorist Hunt Server
	testTimeout     = 1 * time.Second
)

func TestGetBeacon(t *testing.T) {
	if _, err := GetServerReport(testIntegHost, testPort, testTimeout); err != nil {
		t.Fatal(err)
	}
}

func TestGetBeaconFailure(t *testing.T) {
	if _, err := GetServerReport(testInvalidHost, testPort, testTimeout); err == nil {
		t.Fail()
	}
}
