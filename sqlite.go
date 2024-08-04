package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	_ "github.com/glebarez/go-sqlite"
)

var sqliteDb *sql.DB

func sqliteConnect() {
	connStr		:= fmt.Sprintf(os.Getenv("DB_FILE"))
	conn, err	:= sql.Open("sqlite", connStr)
	if err != nil {
		panic(err)
	}

	sqliteDb = conn
}

func sqliteClose() {
	err := sqliteDb.Close()
	if err != nil {
		panic(err)
	}
}

func sqliteInitialised(key string) bool {
	var table string
	switch key {
		case "COUNTRY": table = "ipv4_country"
		case "ASN": 	table = "ipv4_asn"
		case "CITY": 	table = "ipv4_city"
	}

	schema := sqliteGetOptionalSchema()

	var total int
	sqlString := fmt.Sprintf(`
		SELECT		COUNT(*) AS "total"
		
		FROM		%s"%s"
		`,
		schema, table)
	row := sqliteDb.QueryRow(sqlString)
	if err := row.Scan(&total); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}

		panic(err)
	}

	return total > 0
}

func sqliteIp(ip net.IP) *Ip {
	ipString	:= ip.String();
	ipVersion	:= getIpVersion(ipString)
	ipNumber	:= sqliteGetIpNumber(ipVersion, ipString)
	ipStruct	:= NewIp(ipString, ipVersion)
	schema		:= sqliteGetOptionalSchema()

	if hasCityDatabase() {
		sqlString := fmt.Sprintf(`
			SELECT		"country_code", 
						"state1", 
						"state2", 
						"city", 
						"postcode", 
						"latitude", 
						"longitude", 
						"timezone"
						
			FROM		%s"ipv%d_city"
			 
			WHERE		"ip_number_start" <= ? 

			ORDER BY	"ip_number_start" DESC
				 
			LIMIT		1`,
			schema, ipVersion)

		row := sqliteDb.QueryRow(sqlString, ipNumber)
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
			sqlString := fmt.Sprintf(`
				SELECT		"country_code" 
							
				FROM		%s"ipv%d_country"
				 
				WHERE		"ip_number_start" <= ? 
				
				ORDER BY	"ip_number_start" DESC
					 
				LIMIT		1`,
				schema, ipVersion)
			row := sqliteDb.QueryRow(sqlString, ipNumber)
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
		sqlString := fmt.Sprintf(`
			SELECT		"as_number",
						"as_organisation" 
						
			FROM		%s"ipv%d_asn"
			 
			WHERE		"ip_number_start" <= ? 

			ORDER BY	"ip_number_start" DESC
						 
			LIMIT		1`,
			schema, ipVersion)
		row := sqliteDb.QueryRow(sqlString, ipNumber)
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

func sqliteQueryMaxVersion(table string, ipVersion int) int {
	var version int

	schema := sqliteGetOptionalSchema()
	table = strings.Replace(table, "ip_", "ipv" + strconv.Itoa(ipVersion) + "_", 1)

	sqlString := fmt.Sprintf(`SELECT "db_version" FROM %s"%s" WHERE "ip_version" = ? ORDER BY "db_version" DESC LIMIT 1`, schema, table)
	row := sqliteDb.QueryRow(sqlString, ipVersion)
	if err := row.Scan(&version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0
		}

		panic(err)
	}

	return version
}

func sqliteDropOld(table string, ipVersion int, dbVersion int) {
	schema := sqliteGetOptionalSchema()
	table = strings.Replace(table, "ip_", "ipv" + strconv.Itoa(ipVersion) + "_", 1)

	fmt.Printf(`Dropping old SQLite data: %s"%s" ipv%d... `, schema, table, ipVersion)
	sqlString := fmt.Sprintf(`DELETE FROM %s"%s" WHERE "ip_version" = ? AND "db_version" < ?`, schema, table)
	_, err := sqliteDb.Exec(sqlString, ipVersion, dbVersion)
	if err != nil {
		panic(err)
	}
	fmt.Println("Complete")
}

func sqliteSaveCountries(countries []IpCountry) {
	var params []any

	schema := sqliteGetOptionalSchema()

	sqlString := fmt.Sprintf(
		`INSERT INTO %s"ipv%d_country" (
			"ip_range_start", 
			"ip_range_end", 
			"ip_number_start", 
			"ip_number_end", 
			"country_code", 
			"ip_version", 
			"db_version"
		) VALUES `,
		schema, countries[0].IpVersion)

	for _, country := range countries {
		ipNumberStart	:= sqliteGetIpNumber(country.IpVersion, country.IpRangeStart)
		ipNumberEnd		:= sqliteGetIpNumber(country.IpVersion, country.IpRangeEnd)

		sqlString += `(?, ?, ?, ?, ?, ?, ?), `
		params = append(params,
			country.IpRangeStart,
			country.IpRangeEnd,
			ipNumberStart,
			ipNumberEnd,
			country.CountryCode,
			country.IpVersion,
			country.DbVersion,
		)
	}
	sqlString = sqlString[0:len(sqlString) - 2]
	sqlString = fixPostgresVars(sqlString)

	stmt, err := sqliteDb.Prepare(sqlString)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		panic(err)
	}
}

