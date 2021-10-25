-- +migrate Up
CREATE TABLE "public"."config" (
    id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    key character varying UNIQUE NOT NULL,
    value character varying,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
)
;

-- +migrate Down
DROP TABLE "public"."config";