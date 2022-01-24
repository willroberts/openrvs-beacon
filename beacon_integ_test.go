//go:build integration

package beacon

import (
	"testing"
	"time"
)

var (
	testIntegHost = "64.225.54.237" // openrvs.org
	testPort      = 7776            // Classic Maps Terrorist Hunt Server
	testTimeout   = 5 * time.Second
)

func TestGetBeacon(t *testing.T) {
	if _, err := GetServerReport(testIntegHost, testPort, testTimeout); err != nil {
		t.Fatal(err)
	}
}
