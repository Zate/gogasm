// Contains all the information known about IP Ranges for Atlas Live servers.
// It should grab this information from the local server binaries where possible so it will be kept up to date
// Downside is we wont know what realm something belongs to really until after we have connected to it, we should make sure
// we do some kind of validation after it's all built into structs so that at anytime a server says it's of a diff realm type
// to what we have, we should flag it for review and or perhaps automatically delete the entry for that one and create it in
// the right place?  Who knows, lets go with this for now.

package main

// NAPVEServers - NAPVE servers first
