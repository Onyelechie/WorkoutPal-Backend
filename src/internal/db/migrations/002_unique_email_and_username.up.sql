-- Extension for case-sensitive text types
CREATE EXTENSION IF NOT EXISTS citext;

-- Make username and email columns case-sensitive
ALTER TABLE users
ALTER COLUMN email TYPE citext,
ALTER COLUMN username TYPE citext;

-- Make username and email columns unique
ALTER TABLE users ADD CONSTRAINT users_username_key UNIQUE (username);
ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);