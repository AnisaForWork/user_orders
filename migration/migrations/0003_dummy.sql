-- +goose Up
-- +goose StatementBegin
 

CREATE TABLE  IF NOT EXISTS prchecks( 
    filename varchar(45) NOT NULL,
    barcode  varchar(10) NOT NULL,
    CONSTRAINT u_pkey PRIMARY KEY (filename),
    CONSTRAINT products_files_fk 
    FOREIGN KEY (barcode)  REFERENCES products (barcode) 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE  IF EXISTS prchecks;
-- +goose StatementEnd
