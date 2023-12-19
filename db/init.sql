CREATE DATABASE IF NOT EXISTS jason;

USE jason;

DROP TABLE IF EXISTS tasks;

CREATE TABLE tasks (
  user_id INT NOT NULL,
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(150) NOT NULL,
  body VARCHAR(500) DEFAULT NULL,
  due DATETIME DEFAULT NULL,
  priority TINYINT DEFAULT 3,
  completed BOOL DEFAULT false
);

DROP TABLE IF EXISTS users;

CREATE TABLE users (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  first_name VARCHAR(40) NOT NULL,
  last_name VARCHAR(40) NOT NULL,
  password VARCHAR(64) NOT NULL,
  email VARCHAR(64) NOT NULL,
  time_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  account_type VARCHAR(10) DEFAULT "standard",
  encoded_api_key VARCHAR(64) NOT NULL
);

ALTER TABLE tasks
  ADD CONSTRAINT fk_tasks_user
  FOREIGN KEY (user_id)
  REFERENCES users(id);

ALTER TABLE users
  ADD CONSTRAINT uc_email 
  UNIQUE (email);

ALTER TABLE users
  ADD CONSTRAINT uc_full_name 
  UNIQUE (first_name, last_name);

