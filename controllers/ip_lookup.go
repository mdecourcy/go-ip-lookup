package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"ip_lookup/models"
	"log"
	"net"
	"os"

	"github.com/ip2location/ip2location-go"
	"github.com/oschwald/geoip2-golang"
)

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func maxMindLookup() {
	db, err := geoip2.Open("./data/db/GeoLite2-City.mmdb")
	var result_slice []string

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	file, err := os.Open("./data/ip_addresses.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var state string
		ip := net.ParseIP(scanner.Text())
		record, err := db.City(ip)
		if err != nil {
			log.Fatal(err)
		}
		if len(record.Subdivisions) > 0 {
			state = record.Subdivisions[0].Names["en"]
		}
		result_string := fmt.Sprintf("%v:%v:%v:%v", record.Country.IsoCode, state, record.City.Names["en"], scanner.Text())
		result_slice = append(result_slice, result_string)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	writeLines(result_slice, "./data/output/output.csv")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func SingleIPLookup(ipAddress string) (jsonResult string) {
	db, err := ip2location.OpenDB("./data/db/IP2LOCATION-LITE-DB9.BIN")

	check(err)

	results, err := db.Get_all(ipAddress)

	if err != nil {
		return ""
	}

	if results.Country_short == "Invalid IP address." {
		return "Invalid IP address."
	}

	resp := models.IPResposne{}

	resp.IP = ipAddress
	resp.Geo.Country = results.Country_short
	resp.Geo.Region = results.Region
	resp.Geo.City = results.City
	resp.Geo.PostalCode = results.Zipcode
	resp.Geo.Latitude = results.Latitude
	resp.Geo.Longitude = results.Longitude

	jsonString, err := json.Marshal(resp)

	return string(jsonString)

}

func IP2LocLookup() {
	var result_slice []string
	db, err := ip2location.OpenDB("./data/db/IP2LOCATION-LITE-DB9.BIN")
	if err != nil {
		fmt.Print(err)
		return
	}

	file, err := os.Open("./data/ip_addresses.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		ip := scanner.Text()
		results, err := db.Get_all(ip)

		if err != nil {
			fmt.Print(err)
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		result_string := fmt.Sprintf("%v:%v:%v:%v:%v", results.Country_short, results.Region, results.City, results.Zipcode, scanner.Text())
		result_slice = append(result_slice, result_string)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	writeLines(result_slice, "./data/output/ip2Loc.csv")

}
