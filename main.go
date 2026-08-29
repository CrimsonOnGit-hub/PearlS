package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const PORT = "7777"

func main() {
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

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}

	request := string(buf[:n])
	lines := strings.Split(request, "\n")

	firstLine := strings.TrimSpace(lines[0])

	if firstLine == "PEARL/1.0" {
		worldId := ""
		contentLen := 0
		bodyStart := -1

		for i := 1; i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])

			if line == "" {
				bodyStart = i + 1
				break
			}

			if strings.HasPrefix(line, "worldId:") {
				worldId = strings.TrimSpace(line[len("worldId:"):])
			}

			if strings.HasPrefix(line, "Content-Length:") {
				lenStr := strings.TrimSpace(line[len("Content-Length:"):])
				n, _ := strconv.Atoi(lenStr)
				contentLen = n
			}
		}

		body := ""
		if bodyStart != -1 && bodyStart < len(lines) {
			body = strings.Join(lines[bodyStart:], "\n")
			if len(body) > contentLen {
				body = body[:contentLen]
			}
		}

		fmt.Println("got pearl for world " + worldId + " content: " + body)

		respBody := "Pearl received for world " + worldId
		resp := "PEARL/1.0 200\r\nContent-Length: " + strconv.Itoa(len(respBody)) + "\r\n\r\n" + respBody
		conn.Write([]byte(resp))

	} else {
		path := "/"
		parts := strings.Split(firstLine, " ")
		if len(parts) > 1 {
			path = parts[1]
		}

		page := ""
		if path == "/pearldashboard" {
			page = "<!DOCTYPE html><html><body><p>dashboard goes here (WIP)</p></body></html>"
		} else {
			page = "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<meta charset=\"UTF-8\">\n<title>Not a Pearl client</title>\n</head>\n<body>\n    <p>Error: not a Pearl client.</p>\n    <p>This server speaks the pearls:// protocol, not plain HTTP. Connect using the Pearl protocol in JavaScript instead.</p>\n    <p><a href=\"/pearldashboard\">Open the Pearl dashboard</a></p>\n\n    <script>\n        console.log(\"Endermite spawned\");\n    </script>\n</body>\n</html>\n"
		}

		resp := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nContent-Length: " + strconv.Itoa(len(page)) + "\r\nConnection: close\r\n\r\n" + page
		conn.Write([]byte(resp))
	}
}
