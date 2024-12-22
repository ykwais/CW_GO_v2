
CREATE TABLE if not exists Users (
                       id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                       login VARCHAR(50) UNIQUE NOT NULL, -- логин
                       password_hash TEXT NOT NULL

);