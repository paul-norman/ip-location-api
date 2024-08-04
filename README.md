# IP Location API

This is simple system to automatically load any of the data from the excellent [ip-location-db](https://github.com/sapics/ip-location-db) project into a chosen database and expose it as a basic API for IP location lookups. This is really only meant for personal / single project use, it doesn't support items required for multiple user access *(i.e. multiple API keys)*, and **responsible compliance of licences is left up to the user**.

It's written in [Go](https://go.dev/) to allow it to compile to many platforms and run from a single binary.

Database support has been added for [MMDB](https://support.maxmind.com/hc/en-us/articles/4408216157723-Database-Formats), [PostgreSQL](https://www.postgresql.org/), [MySQL](https://www.mysql.com/) / [MariaDB](https://mariadb.org/) and [SQLite](https://www.sqlite.org/index.html) *(ipv4 only - see below)*. MMDB is the most optimised for the task, but since many developers use the other databases already, it might be convenient to have the data available directly within those *(without necessarily needing the API at all)*.

## API Usage

The only important route is the IP lookup: `/ip/{ip}`, e.g. `/ip/40.45.124.54`:

```json
{
	"ip": "42.45.124.54",
	"ip_version": 4,
	"found_country": true,
	"found_city": true,
	"found_asn": true,
	"country_code": "KR",
	"state": "Seoul",
	"state_2": "",
	"city": "Seoul (Eulji-ro)",
	"postcode": "",
	"lat": 37.566,
	"lon": 126.993,
	"timezone": "",
	"as_number": 9644,
	"as_organisation": "SK Telecom",
	"ms_taken": 0,
	"μs_taken": 224
}
```

This route accepts IPv4 and IPv6 strings in any format that [Go](https://pkg.go.dev/net#ParseIP) will support.

There are two more routes, but these **only run with an API key defined**:

- `/random/{ipVersion}`, e.g. `/random/6`
  - return the above result for a random IP
- `/benchmark/{ipVersion}/{times}`, e.g. `/benchmark/4/500`
  - run `{times}` number of lookups of randomly generated IP addresses

## Installation

Download the most suitable build for your system from the `releases`. This will be an executable file. Place this file in a directory, e.g. `/var/www/ip-location-api`, and ensure that it's executable:

```Shell
sudo mkdir -p /var/www/ip-location-api
sudo wget -q -O /var/www/ip-location-api/ip-location-api https://github.com/paul-norman/ip-location-api/releases/download/1.0.0/ip-location-api-linux-x64.bin 
sudo chmod +x /var/www/ip-location-api/ip-location-api
``` 

Either make this folder writable for the user that will run the system, or create a downloads folder and make that writable:

```Shell
sudo mkdir /var/www/ip-location-api/downloads
sudo chmod 0777 /var/www/ip-location-api/downloads
```

The system will download updates into this folder and process them.

Create a `.env` file in the main directory containing your required settings *(see next section)*:

```Shell
sudo nano /var/www/ip-location-api/.env
```

Start the system and wait for it to update:

```Shell
cd /var/www/ip-location-api
./ip-location-api
```

## Configuration

Configuration is handled by an environmental (`.env`) file. Example files for each database type exist in the `.env.sample` directory. All database types share some information:

```Dotenv
SERVER_HOST=127.0.0.1
SERVER_PORT=8081

API_KEY=

COUNTRY=dbip-country
CITY=dbip-city
ASN=dbip-asn

UPDATE_TIME=01:30
```

If you wish to expose the system without a reverse proxy, you may wish to update `SERVER_HOST` to `0.0.0.0`.

`API_KEY` allows a very basic protection of the system to be applied, a header named `API_KEY` with a matching value must be passed if this variable is populated. If left blank, the API is open.

`COUNTRY`, `CITY` and `ASN` are the databases that will be loaded. **If you don't need cities or ASNs, just leave them blank.** The values / names used should mirror the directory values found in the [ip-location-db](https://github.com/sapics/ip-location-db) project:

### Allowed `COUNTRY` values

- asn-country
- dbip-country
- geo-asn-country
- geo-whois-asn-country
- geolite2-country
- iptoasn-country
- webnet77-country

### Allowed `ASN` values

- asn
- dbip-asn
- geolite2-asn
- iptoasn-asn

### Allowed `CITY` values

- dbip-city
- geolite2-city

`UPDATE_TIME` is optional, but if present *(and in standard HH:MM format)*, it will check for / download / reload new data every 24 hours at the time specified.

### MMDB

The MMDB adaption doesn't need any initialisation, it just needs to be told to use that format:

```Dotenv
DB_TYPE=mmdb
```

### PostgreSQL

The PostgreSQL adaption will create 3 tables *(`ip_country`, `ip_asn` and `ip_city`)* in a database of your choice:

```Dotenv
DB_TYPE=postgres
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=db_username
DB_PASS=db_password
DB_NAME=db_name
DB_SCHEMA=ip
```

### MySQL

The MySQL adaption will create 3 tables *(`ip_country`, `ip_asn` and `ip_city`)* in a database of your choice:

```Dotenv
DB_TYPE=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=db_username
DB_PASS=db_password
DB_NAME=ip
```

### SQLite

**SQLite currently DOES NOT work with IPv6 because it lacks a datatype which can easily compare numbers of the required size.** I don't personally ever use SQLite, so when I realised that it had no support for larger numeric fields I didn't push any further. If anyone would like to fix this, it would be gratefully received! 

The SQLite adaption will create 6 tables *(ipv4 and ipv6 are separated)* *(`ipv4_country`, `ipv6_country_ipv6`, `ipv4_asn`, `ipv6_asn`, `ipv4_city` and `ipv6_city`)* in a database of your choice:

```Dotenv
DB_TYPE=sqlite
DB_USER=path/to/db.db
DB_SCHEMA=
```

Since schemas are supported by SQLite, their use is optional. `DB_TYPE` may also be set to `:memory:`.

## Install as a service

This is different for every system, but I'm going to assume that Linux is a popular choice and cover Systemd here.

```Shell
sudo nano /etc/systemd/system/ip-location-api.service
```

```
[Unit]
Description=Vikunja
After=syslog.target
After=network.target
# Depending on how you configured the system, you may want to uncomment these:
#Requires=postgresql.service
#Requires=mysql.service
#Requires=mariadb.service
#Requires=sqlite.service

[Service]
RestartSec=2s
Type=simple
WorkingDirectory=/var/www/ip-location-api
ExecStart=/var/www/ip-location-api/ip-location-api
Restart=always

[Install]
WantedBy=multi-user.target
```

```Shell
sudo systemctl enable ip-location-api
```

## Updates

The system will update whenever it restarts *(if the data is missing)*. It will attempt to do this in a way that minimises downtime, but may still have a second or so of outage for MMDB files *(when they reload)*. If the `UPDATE_TIME` has been specified, the system will keep itself up-to-date every 24 hours.

## Reverse Proxy

**This is optional**, you only need this if you are running this service on a different machine from your codebase *(and you actually want the API functionality, not just the database)*.

Any reverse proxy can handle this task, but for [Nginx](https://nginx.org/en/) the config might look like:

```Nginx
server {
	listen		443 ssl http2;
	listen		[::]:443 ssl http2;
	server_name	location-api.yoursite.com;

	include		ssl.d/yoursite.com.conf;
	
	location / {	
		proxy_http_version	1.1;
		proxy_cache_bypass	$http_upgrade;
		proxy_set_header	Upgrade $http_upgrade;
		proxy_set_header	Connection 'upgrade';
		proxy_set_header	Host $host;
		proxy_set_header	X-Real-IP $remote_addr;
		proxy_set_header	X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header	X-Forwarded-Proto $scheme;
		proxy_pass		http://127.0.0.1:8081/;
	}
}
```

## Benchmarks

These are very unscientific benchmarks running on my home PC. They use the built-in benchmarking route to test with. They are only meant to provide an idea of relative performance.

Over 10,000 loops, `/benchmark/{ipVersion}/10000`, the average lookup operation time:

| DB        | IPv4 (μs) | IPv6 (μs) |
|-----------|-----------|-----------|
| MMDB      | 2-4       | 2-4       |
| Postgres  | 220-280   | 220-280   |
| MySQL     | 220-280   | 220-280   |
| SQLite    | 30-60     | N/A       |

More interestingly, for a single IP, `/benchmark/{ipVersion}/1`:

| DB        | IPv4 (μs) | IPv6 (μs) |
|-----------|-----------|-----------|
| MMDB      | 200-1000  | 200-1000  |
| Postgres  | 600-1000  | 600-1000  |
| MySQL     | 600-1000  | 600-1000  |
| SQLite    | 200-1000  | N/A       |

**1000μs is still only 0.001 seconds, so all are acceptably quick.**

## Possible Future Improvements / Enhancements

- [ ] Make the webserver optional
- [ ] Add ready to use Docker examples
- [ ] Return licence info with the API results *(if required)*
- [ ] Improve my sloppy Go code
- [ ] Add proper tests
- [ ] Load in proper testing / benchmarking data *(currently just randomly generated ips)*