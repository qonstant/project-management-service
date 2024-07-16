CREATE TYPE "task_priority" AS ENUM (
  'low',
  'medium',
  'high'
);

CREATE TYPE "task_status" AS ENUM (
  'new',
  'in_progress',
  'completed'
);

CREATE TABLE "users" (
  "id" BIGSERIAL PRIMARY KEY,
  "full_name" varchar(255) NOT NULL,
  "email" varchar(320) UNIQUE NOT NULL,
  "registration_date" timestamp NOT NULL DEFAULT (now()),
  "role" varchar(50) NOT NULL
);
