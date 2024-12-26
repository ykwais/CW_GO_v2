

CREATE OR REPLACE VIEW VehicleDetails AS
SELECT
    v.id AS vehicle_id,
    vt.type_name AS vehicle_type,
    b.brand_name,
    m.name AS model_name,
    v.color,
    (vt.price_per_day + m.price_per_day) AS total_price_per_day,
    vp.photo_url
FROM Vehicles v
         JOIN VehicleTypes vt ON v.type_id = vt.id
         JOIN Brands b ON v.brand_id = b.id
         JOIN Models m ON v.model_id = m.id
         LEFT JOIN LATERAL (
    SELECT photo_url
    FROM VehiclePhotos
    WHERE vehicle_id = v.id
    LIMIT 1
    ) vp ON true;


CREATE OR REPLACE VIEW AdminOverview AS
SELECT
    u.username AS login,
    u.email AS user_email,
    u.real_name AS user_real_name,
    br.brand_name AS brand_name,
    m.name AS model_name,
    b.date_begin AS booking_start_date,
    b.date_end AS booking_end_date,
    (vt.price_per_day + m.price_per_day) AS total_price_per_day
FROM Bookings b
         JOIN Users u ON b.user_id = u.id
         JOIN Vehicles v ON b.vehicle_id = v.id
         JOIN VehicleTypes vt ON v.type_id = vt.id
         JOIN Brands br ON v.brand_id = br.id
         JOIN Models m ON v.model_id = m.id;


CREATE OR REPLACE VIEW UserDetails AS
SELECT
    u.id as user_id,
    u.username,
    u.email,
    u.real_name,
    u.created_at,
    COUNT(b.id) AS total_bookings
FROM
    users u
        LEFT JOIN bookings b ON u.id = b.user_id
WHERE
    u.role != 'admin'
GROUP BY
    u.id, u.username, u.email, u.real_name, u.role, u.created_at;