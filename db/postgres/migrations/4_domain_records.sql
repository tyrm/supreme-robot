-- +migrate Up
CREATE TYPE dns_record_types AS ENUM ('A', 'AAAA', 'CNAME', 'TXT', 'NS', 'MX', 'SRV', 'SOA', 'CAA', 'TLSA');

CREATE TABLE "public"."domain_records" (
    id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    name character varying NOT NULL,
    domain_id uuid NOT NULL,
    type dns_record_types NOT NULL,
    value character varying NOT NULL,
    ttl integer NOT NULL,
    priority integer,
    port integer,
    weight integer,
    refresh integer,
    retry integer,
    expire integer,
    mbox character varying,
    tag character varying,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    PRIMARY KEY ("id"),
    FOREIGN KEY (domain_id) REFERENCES domains (id) ON DELETE CASCADE
)
;

-- +migrate Down
DROP TABLE "public"."domain_records";