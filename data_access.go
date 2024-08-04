package main

import (
	"embed"
	"net"
	"os"
)

//go:embed structure/*.sql
var dbStructures embed.FS

func dbConnect() {
	switch os.Getenv("DB_TYPE") {
		case "postgres": 	postgresConnect()
		case "mysql": 		mysqlConnect()
		case "sqlite": 		sqliteConnect()
		case "mmdb": 		mmdbConnect()
	}
}

func dbClose() {
	switch os.Getenv("DB_TYPE") {
		case "postgres": 	postgresClose()
		case "mysql": 		mysqlClose()
		case "sqlite": 		sqliteClose()
		case "mmdb": 		mmdbClose()
	}
}

func dbInitialised(key string) bool {
	switch os.Getenv("DB_TYPE") {
		case "postgres": 	return postgresInitialised(key)
		case "mysql": 		return mysqlInitialised(key)
		case "sqlite": 		return sqliteInitialised(key)
		case "mmdb": 		return mmdbInitialised(key)
	}

	return false
}

func dbIp(ip net.IP) *Ip {
	switch os.Getenv("DB_TYPE") {
		case "postgres":	return postgresIp(ip)
		case "mysql": 		return mysqlIp(ip)
		case "sqlite": 		return sqliteIp(ip)
		case "mmdb":		return mmdbIp(ip)
	}

	return NewIp("0.0.0.0", 4)
}

func dbDropOld(table string, ipVersion int, dbVersion int) {
	switch os.Getenv("DB_TYPE") {
		case "postgres":	postgresDropOld(table, ipVersion, dbVersion)
		case "mysql": 		mysqlDropOld(table, ipVersion, dbVersion)
		case "sqlite": 		sqliteDropOld(table, ipVersion, dbVersion)
		case "mmdb":		mmdbSaveRestart(table, ipVersion)
	}
}

func dbSaveCountries(countries []IpCountry) {
	switch os.Getenv("DB_TYPE") {
		case "postgres":	postgresSaveCountries(countries)
		case "mysql": 		mysqlSaveCountries(countries)
		case "sqlite": 		sqliteSaveCountries(countries)
		case "mmdb":		mmdbSaveCountries(countries)
	}
}

func dbSaveASNs(ASNs []IpASN) {
	switch os.Getenv("DB_TYPE") {
		case "postgres":	postgresSaveASNs(ASNs)
		case "mysql": 		mysqlSaveASNs(ASNs)
		case "sqlite": 		sqliteSaveASNs(ASNs)
		case "mmdb":		mmdbSaveASNs(ASNs)
	}
}

func dbSaveCities(cities []IpCity) {
	switch os.Getenv("DB_TYPE") {
		case "postgres":	postgresSaveCities(cities)
		case "mysql": 		mysqlSaveCities(cities)
		case "sqlite": 		sqliteSaveCities(cities)
		case "mmdb":		mmdbSaveCities(cities)
	}
}

func dbFile() {
	switch os.Getenv("DB_TYPE") {
		case "postgres":	postgresFile("structure/postgres.sql")
		case "mysql": 		mysqlFile("structure/mysql.sql")
		case "sqlite": 		sqliteFile("structure/sqlite.sql")
	}
}

func dbQueryMaxVersion(table string, ipVersion int) int {
	switch os.Getenv("DB_TYPE") {
		case "postgres":	return postgresQueryMaxVersion(table, ipVersion)
		case "mysql": 		return mysqlQueryMaxVersion(table, ipVersion)
		case "sqlite": 		return sqliteQueryMaxVersion(table, ipVersion)
	}

	return 0
}