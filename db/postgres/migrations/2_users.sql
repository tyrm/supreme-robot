-- +migrate Up
CREATE TABLE "public"."users" (
    id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    username character varying NOT NULL,
    password character varying NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    PRIMARY KEY ("id"),
    UNIQUE ("username", "deleted_at")
)
;

CREATE TABLE "public"."group_membership" (
    id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL,
    group_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    PRIMARY KEY ("id"),
    UNIQUE ("user_id", "group_id", "deleted_at"),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
)
;

INSERT INTO "public"."users"("id", "username", "password")
    VALUES('8c504483-1e11-4243-b6c8-14499877a641', 'admin', '$2a$14$mmOFu7eOyQUFC0S/gopbDeJKcADiUx7QleU85WW7FnnCiXNgENb1G')
;

INSERT INTO "public"."group_membership"("user_id", "group_id")
    VALUES('8c504483-1e11-4243-b6c8-14499877a641', '71df8f2b-f293-4fde-93b1-e40dbe5c97ea')
;

-- +migrate Down
DROP TABLE "public"."group_membership";
DROP TABLE "public"."users";