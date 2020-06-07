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
	BeaconBufferSize = 4096      // Most responses are under 2048 bytes, data loss begins at 1023.
	sep              = '¶'       // "Pilcrow Sign". Red Storm used this as a field separator.
	header           = "rvnshld" // Start of header line in UDP response.
	enabled          = "1"
	disabled         = "0"
)

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

// GetServerReport handles the UDP connection to the server's beacon port.
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
	buf := make([]byte, BeaconBufferSize) // Most responses are under 2048 bytes.
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
// into a serverResponse object.
func ParseServerReport(ip string, report []byte) (*ServerReport, error) {
	r := &ServerReport{IPAddress: ip}
	for _, line := range bytes.Split(report, []byte{sep}) {
		// Skip the header line, no useful info to parse.
		if strings.HasPrefix(string(line), header) {
			continue
		}

		// These two iterations convert ASCII bytes to UTF-8. If we do something
		// like string(keyBytes) instead, non-ASCII character will be converted
		// into '�'.
		keyBytes := line[0:2]
		key := ""
		for _, b := range keyBytes {
			key += string(b)
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
