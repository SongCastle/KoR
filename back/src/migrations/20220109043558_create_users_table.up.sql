CREATE TABLE IF NOT EXISTS users (
    id int AUTO_INCREMENT NOT NULL,
    login VARCHAR (255) UNIQUE NOT NULL,
    password VARCHAR (255) NOT NULL,
    email VARCHAR (255),
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    PRIMARY KEY (id)
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4;
