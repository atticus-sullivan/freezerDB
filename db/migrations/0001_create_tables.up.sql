START TRANSACTION;
CREATE TABLE IF NOT EXISTS categories (
  name VARCHAR(255) NOT NULL,
  CONSTRAINT pk_cat PRIMARY KEY (name)
);
CREATE TABLE IF NOT EXISTS item_types (
  name varchar(255) NOT NULL,
  category_name varchar(255) NOT NULL,
  CONSTRAINT pk_item_types PRIMARY KEY (name),
  CONSTRAINT fk_item_types_cat FOREIGN KEY (category_name) REFERENCES categories(name)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
);
CREATE TABLE IF NOT EXISTS freezer_items (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  identifier varchar(255) NOT NULL,
  amount varchar(255) NOT NULL,
  misc varchar(255) NOT NULL,
  item_name varchar(255) NOT NULL,
  CONSTRAINT pk_freezer_items PRIMARY KEY (id),
  CONSTRAINT fk_freezer_items_name FOREIGN KEY (item_name) REFERENCES item_types(name)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
	);
COMMIT;
