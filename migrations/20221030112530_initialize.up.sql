create table if not exists configs
(
    id        bigserial primary key,
    service   varchar(255) NOT NULL,
    version   int default 0,
    used      boolean default true,
    data      json
);
create index if not exists ix_service on configs (service);