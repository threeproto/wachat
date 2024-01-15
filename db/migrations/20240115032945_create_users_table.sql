-- migrate:up
create table if not exists users (
  user_id integer primary key,
  name text not null,
  selected boolean not null,
);

-- migrate:down
drop table if exists users;

