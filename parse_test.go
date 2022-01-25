package beacon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testHost = "127.0.0.1"

func TestParseServerReport(t *testing.T) {
	r, err := ParseServerReport(testHost, testBeaconFixture)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure values are expected.
	assert.Equal(t, r.ServerName, "Classic Maps | Terrorist Hunt")
	assert.Equal(t, r.IPAddress, testHost)
	assert.Equal(t, r.Port, 6776)
	assert.Equal(t, r.BeaconPort, 7776)
	assert.Equal(t, r.InternetServer, true)
	assert.Equal(t, r.Dedicated, true)
	assert.Equal(t, r.PunkbusterEnabled, false)
	assert.Equal(t, r.Locked, false)
	assert.Equal(t, r.MaxPlayers, 8)
	assert.Equal(t, r.NumPlayers, 0) // TODO: Generate another fixture while connected.
	assert.Equal(t, r.GameVersion, "PATCH 1.60 (build 412)")
	assert.Equal(t, r.ModName, "RavenShield")
	assert.Equal(t, len(r.OptionsList), 0)
	assert.Equal(t, r.LobbyServerID, 0)
	assert.Equal(t, r.GroupID, 0)
	assert.Equal(t, r.AIBackup, false)
	assert.Equal(t, r.AutoTeamBalance, false)
	assert.Equal(t, r.BombTimer, 0)
	assert.Equal(t, len(r.ConnectedPlayerKills), 0)
	assert.Equal(t, len(r.ConnectedPlayerLatencies), 0)
	assert.Equal(t, len(r.ConnectedPlayerNames), 0)
	assert.Equal(t, len(r.ConnectedPlayerTimes), 0)
	assert.Equal(t, r.CurrentMap, "Streets")
	assert.Equal(t, r.CurrentMode, "RGM_TerroristHuntCoopMode")
	assert.Equal(t, r.ForceFirstPerson, true)
	assert.Equal(t, r.FriendlyFire, false)
	assert.Equal(t, len(r.MapRotation), 5)
	assert.Equal(t, len(r.ModeRotation), 5)
	assert.Equal(t, r.NumTerrorists, 35)
	assert.Equal(t, r.RadarAllowed, true)
	assert.Equal(t, r.RotateMapOnSuccess, true)
	assert.Equal(t, r.RoundsPerMatch, 5)
	assert.Equal(t, r.TeamNamesVisible, true)
	assert.Equal(t, r.TeamkillPenalty, false)
	assert.Equal(t, r.TimeBetweenRounds, 45)
	assert.Equal(t, r.TimePerRound, 900)
	assert.Equal(t, r.MOTD, "") // Not supplied, OpenRVS 1.5 only.
}

func TestValToIntSlice(t *testing.T) {
	value := "/1/2/3"
	expected := []int{1, 2, 3}
	actual, err := valToIntSlice(value)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, expected, actual)
}
