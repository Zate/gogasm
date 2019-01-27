package main

import (
	"encoding/binary"
	"encoding/hex"
	"log"
	"net"
	"strconv"
	"time"
	// "github.com/zate/gogasm/rcon"
)

var debug = false
var colorize = false

// CheckNoError standard error checking function
func CheckNoError(err error) bool {
	if err != nil {
		log.Println("Error: ", err)
		return false
	}
	return true
}

// Colorize wraps color around a string
func Colorize(s string) string {
	return "\033[1;31m" + s + "\033[0m"
}

// MyHexDump takes in array of bytes and and inter and returns a string
func MyHexDump(arr []byte, s int) string {
	var b = make([]byte, s)
	for i := 0; i < s; i++ {
		b[i] = arr[i]
	}
	if colorize {
		return Colorize(hex.Dump(b))
	}
	return hex.Dump(b)
}

// SendPacket sends a packet,duh.
func SendPacket(conn net.Conn, arr []byte, timeout time.Duration) (int, []byte) {
	_, err := conn.Write(arr)
	if err != nil {
		return 0, nil
	}
	buffer := make([]byte, 1500)
	conn.SetReadDeadline(time.Now().Add(timeout))
	n, err := conn.Read(buffer)
	if n == 0 || err != nil {
		conn.SetReadDeadline(time.Now().Add(timeout))
		n, err = conn.Read(buffer)
	}
	if n == 0 || err != nil {
		return 0, nil
	}
	return n, buffer
}

func stripCtlAndExtFromBytes(str string) string {
	b := make([]byte, len(str))
	var bl int
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 32 && c < 127 {
			b[bl] = c
			bl++
		}
	}
	return string(b[:bl])
}

// GetString used to pull strings from array of bytes returned from SendPacket
func GetString(arr []byte, index int) (string, int) {
	data := ""
	for i := index; i < len(arr); i++ {
		index = i
		if arr[i] == 0x00 {
			break
		} else {
			data = data + string(arr[i])
		}
	}
	index++
	data = stripCtlAndExtFromBytes(data)
	return data, index
}

// GetUInt16 converts array of bytes to uint16
func GetUInt16(arr []byte, index int) (uint16, int) {
	num1 := arr[index]
	index++
	num2 := arr[index]
	index++
	num := uint16(num1) | uint16(num2)<<8
	return num, index
}

// GetUInt32 converts array of bytes to uint32
func GetUInt32(arr []byte, index int) (uint32, int) {
	num1 := arr[index]
	index++
	num2 := arr[index]
	index++
	num3 := arr[index]
	index++
	num4 := arr[index]
	index++
	num := uint32(num4)<<24 | uint32(num3)<<16 | uint32(num2)<<8 | uint32(num1)
	return num, index
}

// CheckHeader to see if it's the right byte fotmat
func CheckHeader(hdr byte, chk byte) bool {
	if hdr != chk {
		log.Printf("Header was 0x%x instead of 0x%x\n", hdr, chk)
		return false
	}
	return true
}

// ServerPing sends a Query Protocol Ping packet to the server and looks at
// the response and returns a bool of alive or not
func ServerPing(server, port string) (alive bool) {
	a2sPing := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x56, 0xFF, 0xFF, 0xFF, 0xFF}
	timeout := 1500 * time.Millisecond
	sp := server + ":" + port

	conn, err := net.DialTimeout("udp", sp, timeout)
	if err != nil {
		return false
	}

	defer conn.Close()

	ret, err := conn.Write(a2sPing)
	if ret == 0 || err != nil {
		// Need to handle this better one day
		return false
	}
	BytesReceived := make([]byte, 1500)
	conn.SetReadDeadline(time.Now().Add(timeout))
	n, err := conn.Read(BytesReceived)

	if BytesReceived == nil || n == 0 {
		// Need to handle this better one day
		return false
	}

	if !CheckHeader(BytesReceived[4], 0x41) {
		// Need to handle this better one day
		return false
	}
	return true
}

