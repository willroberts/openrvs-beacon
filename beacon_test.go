package beacon

import (
	"testing"
)

var (
	testHost            = "127.0.0.1"
	testBeaconFixture   = []byte("rvnshld 6776 KEYWORD  �P1 6776 �E1 Streets �I1 Classic Maps | Terrorist Hunt �F1 RGM_TerroristHuntCoopMode �A1 8 �G1 0 �H1 1 �L1  �M1  �N1  �O1  �B1 0 �Q1 5 �R1 900 �S1 45 �W1 1 �X1 1 �Y1 0 �Z1 0 �A2 0 �D2 PATCH 1.60 (build 412) �B2 1 �E2 0 �F2 0 �G2 7776 �H2 35 �I2 0 �J2 1 �K2 1 �L2 RavenShield �L3 0 �K1 /Streets/Training/Island_Dawn/Import_Export/Prison �J1 /RGM_TerroristHuntCoopMode/RGM_TerroristHuntCoopMode/RGM_TerroristHuntCoopMode/RGM_TerroristHuntCoopMode/RGM_TerroristHuntCoopMode///////////////////////////")
	testTooSmallFixture = []byte("")
	testNoHeaderFixture = []byte("1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890")
)

func TestParseBeacon(t *testing.T) {
	if _, err := ParseServerReport(testHost, testBeaconFixture); err != nil {
		t.Fatal(err)
	}
}

func TestValidateSuccess(t *testing.T) {
	if err := validateServerReport(testBeaconFixture); err != nil {
		t.Fatal(err)
	}
}
func TestValidateReportTooSmall(t *testing.T) {
	if err := validateServerReport(testTooSmallFixture); err != ErrNotABeacon {
		t.Fatal(err)
	}
}
func TestValidateInvalidHeader(t *testing.T) {
	if err := validateServerReport(testNoHeaderFixture); err != ErrNotABeacon {
		t.Fatal(err)
	}
}
