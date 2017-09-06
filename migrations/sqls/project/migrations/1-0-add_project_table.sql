CREATE TABLE projects (
    id char(32) NOT NULL,
    name VARCHAR(256) NOT NULL,
    status TINYINT DEFAULT 0, # -1 delete
    description VARCHAR(512),
    owner_id char(32) NOT NULL,
    created_at TIMESTAMP NULL DEFAULT NULL,
    updated_at TIMESTAMP NULL DEFAULT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY(`id`)
);