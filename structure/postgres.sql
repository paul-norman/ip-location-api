CREATE TABLE IF NOT EXISTS "${schema}"."ip_city"
(
	"ip_range_start"	inet not null,
	"ip_range_end"		inet not null,
	"country_code"		varchar,
	"state1"			varchar,
	"state2"			varchar,
	"city"				varchar,
	"postcode"			varchar,
	"latitude"			numeric,
	"longitude"			numeric,
	"timezone"			varchar,
	"ip_version"		int DEFAULT 4,
	"db_version"		int DEFAULT 1
);

CREATE INDEX IF NOT EXISTS "I:${schema}:ip_city:ip_range_start" ON "${schema}"."ip_city" ("ip_range_start");
CREATE INDEX IF NOT EXISTS "I:${schema}:ip_city:ip_range_end" ON "${schema}"."ip_city" ("ip_range_end");
CREATE INDEX IF NOT EXISTS "I:${schema}:ip_city:country_code" ON "${schema}"."ip_city" ("country_code");
CREATE INDEX IF NOT EXISTS "I:${schema}:ip_city:ip_version" ON "${schema}"."ip_city" ("ip_version");
CREATE INDEX IF NOT EXISTS "I:${schema}:ip_city:db_version" ON "${schema}"."ip_city" ("db_version");

CREATE TABLE IF NOT EXISTS "${schema}"."ip_asn"
(
	"ip_range_start"	inet not null,
	"ip_range_end"		inet not null,
	"as_number"			varchar,
	"as_organisation"	varchar,
	"ip_version"		int DEFAULT 4,
	"db_version"		int DEFAULT 1
);

CREATE INDEX IF NOT EXISTS "I:${schema}:ip_asn:ip_range_start" ON "${schema}"."ip_asn" ("ip_range_start");
CREATE INDEX IF NOT EXISTS "I:${schema}:ip_asn:ip_range_end" ON "${schema}"."ip_asn" ("ip_range_end");
CREATE INDEX IF NOT EXISTS "I:${schema}:ip_asn:ip_version" ON "${schema}"."ip_asn" ("ip_version");
CREATE INDEX IF NOT EXISTS "I:${schema}:ip_asn:db_version" ON "${schema}"."ip_asn" ("db_version");

CREATE TABLE IF NOT EXISTS "${schema}"."ip_country"
(
	"ip_range_start"	inet not null,
	"ip_range_end"		inet not null,
	"country_code"		varchar,
	"ip_version"		int DEFAULT 4,
	"db_version"		int DEFAULT 1
);

CREATE INDEX IF NOT EXISTS "I:${schema}:ip_country:ip_range_start" ON "${schema}"."ip_country" ("ip_range_start");
CREATE INDEX IF NOT EXISTS "I:${schema}:ip_country:ip_range_end" ON "${schema}"."ip_country" ("ip_range_end");
CREATE INDEX IF NOT EXISTS "I:${schema}:ip_country:country_code" ON "${schema}"."ip_country" ("country_code");
CREATE INDEX IF NOT EXISTS "I:${schema}:ip_country:ip_version" ON "${schema}"."ip_country" ("ip_version");
CREATE INDEX IF NOT EXISTS "I:${schema}:ip_country:db_version" ON "${schema}"."ip_country" ("db_version");