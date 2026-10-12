-- +goose Up
ALTER TABLE martin_report_products RENAME to report_products;
ALTER TABLE martin_report_list RENAME to report_list;

-- +goose Down
ALTER TABLE report_products RENAME to martin_report_products;
ALTER TABLE report_list RENAME to martin_report_list;
