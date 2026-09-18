-- +goose Up
ALTER TABLE martin_report_products
RENAME COLUMN eng TO description;

-- +goose Down
ALTER TABLE martin_report_products
RENAME COLUMN description TO eng;