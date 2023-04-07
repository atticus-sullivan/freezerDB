START TRANSACTION;
ALTER TABLE categories
ADD CONSTRAINT chk_cat_no_empty CHECK (name NOT IN (
		""
));
COMMIT;
