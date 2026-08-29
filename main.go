package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// port pearl server runs on, might change this later idk
const PORT = "7777"

func main() {
	// Seed the random number generator so facts change every time
	rand.Seed(time.Now().UnixNano())

	ln, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		fmt.Println("couldnt start server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("pearlS is running on port", PORT)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("connection failed:", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)

	line1, err := r.ReadString('\n')
	if err != nil {
		return
	}
	line1 = strings.TrimSpace(line1)

	// check if its actually a pearl request or just some random http thing
	if line1 == "PEARL/1.0" {
		doPearlStuff(conn, r)
	} else {
		// not a pearl client, probably a browser or something. send fallback page
		doHttpFallback(conn, r, line1)
	}
}

func doPearlStuff(conn net.Conn, r *bufio.Reader) {
	worldId := ""
	contentLen := 0

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimSpace(line)

		if line == "" {
			break // done reading headers
		}

		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			continue
		}

		if parts[0] == "worldId" {
			worldId = parts[1]
		}
		if parts[0] == "Content-Length" {
			n, err := strconv.Atoi(parts[1])
			if err == nil {
				contentLen = n
			}
		}
	}

	content := make([]byte, contentLen)
	r.Read(content) // TODO probably should check the error here

	fmt.Println("got a pearl for world", worldId, "-", string(content))

	// this is where the actual game logic should go eventually.
	// for now it just says it got it
	resp := "Pearl received for world " + worldId
	sendPearlResponse(conn, 200, resp)
}

func sendPearlResponse(conn net.Conn, status int, body string) {
	fmt.Fprintf(conn, "PEARL/1.0 %d\r\n", status)
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprintf(conn, "\r\n")
	fmt.Fprintf(conn, body)
}

func doHttpFallback(conn net.Conn, r *bufio.Reader, reqLine string) {
	// reqLine looks like "GET /somepath HTTP/1.1", just want the path part
	pieces := strings.Fields(reqLine)
	path := "/"
	if len(pieces) >= 2 {
		path = pieces[1]
	}

	var page string

	if path == "/pearldashboard" {
		// Define the trivia facts
		facts := []string{
			"Pearl was originally called 'EnderManAuth' (EMAuth, as we called it) and it wasnt a server software. It was a simple SSE auth system before we killeed it off and switched to pearl.",
			"Pearl was originally made in minecraft MakeCode PXT.",
			"PearlS runs entirely on a custom TCP protocol (PEARL/1.0), bypassing HTTP overhead for faster streams.",
			"MakeCode extensions compile down to static TypeScript (STS), which is why Pearl's early block days were so strict!",
		}
		
		// Pick a random fact
		fact := facts[rand.Intn(len(facts))]

		// Get real server stats
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		
		// NumGoroutine gives us a rough estimate of active connections + background tasks
		activeRoutines := runtime.NumGoroutine()

		page = fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>PearlS Dashboard</title>
	<style>
		body { font-family: system-ui, sans-serif; background: #18181b; color: #e4e4e7; padding: 2rem; }
		.card { background: #27272a; padding: 2rem; border-radius: 8px; max-width: 600px; margin: 0 auto; box-shadow: 0 4px 15px rgba(0,0,0,0.5); }
		h1 { color: #a78bfa; margin-top: 0; }
		.stat { font-size: 1.2rem; margin: 12px 0; border-bottom: 1px solid #3f3f46; padding-bottom: 8px; }
		.value { float: right; font-weight: bold; }
		.online { color: #4ade80; }
		.fact-box { margin-top: 30px; padding: 15px; background: #3f3f46; border-left: 4px solid #a78bfa; border-radius: 4px; }
	</style>
</head>
<body>
	<div class="card">
		<h1>🔮 PearlS Dashboard</h1>
		
		<div class="stat">Server Status: <span class="value online">Online</span></div>
		<div class="stat">Port: <span class="value">%s</span></div>
		<div class="stat">Active Threads: <span class="value">%d</span></div>
		<div class="stat">Memory Allocated: <span class="value">%v MB</span></div>
		
		<div class="fact-box">
			<strong>Did you know?</strong><br>
			<i>%s</i>
		</div>
	</div>
</body>
</html>`, PORT, activeRoutines, m.Alloc/1024/1024, fact)

	} else {
		page = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Not a Pearl client</title>
</head>
<body>
    <p>Error: not a Pearl client.</p>
    <p>This server speaks the pearls:// protocol, not plain HTTP. Connect using the Pearl protocol in JavaScript instead.</p>
    <p><a href="/pearldashboard">Open the Pearl dashboard</a></p>

    <script>
        console.log("Endermite spawned");
    </script>
</body>
</html>`
	}

	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Type: text/html\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(page))
	fmt.Fprintf(conn, "Connection: close\r\n")
	fmt.Fprintf(conn, "\r\n")
	fmt.Fprintf(conn, page)
}