// CheckStatus sends a Server Query Protocol request and prses the response
func CheckStatus(g Grids) Info {
	a2sInfo := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x54, 0x53, 0x6F, 0x75, 0x72, 0x63, 0x65, 0x20, 0x45, 0x6E, 0x67, 0x69, 0x6E, 0x65, 0x20, 0x51, 0x75, 0x65, 0x72, 0x79, 0x00}
	a2sRules := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x56, 0xFF, 0xFF, 0xFF, 0xFF}
	//a2sPlayer := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x55, 0xFF, 0xFF, 0xFF, 0xFF}
	seconds := 3
	timeout := time.Duration(seconds) * time.Second
	sp := g.Config.AtlasIP + ":" + g.Config.AtlasQueryPort

	Conn, err := net.DialTimeout("udp", sp, timeout)
	if err != nil {
		return g.Info
	}

	defer Conn.Close()

	// Get Info
	start := time.Now()
	n, BytesReceived := SendPacket(Conn, a2sInfo, timeout)
	t := time.Now()
	elapsed1 := t.Sub(start)

	if BytesReceived == nil || n == 0 {
		start = time.Now()
		n, BytesReceived = SendPacket(Conn, a2sInfo, timeout)
		t = time.Now()
		elapsed1 = t.Sub(start)
	}

	if BytesReceived == nil || n == 0 {
		log.Printf("%v Received no data! for A2S_INFO", g.Grid)
		return g.Info
	}

	if !CheckHeader(BytesReceived[4], 0x49) {
		log.Printf("%v Got wrong header response for A2S_INFO", g.Grid)
		return g.Info
	}

	var sPtr int
	var info string
	sPtr = 5
	g.Info.Name, sPtr = GetString(BytesReceived, sPtr)
	g.Info.Map, sPtr = GetString(BytesReceived, sPtr)
	g.Info.Folder, sPtr = GetString(BytesReceived, sPtr)
	g.Info.Game, sPtr = GetString(BytesReceived, sPtr)
	g.Info.ID, sPtr = GetUInt16(BytesReceived, sPtr)
	g.Info.Players = BytesReceived[sPtr]
	sPtr++
	g.Info.MaxPlayers = BytesReceived[sPtr]
	sPtr++
	g.Info.Bot = BytesReceived[sPtr]
	sPtr++
	g.Info.ServerType = BytesReceived[sPtr]
	sPtr++
	g.Info.Environment = BytesReceived[sPtr]
	sPtr++
	g.Info.Visibility = BytesReceived[sPtr]
	sPtr++
	g.Info.Vac = BytesReceived[sPtr]
	sPtr++
	g.Info.Version, sPtr = GetString(BytesReceived, sPtr)

	if n > sPtr {
		// EDF
		edf := BytesReceived[sPtr]
		sPtr++

		// PORT
		if edf&0x80 != 0 {
			//var port uint16
			g.Info.Port, _ = GetUInt16(BytesReceived, sPtr)
		}

		// STEAMID
		if edf&0x10 != 0 {
			sPtr += 8
		}

		// Keywords
		if edf&0x20 != 0 {
			g.Info.KeyWords, sPtr = GetString(BytesReceived, sPtr)
		}
	}
	// Get Rules
	sPtr = 5

	start = time.Now()
	n, BytesReceived = SendPacket(Conn, a2sRules, timeout)
	t = time.Now()
	elapsed2 := t.Sub(start)

	if BytesReceived == nil || n == 0 {
		start = time.Now()
		n, BytesReceived = SendPacket(Conn, a2sRules, timeout)
		t = time.Now()
		elapsed2 = t.Sub(start)
	}
	if BytesReceived == nil || n == 0 {
		log.Printf("%v Received no data! for A2S_RULES - Pre Challenge", g.Grid)
		return g.Info
	}

	if !CheckHeader(BytesReceived[4], 0x41) {
		log.Printf("%v Got wrong header response for A2S_RULES", g.Grid)
		return g.Info
	}

	// Challenge number
	var chnum uint32
	chnum, sPtr = GetUInt32(BytesReceived, sPtr)

	a2sRules[5] = byte(chnum)
	a2sRules[6] = byte(chnum >> 8)
	a2sRules[7] = byte(chnum >> 16)
	a2sRules[8] = byte(chnum >> 24)

	start = time.Now()
	n, BytesReceived = SendPacket(Conn, a2sRules, timeout)
	t = time.Now()
	elapsed3 := t.Sub(start)

	if BytesReceived == nil || n == 0 {
		log.Printf("%v Received no data! for A2S_RULES Post Challenge", g.Grid)
		return g.Info
	}

	elapsed := (elapsed1 + elapsed2 + elapsed3) / 3
	g.Info.Ping = int(elapsed) / 1000000
	if !CheckHeader(BytesReceived[4], 0x45) {
		return g.Info
	}

	// reset sPtr
	sPtr = 5
	var rules uint16
	rules, sPtr = GetUInt16(BytesReceived, sPtr)
	rulesMap := make(map[string]string)
	if rules > 0 {
	}

	for i := uint16(0); i < rules; i++ {
		// Name
		info, sPtr = GetString(BytesReceived, sPtr)
		// Value
		val := ""
		val, sPtr = GetString(BytesReceived, sPtr)
		rulesMap[info] = val
	}
	return g.Info
}

