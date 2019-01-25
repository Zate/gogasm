package main

type config struct {
	AtlasIP string `env:"ATLASIP" envDefault:"159.203.52.169"` // ATLASIP=159.203.52.169
	//AtlasMaxPlayers		int			`env:"ATLASMAXPLAYERS" envDefault:"10"` // MAXPLAYERS=10
	//AtlasGamePort1		int			`env:"ATLASGAMEPORT1" envDefault:"27005"` // GAMEPORT1=27005
	//AtlasGamePort2		int			`env:"ATLASGAMEPORT2" envDefault:"27006"`// GAMEPORT2=27006
	AtlasQueryPort string `env:"ATLASQUERYPORT" envDefault:"27015"` // ATLASQUERYPORT=27015
	AtlasRCONPort  string `env:"ATLASRCONPORT" envDefault:"27025"`  // RCONPORT=27025
	//AtlasSeamlessPort	int			`env:"ATLASSEAMLESSPORT envDefault:"27020"` // SEAMLESSPORT=27020
	AtlasAdminPass string `env:"ATLASADMINPASS" envDefault:"changeme"` // ADMINPASS=changeme
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