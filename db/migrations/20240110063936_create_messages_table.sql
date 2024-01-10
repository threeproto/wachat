-- migrate:up
create table if not exists messages (
  message_id integer primary key,
  user_name text not null,
  content text not null,
  timestamp integer not null,
  message_hash text not null,
  is_stored boolean not null default false
);

-- migrate:down
drop table if exists messages;
