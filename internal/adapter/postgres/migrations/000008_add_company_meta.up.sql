alter table companies
add column country_code text not null default '',
add column market_place_code text not null default '';

create index companies_country_code on companies(country_code);