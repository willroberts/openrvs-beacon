package beacon

import "testing"

func TestParseServerReport(t *testing.T) {
	if _, err := ParseServerReport(testHost, testBeaconFixture); err != nil {
		t.Fatal(err)
	}
}
