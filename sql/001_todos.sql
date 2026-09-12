-- Database and table used by every Todo API implementation in this repository.
-- Safe to execute more than once: it does not drop existing data.

CREATE DATABASE IF NOT EXISTS todos
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE todos;

CREATE TABLE IF NOT EXISTS todos (
  id BIGINT NOT NULL AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  title VARCHAR(255) NOT NULL,
  description TEXT NULL,
  completed TINYINT(1) NOT NULL DEFAULT 0,
  due_date DATE NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  INDEX idx_todos_user_id (user_id),
  INDEX idx_todos_completed (completed)
) ENGINE=InnoDB;
