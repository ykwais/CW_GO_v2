

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