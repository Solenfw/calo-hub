-- +goose Up
ALTER TABLE martin_report_products
    DROP CONSTRAINT martin_report_products_code_fkey;
-- +goose Down
ALTER TABLE martin_report_products
    ADD CONSTRAINT martin_report_products_code_fkey FOREIGN KEY (code) REFERENCES products(code) ON DELETE CASCADE;

