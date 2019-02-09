// Contains all the information known about IP Ranges for Atlas Live servers.
// It should grab this information from the local server binaries where possible
// so it will be kept up to date Downside is we wont know what realm something
// belongs to really until after we have connected to it, we should make sure we
// do some kind of validation after it's all built into structs so that at
// anytime a server says it's of a diff realm type to what we have, we should
// flag it for review and or perhaps automatically delete the entry for that one
// and create it in the right place?  Who knows, lets go with this for now.

package main

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"
)

// Need something that has the base servers for each realm that we can lookup.
// Hard code it for now, work out how to make it dynamic and bulltproof later

// LiveRealms is just an array of the names of each of the 4 realms.
var LiveRealms = []string{"napve", "napvp", "eupve", "eupvp"}

// InitLive populates the AtlasSearcer struct with a Realm for each live realm.
// Should only run on server startup.
func initLive() {
	var l AtlasServers
	for _, r := range LiveRealms {
		r := InitLiveRealm(r)
		l.AddOfficialRealm(r)
	}
	PrettyPrint(l)
}

// InitOfficialBaseIPs will build out the initial listing of base IPs for each
// official realm.  Defaults to NAPVE cos fuck PVP :-P
func InitOfficialBaseIPs(r string) string {
	r = strings.ToLower(r)

	switch r {
	case "napve":
		return "37.10.126.130"
	case "napvp":
		return "37.10.127.123"
	case "eupve":
		return "46.251.238.2"
	case "eupvp":
		return "46.251.238.59"
	default:
		return "37.10.126.130"
	}
}

// InitLiveRealm populates the defauly data into a realm struct
func InitLiveRealm(realm string) Realm {
	b := InitOfficialBaseIPs(realm)
	var r Realm
	r.RealmName = realm
	gridSize := 15
	portsPerServer := 4
	portCount := 1

	for i := 0; i < (gridSize * gridSize); i++ {
		var g Grids
		ipcount := (i / portsPerServer)
		portNum := 57554 + portCount
		grid := toGrid(i, gridSize)
		g.Grid = grid
		g.Config.AtlasIP = incIP(b, ipcount)
		g.Config.AtlasQueryPort = strconv.Itoa(portNum)
		portCount++
		portCount++
		if portCount == 9 {
			portCount = 1
		}
		// Need to calculate what the ip should be based off a base IP, 4
		// servers per IP.  lets try and use ipcount to count from what the base
		// server for this realm is.  likely to break if they change subnets,
		// but unfuck that when we get to it Also not going to convert it to a
		// net.ip or anything just yet.  This means we could easily break things
		// by say incrementing beyond a subnet boundary, but lets get it working
		// right first.

		r.AddGrid(g)

	}
	return r
}

func getLive() {
	for {
		realminfo := make(chan Realm)
		var wg sync.WaitGroup

		wg.Add(len(LiveRealms))
		for _, realm := range LiveRealms {
			go func(realm string) {
				defer wg.Done()
				var r Realm
				r.RealmName = realm
				for _, g := range r.Grids {
					g, err := CheckStatus(g)
					if err != nil {
						log.Printf("%v failed to get status\n", g.Grid)
					} else {
						realminfo <- r
					}

					gj, err := json.Marshal(g)
					if err != nil {
						log.Fatalf("Marshaling live servers to json failed: %v", err)
					}
					PrettyPrint(string(gj))
				}
			}(realm)
		}
		go func() {
			for info := range realminfo {
				ri, err := json.Marshal(info)
				if err != nil {
					log.Fatalf("turning realm info to json failed: %v", err)
				}
				log.Printf("Putting %v in db", info.RealmName)
				PrettyPrint(string(ri))
				// This is where we would insert to db with a post to /dbname
			}
		}()

		wg.Wait()
		// Now synch it to the db

	}
}

// GetLiveRealm will build a list of all the live servers in the grid for
// the realm specified
func GetLiveRealm(r Realm) Realm {

	return r

}
