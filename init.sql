CREATE TABLE IF NOT EXISTS "modules" (
  "id" integer PRIMARY KEY,
  "name" text NOT NULL,
  "form" text NOT NULL,
  "mod" text NOT NULL,
  "code" text
);

CREATE TABLE IF NOT EXISTS "nflow" (
  "id" integer PRIMARY KEY,
  "json" text NOT NULL,
  "name" text NOT NULL,
  "default_js" text
);