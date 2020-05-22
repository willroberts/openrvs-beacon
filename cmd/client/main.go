package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"

	beacon "github.com/ijemafe/openrvs-beacon"
)

func main() {
	servers, err := readServerList()
	if err != nil {
		log.Fatal("failed to read server list:", err)
	}

	reports := make(chan *beacon.ServerReport, len(servers))
	errs := make(chan error, len(servers))

	var wg sync.WaitGroup
	for _, s := range servers {
		wg.Add(1)
		go func(s beacon.Server) {
			b, err := beacon.GetServerReport(s.IP, s.Port+1000)
			if err != nil {
				errs <- err
				wg.Done()
				return
			}

			r, err := beacon.ParseServerReport(s.IP, b)
			if err != nil {
				errs <- err
				wg.Done()
				return
			}
			reports <- r
			wg.Done()
		}(s)
	}
	wg.Wait()
	close(reports)
	close(errs)

	for e := range errs {
		if e != nil {
			log.Println(e)
		}
	}

	for r := range reports {
		fmt.Println("Server:", r.ServerName)
		fmt.Printf("Address: %s:%d\n", r.IPAddress, r.Port)
		fmt.Println("Game Version:", r.GameVersion)
		fmt.Println("Mod Name:", r.ModName)
		fmt.Println("MOTD:", r.MOTD)
		fmt.Println("Current Map:", r.CurrentMap)
		fmt.Println("Current Game Mode:", r.CurrentMode)
		if r.NumTerrorists > 0 {
			fmt.Println("Number of Terrorists:", r.NumTerrorists)
		}
		fmt.Println("Friendly Fire:", r.FriendlyFire)
		fmt.Printf("Active Players: %d out of %d\n", r.NumPlayers, r.MaxPlayers)
		for i := 0; i < len(r.ConnectedPlayerNames); i++ {
			fmt.Printf("- %s (Kills: %d, Ping: %dms)\n",
				r.ConnectedPlayerNames[i],
				r.ConnectedPlayerKills[i],
				r.ConnectedPlayerLatencies[i])
		}
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
