CREATE EXTENSION IF NOT EXISTS unaccent;
alter table companies
    add column ts tsvector generated always as (to_tsvector('simple', name || ' ' || symbol || ' ' || isin || ' ' || coalesce(bio, ''))) STORED;
create index if not exists companies_ts on companies using gin (ts);