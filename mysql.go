package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	mysql "github.com/go-sql-driver/mysql"
)

var mysqlDb *sql.DB

func mysqlConnect() {
	config := mysql.Config{
		User:					os.Getenv("DB_USER"),
		Passwd:					os.Getenv("DB_PASS"),
		Addr:					os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		DBName:					os.Getenv("DB_NAME"),
		Net:					"tcp",
		AllowNativePasswords:	true,
	}

	conn, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		panic(err)
	}

	mysqlDb = conn
}

func mysqlClose() {
	err := mysqlDb.Close()
	if err != nil {
		panic(err)
	}
}

func mysqlInitialised(key string) bool {
	var table string
	switch key {
		case "COUNTRY": table = "ip_country"
		case "ASN": 	table = "ip_asn"
		case "CITY": 	table = "ip_city"
	}

	var total int
	sqlString := fmt.Sprintf("SELECT COUNT(*) AS `total` FROM `%s` ", table)
	row := mysqlDb.QueryRow(sqlString)
	if err := row.Scan(&total); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}

		panic(err)
	}

	return total > 0
}

func mysqlIp(ip net.IP) *Ip {
	ipString	:= ip.String();
	ipVersion	:= getIpVersion(ipString)
	function	:= mysqlGetConversionFunction(ipVersion)
	ipStruct	:= NewIp(ipString, ipVersion)

	if hasCityDatabase() {
		// Could also use: `SET SESSION sql_mode = 'ANSI_QUOTES';` but, meh
		sqlString := fmt.Sprintf(`
			SELECT		`+"`"+`country_code`+"`"+`,
						`+"`"+`state1`+"`"+`,
						`+"`"+`state2`+"`"+`,
						`+"`"+`city`+"`"+`,
						`+"`"+`postcode`+"`"+`,
						`+"`"+`latitude`+"`"+`,
						`+"`"+`longitude`+"`"+`,
						`+"`"+`timezone`+"`"+`
						
			FROM		`+"`"+`ip_city`+"`"+`
			 
			WHERE		`+"`"+`ip_number_start`+"`"+` <= %s(?) 
			
			ORDER BY	`+"`"+`ip_number_start`+"`"+` DESC
			 
			LIMIT 		1`,
			function)

		row := mysqlDb.QueryRow(sqlString, ipString)
		if err := row.Scan(&ipStruct.CountryCode, &ipStruct.State1, &ipStruct.State2, &ipStruct.City, &ipStruct.Postcode, &ipStruct.Latitude, &ipStruct.Longitude, &ipStruct.Timezone); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ipStruct
			}

			panic(err)
		}

		if len(ipStruct.City) > 0 {
			ipStruct.FoundCity = true
		}
	}

	if len(ipStruct.CountryCode) > 0 {
		ipStruct.FoundCountry = true
	} else {
		if hasCountryDatabase() {
			// Could also use: `SET SESSION sql_mode = 'ANSI_QUOTES';` but, meh
			sqlString := fmt.Sprintf(`
				SELECT		`+"`"+`country_code`+"`"+` 
							
				FROM		`+"`"+`ip_country`+"`"+`
				 
				WHERE		`+"`"+`ip_number_start`+"`"+` <= %s(?)
				
				ORDER BY	`+"`"+`ip_number_start`+"`"+` DESC
				 
				LIMIT		1`,
				function)
			row := mysqlDb.QueryRow(sqlString, ipString)
			if err := row.Scan(&ipStruct.CountryCode); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return ipStruct
				}

				panic(err)
			}

			if len(ipStruct.CountryCode) > 0 {
				ipStruct.FoundCountry = true
			}
		}
	}

	if hasASNDatabase() {
		// Could also use: `SET SESSION sql_mode = 'ANSI_QUOTES';` but, meh
		sqlString := fmt.Sprintf(`
			SELECT		`+"`"+`as_number`+"`"+`,
						`+"`"+`as_organisation`+"`"+` 
						
			FROM		`+"`"+`ip_asn`+"`"+`
			 
			WHERE		`+"`"+`ip_number_start`+"`"+` <= %s(?) 

			ORDER BY	`+"`"+`ip_number_start`+"`"+` DESC
			 
			LIMIT		1`,
			function)
		row := mysqlDb.QueryRow(sqlString, ipString)
		if err := row.Scan(&ipStruct.OrganisationNumber, &ipStruct.OrganisationName); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ipStruct
			}

			panic(err)
		}

		if len(ipStruct.OrganisationName) > 0 {
			ipStruct.FoundASN = true
		}
	}

	return ipStruct
}