func toGrid(s, gs int) (grid string) {
	g := (s % gs) + 1
	l := (s / gs)
	letter := string('A' + l)
	grid = letter + strconv.Itoa(g)
	return grid
}

// Grid attempts to map official server space for each port/ip combo
func makeGrid(gridSize, portsPerServer int) {
	count := 0
	portCount := 1
	for i := 0; i < (gridSize*gridSize)+1; i++ {
		ipcount := (count / portsPerServer)
		portNum := 57554 + portCount
		log.Printf("count: %v is %v and ip: %v port: %v", count, toGrid(count, gridSize), ipcount, portNum)
		portCount++
		portCount++
		if portCount == 9 {
			portCount = 1
		}
		count++
	}
	return

}

// AddOfficialRealm takes a []Realm struct and adds it to the base AtlasServers
// struct
func (l *AtlasServers) AddOfficialRealm(realm Realm) []Realm {
	l.Official = append(l.Official, realm)
	return l.Official
}

// AddGrid takes a []Grids struct and adds it to the Realm struct
func (r *Realm) AddGrid(grid Grids) []Grids {
	r.Grids = append(r.Grids, grid)
	return r.Grids
}

// incIP increments the base IP by int and returns a string
func incIP(ip string, i int) string {
	ui := uint32(i)
	s := int2ip(ip2int(net.ParseIP(ip)) + ui)

	//log.Println(s.String())
	return s.String()
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}

// LiveAtlasServers will build a list of all the live servers in the grid for
// the realm specified
func LiveAtlasServers(realm string) {
	b := InitOfficialBaseIPs(realm)
	//log.Println(b)
	var l AtlasServers
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
		//log.Printf("count: %v is %v and ip: %v port: %v", i, grid, ipcount, portNum)
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
	l.AddOfficialRealm(r)
	//PrettyPrint(l)

	for _, v := range l.Official {
		if v.RealmName == realm {
			//PrettyPrint(v)
			for _, g := range v.Grids {
				// up := ServerPing(g.Config.AtlasIP, g.Config.AtlasQueryPort)
				// if up == true {
				g.Info = CheckStatus(g)
				//PrettyPrint(g.Info)
				// }
				log.Printf("%v | %v | %v:%v Pop: %v", g.Grid, g.Info.Name, g.Config.AtlasIP, g.Config.AtlasQueryPort, g.Info.Players)
			}
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	// statusPtr := flag.Bool("status", false, "Check the status of the server")
	// serverPtr := flag.String("s", "", "IP Address of the server to connect to")
	// portPtr := flag.String("p", "", "Port on the server to connect to")
	// pingPtr := flag.Bool("ping", false, "Simply checks if server responds to query on that port")
	// gridPtr := flag.Bool("grid", false, "used to print out possible grid assignments based on 4 ports per IP, consecutive IP's")
	// livePtr := flag.String("live", "", "to run status/ping on live servers")
	// flag.Parse()

	// if len(*livePtr) > 0 {
	// 	// Build out a grid of all the live servers
	// 	realm := *livePtr
	// 	go LiveAtlasServers(realm)
	// 	// os.Exit(0)
	// }

	LiveAtlasServers("napve")
	//initWeb()

}
