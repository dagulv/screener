create table shares(
    "company_id" text not null,
    "date" timestamp not null,
    "open" int not null,
    "high" int not null,
    "low" int not null,
    "close" int not null,
    "volume" bigint not null,
    "average" int not null,
    primary key("company_id", "date")
);