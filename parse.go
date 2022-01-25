package beacon

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

const keySize = 2

var ErrInvalidLine = errors.New("invalid report line")

// ParseServerReport parses the bytes from the game server into a ServerReport.
func ParseServerReport(ip string, report []byte) (*ServerReport, error) {
	r := &ServerReport{IPAddress: ip}
	lines := bytes.Split(report, []byte{pilcrow})

	// Skip the first line containing the header.
	for _, line := range lines[1:] {
		if err := validateLine(line); err != nil {
			continue
		}

		// Parse the key and value. If there is no value, do nothing.
		key, value := lineToKeyVal(line)
		if len(value) == 0 {
			continue
		}

		// Decide which action to take based on the key.
		var err error
		switch key {
		case "A1":
			r.MaxPlayers, err = strconv.Atoi(value)
		case "B1":
			r.NumPlayers, err = strconv.Atoi(value)
		// C1 and D1 are unused.
		case "E1":
			r.CurrentMap = value
		case "F1":
			r.CurrentMode = value
		case "G1":
			r.Locked = valToBool(value)
		case "H1":
			r.Dedicated = valToBool(value)
		case "I1":
			r.ServerName = value
		case "J1":
			r.ModeRotation = parseModeRotation(value)
		case "K1":
			r.MapRotation = strings.Split(value, "/")[1:]
		case "L1":
			r.ConnectedPlayerNames = strings.Split(value, "/")[1:]
		case "M1":
			r.ConnectedPlayerTimes = strings.Split(value, "/")[1:]
		case "N1":
			r.ConnectedPlayerLatencies, err = valToIntSlice(value)
		case "O1":
			r.ConnectedPlayerKills, err = valToIntSlice(value)
		case "P1":
			r.Port, err = strconv.Atoi(value)
		case "Q1":
			r.RoundsPerMatch, err = strconv.Atoi(value)
		case "R1":
			r.TimePerRound, err = strconv.Atoi(value)
		case "S1":
			r.TimeBetweenRounds, err = strconv.Atoi(value)
		case "T1":
			r.BombTimer, err = strconv.Atoi(value)
		// U1 and V1 are unused.
		case "W1":
			r.TeamNamesVisible = valToBool(value)
		case "X1":
			r.InternetServer = valToBool(value)
		case "Y1":
			r.FriendlyFire = valToBool(value)
		case "Z1":
			r.AutoTeamBalance = valToBool(value)
		case "A2":
			r.TeamkillPenalty = valToBool(value)
		case "B2":
			r.RadarAllowed = valToBool(value)
		case "C2":
			r.OptionsList = value
		case "D2":
			r.GameVersion = value
		case "E2":
			r.LobbyServerID, err = strconv.Atoi(value)
		case "F2":
			r.GroupID, err = strconv.Atoi(value)
		case "G2":
			r.BeaconPort, err = strconv.Atoi(value)
		case "H2":
			r.NumTerrorists, err = strconv.Atoi(value)
		case "I2":
			r.AIBackup = valToBool(value)
		case "J2":
			r.RotateMapOnSuccess = valToBool(value)
		case "K2":
			r.ForceFirstPerson = valToBool(value)
		case "L2":
			r.ModName = value
		case "L3":
			r.PunkbusterEnabled = valToBool(value)
		// O2 through O9 are reserved for OpenRVS custom fields.
		case "O2":
			r.MOTD = value
		default:
			// Unused or invalid key.
		}

		// Some case statement branches can return an error.
		if err != nil {
			return &ServerReport{}, err
		}
	}

	return r, nil
}

func validateLine(line []byte) error {
	if len(line) < keySize {
		return ErrInvalidLine
	}
	return nil
}

// Removes traling whitespace and returns the key and value from the given line.
func lineToKeyVal(line []byte) (string, string) {
	key := string(line[0:keySize])
	val := string(bytes.Trim(line[keySize+1:], "\x20"))
	return key, val
}

// Some values (kill counts, latencies, etc.) are lists separated by forward slashes.
// This converts them to slices of integers.
func valToIntSlice(value string) ([]int, error) {
	ints := make([]int, strings.Count(value, "/"))
	for i, strVal := range strings.Split(value, "/")[1:] {
		intVal, err := strconv.Atoi(strVal)
		if err != nil {
			return []int{}, err
		}
		ints[i] = intVal
	}
	return ints, nil
}

func valToBool(value string) bool {
	return value == enabled
}

// Unlike the Map Rotation, the Mode Rotation always includes 32 fields, regardless of
// how many actually contain values. Each field is prefixed with '/'.
func parseModeRotation(value string) []string {
	modes := make([]string, 0)
	for _, m := range strings.Split(value, "/")[1:] {
		if m != "" {
			modes = append(modes, m)
		}
	}
	return modes
}
