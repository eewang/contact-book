USE contact_book;

CREATE TABLE IF NOT EXISTS persons(
  id       INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
  name     VARCHAR(255),
  notes    MEDIUMTEXT,
  group_id INT
);

CREATE TABLE IF NOT EXISTS groups(
  id               INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
  name             VARCHAR(255)
);

ALTER TABLE persons ADD FOREIGN KEY (group_id) REFERENCES groups(id);
