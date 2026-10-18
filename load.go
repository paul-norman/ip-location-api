package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

func loadCheckInitialised() (bool, []string) {
	initialised := true
	var missing []string

	if hasCountryDatabase() {
		if !dbInitialised("COUNTRY") {
			missing		= append(missing, "COUNTRY")
			initialised	= false
		}
	}

	if hasASNDatabase() {
		if !dbInitialised("ASN") {
			missing		= append(missing, "ASN")
			initialised	= false
		}
	}

	if hasCityDatabase() {
		if !dbInitialised("CITY") {
			missing		= append(missing, "CITY")
			initialised	= false
		}
	}

	return initialised, missing
}

func loadData(dataToLoad []DataToLoad) {
	for _, item := range dataToLoad {
		switch item.Download.Type {
			case "CITY":	loadCities(item)
			case "ASN":		loadASNs(item)
			case "COUNTRY":	loadCountries(item)
		}
	}
}

func loadCities(dataToLoad DataToLoad) {
	csvFile, err := os.Open(dataToLoad.Path)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	csvFileReader := csv.NewReader(csvFile)

	version 	:= dbQueryMaxVersion("ip_city", dataToLoad.Version) + 1
	cities		:= []IpCity{}
	numSaved	:= 0;

	fmt.Println("rebuilding: ip_city ipv", dataToLoad.Version)
	fmt.Print("\033[s") // Save the cursor position
	for {
		record, err := csvFileReader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
		lat, _ := strconv.ParseFloat(record[7], 64)
		lon, _ := strconv.ParseFloat(record[8], 64)

		cities = append(cities, IpCity{ record[0], record[1], record[2], record[3], record[4], record[5], record[6], lat, lon, record[9], dataToLoad.Version, version })

		if len(cities) == 100 {
			dbSaveCities(cities)
			logEntriesConditionally(&numSaved, &cities)
			cities = []IpCity{}
		}
	}

	if len(cities) > 0 {
		dbSaveCities(cities)
		numSaved += len(cities)
	}
	logEntries(numSaved)

	dbDropOld("ip_city", dataToLoad.Version, version)
}

func loadASNs(dataToLoad DataToLoad) {
	csvFile, err := os.Open(dataToLoad.Path)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	csvFileReader := csv.NewReader(csvFile)

	version 	:= dbQueryMaxVersion("ip_asn", dataToLoad.Version) + 1
	ASNs		:= []IpASN{}
	numSaved	:= 0;

	fmt.Println("rebuilding: ip_asn ipv", dataToLoad.Version)
	fmt.Print("\033[s") // Save the cursor position
	for {
		record, err := csvFileReader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
		asn, _ := strconv.Atoi(record[2])
		ASNs = append(ASNs, IpASN{ record[0], record[1], asn, record[3], dataToLoad.Version, version })

		if len(ASNs) == 100 {
			dbSaveASNs(ASNs)
			logEntriesConditionally(&numSaved, &ASNs)
			ASNs = []IpASN{}
		}
	}

	if len(ASNs) > 0 {
		dbSaveASNs(ASNs)
		numSaved += len(ASNs)
	}
	logEntries(numSaved)

	dbDropOld("ip_asn", dataToLoad.Version, version)
}

func loadCountries(dataToLoad DataToLoad) {
	csvFile, err := os.Open(dataToLoad.Path)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	csvFileReader := csv.NewReader(csvFile)

	version		:= dbQueryMaxVersion("ip_country", dataToLoad.Version) + 1
	countries	:= []IpCountry{}
	numSaved	:= 0

	fmt.Println("rebuilding: ip_country ipv", dataToLoad.Version)
	fmt.Print("\033[s") // Save the cursor position
	for {
		record, err := csvFileReader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
		countries = append(countries, IpCountry{ record[0], record[1], record[2], dataToLoad.Version, version })

		if len(countries) == 100 {
			dbSaveCountries(countries)
			logEntriesConditionally(&numSaved, &countries)
			countries = []IpCountry{}
		}
	}

	if len(countries) > 0 {
		dbSaveCountries(countries)
		numSaved += len(countries)
	}
	logEntries(numSaved)

	dbDropOld("ip_country", dataToLoad.Version, version)
}

func loadDbStructure() {
	dbFile()
}

func logEntries(numSaved int) {
	fmt.Print("\033[u\033[K") // Restore the cursor position and clear the line
	fmt.Printf("Saved: %d entries\n", numSaved)
}

// Don't spam the logs, only display every message after `logFreq` records
func logEntriesConditionally[T any](numSaved *int, list *[]T) {
	logFreq := getLogFrequency()

	previous	:= *numSaved
	*numSaved	+= len(*list)

	if *numSaved / logFreq > previous / logFreq {
		logEntries(*numSaved)
	}
}