CREATE TABLE IF NOT EXISTS "beers"(
    id bigserial NOT NULL,
    name character varying(150) NOT NULL,
    brewery character varying(150) NOT NULL,
    country character varying(100) NOT NULL,
    price numeric(15, 9) NOT NULL,
    currency character varying(3) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE "beers" OWNER to postgres;