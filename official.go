// Contains all the information known about IP Ranges for Atlas Live servers.
// It should grab this information from the local server binaries where possible
// so it will be kept up to date Downside is we wont know what realm something
// belongs to really until after we have connected to it, we should make sure we
// do some kind of validation after it's all built into structs so that at
// anytime a server says it's of a diff realm type to what we have, we should
// flag it for review and or perhaps automatically delete the entry for that one
// and create it in the right place?  Who knows, lets go with this for now.

package main

import "strings"

// Need something that has the base servers for each realm that we can lookup.
// Hard code it for now, work out how to make it dynamic and bulltproof later

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
