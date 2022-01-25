# openrvs-beacon

Library and client for the OpenRVS Beacon Server

## Documentation

Complete documentation can be found [here](https://godoc.org/github.com/willroberts/openrvs-beacon)

## Example client

There is an example of using this library in the `cmd/client` directory:

```
$ go run main.go -ip 64.225.54.237 -port 7776
Server: Classic Maps | Terrorist Hunt
Address: 64.225.54.237:6776
Game Version: PATCH 1.60 (build 412)
Mod Name: RavenShield
MOTD: 
Current Map: Streets
Current Game Mode: RGM_TerroristHuntCoopMode
Number of Terrorists: 35
Friendly Fire: false
Active Players: 0 out of 8
```