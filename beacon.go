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
//     - Player Names (20-320 bytes for 1-16 players)
//     - Player Times (5-80 bytes for 1-16 players)
//     - Player Pings (5-80 bytes for 1-16 players)
//     - Player Kills (4-64 bytes for 1-16 players)
//     - Current Players (2 bytes)
//     - Rounds Per Match (2 bytes, more above 99)
//     - Time Per Round (4 bytes, more above 2.75 hours)
//     - Time Between Rounds (2 bytes, more above 99s)
//     - Bomb Timer (2 bytes, more above 99s)
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
	"time"
)

const (
	beaconBufferSize  = 4096 // 4kb. Most responses are under 2kb, data loss begins at 1kb.
	minimumReportSize = 128  // High enough to avoid invalid data, but lower than actual values (320+).

	sep          = '¶'       // "Pilcrow Sign". Red Storm used this as a field separator.
	beaconHeader = "rvnshld" // Start of header line in UDP response.

	enabled  = "1"
	disabled = "0"
)

// ErrNotABeacon indicates a valid UDP response which is not from OpenRVS.
var ErrNotABeacon = fmt.Errorf("error: not a valid openrvs beacon")

// GetServerReport will send "REPORT" over UDP to the given server and return the response bytes.
// If the response does not have a 'rvnshld' header, ErrNotABeacon will be returned.
// NOTE: The port in question is the beacon port and not the game server port. The beacon port is typically the game
// server port plus 1000.
func GetServerReport(ip string, port int, timeout time.Duration) ([]byte, error) {
	b, err := sendUDP("REPORT", ip, port, timeout)
	if err != nil {
		return []byte{}, err
	}

	if err := validateServerReport(b); err != nil {
		return []byte{}, err
	}

	return b, nil
}

func validateServerReport(report []byte) error {
	if len(report) < minimumReportSize {
		return ErrNotABeacon
	}
	if !bytes.HasPrefix(report, []byte(beaconHeader)) {
		return ErrNotABeacon
	}
	return nil
}
