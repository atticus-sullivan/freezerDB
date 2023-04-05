START TRANSACTION;
CREATE TABLE IF NOT EXISTS categories (
  name VARCHAR(255) NOT NULL,
  PRIMARY KEY (name)
);
CREATE TABLE IF NOT EXISTS item_types (
  name varchar(255) NOT NULL,
  category_name varchar(255) NOT NULL,
  PRIMARY KEY (name),
  FOREIGN KEY (category_name) REFERENCES categories(name)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
);
CREATE TABLE IF NOT EXISTS freezer_items (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  identifier varchar(255),
  amount varchar(255),
  misc varchar(255),
  item_name varchar(255),
  PRIMARY KEY (id),
  FOREIGN KEY (item_name) REFERENCES item_types(name)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
	);
COMMIT;
