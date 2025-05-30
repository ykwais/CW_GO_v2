
CREATE OR REPLACE FUNCTION register_user( --вроде норм
    p_username VARCHAR(50),
    p_password_hash TEXT,
    p_email VARCHAR(100),
    p_real_name VARCHAR(100),
    p_role VARCHAR(7) DEFAULT 'client'
)
    RETURNS BIGINT

    LANGUAGE plpgsql
AS $$
DECLARE
    v_user_id BIGINT;
BEGIN
    INSERT INTO users (username, password_hash, email, real_name, role)
    VALUES (p_username, p_password_hash, p_email, p_real_name, p_role)
    RETURNING id INTO v_user_id;
    RETURN v_user_id;
END;
$$ ;





CREATE OR REPLACE FUNCTION login_user(--вроде норм
    p_username VARCHAR(50)
) RETURNS TABLE(id BIGINT, login VARCHAR(50), pass_hash Text)
    LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
        SELECT users.id, username, password_hash
        FROM users
        WHERE username = p_username;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Invalid username';
    END IF;
END;
$$;



CREATE OR REPLACE FUNCTION get_available_vehicles( --вроде норм
    p_date_begin DATE,
    p_date_end DATE
)
    RETURNS TABLE (
                      veh_id BIGINT,
                      brand_name VARCHAR,
                      model_name VARCHAR,
                      price_per_day MONEY,
                      photo_url TEXT
                  ) AS $$
BEGIN
    RETURN QUERY
        SELECT
            vd.vehicle_id,
            vd.brand_name,
            vd.model_name,
            vd.total_price_per_day,
            vd.photo_url
        FROM VehicleDetails vd
        WHERE vd.vehicle_id NOT IN (
            SELECT Bookings.vehicle_id
            FROM Bookings
            WHERE (p_date_begin <= date_end AND p_date_end >= date_begin)
        );
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION book_vehicle( --вроде норм
    p_user_id BIGINT,
    p_vehicle_id BIGINT,
    p_date_begin DATE,
    p_date_end DATE
) RETURNS BIGINT AS $$
DECLARE
    v_booking_id BIGINT;
BEGIN

    IF EXISTS (
        SELECT 1
        FROM bookings
        WHERE vehicle_id = p_vehicle_id
          AND (
            (p_date_begin BETWEEN date_begin AND date_end)
                OR (p_date_end BETWEEN date_begin AND date_end)
                OR (date_begin BETWEEN p_date_begin AND p_date_end)
            )
    ) THEN
        RAISE EXCEPTION 'Vehicle is not available for the selected dates';
    END IF;


    INSERT INTO bookings (user_id, vehicle_id, date_begin, date_end)
    VALUES (p_user_id, p_vehicle_id, p_date_begin, p_date_end)
    RETURNING id INTO v_booking_id;


    UPDATE vehicles
    SET status = 'rented'
    WHERE id = p_vehicle_id;


    RETURN v_booking_id;
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION cancel_booking(--вроде ок
    p_user_id BIGINT,
    p_vehicle_id BIGINT
) RETURNS VOID AS $$
BEGIN

    IF NOT EXISTS (
        SELECT 1
        FROM bookings
        WHERE user_id = p_user_id
          AND vehicle_id = p_vehicle_id
    ) THEN
        RAISE EXCEPTION 'No booking found for this user and vehicle';
    END IF;


    DELETE FROM bookings
    WHERE user_id = p_user_id AND vehicle_id = p_vehicle_id;


    UPDATE vehicles
    SET status = 'available'
    WHERE id = p_vehicle_id;
END;
$$ LANGUAGE plpgsql;




CREATE OR REPLACE FUNCTION log_user_action() --вроде ок
    RETURNS TRIGGER AS $$
BEGIN

    IF TG_TABLE_NAME = 'users' THEN
        INSERT INTO ActionLogs (user_id, action_type, details)
        VALUES (NEW.id, 'registration', 'New user registered: ' || NEW.username);
    END IF;

    IF TG_TABLE_NAME = 'bookings' THEN
        INSERT INTO ActionLogs (user_id, action_type, details)
        VALUES (NEW.user_id, 'booking', 'Vehicle booked: ' || NEW.vehicle_id || ' from ' || NEW.date_begin || ' to ' || NEW.date_end);
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION delete_user(p_user_id BIGINT)--
    RETURNS VOID AS $$
BEGIN

    UPDATE vehicles
    SET status = 'available'
    WHERE id IN (
        SELECT vehicle_id
        FROM bookings
        WHERE user_id = p_user_id
    );


    DELETE FROM bookings
    WHERE user_id = p_user_id;


    DELETE FROM users
    WHERE id = p_user_id;


    INSERT INTO ActionLogs (user_id, action_type, details)
    VALUES (p_user_id, 'deletion', 'User deleted along with their bookings.');
END;
$$ LANGUAGE plpgsql;




CREATE OR REPLACE FUNCTION get_vehicle_photos_table(_id BIGINT) -- вроде норм
    RETURNS TABLE(photo_url TEXT) AS $$
BEGIN
    RETURN QUERY
        SELECT VehiclePhotos.photo_url
        FROM VehiclePhotos
        WHERE vehicle_id = _id;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION get_user_bookings(p_user_id BIGINT)--вроде ок
    RETURNS TABLE (
                      vehicle_id BIGINT,
                      brand_name VARCHAR,
                      model_name VARCHAR,
                      date_begin DATE,
                      date_end DATE
                  ) AS $$
BEGIN
    RETURN QUERY
        SELECT
            b.vehicle_id As vehicle_id,
            br.brand_name,
            m.name AS model_name,
            b.date_begin,
            b.date_end
        FROM bookings b
                 JOIN vehicles v ON b.vehicle_id = v.id
                 JOIN brands br ON v.brand_id = br.id
                 JOIN models m ON v.model_id = m.id
        WHERE b.user_id = p_user_id;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION get_admin_overview() -- вроде ок
    RETURNS TABLE (
                      login VARCHAR(50),
                      user_email VARCHAR(100),
                      user_real_name VARCHAR(100),
                      brand_name VARCHAR(50),
                      model_name VARCHAR(50),
                      booking_start_date Date,
                      booking_end_date Date,
                      total_price_per_day money
                  ) AS $$
BEGIN
    RETURN QUERY SELECT * FROM AdminOverview;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION get_user_details() --вроде ок
    RETURNS TABLE (
                      user_id BIGINT,
                      username VARCHAR(50),
                      email VARCHAR(100),
                      real_name VARCHAR(50),
                      created_at TIMESTAMP,
                      total_bookings BIGINT
                  ) AS $$
BEGIN
    RETURN QUERY SELECT * FROM UserDetails;
END;
$$ LANGUAGE plpgsql;