
CREATE TABLE if not exists users (
                       id BIGINT GENERATED ALWAYS AS IDENTITY Primary key,
                       username Varchar(50) UNIQUE NOt NULL,
                       password_hash TEXT NOT NULL,
                       email VARCHAR(100) NOT NULL,
                       real_name VARCHAR(100) NOT NULL,
                       role VARCHAR(7) NOT NULL CHECK (role IN ('client', 'admin')) DEFAULT 'client',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE if not exists ActionLogs (
                            id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                            user_id BIGINT NOT NULL,
                            action_type VARCHAR(50) NOT NULL,
                            action_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            details TEXT
);

CREATE TABLE if not exists VehicleTypes (
                              id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                              type_name VARCHAR(50) NOT NULL,
                              price_per_day MONEY NOT NULL CHECK (price_per_day > 0::money)
);

DO $$
    BEGIN
    IF (SELECT COUNT(*) FROM VehicleTypes) = 0 THEN
        COPY VehicleTypes (type_name, price_per_day)
            FROM '/data_for_lab_2/test_copy/vehicle_types.csv'
            WITH (FORMAT csv, HEADER true);
    END IF;
END $$;


CREATE TABLE if not exists Brands (
                        id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                        brand_name VARCHAR(50) NOT NULL UNIQUE
);

DO $$
    BEGIN
        IF (SELECT COUNT(*) FROM Brands) = 0 THEN
            COPY Brands (brand_name)
                FROM '/data_for_lab_2/test_copy/brands.csv'
                WITH (FORMAT csv, HEADER true);
        END IF;
END $$;


CREATE TABLE if not exists Models (
                        id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                        name VARCHAR(50) NOT NULL,
                        price_per_day MONEY NOT NULL CHECK (price_per_day > 0::money)
);

DO $$
    BEGIN
        IF (SELECT COUNT(*) FROM Models) = 0 THEN
            COPY Models (name, price_per_day)
                FROM '/data_for_lab_2/test_copy/models.csv'
                WITH (FORMAT csv, HEADER true);
        END IF;
END $$;

CREATE TABLE if not exists Vehicles (
                          id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                          type_id BIGINT NOT NULL,
                          model_id BIGINT NOT NULL,
                          brand_id BIGINT NOT NULL,
                          color VARCHAR(50) NOT NULL,
                          status VARCHAR(20) DEFAULT 'available' CHECK (status IN ('available', 'rented')),
                          FOREIGN KEY (type_id) REFERENCES VehicleTypes (id) ON DELETE CASCADE,
                          FOREIGN KEY (model_id) REFERENCES Models (id) ON DELETE CASCADE,
                          FOREIGN KEY (brand_id) REFERENCES Brands (id) ON DELETE CASCADE
);

DO $$
    BEGIN
        IF (SELECT COUNT(*) FROM Vehicles) = 0 THEN
            COPY Vehicles (type_id, model_id, brand_id, color, status)
                FROM '/data_for_lab_2/test_copy/vehicle.csv'
                WITH (FORMAT csv, HEADER true);
        END IF;
    END $$;


CREATE TABLE if not exists VehiclePhotos (
                               id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                               vehicle_id BIGINT NOT NULL,
                               photo_url TEXT NOT NULL,
                               FOREIGN KEY (vehicle_id) REFERENCES Vehicles (id) ON DELETE CASCADE
);

DO $$
    BEGIN
        IF (SELECT COUNT(*) FROM VehiclePhotos) = 0 THEN
            COPY VehiclePhotos (vehicle_id, photo_url)
                FROM '/data_for_lab_2/test_copy/vehicles_photos.csv'
                WITH (FORMAT csv, HEADER true);
        END IF;
    END $$;

CREATE TABLE if not exists Bookings (
                          id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                          user_id BIGINT NOT NULL,
                          vehicle_id BIGINT NOT NULL,
                          date_begin DATE NOT NULL,
                          date_end DATE NOT NULL,
                          FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
                          FOREIGN KEY (vehicle_id) REFERENCES Vehicles (id) ON DELETE CASCADE
);
