CREATE TABLE "Products"(
	"id" BIGSERIAL NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"price" NUMERIC(12,2) NOT NULL,
	CONSTRAINT "unique_name" UNIQUE ("name"),
	CONSTRAINT "primary_key" PRIMARY KEY ("id")
)