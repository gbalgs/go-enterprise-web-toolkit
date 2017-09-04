CREATE TABLE user IF NOT EXISTS (
    id char(32) NOT NULL,
    email VARCHAR(256),
    phone VARCHAR(32),
    password VARCHAR(256) NOT NULL,
    name VARCHAR(256) not null,
    status tinyint default -1, # -1 for inactive, 0 for delete, 1 for normal
    avatar VARCHAR(256),
    actived_date TIMESTAMP,
    created_date TIMESTAMP NOT NULL,
    updated_date TIMESTAMP,
    deleted_date TIMESTAMP
    PRIMARY KEY(`id`)
);

CREATE TABLE user_profile IF NOT EXISTS (
    user_id CHAR(32) NOT NULL,
    avatar VARCHAR(256),
    sex CHAR(8),
    address1 VARCHAR(1024),
    address2 VARCHAR (1024),
    created_date TIMESTAMP NOT NULL,
    updated_date TIMESTAMP,
    deleted_date TIMESTAMP
    PRIMARY KEY(`user_id`)
)