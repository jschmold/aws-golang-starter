CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA common;
CREATE SCHEMA accounts; 


---
-- Base type for all "standalone" entities that need timestamps
---
CREATE TABLE common.timestamps (
	created_at TIMESTAMPTZ NOT NULL DEFAULT Now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT Now(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE common.verifications (
  created_at TIMESTAMPTZ  NOT NULL DEFAULT Now(),
  code VARCHAR(32) NOT NULL DEFAULT md5(uuid_generate_v4()::text),
  verified BOOLEAN NOT NULL DEFAULT false
);

---
-- Section: Accounts
--   Logins, organizations, roles, etc
---

CREATE TABLE accounts.users (
  id UUID PRIMARY KEY
    DEFAULT uuid_generate_v4()
    NOT NULL,

  name VARCHAR(256) DEFAULT '',

  email VARCHAR UNIQUE NOT NULL,

  email_verified BOOLEAN
    NOT NULL
    DEFAULT false,

  password bytea

) INHERITS (common.timestamps);


CREATE TABLE accounts.password_resets (
  id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),

  user_id UUID NOT NULL
    REFERENCES accounts.users (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,

  email citext NOT NULL

) INHERITS (common.verifications);

CREATE TABLE accounts.user_confirmations (
  id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),

	created_at TIMESTAMPTZ DEFAULT Now() NOT NULL,

  user_id UUID NOT NULL
    REFERENCES accounts.users (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE

) INHERITS (common.verifications);