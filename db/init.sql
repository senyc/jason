CREATE DATABASE IF NOT EXISTS jason;

USE jason;

DROP TABLE IF EXISTS tasks;

CREATE TABLE tasks (
  user_id UUID NOT NULL,
  id INT NOT NULL,
  PRIMARY KEY (user_id, id),
  title VARCHAR(150) NOT NULL,
  body VARCHAR(500) DEFAULT NULL,
  due TIMESTAMP DEFAULT NULL,
  time_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  priority TINYINT DEFAULT 0,
  completed BOOL DEFAULT false,
  completed_date DATETIME DEFAULT NULL
);

DROP TABLE IF EXISTS users;

CREATE TABLE users (
  id UUID DEFAULT UUID() PRIMARY KEY,
  first_name VARCHAR(40) DEFAULT NULL,
  last_name VARCHAR(40) DEFAULT NULL,
  password VARCHAR(100) NOT NULL,
  email VARCHAR(64) NOT NULL,
  time_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  account_type VARCHAR(10) DEFAULT NULL,
  api_key VARCHAR(64) DEFAULT NULL,
  added_tasks INT DEFAULT 0 NOT NULL,
  deleted_tasks INT DEFAULT 0 NOT NULL
);

DELIMITER //

CREATE TRIGGER update_added_tasks 
AFTER INSERT ON tasks
FOR EACH ROW
UPDATE users
SET added_tasks = added_tasks + 1
WHERE id = NEW.user_id;

//

CREATE TRIGGER update_deleted_tasks 
AFTER DELETE ON tasks
FOR EACH ROW
UPDATE users
SET deleted_tasks = deleted_tasks + 1
WHERE id = OLD.user_id;

//

DELIMITER ;

ALTER TABLE tasks
  ADD CONSTRAINT fk_tasks_user
  FOREIGN KEY (user_id)
  REFERENCES users(id);

ALTER TABLE users
  ADD CONSTRAINT uc_email 
  UNIQUE (email);
