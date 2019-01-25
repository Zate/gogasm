package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/caarlos0/env"
	// "github.com/zate/gogasm/rcon"
)

var debug = false
var colorize = false

// CheckNoError standard error checking function
func CheckNoError(err error) bool {
	if err != nil {
		fmt.Println("Error: ", err)
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
	if debug {
		fmt.Fprintln(os.Stderr, "Writing...")
	}
	ret, err := conn.Write(arr)
	if debug {
		fmt.Fprintf(os.Stderr, MyHexDump(arr, ret))
	}
	if CheckNoError(err) {
		if debug {
			fmt.Fprintf(os.Stderr, "Wrote %d bytes\n", ret)
		}
		buffer := make([]byte, 1500)
		if debug {
			fmt.Fprintln(os.Stderr, "Reading...")
		}
		conn.SetReadDeadline(time.Now().Add(timeout))
		n, err := conn.Read(buffer)
		if CheckNoError(err) {
			if debug {
				fmt.Fprintf(os.Stderr, "Read %d bytes\n", n)
			}
			if debug {
				fmt.Fprintf(os.Stderr, MyHexDump(buffer, n))
			}
			return n, buffer
		}
		return 0, nil
	}
	return 0, nil
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
		fmt.Fprintf(os.Stderr, "Header was 0x%x instead of 0x%x\n", hdr, chk)
		return false
	}
	return true
}

