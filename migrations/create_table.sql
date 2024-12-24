
CREATE TABLE if not exists Users (
                       id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                       login VARCHAR(50) UNIQUE NOT NULL, -- логин
                       password_hash TEXT NOT NULL

);

DO $$
    BEGIN
    IF (SELECT COUNT(*) FROM Users) = 0 THEN
        COPY Users (login, password_hash)
            FROM '/data_for_lab_2/test_copy/users.csv'
            WITH (FORMAT csv, HEADER true);
    END IF;
END $$;