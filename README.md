# openrvs-beacon

Library and client for the OpenRVS Beacon Server

## example client

There is an example of using this library in `cmd/client/main.go` which reads from a local list of server addresses. The output for each server looks like this:

```
Server: ALLR6 | Original Maps Only
Address: 104.243.46.138:9777
Game Version: PATCH 1.60 (build 412)
Mod Name: RavenShield
Current Map: Mountain_High
Current Game Mode: RGM_TerroristHuntCoopMode
Number of Terrorists: 50
Friendly Fire: true
Active Players: 5 out of 8
- (Srgt)ThriceQC (Kills: 3, Ping: 62ms)
- R6_Pride (Kills: 4, Ping: 94ms)
- 1 (Kills: 16, Ping: 140ms)
- SusanSarandonJr (Kills: 0, Ping: 31ms)
- (Srgt)McAGravQc (Kills: 18, Ping: 46ms)
```
