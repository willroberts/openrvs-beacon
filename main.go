package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

var (
	timeout = 3 * time.Second // Time to wait for beacon response before moving on.
)

// serverReport is the response object for the server beacon.
type serverReport struct {
	IPAddress                string
	Port                     int
	CurrentMap               string
	Name                     string
	CurrentMode              string
	MaxPlayers               int
	Locked                   bool
	Dedicated                bool
	ConnectedPlayerNames     []string
	ConnectedPlayerTimes     []string
	ConnectedPlayerLatencies []int
	ConnectedPlayerKills     []int
	GameMode                 string
	RoundsPerMatch           int
	TimePerRound             int
	TimeBetweenRounds        int
	BombTimer                int
	TeamNamesVisible         bool
	InternetServer           bool
	FriendlyFire             bool
	AutoTeamBalance          bool
	TeamkillPenalty          bool
	GameVersion              string
	RadarAllowed             bool
	UnknownE2                int // Unknown (E2). Appears to always be set to 0.
	UnknownF2                int // Unknown (F2). Appears to always be set to 0.
	BeaconPort               int
	NumTerrorists            int
	AIBackup                 bool
	RotateMapOnSuccess       bool
	ForceFirstPerson         bool
	ModName                  string
	UnknownL3                int // Unknown (L3). Appears to always be set to 0.
	MapRotation              []string
	ModeRotation             []string
}

// server is an endpoint for us to check.
type server struct {
	IP   string
	Port int // This is the GAME SERVER port.
}

func main() {
	servers, err := readServerList()
	if err != nil {
		log.Fatal("failed to read server list:", err)
	}

	for _, s := range servers {
		queryPort := s.Port + 1000
		log.Printf("Querying %s:%d", s.IP, s.Port)
		b, err := getServerReport(s.IP, queryPort)
		if err != nil {
			log.Println("failed to read from server:", err)
			fmt.Println()
			continue
		}

		r, err := parseServerReport(s.IP, b)
		if err != nil {
			log.Println("failed to parse report:", err)
			fmt.Println()
			continue
		}

		fmt.Println("Server:", r.Name)
		fmt.Printf("Address: %s:%d\n", r.IPAddress, r.Port)
		fmt.Println("Game Version:", r.GameVersion)
		fmt.Println("Mod Name:", r.ModName)
		fmt.Println("Current Map:", r.CurrentMap)
		fmt.Println("Current Game Mode:", r.CurrentMode)
		if r.NumTerrorists > 0 {
			fmt.Println("Number of Terrorists:", r.NumTerrorists)
		}
		fmt.Println("Friendly Fire:", r.FriendlyFire)
		fmt.Printf("Active Players: %d out of %d\n", len(r.ConnectedPlayerNames), r.MaxPlayers)
		fmt.Println()
	}
}

// readServerList checks serverlist.example and parses the contents. This text was copied from RVSGaming.org.
func readServerList() ([]server, error) {
	b, err := ioutil.ReadFile("servers.example")
	if err != nil {
		return nil, err
	}
	var servers []server

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		port, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		servers = append(servers, server{
			IP:   fields[0],
			Port: port,
		})
	}

	return servers, nil
}

// getServerReport handles the UDP connection to the server's beacon port.
func getServerReport(ip string, port int) ([]byte, error) {
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{IP: net.ParseIP(ip), Port: port})
	if err != nil {
		return nil, err
	}
	if _, err = conn.Write([]byte("REPORT")); err != nil {
		return nil, err
	}
	buf := make([]byte, 65536)
	conn.SetReadDeadline(time.Now().Add(timeout))
	if _, err = conn.Read(buf); err != nil {
		return nil, err
	}
	b, err := bytes.Trim(buf, "\x00"), nil // 477 bytes remain.
	if err != nil {
		return nil, err
	}
	return b, nil
}

// parseServerReport reads the bytestream from the game server and parses it into a serverResponse object.
func parseServerReport(ip string, report []byte) (*serverReport, error) {
	r := &serverReport{IPAddress: ip}
	for _, line := range bytes.Split(report, []byte{sep}) {
		// Skip the header line, no useful info to parse.
		if strings.HasPrefix(string(line), header) {
			continue
		}

		key := string(line[0:2])
		value := string(bytes.Trim(line[3:], "\x20"))

		// Case statements are brittle, but there's no risk of this code changing.
		var err error
		switch key {
		case "P1":
			r.Port, err = strconv.Atoi(value)
		case "E1":
			r.CurrentMap = value
		case "I1":
			r.Name = value
		case "F1":
			r.CurrentMode = value
		case "A1":
			r.MaxPlayers, err = strconv.Atoi(value)
		case "G1":
			if value == enabled {
				r.Locked = true
			}
		case "H1":
			if value == enabled {
				r.Dedicated = true
			}
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
		case "B1":
			if value == enabled {
				r.GameMode = "coop"
				break
			}
			r.GameMode = "adver"
		case "Q1":
			r.RoundsPerMatch, err = strconv.Atoi(value)
		case "R1":
			r.TimePerRound, err = strconv.Atoi(value)
		case "S1":
			r.TimeBetweenRounds, err = strconv.Atoi(value)
		case "T1":
			r.BombTimer, err = strconv.Atoi(value)
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
		case "D2":
			r.GameVersion = value
		case "B2":
			if value == enabled {
				r.RadarAllowed = true
			}
		case "E2":
			// Unknown field.
			if value != disabled {
				log.Println("E2 is nonzero!")
			}
			break
		case "F2":
			// Unknown field.
			if value != disabled {
				log.Println("F2 is nonzero!")
			}
			break
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
			// Unknown field.
			if value != disabled {
				log.Println("L3 is nonzero!")
			}
			break
		case "K1":
			r.MapRotation = strings.Split(value, "/")[1:]
		case "J1":
			// ModeRotation includes "/" separators for every slot, not every mode. Omit empty values.
			modes := make([]string, 0)
			for _, m := range strings.Split(value, "/")[1:] {
				if m != "" {
					modes = append(modes, m)
				}
			}
			r.ModeRotation = modes
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
