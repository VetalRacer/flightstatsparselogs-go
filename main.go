package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var input string

type Item struct {
	FullString       string
	Flight           string
	Date             string
	TypeRequest      string
	TypeOwnerRequest string
}

func main() {
	var dir string
	var dirSet string

	flag.StringVar(&dirSet, "d", "current", "Set directory")
	flag.Parse()

	switch dirSet {
	case "echo":
		dir = "E:/core-log/"
	default:
		dir = "."
	}

	files, err := readDir("flightdata_flightstats.log", dir)
	if err != nil {
		fmt.Printf("findFiles: %s %s", err, dir)
		_, _ = fmt.Scanf("%v", &input)
		os.Exit(1)
	}

	lines, err := readLines(dir, files)
	if err != nil {
		fmt.Printf("readLines: %s", err)
		_, _ = fmt.Scanf("%v", &input)
		os.Exit(1)
	}

	if err := writeToFile(parseItem(lines)); err != nil {
		fmt.Printf("writeToFile: %s", err)
		_, _ = fmt.Scanf("%v", &input)
		os.Exit(1)
	} else {

		fmt.Println("Success! Please find FS.csv file in current folder.")
		_, _ = fmt.Scanf("%v", &input)
	}
}

func readDir(mask, dir string) ([]string, error) {
	var Files []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), mask) {
			Files = append(Files, file.Name())
		}
	}
	if Files != nil {
		return Files, nil
	}
	return Files, fmt.Errorf("Unable to find files in dir: ")
}

func writeToFile(data []Item) error {
	if data == nil {
		return fmt.Errorf("no data found")
	}

	file, err := os.Create("FS_statistic.csv")
	if err != nil {
		log.Fatal("Unable to create file:", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	header := []string{"The whole request", "Flight", "DateRequest", "TypeRequest", "TypeOwnerRequest"}
	write(header, file)

	for _, v := range data {

		s := []string{
			v.FullString,
			v.Flight,
			v.Date,
			v.TypeRequest,
			v.TypeOwnerRequest,
		}

		write(s, file)
	}
	return nil
}

func parseItem(lines []string) []Item {
	var dataItem []Item

	for _, v := range lines {

		var typeOwnerRequest, typeRequest string
		var flight []byte

		reDate := regexp.MustCompile(`^[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]`)
		reFlight := regexp.MustCompile(`(v2/json/flight/status|v1/json/flight)/(?P<flightId>[A-Za-z0-9]+)/(?P<flightNum>[0-9]+)/(?:arriving|arr)/(?P<flightDate>[0-9]{4}/[0-9]{1,2}/[0-9]{1,2})`)

		templateFlight := "FlightNum: $flightId$flightNum | Date: $flightDate"

		for _, submatches := range reFlight.FindAllStringSubmatchIndex(v, -1) {
			flight = reFlight.ExpandString(flight, templateFlight, v, submatches)
		}

		date := reDate.FindAllString(v, -1)

		switch strings.Contains(v, "schedules") {
		case true:
			typeRequest = "Schedules"
		case false:
			typeRequest = "Status"
		}

		switch strings.Contains(v, "user - echo") {
		case true:
			typeOwnerRequest = "BookingRequest"
		case false:
			typeOwnerRequest = "UsersRequest"
		}
		addItem(v, string(flight), date[0], typeRequest, typeOwnerRequest, &dataItem)

	}
	return dataItem
}

func addItem(fullString, flight, date, typeR, typeOR string, data *[]Item) {
	var item = new(Item)
	item.FullString = fullString
	item.Flight = flight
	item.Date = date
	item.TypeRequest = typeR
	item.TypeOwnerRequest = typeOR

	*data = append(*data, *item)
}

func readLines(dir string, files []string) ([]string, error) {
	var lines []string

	if dir == "." {
		dir = ""
	}

	for _, filePath := range files {
		file, err := os.Open(dir + filePath)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "result") {
				lines = append(lines, scanner.Text())
			}
		}

		if file.Close() != nil {
			fmt.Print(err)
		}

	}

	return lines, nil
}

//Write to CSV file
func write(s []string, f *os.File) {

	w := csv.NewWriter(f)

	if err := w.Write(s); err != nil {
		fmt.Printf("Error writing csv: %v\n", err)
	}
	w.Flush()
}
