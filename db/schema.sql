CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE messages (
  message_id integer primary key,
  user_name text not null,
  content text not null,
  timestamp integer not null,
  message_hash text not null,
  is_stored boolean not null default false
, waku_timestamp integer not null default 0);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20240110063936'),
  ('20240110091221');
