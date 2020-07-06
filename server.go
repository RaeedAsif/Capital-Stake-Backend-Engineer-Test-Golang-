package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"

	covidlib "Capital-Stake-Test/lib"
)

var (
	reports = covidlib.Fetch("data/covid_final_data.csv")
)

func main() {
	var addr string
	var network string
	flag.StringVar(&addr, "e", ":4040", "service endpoint [ip addr or socket path]")
	flag.StringVar(&network, "n", "tcp", "network protocol [tcp,unix]")
	flag.Parse()

	// validate supported network protocols
	switch network {
	case "tcp", "tcp4", "tcp6", "unix":
	default:
		log.Fatalln("unsupported network protocol:", network)
	}

	// create a listener for provided network and host address
	ln, err := net.Listen(network, addr)
	if err != nil {
		log.Fatal("failed to create listener:", err)
	}
	defer ln.Close()
	log.Println("**** Capital Stake Backend engineer test server ***")
	log.Printf("Service started: (%s) %s\n", network, addr)

	// connection-loop - handle incoming requests
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			if err := conn.Close(); err != nil {
				log.Println("failed to close listener:", err)
			}
			continue
		}
		log.Println("Connected to", conn.RemoteAddr())

		go handleConnection(conn)
	}
}

//handleConnection :to handle incoming connection
func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("error closing connection:", err)
		}
	}()

	if _, err := conn.Write([]byte("Connected...\nUsage: GET Date or Region\n")); err != nil {
		log.Println("error writing:", err)
		return
	}

	// loop to stay connected with client until client breaks connection
	for {
		// buffer for client command
		cmdLine := make([]byte, (1024 * 4))
		n, err := conn.Read(cmdLine)
		if n == 0 || err != nil {
			log.Println("connection read error:", err)
			return
		}
		param := parseCommand(string(cmdLine[0:n]))

		if param == "" {
			if _, err := conn.Write([]byte("Invalid command\n")); err != nil {
				log.Println("failed to write:", err)
				return
			}
			continue
		}

		result := covidlib.Query(reports, param)
		if len(result) == 0 {
			if _, err := conn.Write([]byte("Nothing found\n")); err != nil {
				log.Println("failed to write:", err)
			}
			continue
		}
		// send each currency info as a line to the client wiht fmt.Fprintf()

		rsp, _ := json.Marshal(result)
		_, jrr := conn.Write([]byte(

			fmt.Sprintf(
				"{ \"response\":  %s }\n",
				string(rsp),
			),
		))
		if jrr != nil {
			log.Println("failed to write response:", err)
			return
		}

		// execute command
	}
}

//parseCommand :to parse incoming input
func parseCommand(cmdLine string) (param string) {
	var result = map[string]map[string]string{}
	json.Unmarshal([]byte(cmdLine), &result)

	// Print the data type of result variable
	switch {
	case result["query"]["region"] != "":
		param = result["query"]["region"]
	case result["query"]["date"] != "":
		param = result["query"]["date"]
	default:
		param = ""
	}

	return
}
