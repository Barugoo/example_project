CREATE TABLE orders (
    id serial primary key,
    email text not null,
    created_at timestamp not null
);