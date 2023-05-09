package main

//testing

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	reports = Fetch("data/covid_final_data.csv")
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

// ”'Library for setting and hosting server”'
// handleConnection :to handle incoming connection
func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("error closing connection:", err)
		}
	}()

	if _, err := conn.Write([]byte("Connected...\nUsage:  Send Date or Region in JSON fromat\n")); err != nil {
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

		result := Query(reports, param)
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

// parseCommand :to parse incoming input
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

type covidstats struct {
	Date                     string `json:"date"`
	CumulativeTestPostive    string `json:"positive"`
	CumulativeTestsPerformed string `json:"tests"`
	Expired                  string `json:"expired"`
	StillAdmitted            string `json:"admitted"`
	Discharged               string `json:"discharged"`
	Region                   string `json:"region"`
}

// ”'Library for parsing and arranging data”'
// Fetch :function to load and fetch csv file into table
func Fetch(path string) []covidstats {
	table := make([]covidstats, 0)
	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err.Error())
		}
		c := covidstats{
			Date:                     DateFormat(row[4]),
			CumulativeTestPostive:    row[2],
			CumulativeTestsPerformed: row[3],
			Expired:                  row[6],
			StillAdmitted:            row[10],
			Discharged:               row[5],
			Region:                   row[9],
		}
		table = append(table, c)
	}
	return table
}

// Query :function to query the table and return results
func Query(table []covidstats, filter string) []covidstats {
	result := make([]covidstats, 0)
	filter = strings.ToLower(filter)
	for _, cov := range table {
		region := strings.ToLower(cov.Region)
		if cov.Date == filter || region == filter {
			result = append(result, cov)
		}
	}
	return result
}

// DateFormat: to format date to said format
func DateFormat(DateString string) string {
	var strDate = string(DateString[2])
	switch {
	case strDate == "-":
		myDate, _ := time.Parse("02-Jan-06", DateString)
		str := myDate.Format("2006-01-02")
		return str
	case strDate == "/":
		myDate, _ := time.Parse("02/01/2006", DateString)
		str := myDate.Format("2006-01-02")
		return str
	default:
		return DateString
	}
}
