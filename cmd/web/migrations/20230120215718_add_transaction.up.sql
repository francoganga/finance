SET statement_timeout = 0;

--bun:split

CREATE TABLE transactions (
    id serial primary key,
    "date" date not null,
    code varchar not null,
    description varchar not null,
    amount int not null,
    balance int not null
);
