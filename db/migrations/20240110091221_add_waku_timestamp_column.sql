-- migrate:up
alter table messages add column waku_timestamp integer not null default 0;

-- migrate:down

