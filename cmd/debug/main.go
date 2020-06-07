package main

import (
	"log"
	"time"

	beacon "github.com/ijemafe/openrvs-beacon"
)

const targetIP = "72.251.228.169"
const targetPort = 7777

func main() {
	// test cut off game modes
	// update: this is caused by oversized UDP beacons. a safe limit is between
	// 512 and 548 bytes. Many beacons are over 1700 bytes, resulting in lost
	// data.
	b, err := beacon.GetServerReport(targetIP, targetPort+1000, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	r, err := beacon.ParseServerReport(targetIP, b)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("map count:", len(r.MapRotation))
	log.Println("mode count:", len(r.ModeRotation))
}
