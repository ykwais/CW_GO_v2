
CREATE OR REPLACE FUNCTION register_client(
    user_name VARCHAR(50),
    pass_hash TEXT
) RETURNS VOID
    LANGUAGE plpgsql
AS $$
BEGIN
INSERT INTO Users (login, password_hash)
VALUES (user_name, pass_hash);
END;
$$;