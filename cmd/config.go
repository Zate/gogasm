package main

import (
	"encoding/json"
	"fmt"
)

// PrettyPrint .. well it prints things pretty see..
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

// AtlasServers struct contains all the knowledge about unofficial and official
// servers
type AtlasServers struct {
	Official   []Realm `json:"official"`
	Unofficial []Realm `json:"unofficial"`
}

// Realm contains a name (such as NAPVE) and a collection of grid servers
type Realm struct {
	RealmName string  `json:"realm"`
	Grids     []Grids `json:"grids"`
}

// Grids is a struct for a grid server (such as A1) containing the QueryInfo and
// Config structs
type Grids struct {
	Grid   string `json:"Grid"`
	Info   `json:"info"`
	Config `json:"config"`
}

// Config struct contains all the information we will need to interact with a
// server/grid
type Config struct {
	AtlasIP string `env:"ATLASIP" envDefault:"159.203.52.169"` // ATLASIP=159.203.52.169
	//AtlasMaxPlayers		int			`env:"ATLASMAXPLAYERS" envDefault:"10"` // MAXPLAYERS=10
	AtlasGamePort     int    `env:"ATLASGAMEPORT1" envDefault:"27005"`    // GAMEPORT1=27005
	AtlasGamePortAlt  int    `env:"ATLASGAMEPORT2" envDefault:"27006"`    // GAMEPORT2=27006
	AtlasQueryPort    string `env:"ATLASQUERYPORT" envDefault:"27015"`    // ATLASQUERYPORT=27015
	AtlasRCONPort     string `env:"ATLASRCONPORT" envDefault:"27025"`     // RCONPORT=27025
	AtlasSeamlessPort int    `env:"ATLASSEAMLESSPORT" envDefault:"27020"` // SEAMLESSPORT=27020
	AtlasAdminPass    string `env:"ATLASADMINPASS" envDefault:"changeme"` // ADMINPASS=changeme
	//AtlasRCON 			bool		`env:"ATLASRCON" envDefault:false` // RCON=false
	//AtlasResPlayers		int			`env:"ATLASRESPLAYERS" envDefault:"0"` // RESPLAYERS=0
	// SLOG=-log
	// ALLHOME=-ForceAllHomeServer
	// MAP=Ocean
	// SVRX=0
	// SVRY=0
	// Home         string        `env:"HOME"`
	// Port         int           `env:"PORT" envDefault:"3000"`
	// IsProduction bool          `env:"PRODUCTION"`
	// Hosts        []string      `env:"HOSTS" envSeparator:":"`
	// Duration     time.Duration `env:"DURATION"`
	// TempFolder   string        `env:"TEMP_FOLDER" envDefault:"${HOME}/tmp" envExpand:"true"`
}

// Info contains all the information returned from a QueryServer request to an
// Atlas Server
type Info struct {
	Name        string // NAME: Atlas_D6 - (v16.14)
	Map         string // MAP: Ocean
	Folder      string // FOLDER: atlas
	Game        string // GAME: ATLAS
	ID          uint16 // ID: 0
	Players     byte   // PLAYERS: 26
	MaxPlayers  byte   // MAXPLAYERS: 150
	Bot         byte   // BOTS: 0
	ServerType  byte   // SERVERTYPE: d
	Environment byte   // ENVIRONMENT: w
	Visibility  byte   // VISIBILITY: 0
	Vac         byte   // VAC: 1
	Version     string // VERSION: 1.0.0.0
	Port        uint16 // PORT: 5759
	KeyWords    string // KEYWORDS: @,OWNINGID:90122942757731332,OWNINGNAME:90122942757731332,NUMOPENPUBCONN:124,P2PADDR:90122942757731332,P2PPORT:5759,NONATLAS_i:0
	Ping        int    // PING: 96
}

// 	Config      config // config struct

// }

// type realm struct {
// 	Name    string
// 	Servers []server
// }

// type server struct {
// 	Grid   string
// 	config config
// 	info   serverInfo
// }

// //type providers []provider

// //type provider []realm

// type realm []server

// type providers struct {
// 	provider []provider
// }

// type provider struct {
// 	realm
// }

// type name string

// type liveServers struct {
// 	// Each server should be able to have config and serverInfo
// 	Realm map[name]server
// }
