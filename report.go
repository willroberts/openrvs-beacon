package beacon

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
