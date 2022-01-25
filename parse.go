package beacon

import (
	"bytes"
	"errors"
	"log"
	"strconv"
	"strings"
)

const minimumLineSize = 2 // Based on key size; value can be empty.

var ErrInvalidLine = errors.New("invalid report line")

func validateLine(line []byte) error {
	if len(line) < minimumLineSize {
		return ErrInvalidLine
	}

	return nil
}

// Removes traling whitespace and returns the key and value from the given line.
func lineToKeyVal(line []byte) (string, string) {
	key := string(line[0:2])
	val := string(bytes.Trim(line[3:], "\x20"))
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
		// Case statements can be brittle, but there's no risk of this code changing.
		var err error
		switch key {
		case "A1":
			r.MaxPlayers, err = strconv.Atoi(value)
		case "B1":
			r.NumPlayers, err = strconv.Atoi(value)
		// No C1 or D1.
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
		case "J1":
			// ModeRotation always includes 32 "/" separators, not just one for
			// each mode. Omit empty strings from ModeRotation.
			modes := make([]string, 0)
			fields := strings.Split(value, "/")
			for _, m := range fields[1:] {
				if m != "" {
					modes = append(modes, m)
				}
			}
			r.ModeRotation = modes
			// Note: Mode rotation is the last thing to arrive over UDP. If it
			// is missing any placeholder '/' characters, data loss occurred.
			if len(fields)-1 != 32 && r.Port != 0 {
				log.Printf("warning: data loss occurred for server %s:%d (received %d bytes)",
					ip, r.Port, len(report))
			}
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
		// No U1 or V1.
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
		// O2 and above are OpenRVS custom fields.
		case "O2":
			r.MOTD = value
		default:
			log.Println("unknown key:", key)
		}

		// Several case statement branches write to err before breaking. Check it now.
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}
