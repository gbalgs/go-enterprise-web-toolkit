CREATE TABLE users (
    id char(32) NOT NULL,
    name VARCHAR(256) NOT NULL UNIQUE,
    email VARCHAR(256) UNIQUE,
    phone VARCHAR(32) UNIQUE,
    password VARCHAR(256) NOT NULL,
    status TINYINT DEFAULT -1, # -1 inactive , 0 for delete, 1 for normal
    activated_at DATETIME,
    created_at TIMESTAMP NULL DEFAULT NULL,
    updated_at TIMESTAMP NULL DEFAULT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY(`id`)
);