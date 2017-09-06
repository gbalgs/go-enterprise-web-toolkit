CREATE TABLE user_profiles  (
    id char(32) NOT NULL,
    user_id CHAR(32) NOT NULL UNIQUE,
    avatar VARCHAR(1024),
    sex CHAR(8),
    address1 VARCHAR(1024),
    address2 VARCHAR (1024),
    created_at TIMESTAMP NULL DEFAULT NULL,
    updated_at TIMESTAMP NULL DEFAULT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY(`id`)
);