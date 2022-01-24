//go:build integration

package beacon

import (
	"testing"
	"time"
)

var (
	testPort    = 7776
	testTimeout = 5 * time.Second
)

func TestGetBeacon(t *testing.T) {
	_, err := GetServerReport(testHost, testPort, testTimeout)
	if err != nil {
		t.Fatal(err)
	}
}
