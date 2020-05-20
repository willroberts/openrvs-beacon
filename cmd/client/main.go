package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	beacon "github.com/ijemafe/openrvs-beacon"
)

func main() {
	servers, err := readServerList()
	if err != nil {
		log.Fatal("failed to read server list:", err)
	}

	for _, s := range servers {
		b, err := beacon.GetServerReport(s.IP, s.Port+1000)
		if err != nil {
			log.Println("failed to read from server:", err)
			fmt.Println()
			continue
		}

		r, err := beacon.ParseServerReport(s.IP, b)
		if err != nil {
			log.Println("failed to parse report:", err)
			fmt.Println()
			continue
		}

		fmt.Println("Server:", r.ServerName)
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
func readServerList() ([]beacon.Server, error) {
	b, err := ioutil.ReadFile("servers.example")
	if err != nil {
		return nil, err
	}
	var servers []beacon.Server

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		port, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		servers = append(servers, beacon.Server{
			IP:   fields[0],
			Port: port,
		})
	}

	return servers, nil
}