func mysqlQueryMaxVersion(table string, ipVersion int) int {
	var version int

	sqlString := fmt.Sprintf("SELECT `db_version` FROM `%s` WHERE `ip_version` = ? ORDER BY `db_version` DESC LIMIT 1", table)
	row := mysqlDb.QueryRow(sqlString, ipVersion)
	if err := row.Scan(&version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0
		}

		panic(err)
	}

	return version
}

func mysqlDropOld(table string, ipVersion int, dbVersion int) {
	fmt.Printf("Dropping old MySQL data: `%s` ipv%d... ", table, ipVersion)
	sqlString := fmt.Sprintf("DELETE FROM `%s` WHERE `ip_version` = ? AND `db_version` < ?", table)
	_, err := mysqlDb.Exec(sqlString, ipVersion, dbVersion)
	if err != nil {
		panic(err)
	}
	fmt.Println("Complete")
}

func mysqlSaveCountries(countries []IpCountry) {
	var params []any

	function := mysqlGetConversionFunction(countries[0].IpVersion)
	sqlString := "INSERT INTO `ip_country` (`ip_range_start`, `ip_range_end`, `ip_number_start`, `ip_number_end`, `country_code`, `ip_version`, `db_version`) VALUES "
	for _, country := range countries {
		sqlString += `(?, ?, ` + function + `(?), ` + function + `(?), ?, ?, ?), `
		params = append(params, country.IpRangeStart, country.IpRangeEnd, country.IpRangeStart, country.IpRangeEnd, country.CountryCode, country.IpVersion, country.DbVersion)
	}
	sqlString = sqlString[0:len(sqlString) - 2]

	stmt, err := mysqlDb.Prepare(sqlString)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		panic(err)
	}
}

func mysqlSaveASNs(ASNs []IpASN) {
	var params []any

	function := mysqlGetConversionFunction(ASNs[0].IpVersion)
	sqlString := "INSERT INTO `ip_asn` (`ip_range_start`, `ip_range_end`, `ip_number_start`, `ip_number_end`, `as_number`, `as_organisation`, `ip_version`, `db_version`) VALUES "
	for _, asn := range ASNs {
		sqlString += `(?, ?, ` + function + `(?), ` + function + `(?), ?, ?, ?, ?), `
		params = append(params, asn.IpRangeStart, asn.IpRangeEnd, asn.IpRangeStart, asn.IpRangeEnd, asn.AsNumber, asn.AsOrganisation, asn.IpVersion, asn.DbVersion)
	}
	sqlString = sqlString[0:len(sqlString) - 2]

	stmt, err := mysqlDb.Prepare(sqlString)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		panic(err)
	}
}

func mysqlSaveCities(cities []IpCity) {
	var params []any

	function := mysqlGetConversionFunction(cities[0].IpVersion)
	sqlString := "INSERT INTO `ip_city` (`ip_range_start`, `ip_range_end`, `ip_number_start`, `ip_number_end`, `country_code`, `state1`, `state2`, `city`, `postcode`, `latitude`, `longitude`, `timezone`, `ip_version`, `db_version`) VALUES "
	for _, city := range cities {
		sqlString += `(?, ?, ` + function + `(?), ` + function + `(?), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?), `
		params = append(params, city.IpRangeStart, city.IpRangeEnd, city.IpRangeStart, city.IpRangeEnd, city.CountryCode, city.State1, city.State2, city.City, city.Postcode, city.Latitude, city.Longitude, city.Timezone, city.IpVersion, city.DbVersion)
	}
	sqlString = sqlString[0:len(sqlString) - 2]

	stmt, err := mysqlDb.Prepare(sqlString)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		panic(err)
	}
}

func mysqlFile(sqlPath string) {
	sqlBytes, err := os.ReadFile(sqlPath)
	if err != nil {
		panic(err)
	}
	sqlString := string(sqlBytes)

	sqlStatements := strings.Split(sqlString, ";")

	transaction, err := mysqlDb.Begin()
	if err != nil {
		panic(err)
	}
	defer transaction.Rollback()

	for _, sqlStatement := range sqlStatements {
		sqlStatement = strings.Trim(sqlStatement, " 	")
		if len(sqlStatement) > 0 {
			_, err := transaction.Exec(sqlStatement)
			if err != nil {

				panic(err)
			}
		}
	}

	err = transaction.Commit()
	if err != nil {
		panic(err)
	}
}

// Probably don't actually need to do this
func mysqlGetConversionFunction(ipVersion int) string {
	function := "INET_ATON"
	if ipVersion == 6 {
		function = "INET6_ATON"
	}

	return function
}