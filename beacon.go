package beacon

import (
	"bytes"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	sep      = '¶'       // "Pilcrow Sign". Red Storm used this as a field separator.
	header   = "rvnshld" // Start of header line in UDP response.
	enabled  = "1"
	disabled = "0"
)

var timeout = 3 * time.Second // Time to wait for beacon response before moving on.

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

// Server is an endpoint for us to check.
type Server struct {
	IP   string
	Port int // This is the GAME SERVER port.
}

// GetServerReport handles the UDP connection to the server's beacon port.
func GetServerReport(ip string, port int) ([]byte, error) {
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{IP: net.ParseIP(ip), Port: port})
	if err != nil {
		return nil, err
	}
	if _, err = conn.Write([]byte("REPORT")); err != nil {
		return nil, err
	}
	buf := make([]byte, 4096) // Most responses are under 512 bytes.
	conn.SetReadDeadline(time.Now().Add(timeout))
	if _, err = conn.Read(buf); err != nil {
		return nil, err
	}
	b, err := bytes.Trim(buf, "\x00"), nil
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ParseServerReport reads the bytestream from the game server and parses it into a serverResponse object.
func ParseServerReport(ip string, report []byte) (*ServerReport, error) {
	r := &ServerReport{IPAddress: ip}
	for _, line := range bytes.Split(report, []byte{sep}) {
		// Skip the header line, no useful info to parse.
		if strings.HasPrefix(string(line), header) {
			continue
		}

		key := string(line[0:2])
		value := string(bytes.Trim(line[3:], "\x20"))

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
			// ModeRotation includes "/" separators for every slot, not every mode. Omit empty values.
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
