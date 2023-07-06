-- +goose Up
-- +goose StatementBegin
CREATE TABLE  IF NOT EXISTS products(
    barcode varchar(10) NOT NULL,
    name varchar(255) NOT NULL,
    desc Text NOT NULL ,
    cost int NOT NULL , 
    userId int NOT NULL ,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT u_pkey PRIMARY KEY (barcode),
    CONSTRAINT products_custonmers_fk 
    FOREIGN KEY (userId)  REFERENCES user (id) ON DELETE,
    CONSTRAINT u_email_UNQ UNIQUE (email),
    CONSTRAINT u_login_UNQ UNIQUE (login)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE  IF EXISTS users;
-- +goose StatementEnd
