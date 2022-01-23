// Package beacon provides a library for interacting with Rainbow Six 3: Raven
// Shield game servers.
//
// Raven Shield game servers communicate with UDP for all purposes. Servers
// listen on UDP ports for clients, clients join games with UDP, and there is
// communication between clients and servers over UDP as well. The format used
// in this communication is called a "beacon". A beacon is a stream of text
// separated by named markers containing pilcrow (¶) signs. Each segment follows
// a specific marker, so the format is similar to an ordered map of string keys
// to string values.
//
// The UDP beacon structure looks like this (in exact order of appearance):
//
//     - Port (4-5 bytes)
//     - Map Name (up to 32 bytes)
//     - Server Name (up to 32 bytes)
//     - Current Game Mode (15-25 bytes)
//     - Maximum Players (2 bytes)
//     - Locked? (1 byte)
//     - Dedicated? (1 byte)
//     - Player Names (0 or 20-320 bytes for 1-16 players)
//     - Player Times (0 or 5-80 bytes for 1-16 players)
//     - Player Pings (0 or 5-80 bytes for 1-16 players)
//     - Player Kills (0 or 4-64 bytes for 1-16 players)
//     - Current Players (2 bytes)
//     - Rounds Per Match (2 bytes, more above 99)
//     - Time Per Round (4 bytes, more above 2.75 hours)
//     - Time Between Rounds (2 bytes, more above 1m39s)
//     - Bomb Timer (2 bytes, more above 1m39s)
//     - Team Names Visible? (1 byte)
//     - Internet Server? (1 byte)
//     - Friendly Fire? (1 byte)
//     - Auto Team Balance? (1 byte)
//     - Teamkill Penalty? (1 byte)
//     - Game Version: (22-32 bytes)
//     - Radar Allowed? (1 byte)
//     - Lobby Server ID (1 byte)
//     - Group ID (1 byte)
//     - Beacon Port (4-5 bytes)
//     - Num Terrorists (2 bytes)
//     - AI Backup? (1 byte)
//     - Rotate Map? (1 byte)
//     - Force First Person? (1 byte)
//     - Mod Name (9-12 bytes)
//     - Punkbuster? (1 byte)
//     - Map Rotation (10-831 bytes for 1-32 maps)
//     - Mode Rotation (47-832 bytes for 1-32 modes)
//     - MOTD (up to 60 bytes, only in OpenRVS 1.5+)
//
// Each component is preceded by a 3-byte marker. Overall, without the map and
// mode rotations, and with no players connected, a server should be able to fit
// all data within 320 bytes.
//
// Depending on the network and OS this code runs on, UDP data loss may occur at
// different points. The safe limit is generally 512 bytes, and data loss could
// occur at higher values. In local development, data loss begins at 1024 bytes.
// When data loss occurs, the list of game modes and the MOTD will be cut off.
//
// In order to avoid data loss, we may need to fragment beacons across multiple
// UDP messages, such as moving the map rotation and/or mode rotation to
// separate beacons. 1024 bytes is enough space to fit the base data, 8 sets of
// connected player data, and 431 bytes of map rotation data. Some servers will
// exceed this map rotation data length, and PVP servers will have 16 sets of
// connected player data, so excising the game mode rotation may not be enough
// to get responses under 1024 bytes.
package beacon

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	beaconBufferSize = 4096      // 4kb. Most responses are under 2kb, data loss begins at 1kb.
	sep              = '¶'       // "Pilcrow Sign". Red Storm used this as a field separator.
	header           = "rvnshld" // Start of header line in UDP response.
	enabled          = "1"
	disabled         = "0"
)

// ErrNotABeacon indicates a valid UDP response which is not from OpenRVS.
var ErrNotABeacon = fmt.Errorf("error: response was not an openrvs beacon")

// ServerReport is the response object from the game server's beacon port.
type ServerReport struct {

	// Server settings.

	ServerName        string
	IPAddress         string
	Port              int
	BeaconPort        int
	InternetServer    bool
	Dedicated         bool
	PunkbusterEnabled bool
	Locked            bool
	MaxPlayers        int
	NumPlayers        int
	GameVersion       string
	ModName           string
	OptionsList       string // The beacon does not seem to return this value.
	LobbyServerID     int    // Ubisoft-specific. Always 0.
	GroupID           int    // Ubisoft-specific. Always 0.

	// Game settings.

	AIBackup                 bool
	AutoTeamBalance          bool
	BombTimer                int
	ConnectedPlayerKills     []int
	ConnectedPlayerLatencies []int
	ConnectedPlayerNames     []string
	ConnectedPlayerTimes     []string
	CurrentMap               string
	CurrentMode              string
	ForceFirstPerson         bool
	FriendlyFire             bool
	MapRotation              []string
	ModeRotation             []string
	NumTerrorists            int
	RadarAllowed             bool
	RotateMapOnSuccess       bool
	RoundsPerMatch           int
	TeamNamesVisible         bool
	TeamkillPenalty          bool
	TimeBetweenRounds        int
	TimePerRound             int

	// OpenRVS custom fields.

	MOTD string
}

// GetServerReport handles the UDP connection to the server's beacon port and
// retrieves the report bytes. Note that the port in question is the beacon port
// and not the game server port. The beacon port is typically the gamer server
// port plus 1000.
func GetServerReport(ip string, port int, timeout time.Duration) ([]byte, error) {
	// "Connect" to the remote UDP port.
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{IP: net.ParseIP(ip), Port: port})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Send a REPORT request.
	conn.SetReadDeadline(time.Now().Add(timeout))
	if _, err = conn.Write([]byte("REPORT")); err != nil {
		return nil, err
	}

	// Try to read the REPORT response into a buffer.
	buf := make([]byte, beaconBufferSize) // Most responses are under 2048 bytes.
	if _, err = conn.Read(buf); err != nil {
		return nil, err
	}

	// Validate the response.
	if !bytes.HasPrefix(buf, []byte(header)) {
		return nil, ErrNotABeacon
	}

	// Remove empty bytes from the end of the buffer.
	b, err := bytes.Trim(buf, "\x00"), nil
	if err != nil {
		return nil, err
	}

	return b, nil
}

// ParseServerReport reads the bytestream from the game server and parses it
// into a ServerReport.
func ParseServerReport(ip string, report []byte) (*ServerReport, error) {
	r := &ServerReport{IPAddress: ip}
	for _, line := range bytes.Split(report, []byte{sep}) {
		// Skip the header line, no useful info to parse.
		if strings.HasPrefix(string(line), header) {
			continue
		}

		// These two iterations convert ASCII bytes to UTF-8. If we do something
		// like string(keyBytes) instead, non-ASCII characters will be converted
		// into '�'.
		keyBytes := line[0:2]
		key := ""
		for _, b := range keyBytes {
			key += string(b)
		}
		// A valid line must contain both key and value.
		if len(line) < 3 {
			continue
		}
		valueBytes := bytes.Trim(line[3:], "\x20")
		value := ""
		for _, b := range valueBytes {
			value += string(b)
		}

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
			break
		}

		if err != nil {
			return nil, err
		}
	}

	return r, nil
}
