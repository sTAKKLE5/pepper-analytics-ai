CREATE SEQUENCE IF NOT EXISTS journal_entries_id_seq;

-- Table Definition
CREATE TABLE "public"."journal_entries" (
    "id" int4 NOT NULL DEFAULT nextval('journal_entries_id_seq'::regclass),
    "plant_id" int4 NOT NULL,
    "title" varchar(200) NOT NULL,
    "description" text,
    "entry_type" varchar(50),
    "image_path" varchar(255),
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp,
    "entry_date" date NOT NULL,
    PRIMARY KEY ("id")
);

-- This script only contains the table creation statements and does not fully represent the table in the database. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS plants_id_seq;

-- Table Definition
CREATE TABLE "public"."plants" (
    "id" int4 NOT NULL DEFAULT nextval('plants_id_seq'::regclass),
    "name" varchar(100) NOT NULL,
    "species" varchar(100),
    "health" varchar(20) CHECK ((health)::text = ANY ((ARRAY['Excellent'::character varying, 'Good'::character varying, 'Fair'::character varying, 'Poor'::character varying])::text[])),
    "growth_stage" varchar(20) CHECK ((growth_stage)::text = ANY (ARRAY[('Seed'::character varying)::text, ('Seedling'::character varying)::text, ('Vegetative'::character varying)::text, ('Flowering'::character varying)::text, ('Fruiting'::character varying)::text])),
    "planting_date" date NOT NULL,
    "image_path" varchar(255),
    "notes" text,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp,
    "last_fertilized_at" timestamp,
    "last_watered_at" timestamp,
    "is_cross" bool DEFAULT false,
    "generation" varchar(50),
    "is_harvested" bool NOT NULL DEFAULT false,
    "harvested_at" timestamptz,
    PRIMARY KEY ("id")
);

ALTER TABLE "public"."journal_entries" ADD FOREIGN KEY ("plant_id") REFERENCES "public"."plants"("id") ON DELETE CASCADE;


-- Indices
CREATE INDEX idx_journal_plant_id ON public.journal_entries USING btree (plant_id);
CREATE INDEX idx_journal_created_at ON public.journal_entries USING btree (created_at);


-- Indices
CREATE INDEX idx_plants_planting_date ON public.plants USING btree (planting_date);
CREATE INDEX idx_plants_health ON public.plants USING btree (health);
CREATE INDEX idx_plants_growth_stage ON public.plants USING btree (growth_stage);
