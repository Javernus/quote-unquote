CREATE TABLE quote (
    id uuid primary key,
    message varchar(256) not null,
    person varchar(256) not null,
    ip inet not null,
    created_at timestamptz not null,
    updated_at timestamptz not null
);

CREATE INDEX ON quote (created_at);