func sqliteSaveASNs(ASNs []IpASN) {
	var params []any

	schema := sqliteGetOptionalSchema()

	sqlString := fmt.Sprintf(
		`INSERT INTO %s"ipv%d_asn" (
			"ip_range_start", 
			"ip_range_end", 
			"ip_number_start", 
			"ip_number_end", 
			"as_number", 
			"as_organisation", 
			"ip_version", 
			"db_version"
		) VALUES `,
		schema, ASNs[0].IpVersion)

	for _, asn := range ASNs {
		ipNumberStart	:= sqliteGetIpNumber(asn.IpVersion, asn.IpRangeStart)
		ipNumberEnd		:= sqliteGetIpNumber(asn.IpVersion, asn.IpRangeEnd)

		sqlString += `(?, ?, ?, ?, ?, ?, ?, ?), `
		params = append(params,
			asn.IpRangeStart,
			asn.IpRangeEnd,
			ipNumberStart,
			ipNumberEnd,
			asn.AsNumber,
			asn.AsOrganisation,
			asn.IpVersion,
			asn.DbVersion,
		)
	}
	sqlString = sqlString[0:len(sqlString) - 2]
	sqlString = fixPostgresVars(sqlString)

	stmt, err := sqliteDb.Prepare(sqlString)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		panic(err)
	}
}

func sqliteSaveCities(cities []IpCity) {
	var params []any

	schema := sqliteGetOptionalSchema()

	sqlString := fmt.Sprintf(
		`INSERT INTO %s"ipv%d_city" (
			"ip_range_start", 
			"ip_range_end", 
			"ip_number_start", 
			"ip_number_end", 
			"country_code", 
			"state1", 
			"state2", 
			"city", 
			"postcode", 
			"latitude", 
			"longitude", 
			"timezone", 
			"ip_version", 
			"db_version"
		) VALUES `,
		schema, cities[0].IpVersion)

	for _, city := range cities {
		ipNumberStart	:= sqliteGetIpNumber(city.IpVersion, city.IpRangeStart)
		ipNumberEnd		:= sqliteGetIpNumber(city.IpVersion, city.IpRangeEnd)

		sqlString += `(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?), `
		params = append(params,
			city.IpRangeStart,
			city.IpRangeEnd,
			ipNumberStart,
			ipNumberEnd,
			city.CountryCode,
			city.State1,
			city.State2,
			city.City,
			city.Postcode,
			city.Latitude,
			city.Longitude,
			city.Timezone,
			city.IpVersion,
			city.DbVersion,
		)
	}
	sqlString = sqlString[0:len(sqlString) - 2]
	sqlString = fixPostgresVars(sqlString)

	stmt, err := sqliteDb.Prepare(sqlString)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		panic(err)
	}
}

func sqliteFile(sqlPath string) {
	sqlBytes, err := os.ReadFile(sqlPath)
	if err != nil {
		panic(err)
	}
	sqlString := string(sqlBytes)
	if len(os.Getenv("DB_SCHEMA")) > 0 {
		sqlString = strings.Replace(sqlString, "${schema}", os.Getenv("DB_SCHEMA"), -1)
	} else {
		sqlString = strings.Replace(sqlString, `"${schema}".`, "", -1)
		sqlString = strings.Replace(sqlString, `${schema}:`, "", -1)
	}

	sqlStatements := strings.Split(sqlString, ";")

	transaction, err := sqliteDb.Begin()
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

func sqliteGetOptionalSchema() string {
	schema := ""
	if len(os.Getenv("DB_SCHEMA")) > 0 {
		schema = fmt.Sprintf(`"%s".`, os.Getenv("DB_SCHEMA"))
	}

	return schema
}

func sqliteGetIpNumber(ipVersion int, ipString string) any {
	var ipNumber any

	if ipVersion == 4 {
		ipNumber = ipv4ToNumber(ipString)
	} else {
		ipNumber = ipv6ToNumber(ipString)
	}

	return ipNumber
}