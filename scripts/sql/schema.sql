-- storing image meta
CREATE TABLE images (
    id              serial,
    title           varchar,
    orig_filename   varchar,
    filepath        varchar,
    height          integer,
    width           integer,
    hash            varchar,
    created_at      timestamp,
    PRIMARY KEY(id)
);