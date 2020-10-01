create extension if not exists "pgcrypto";

create table if not exists stories (
    id uuid primary key default gen_random_uuid(),
    title  varchar(100) not null,
    body  varchar(100000) not null,
    viewCount bigint default 0,
    upVotes bigint default 0,
    downVotes bigint default 0,
    createdAt timestamp without time zone default (now() at time zone 'utc'),
    updatedAt timestamp without time zone default (now() at time zone 'utc'),
    CHECK (title <> ''),
    CHECK (body <> '')
);