create table sectors (
    id text primary key,
    name text not null unique
);

create table currencies (
    id text primary key,
    name text not null unique
);

create table companies (
    id text primary key,
    name text not null,
    bio text,
    symbol text not null unique,
    isin text not null,
    "currencyId" text not null references currencies(id)
        on update cascade
        on delete cascade,
    "sectorId" text not null references sectors(id)
        on update cascade
        on delete cascade,
    "orderbookId" text not null
);
create index companies_name on companies(name);
create index companies_isin on companies(isin);