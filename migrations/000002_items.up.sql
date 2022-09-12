CREATE TABLE items (
    id serial primary key,
    name text not null,
    color text not null,
    created_at timestamp not null
);

INSERT INTO items (name, color, created_at)
VALUES 
    ('go_course_basic', 'blue', NOW()),
    ('go_course_regular', 'blue', NOW()),
    ('go_course_advanced', 'blue', NOW());