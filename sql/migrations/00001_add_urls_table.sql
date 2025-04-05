-- +goose Up
-- +goose StatementBegin
create table urls (
    id bigserial primary key,
    url varchar(255) not null,
    short varchar(4) not null unique,
    created_at timestamp not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table urls;
-- +goose StatementEnd
