-- +migrate Up
CREATE TABLE "public"."domains" (
    id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    domain character varying NOT NULL,
    owner_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    PRIMARY KEY ("id"),
    UNIQUE ("domain", "deleted_at"),
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
)
;

-- +migrate Down
DROP TABLE "public"."domains";