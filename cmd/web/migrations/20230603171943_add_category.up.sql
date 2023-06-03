SET statement_timeout = 0;

--bun:split

CREATE TABLE category (
    id serial primary key,
    name varchar not null
)

--bun:split

ALTER TABLE transactions
ADD COLUMN category_id int references category(id);
