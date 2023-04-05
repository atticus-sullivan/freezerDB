-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back

START TRANSACTION;
DROP TABLE IF EXISTS freezer_item;
DROP TABLE IF EXISTS item_type;
DROP TABLE IF EXISTS category;
COMMIT;
