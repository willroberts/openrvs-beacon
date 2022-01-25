package beacon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseServerReport(t *testing.T) {
	r, err := ParseServerReport(testHost, testBeaconFixture)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure values are expected.
	assert.Equal(t, r.ServerName, "Classic Maps | Terrorist Hunt", "failed to parse ServerName")
	assert.Equal(t, r.IPAddress, "127.0.0.1", "failed to parse IPAddress")
	assert.Equal(t, r.Port, 6776, "failed to parse Port")
	assert.Equal(t, r.BeaconPort, 7776, "failed to parse BeaconPort")
	assert.Equal(t, r.InternetServer, true, "failed to parse InternetServer")
	assert.Equal(t, r.Dedicated, true, "failed to parse Dedicated")
	assert.Equal(t, r.PunkbusterEnabled, false, "failed to parse PunkbusterEnabled")
	assert.Equal(t, r.Locked, false, "failed to parse Locked")
	assert.Equal(t, r.MaxPlayers, 8, "failed to parse MaxPlayers")
	assert.Equal(t, r.NumPlayers, 0, "failed to parse NumPlayers")
	assert.Equal(t, r.GameVersion, "PATCH 1.60 (build 412)", "failed to parse GameVersion")
	assert.Equal(t, r.ModName, "RavenShield", "failed to parse ModName")
	assert.Equal(t, len(r.OptionsList), 0, "failed to parse OptionsList")
	assert.Equal(t, r.LobbyServerID, 0, "failed to parse LobbyServerID")
	assert.Equal(t, r.GroupID, 0, "failed to parse GroupID")
	assert.Equal(t, r.AIBackup, false, "failed to parse AIBackup")
	assert.Equal(t, r.AutoTeamBalance, false, "failed to parse AutoTeamBalance")
	assert.Equal(t, r.BombTimer, 0, "failed to parse BombTimer")
	assert.Equal(t, len(r.ConnectedPlayerKills), 0, "failed to parse ConnectedPlayerKills")
	assert.Equal(t, len(r.ConnectedPlayerLatencies), 0, "failed to parse ConnectedPlayerLatencies")
	assert.Equal(t, len(r.ConnectedPlayerNames), 0, "failed to parse ConnectedPlayerNames")
	assert.Equal(t, len(r.ConnectedPlayerTimes), 0, "failed to parse ConnectedPlayerTimes")
	assert.Equal(t, r.CurrentMap, "Streets", "failed to parse CurrentMap")
	assert.Equal(t, r.CurrentMode, "RGM_TerroristHuntCoopMode", "failed to parse CurrentMode")
	assert.Equal(t, r.ForceFirstPerson, true, "failed to parse ForceFirstPerson")
	assert.Equal(t, r.FriendlyFire, false, "failed to parse FriendlyFire")
	assert.Equal(t, len(r.MapRotation), 5, "failed to parse MapRotation")
	assert.Equal(t, len(r.ModeRotation), 5, "failed to parse ModeRotation")
	assert.Equal(t, r.NumTerrorists, 35, "failed to parse NumTerrorists")
	assert.Equal(t, r.RadarAllowed, true, "failed to parse RadarAllowed")
	assert.Equal(t, r.RotateMapOnSuccess, true, "failed to parse RotateMapOnSuccess")
	assert.Equal(t, r.RoundsPerMatch, 5, "failed to parse RoundsPerMatch")
	assert.Equal(t, r.TeamNamesVisible, true, "failed to parse TeamNamesVisible")
	assert.Equal(t, r.TeamkillPenalty, false, "failed to parse TeamkillPenalty")
	assert.Equal(t, r.TimeBetweenRounds, 45, "failed to parse TimeBetweenRounds")
	assert.Equal(t, r.TimePerRound, 900, "failed to parse TimePerRound")
	assert.Equal(t, r.MOTD, "", "failed to parse MOTD") // Not supplied, OpenRVS 1.5 only.
}
