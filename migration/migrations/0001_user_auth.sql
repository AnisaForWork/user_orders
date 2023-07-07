-- +goose Up
-- +goose StatementBegin
CREATE TABLE  IF NOT EXISTS users(
    id int NOT NULL AUTO_INCREMENT,
    login varchar(40) NOT NULL,
    fullName varchar(75) NOT NULL ,
    email varchar(255) NOT NULL , 
    password BLOB NOT NULL ,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT u_pkey PRIMARY KEY (id),
    CONSTRAINT u_email_UNQ UNIQUE (email),
    CONSTRAINT u_login_UNQ UNIQUE (login)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE  IF EXISTS users;
-- +goose StatementEnd
