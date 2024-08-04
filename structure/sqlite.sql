CREATE TABLE IF NOT EXISTS "${schema}"."ipv4_city" (
	"ip_range_start"	TEXT NOT NULL,
	"ip_range_end"		TEXT NOT NULL,
	"ip_number_start"	INTEGER NOT NULL,
	"ip_number_end"		INTEGER NOT NULL,
	"country_code"		TEXT NOT NULL,
	"state1"			TEXT NOT NULL,
	"state2"			TEXT NOT NULL,
	"city"				TEXT NOT NULL,
	"postcode"			TEXT NOT NULL,
	"latitude"			REAL NOT NULL,
	"longitude"			REAL NOT NULL,
	"timezone"			TEXT NOT NULL,
	"ip_version"		INTEGER NOT NULL DEFAULT 4,
	"db_version"		INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS "I:${schema}:ipv4_city:ip_number_start" ON "${schema}"."ipv4_city" ("ip_number_start");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv4_city:ip_number_end" ON "${schema}"."ipv4_city" ("ip_number_end");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv4_city:country_code" ON "${schema}"."ipv4_city" ("country_code");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv4_city:ip_version" ON "${schema}"."ipv4_city" ("ip_version");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv4_city:db_version" ON "${schema}"."ipv4_city" ("db_version");

CREATE TABLE IF NOT EXISTS "${schema}"."ipv6_city" (
	"ip_range_start"	TEXT NOT NULL,
	"ip_range_end"		TEXT NOT NULL,
	"ip_number_start"	NUMERIC NOT NULL,
	"ip_number_end"		NUMERIC NOT NULL,
	"country_code"		TEXT NOT NULL,
	"state1"			TEXT NOT NULL,
	"state2"			TEXT NOT NULL,
	"city"				TEXT NOT NULL,
	"postcode"			TEXT NOT NULL,
	"latitude"			REAL NOT NULL,
	"longitude"			REAL NOT NULL,
	"timezone"			TEXT NOT NULL,
	"ip_version"		INTEGER NOT NULL DEFAULT 4,
	"db_version"		INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS "I:${schema}:ipv6_city:ip_number_start" ON "${schema}"."ipv6_city" ("ip_number_start");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv6_city:ip_number_end" ON "${schema}"."ipv6_city" ("ip_number_end");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv6_city:country_code" ON "${schema}"."ipv6_city" ("country_code");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv6_city:ip_version" ON "${schema}"."ipv6_city" ("ip_version");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv6_city:db_version" ON "${schema}"."ipv6_city" ("db_version");

CREATE TABLE IF NOT EXISTS "${schema}"."ipv4_asn" (
	"ip_range_start"	TEXT NOT NULL,
	"ip_range_end"		TEXT NOT NULL,
	"ip_number_start"	INTEGER NOT NULL,
	"ip_number_end"		INTEGER NOT NULL,
	"as_number"			INTEGER NOT NULL DEFAULT 0,
	"as_organisation"	TEXT NOT NULL,
	"ip_version"		INTEGER NOT NULL DEFAULT 4,
	"db_version"		INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS "I:${schema}:ipv4_asn:ip_number_start" ON "${schema}"."ipv4_asn" ("ip_number_start");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv4_asn:ip_number_end" ON "${schema}"."ipv4_asn" ("ip_number_end");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv4_asn:ip_version" ON "${schema}"."ipv4_asn" ("ip_version");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv4_asn:db_version" ON "${schema}"."ipv4_asn" ("db_version");

CREATE TABLE IF NOT EXISTS "${schema}"."ipv6_asn" (
	"ip_range_start"	TEXT NOT NULL,
	"ip_range_end"		TEXT NOT NULL,
	"ip_number_start"	NUMERIC NOT NULL,
	"ip_number_end"		NUMERIC NOT NULL,
	"as_number"			INTEGER NOT NULL DEFAULT 0,
	"as_organisation"	TEXT NOT NULL,
	"ip_version"		INTEGER NOT NULL DEFAULT 4,
	"db_version"		INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS "I:${schema}:ipv6_asn:ip_number_start" ON "${schema}"."ipv6_asn" ("ip_number_start");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv6_asn:ip_number_end" ON "${schema}"."ipv6_asn" ("ip_number_end");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv6_asn:ip_version" ON "${schema}"."ipv6_asn" ("ip_version");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv6_asn:db_version" ON "${schema}"."ipv6_asn" ("db_version");

CREATE TABLE IF NOT EXISTS "${schema}"."ipv4_country" (
	"ip_range_start"	TEXT NOT NULL,
	"ip_range_end"		TEXT NOT NULL,
	"ip_number_start"	INTEGER NOT NULL,
	"ip_number_end"		INTEGER NOT NULL,
	"country_code"		TEXT NOT NULL,
	"ip_version"		INTEGER NOT NULL DEFAULT 4,
	"db_version"		INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS "I:${schema}:ipv4_country:ip_number_start" ON "${schema}"."ipv4_country" ("ip_number_start");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv4_country:ip_number_end" ON "${schema}"."ipv4_country" ("ip_number_end");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv4_country:country_code" ON "${schema}"."ipv4_country" ("country_code");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv4_country:ip_version" ON "${schema}"."ipv4_country" ("ip_version");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv4_country:db_version" ON "${schema}"."ipv4_country" ("db_version");

CREATE TABLE IF NOT EXISTS "${schema}"."ipv6_country" (
	"ip_range_start"	TEXT NOT NULL,
	"ip_range_end"		TEXT NOT NULL,
	"ip_number_start"	NUMERIC NOT NULL,
	"ip_number_end"		NUMERIC NOT NULL,
	"country_code"		TEXT NOT NULL,
	"ip_version"		INTEGER NOT NULL DEFAULT 4,
	"db_version"		INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS "I:${schema}:ipv6_country:ip_number_start" ON "${schema}"."ipv6_country" ("ip_number_start");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv6_country:ip_number_end" ON "${schema}"."ipv6_country" ("ip_number_end");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv6_country:country_code" ON "${schema}"."ipv6_country" ("country_code");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv6_country:ip_version" ON "${schema}"."ipv6_country" ("ip_version");
CREATE INDEX IF NOT EXISTS "I:${schema}:ipv6_country:db_version" ON "${schema}"."ipv6_country" ("db_version");