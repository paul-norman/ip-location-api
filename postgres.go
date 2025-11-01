package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

var pgDb *sql.DB

func postgresConnect() {
	connStr		:= fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	conn, err	:= sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	pgDb = conn
}

func postgresClose() {
	err := pgDb.Close()
	if err != nil {
		panic(err)
	}
}

func postgresInitialised(key string) bool {
	var table string
	switch key {
		case "COUNTRY": table = "ip_country"
		case "ASN": 	table = "ip_asn"
		case "CITY": 	table = "ip_city"
	}

	var total int
	sqlString := fmt.Sprintf(`
		SELECT		COUNT(*) AS "total"
		
		FROM		"%s"."%s"
		`,
		os.Getenv("DB_SCHEMA"), table)
	row := pgDb.QueryRow(sqlString)
	if err := row.Scan(&total); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}

		panic(err)
	}

	return total > 0
}

func postgresIp(ip net.IP) *Ip {
	ipString	:= ip.String();
	ipVersion	:= getIpVersion(ipString)
	ipStruct	:= NewIp(ipString, ipVersion)

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
						
			FROM		"%s"."ip_city"
			 
			WHERE		"ip_range_start"	<= $1 
			AND			"ip_range_end"		>= $1

			ORDER BY	"ip_range_start" DESC
			 
			LIMIT		1`,
			os.Getenv("DB_SCHEMA"))
		row := pgDb.QueryRow(sqlString, ipString)
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
							
				FROM		"%s"."ip_country"
				 
				WHERE		"ip_range_start"	<= $1
				AND			"ip_range_end"		>= $1
				
				ORDER BY	"ip_range_start" DESC
					 
				LIMIT		1`,
				os.Getenv("DB_SCHEMA"))
			row := pgDb.QueryRow(sqlString, ipString)
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
						
			FROM		"%s"."ip_asn"
			 
			WHERE		"ip_range_start"	<= $1
			AND			"ip_range_end"		>= $1
			
			ORDER BY	"ip_range_start" DESC
			 
			LIMIT		1`,
			os.Getenv("DB_SCHEMA"))
		row := pgDb.QueryRow(sqlString, ipString)
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

func postgresQueryMaxVersion(table string, ipVersion int) int {
	var version int

	sqlString := fmt.Sprintf(`
		SELECT		"db_version"
		
		FROM		"%s"."%s"
		
		WHERE		"ip_version" = $1
		
		ORDER BY	"db_version" DESC
		
		LIMIT 1`,
		os.Getenv("DB_SCHEMA"), table)
	row := pgDb.QueryRow(sqlString, ipVersion)
	if err := row.Scan(&version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0
		}

		panic(err)
	}

	return version
}

func postgresDropOld(table string, ipVersion int, dbVersion int) {
	fmt.Printf("Dropping old Postgres data: `%s`.`%s` ipv%d... ", os.Getenv("DB_SCHEMA"), table, ipVersion)
	sqlString := fmt.Sprintf(`
		DELETE FROM	"%s"."%s" 
		
		WHERE		"ip_version" = $1
		AND			"db_version" < $2`,
		os.Getenv("DB_SCHEMA"), table)
	_, err := pgDb.Exec(sqlString, ipVersion, dbVersion)
	if err != nil {
		panic(err)
	}
	fmt.Println("Complete")
}

func postgresSaveCountries(countries []IpCountry) {
	var params []any

	sqlString := fmt.Sprintf(`
		INSERT INTO "%s"."ip_country" (
			"ip_range_start", 
			"ip_range_end", 
			"country_code", 
			"ip_version", 
			"db_version"
		) VALUES `,
		os.Getenv("DB_SCHEMA"))

	for _, country := range countries {
		sqlString += `($?, $?, $?, $?, $?), `
		params = append(params, country.IpRangeStart, country.IpRangeEnd, country.CountryCode, country.IpVersion, country.DbVersion)
	}
	sqlString = sqlString[0:len(sqlString) - 2]
	sqlString = fixPostgresVars(sqlString)

	stmt, err := pgDb.Prepare(sqlString)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		panic(err)
	}
}

func postgresSaveASNs(ASNs []IpASN) {
	var params []any

	sqlString := fmt.Sprintf(`
		INSERT INTO	"%s"."ip_asn" (
			"ip_range_start", 
			"ip_range_end", 
			"as_number", 
			"as_organisation", 
			"ip_version", 
			"db_version"
		) VALUES `,
		os.Getenv("DB_SCHEMA"))

	for _, asn := range ASNs {
		sqlString += `($?, $?, $?, $?, $?, $?), `
		params = append(params, asn.IpRangeStart, asn.IpRangeEnd, asn.AsNumber, asn.AsOrganisation, asn.IpVersion, asn.DbVersion)
	}
	sqlString = sqlString[0:len(sqlString) - 2]
	sqlString = fixPostgresVars(sqlString)

	stmt, err := pgDb.Prepare(sqlString)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		panic(err)
	}
}

func postgresSaveCities(cities []IpCity) {
	var params []any

	sqlString := fmt.Sprintf(`
		INSERT INTO "%s"."ip_city" (
			"ip_range_start", 
			"ip_range_end", 
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
		os.Getenv("DB_SCHEMA"))
	for _, city := range cities {
		sqlString += `($?, $?, $?, $?, $?, $?, $?, $?, $?, $?, $?, $?), `
		params = append(params, city.IpRangeStart, city.IpRangeEnd, city.CountryCode, city.State1, city.State2, city.City, city.Postcode, city.Latitude, city.Longitude, city.Timezone, city.IpVersion, city.DbVersion)
	}
	sqlString = sqlString[0:len(sqlString) - 2]
	sqlString = fixPostgresVars(sqlString)

	stmt, err := pgDb.Prepare(sqlString)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		panic(err)
	}
}

func fixPostgresVars(sqlString string) string {
	varIncrement := 1;
	for {
		if strings.Contains(sqlString, "$?") {
			sqlString = strings.Replace(sqlString, "$?", "$" + strconv.Itoa(varIncrement), 1)
			varIncrement++
		} else {
			break
		}
	}

	return sqlString
}

func postgresFile(sqlPath string) {
	sqlBytes, err := dbStructures.ReadFile(sqlPath)
	//sqlBytes, err := os.ReadFile(sqlPath)
	if err != nil {
		panic(err)
	}
	sqlString := string(sqlBytes)
	sqlString = strings.Replace(sqlString, "${schema}", os.Getenv("DB_SCHEMA"), -1)

	sqlStatements := strings.Split(sqlString, ";")

	transaction, err := pgDb.Begin()
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