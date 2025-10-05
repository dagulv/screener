create table quarterly_currency_rates(
    fiscal_year int not null,
    quarter smallint not null,
    currency_id text not null references currencies(id)
        on update cascade
        on delete cascade,
    rate real not null,
    primary key (fiscal_year, quarter, currency_id)
);