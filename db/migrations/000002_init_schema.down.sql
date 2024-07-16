-- Drop foreign keys from tasks table
ALTER TABLE "tasks" DROP CONSTRAINT IF EXISTS "tasks_project_id_fkey";
ALTER TABLE "tasks" DROP CONSTRAINT IF EXISTS "tasks_assignee_id_fkey";
ALTER TABLE "tasks" DROP CONSTRAINT IF EXISTS "tasks_priority_fkey";

-- Drop indices from tasks table
DROP INDEX IF EXISTS "tasks_project_id_idx";
DROP INDEX IF EXISTS "tasks_assignee_id_idx";

-- Drop tasks table
DROP TABLE IF EXISTS "tasks";

-- Drop foreign key from projects table
ALTER TABLE "projects" DROP CONSTRAINT IF EXISTS "projects_manager_id_fkey";

-- Drop index from projects table
DROP INDEX IF EXISTS "projects_manager_id_idx";

-- Drop projects table
DROP TABLE IF EXISTS "projects";
