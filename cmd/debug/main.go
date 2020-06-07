// The debug cmd retrieves UDP beacons and tests for missing data.
// Maximum safe size for UDP packets is 512 bytes on Linux and 548 bytes on
// Windows. OpenRVS beacons start to drop data around 1700 bytes.
package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	beacon "github.com/ijemafe/openrvs-beacon"
)

type Input struct {
	IP   string
	Port int
}

func main() {
	//targets := []Input{Input{"162.248.92.181", 7778}} // Broken beacon response
	targets, err := getHostPorts()
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range targets {
		b, err := beacon.GetServerReport(t.IP, t.Port+1000, 5*time.Second)
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = beacon.ParseServerReport(t.IP, b)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func getHostPorts() ([]Input, error) {
	var inputs = make([]Input, 0)
	resp, err := http.Get("http://64.225.54.237/servers")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(bytes.TrimSuffix(b, []byte{'\n'}), []byte{'\n'})
	for i := 1; i < len(lines); i++ {
		fields := bytes.Split(lines[i], []byte{','})
		host := string(fields[1])
		portBytes := fields[2]
		port, err := strconv.Atoi(string(portBytes))
		if err != nil {
			log.Println("atoi error:", err)
			continue
		}
		inputs = append(inputs, Input{IP: host, Port: port})
	}
	return inputs, nil
}
