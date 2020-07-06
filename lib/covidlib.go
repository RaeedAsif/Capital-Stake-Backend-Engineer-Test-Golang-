package covidlib

import (
	"encoding/csv"
	"io"
	"os"
	"strings"
	"time"
)

type covidstats struct {
	Date                     string `json:"date"`
	CumulativeTestPostive    string `json:"positive"`
	CumulativeTestsPerformed string `json:"tests"`
	Expired                  string `json:"expired"`
	StillAdmitted            string `json:"admitted"`
	Discharged               string `json:"discharged"`
	Region                   string `json:"region"`
}

//Fetch :function to load and fetch csv file into table
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

//Query :function to query the table and return results
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

//DateFormat: to format date to said format
func DateFormat(DateString string) string {
	if string(DateString[2]) == "-" {
		myDate, _ := time.Parse("02-Jan-2006", DateString)
		str := myDate.Format("2006-01-02")
		return str
	}
	if string(DateString[2]) == "/" {
		myDate, _ := time.Parse("02/01/2006", DateString)
		str := myDate.Format("2006-01-02")
		return str
	}
	return DateString
}