// CheckStatus sends a Server Query Protocol request and prses the response
func CheckStatus(Cfg config) {
	a2sInfo := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x54, 0x53, 0x6F, 0x75, 0x72, 0x63, 0x65, 0x20, 0x45, 0x6E, 0x67, 0x69, 0x6E, 0x65, 0x20, 0x51, 0x75, 0x65, 0x72, 0x79, 0x00}
	a2sRules := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x56, 0xFF, 0xFF, 0xFF, 0xFF}
	a2sPlayer := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x55, 0xFF, 0xFF, 0xFF, 0xFF}

	server := Cfg.AtlasIP
	port := Cfg.AtlasQueryPort
	seconds := 15
	timeout := time.Duration(seconds) * time.Second

	if debug {
		fmt.Fprintln(os.Stderr, "Opening UDP connection...")
	}
	Conn, err := net.DialTimeout("udp", server+":"+port, timeout)
	if !CheckNoError(err) {
		os.Exit(2)
	}

	defer Conn.Close()

	// Get Info

	if debug {
		fmt.Fprintln(os.Stderr, "Sending A2S_INFO...")
	}

	start := time.Now()
	n, BytesReceived := SendPacket(Conn, a2sInfo, timeout)
	t := time.Now()
	elapsed1 := t.Sub(start)

	if BytesReceived == nil || n == 0 {
		fmt.Fprintln(os.Stderr, "Received no data!")
		os.Exit(2)
	}

	if !CheckHeader(BytesReceived[4], 0x49) {
		os.Exit(2)
	}

	if debug {
		fmt.Fprintf(os.Stderr, "HEADER: 0x%x\n", BytesReceived[4])
	}
	if debug {
		fmt.Fprintf(os.Stderr, "PROTOCOL: 0x%x\n", BytesReceived[5])
	}

	var sPtr int
	var info string
	sPtr = 5
	info, sPtr = GetString(BytesReceived, sPtr)
	fmt.Printf("NAME: %s\n", info)

	info, sPtr = GetString(BytesReceived, sPtr)
	fmt.Printf("MAP: %s\n", info)

	info, sPtr = GetString(BytesReceived, sPtr)
	fmt.Printf("FOLDER: %s\n", info)

	info, sPtr = GetString(BytesReceived, sPtr)
	fmt.Printf("GAME: %s\n", info)

	var id uint16
	id, sPtr = GetUInt16(BytesReceived, sPtr)
	fmt.Printf("ID: %d\n", id)

	fmt.Printf("PLAYERS: %d\n", BytesReceived[sPtr])
	sPtr++

	fmt.Printf("MAXPLAYERS: %d\n", BytesReceived[sPtr])
	sPtr++

	fmt.Printf("BOTS: %d\n", BytesReceived[sPtr])
	sPtr++

	fmt.Printf("SERVERTYPE: %c\n", BytesReceived[sPtr])
	sPtr++

	fmt.Printf("ENVIRONMENT: %c\n", BytesReceived[sPtr])
	sPtr++

	fmt.Printf("VISIBILITY: %d\n", BytesReceived[sPtr])
	sPtr++

	fmt.Printf("VAC: %d\n", BytesReceived[sPtr])
	sPtr++

	info, sPtr = GetString(BytesReceived, sPtr)
	fmt.Printf("VERSION: %s\n", info)

	if n > sPtr {
		// EDF
		edf := BytesReceived[sPtr]
		sPtr++

		// PORT
		if edf&0x80 != 0 {
			var port uint16
			port, _ = GetUInt16(BytesReceived, sPtr)
			fmt.Printf("PORT: %d\n", port)
		}

		// STEAMID
		if edf&0x10 != 0 {
			sPtr += 8
		}

		// Keywords
		if edf&0x20 != 0 {
			info, sPtr = GetString(BytesReceived, sPtr)
			fmt.Printf("KEYWORDS: %s\n", info)
		}
	}

	// Get Rules
	sPtr = 5

	if debug {
		fmt.Fprintln(os.Stderr, "Sending A2S_RULES...")
	}

	start = time.Now()
	n, BytesReceived = SendPacket(Conn, a2sRules, timeout)
	t = time.Now()
	elapsed2 := t.Sub(start)

	if BytesReceived == nil || n == 0 {
		fmt.Fprintln(os.Stderr, "Received no data!")
		os.Exit(2)
	}

	if !CheckHeader(BytesReceived[4], 0x41) {
		os.Exit(2)
	}

	// Challenge number
	var chnum uint32
	chnum, sPtr = GetUInt32(BytesReceived, sPtr)
	if debug {
		fmt.Fprintf(os.Stderr, "Challenge number: %d\n", chnum)
	}

	a2sRules[5] = byte(chnum)
	a2sRules[6] = byte(chnum >> 8)
	a2sRules[7] = byte(chnum >> 16)
	a2sRules[8] = byte(chnum >> 24)

	if debug {
		fmt.Fprintln(os.Stderr, "Sending A2S_RULES...")
	}

	start = time.Now()
	n, BytesReceived = SendPacket(Conn, a2sRules, timeout)
	t = time.Now()
	elapsed3 := t.Sub(start)

	if BytesReceived == nil || n == 0 {
		fmt.Fprintln(os.Stderr, "Received no data!")
		os.Exit(2)
	}

	elapsed := (elapsed1 + elapsed2 + elapsed3) / 3

	fmt.Printf("PING: %d\n", int(elapsed)/1000000)

	if !CheckHeader(BytesReceived[4], 0x45) {
		os.Exit(2)
	}

	// reset sPtr
	sPtr = 5
	var rules uint16
	rules, sPtr = GetUInt16(BytesReceived, sPtr)

	if rules > 0 {
		fmt.Println("RULE LIST:")
	}

	for i := uint16(0); i < rules; i++ {
		// Name
		info, sPtr = GetString(BytesReceived, sPtr)
		// Value
		val := ""
		val, sPtr = GetString(BytesReceived, sPtr)

		fmt.Printf("%s %s\n", info, val)
	}

	// Get Players
	sPtr = 5

	if debug {
		fmt.Fprintln(os.Stderr, "Sending A2S_PLAYER...")
	}
	n, BytesReceived = SendPacket(Conn, a2sPlayer, timeout)

	if BytesReceived == nil || n == 0 {
		fmt.Fprintln(os.Stderr, "Received no data!")
		os.Exit(2)
	}

	if !CheckHeader(BytesReceived[4], 0x41) {
		os.Exit(2)
	}

	// Challenge number
	chnum, sPtr = GetUInt32(BytesReceived, sPtr)
	if debug {
		fmt.Fprintf(os.Stderr, "Challenge number: %d\n", chnum)
	}

	a2sPlayer[5] = byte(chnum)
	a2sPlayer[6] = byte(chnum >> 8)
	a2sPlayer[7] = byte(chnum >> 16)
	a2sPlayer[8] = byte(chnum >> 24)

	if debug {
		fmt.Fprintln(os.Stderr, "Sending A2S_PLAYER...")
	}
	n, BytesReceived = SendPacket(Conn, a2sPlayer, timeout)

	if BytesReceived == nil || n == 0 {
		fmt.Fprintln(os.Stderr, "Received no data!")
		os.Exit(2)
	}

	if !CheckHeader(BytesReceived[4], 0x44) {
		os.Exit(2)
	}

	sPtr = 5
	players := BytesReceived[sPtr]
	sPtr++

	if players > 0 {
		fmt.Println("PLAYER LIST:")
	}

	var score uint32

	for i := 0; i < int(players); i++ {
		// Index (this seems to always be 0, so skipping it)
		sPtr++

		// Name
		info, sPtr = GetString(BytesReceived, sPtr)

		// Score
		score, sPtr = GetUInt32(BytesReceived, sPtr)

		// Duration
		b := []byte{0x00, 0x00, 0x00, 0x00}
		b[0] = BytesReceived[sPtr]
		sPtr++
		b[1] = BytesReceived[sPtr]
		sPtr++
		b[2] = BytesReceived[sPtr]
		sPtr++
		b[3] = BytesReceived[sPtr]
		sPtr++
		var duration float32
		buf := bytes.NewReader(b)
		err := binary.Read(buf, binary.LittleEndian, &duration)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Float conversion failed:", err)
		}

		fmt.Printf("%s %d %.0f\n", info, score, duration)
	}
	// r, err := rcon.Dial(server+":27025", "changeme")
	// r.Write()
	os.Exit(0)
}

func main() {

	Cfg := config{}
	err := env.Parse(&Cfg)

	if CheckNoError(err) == false {
		os.Exit(1)
	}

	fmt.Printf("Server IP: %v\n", Cfg.AtlasIP)
	fmt.Printf("Query Port: %v\n", Cfg.AtlasQueryPort)

	statusPtr := flag.Bool("status", false, "Check the status of the server")
	flag.Parse()
	if statusPtr != nil {
		CheckStatus(Cfg)
	}
	// argsWithProg := os.Args
	// if len(argsWithProg) < 3 {
	// 	fmt.Printf("Usage: %s <server> <port> or set ATLASIP and ATLASQUERYPORT Environment Variables\n", filepath.Base(argsWithProg[0]))
	// 	os.Exit(1)
	// }

	// server := argsWithProg[1]
	// port := argsWithProg[2]
	// if len(argsWithProg) > 3 {
	// 	debug = true
	// }
	// if len(argsWithProg) > 4 {
	// 	colorize = true
	// }

}
