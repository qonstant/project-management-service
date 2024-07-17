CREATE TABLE "projects" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "description" text NOT NULL,
  "start_date" timestamp NOT NULL,
  "end_date" timestamp NOT NULL,
  "manager_id" BIGINT NOT NULL
);

CREATE TABLE "tasks" (
  "id" BIGSERIAL PRIMARY KEY,
  "title" varchar(255) NOT NULL,
  "description" text NOT NULL,
  "priority" task_priority NOT NULL,
  "status" task_status NOT NULL,
  "assignee_id" BIGINT NOT NULL,
  "project_id" BIGINT NOT NULL,
  "creation_date" timestamp NOT NULL DEFAULT (now()),
  "completion_date" timestamp
);

CREATE INDEX ON "projects" ("manager_id");
CREATE INDEX ON "tasks" ("assignee_id");
CREATE INDEX ON "tasks" ("project_id");

ALTER TABLE "projects" ADD FOREIGN KEY ("manager_id") REFERENCES "users" ("id");
ALTER TABLE "tasks" ADD FOREIGN KEY ("assignee_id") REFERENCES "users" ("id");
ALTER TABLE "tasks" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id");