package beacon

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

const keySize = 2

var ErrInvalidLine = errors.New("invalid report line")

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

// ParseServerReport parses the bytes from the game server into a ServerReport.
func ParseServerReport(ip string, report []byte) (*ServerReport, error) {
	r := &ServerReport{IPAddress: ip}
	lines := bytes.Split(report, []byte{pilcrow})

	// Skip the first line containing the header.
	for _, line := range lines[1:] {
		if err := validateLine(line); err != nil {
			return r, err
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
			if value == enabled {
				r.Locked = true
			}
		case "H1":
			if value == enabled {
				r.Dedicated = true
			}
		case "I1":
			r.ServerName = value
		case "J1": // Mode Rotation.
			// Unlike the Map Rotation, the Mode Rotation always includes 32 fields, regardless of
			// how many actually contain values. Each field is prefixed with '/'.
			modes := make([]string, 0)
			for _, m := range strings.Split(value, "/")[1:] {
				if m != "" {
					modes = append(modes, m)
				}
			}
			r.ModeRotation = modes
		case "K1":
			r.MapRotation = strings.Split(value, "/")[1:]
		case "L1":
			r.ConnectedPlayerNames = strings.Split(value, "/")[1:]
		case "M1":
			r.ConnectedPlayerTimes = strings.Split(value, "/")[1:]
		case "N1":
			in := strings.Split(value, "/")[1:]
			out := make([]int, len(in))
			for i, l := range in {
				var v int
				v, err := strconv.Atoi(l)
				if err != nil {
					break
				}
				out[i] = v
			}
			r.ConnectedPlayerLatencies = out
		case "O1":
			in := strings.Split(value, "/")[1:]
			out := make([]int, len(in))
			for i, l := range in {
				var v int
				v, err := strconv.Atoi(l)
				if err != nil {
					break
				}
				out[i] = v
			}
			r.ConnectedPlayerKills = out
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
			if value == enabled {
				r.TeamNamesVisible = true
			}
		case "X1":
			if value == enabled {
				r.InternetServer = true
			}
		case "Y1":
			if value == enabled {
				r.FriendlyFire = true
			}
		case "Z1":
			if value == enabled {
				r.AutoTeamBalance = true
			}
		case "A2":
			if value == enabled {
				r.TeamkillPenalty = true
			}
		case "B2":
			if value == enabled {
				r.RadarAllowed = true
			}
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
			if value == enabled {
				r.AIBackup = true
			}
		case "J2":
			if value == enabled {
				r.RotateMapOnSuccess = true
			}
		case "K2":
			if value == enabled {
				r.ForceFirstPerson = true
			}
		case "L2":
			r.ModName = value
		case "L3":
			if value == enabled {
				r.PunkbusterEnabled = true
			}
		// O2 through O9 are reserved for OpenRVS custom fields.
		case "O2":
			r.MOTD = value
		default:
			// Unused or invalid key.
		}

		// Several case statement branches write to err before breaking. Check it now.
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}
