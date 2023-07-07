-- +goose Up
-- +goose StatementBegin
CREATE TABLE  IF NOT EXISTS products(
    barcode varchar(10) NOT NULL,
    name varchar(60) NOT NULL,
    descr Text NOT NULL ,
    cost int NOT NULL , 
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    userId int NOT NULL ,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT u_pkey PRIMARY KEY (barcode),
    CONSTRAINT products_custonmers_fk 
    FOREIGN KEY (userId)  REFERENCES users (id) 
);

CREATE TABLE IF NOT EXISTS prchecks( 
    filename varchar(45) NOT NULL,
    barcode  varchar(10) NOT NULL,
    CONSTRAINT u_pkey PRIMARY KEY (filename),
    CONSTRAINT products_files_fk 
    FOREIGN KEY (barcode)  REFERENCES products (barcode) 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE  IF EXISTS products;
DROP TABLE  IF EXISTS prchecks;
-- +goose StatementEnd
