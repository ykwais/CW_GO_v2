PGDMP  +                    |            afdb    17.0 (Debian 17.0-1.pgdg120+1)    17.0 (Debian 17.0-1.pgdg120+1) E    �           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                           false            �           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                           false            �           0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                           false            �           1262    16384    afdb    DATABASE     o   CREATE DATABASE afdb WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';
    DROP DATABASE afdb;
                     ykwais    false            
            2615    20492    cw_test    SCHEMA        CREATE SCHEMA cw_test;
    DROP SCHEMA cw_test;
                     ykwais    false                       1255    23324 (   book_vehicle(bigint, bigint, date, date)    FUNCTION     �  CREATE FUNCTION cw_test.book_vehicle(p_user_id bigint, p_vehicle_id bigint, p_date_begin date, p_date_end date) RETURNS bigint
    LANGUAGE plpgsql
    AS $$
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
$$;
 o   DROP FUNCTION cw_test.book_vehicle(p_user_id bigint, p_vehicle_id bigint, p_date_begin date, p_date_end date);
       cw_test               ykwais    false    10                       1255    23325    cancel_booking(bigint, bigint)    FUNCTION     )  CREATE FUNCTION cw_test.cancel_booking(p_user_id bigint, p_vehicle_id bigint) RETURNS void
    LANGUAGE plpgsql
    AS $$
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
$$;
 M   DROP FUNCTION cw_test.cancel_booking(p_user_id bigint, p_vehicle_id bigint);
       cw_test               ykwais    false    10                       1255    23327    delete_user(bigint)    FUNCTION       CREATE FUNCTION cw_test.delete_user(p_user_id bigint) RETURNS void
    LANGUAGE plpgsql
    AS $$
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
$$;
 5   DROP FUNCTION cw_test.delete_user(p_user_id bigint);
       cw_test               ykwais    false    10            �            1255    23330    get_admin_overview()    FUNCTION     {  CREATE FUNCTION cw_test.get_admin_overview() RETURNS TABLE(login character varying, user_email character varying, user_real_name character varying, brand_name character varying, model_name character varying, booking_start_date date, booking_end_date date, total_price_per_day money)
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN QUERY SELECT * FROM AdminOverview;
END;
$$;
 ,   DROP FUNCTION cw_test.get_admin_overview();
       cw_test               ykwais    false    10                       1255    23323 "   get_available_vehicles(date, date)    FUNCTION     �  CREATE FUNCTION cw_test.get_available_vehicles(p_date_begin date, p_date_end date) RETURNS TABLE(veh_id bigint, brand_name character varying, model_name character varying, price_per_day money, photo_url text)
    LANGUAGE plpgsql
    AS $$
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
$$;
 R   DROP FUNCTION cw_test.get_available_vehicles(p_date_begin date, p_date_end date);
       cw_test               ykwais    false    10            �            1255    23329    get_user_bookings(bigint)    FUNCTION     �  CREATE FUNCTION cw_test.get_user_bookings(p_user_id bigint) RETURNS TABLE(vehicle_id bigint, brand_name character varying, model_name character varying, date_begin date, date_end date)
    LANGUAGE plpgsql
    AS $$
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
$$;
 ;   DROP FUNCTION cw_test.get_user_bookings(p_user_id bigint);
       cw_test               ykwais    false    10                        1255    23331    get_user_details()    FUNCTION     8  CREATE FUNCTION cw_test.get_user_details() RETURNS TABLE(user_id bigint, username character varying, email character varying, real_name character varying, created_at timestamp without time zone, total_bookings bigint)
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN QUERY SELECT * FROM UserDetails;
END;
$$;
 *   DROP FUNCTION cw_test.get_user_details();
       cw_test               ykwais    false    10            �            1255    23328     get_vehicle_photos_table(bigint)    FUNCTION       CREATE FUNCTION cw_test.get_vehicle_photos_table(_id bigint) RETURNS TABLE(photo_url text)
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN QUERY
        SELECT VehiclePhotos.photo_url
        FROM VehiclePhotos
        WHERE vehicle_id = _id;
END;
$$;
 <   DROP FUNCTION cw_test.get_vehicle_photos_table(_id bigint);
       cw_test               ykwais    false    10                       1255    23326    log_user_action()    FUNCTION     D  CREATE FUNCTION cw_test.log_user_action() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
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
$$;
 )   DROP FUNCTION cw_test.log_user_action();
       cw_test               ykwais    false    10                       1255    23322    login_user(character varying)    FUNCTION     �  CREATE FUNCTION cw_test.login_user(p_username character varying) RETURNS TABLE(id bigint, login character varying, pass_hash text)
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
 @   DROP FUNCTION cw_test.login_user(p_username character varying);
       cw_test               ykwais    false    10            �            1255    23321 _   register_user(character varying, text, character varying, character varying, character varying)    FUNCTION     �  CREATE FUNCTION cw_test.register_user(p_username character varying, p_password_hash text, p_email character varying, p_real_name character varying, p_role character varying DEFAULT 'client'::character varying) RETURNS bigint
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
$$;
 �   DROP FUNCTION cw_test.register_user(p_username character varying, p_password_hash text, p_email character varying, p_real_name character varying, p_role character varying);
       cw_test               ykwais    false    10            �            1259    23346 
   actionlogs    TABLE     �   CREATE TABLE cw_test.actionlogs (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    action_type character varying(50) NOT NULL,
    action_timestamp timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    details text
);
    DROP TABLE cw_test.actionlogs;
       cw_test         heap r       ykwais    false    10            �            1259    23345    actionlogs_id_seq    SEQUENCE     �   ALTER TABLE cw_test.actionlogs ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME cw_test.actionlogs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);
            cw_test               ykwais    false    10    225            �            1259    23413    bookings    TABLE     �   CREATE TABLE cw_test.bookings (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    vehicle_id bigint NOT NULL,
    date_begin date NOT NULL,
    date_end date NOT NULL
);
    DROP TABLE cw_test.bookings;
       cw_test         heap r       ykwais    false    10            �            1259    23362    brands    TABLE     g   CREATE TABLE cw_test.brands (
    id bigint NOT NULL,
    brand_name character varying(50) NOT NULL
);
    DROP TABLE cw_test.brands;
       cw_test         heap r       ykwais    false    10            �            1259    23370    models    TABLE     �   CREATE TABLE cw_test.models (
    id bigint NOT NULL,
    name character varying(50) NOT NULL,
    price_per_day money NOT NULL,
    CONSTRAINT models_price_per_day_check CHECK ((price_per_day > (0)::money))
);
    DROP TABLE cw_test.models;
       cw_test         heap r       ykwais    false    10            �            1259    23333    users    TABLE     �  CREATE TABLE cw_test.users (
    id bigint NOT NULL,
    username character varying(50) NOT NULL,
    password_hash text NOT NULL,
    email character varying(100) NOT NULL,
    real_name character varying(100) NOT NULL,
    role character varying(7) DEFAULT 'client'::character varying NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_role_check CHECK (((role)::text = ANY ((ARRAY['client'::character varying, 'admin'::character varying])::text[])))
);
    DROP TABLE cw_test.users;
       cw_test         heap r       ykwais    false    10            �            1259    23377    vehicles    TABLE     �  CREATE TABLE cw_test.vehicles (
    id bigint NOT NULL,
    type_id bigint NOT NULL,
    model_id bigint NOT NULL,
    brand_id bigint NOT NULL,
    color character varying(50) NOT NULL,
    status character varying(20) DEFAULT 'available'::character varying,
    CONSTRAINT vehicles_status_check CHECK (((status)::text = ANY ((ARRAY['available'::character varying, 'rented'::character varying])::text[])))
);
    DROP TABLE cw_test.vehicles;
       cw_test         heap r       ykwais    false    10            �            1259    23355    vehicletypes    TABLE     �   CREATE TABLE cw_test.vehicletypes (
    id bigint NOT NULL,
    type_name character varying(50) NOT NULL,
    price_per_day money NOT NULL,
    CONSTRAINT vehicletypes_price_per_day_check CHECK ((price_per_day > (0)::money))
);
 !   DROP TABLE cw_test.vehicletypes;
       cw_test         heap r       ykwais    false    10            �            1259    23435    adminoverview    VIEW     n  CREATE VIEW cw_test.adminoverview AS
 SELECT u.username AS login,
    u.email AS user_email,
    u.real_name AS user_real_name,
    br.brand_name,
    m.name AS model_name,
    b.date_begin AS booking_start_date,
    b.date_end AS booking_end_date,
    (vt.price_per_day + m.price_per_day) AS total_price_per_day
   FROM (((((cw_test.bookings b
     JOIN cw_test.users u ON ((b.user_id = u.id)))
     JOIN cw_test.vehicles v ON ((b.vehicle_id = v.id)))
     JOIN cw_test.vehicletypes vt ON ((v.type_id = vt.id)))
     JOIN cw_test.brands br ON ((v.brand_id = br.id)))
     JOIN cw_test.models m ON ((v.model_id = m.id)));
 !   DROP VIEW cw_test.adminoverview;
       cw_test       v       ykwais    false    233    233    223    231    223    223    231    223    231    229    229    227    227    237    237    237    237    233    233    10            �            1259    23412    bookings_id_seq    SEQUENCE     �   ALTER TABLE cw_test.bookings ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME cw_test.bookings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);
            cw_test               ykwais    false    10    237            �            1259    23361    brands_id_seq    SEQUENCE     �   ALTER TABLE cw_test.brands ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME cw_test.brands_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);
            cw_test               ykwais    false    10    229            �            1259    23369    models_id_seq    SEQUENCE     �   ALTER TABLE cw_test.models ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME cw_test.models_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);
            cw_test               ykwais    false    10    231            �            1259    23440    userdetails    VIEW     e  CREATE VIEW cw_test.userdetails AS
 SELECT u.id AS user_id,
    u.username,
    u.email,
    u.real_name,
    u.created_at,
    count(b.id) AS total_bookings
   FROM (cw_test.users u
     LEFT JOIN cw_test.bookings b ON ((u.id = b.user_id)))
  WHERE ((u.role)::text <> 'admin'::text)
  GROUP BY u.id, u.username, u.email, u.real_name, u.role, u.created_at;
    DROP VIEW cw_test.userdetails;
       cw_test       v       ykwais    false    237    223    223    223    223    223    223    237    10            �            1259    23332    users_id_seq    SEQUENCE     �   ALTER TABLE cw_test.users ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME cw_test.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);
            cw_test               ykwais    false    223    10            �            1259    23400    vehiclephotos    TABLE     |   CREATE TABLE cw_test.vehiclephotos (
    id bigint NOT NULL,
    vehicle_id bigint NOT NULL,
    photo_url text NOT NULL
);
 "   DROP TABLE cw_test.vehiclephotos;
       cw_test         heap r       ykwais    false    10            �            1259    23430    vehicledetails    VIEW     e  CREATE VIEW cw_test.vehicledetails AS
 SELECT v.id AS vehicle_id,
    vt.type_name AS vehicle_type,
    b.brand_name,
    m.name AS model_name,
    v.color,
    (vt.price_per_day + m.price_per_day) AS total_price_per_day,
    vp.photo_url
   FROM ((((cw_test.vehicles v
     JOIN cw_test.vehicletypes vt ON ((v.type_id = vt.id)))
     JOIN cw_test.brands b ON ((v.brand_id = b.id)))
     JOIN cw_test.models m ON ((v.model_id = m.id)))
     LEFT JOIN LATERAL ( SELECT vehiclephotos.photo_url
           FROM cw_test.vehiclephotos
          WHERE (vehiclephotos.vehicle_id = v.id)
         LIMIT 1) vp ON (true));
 "   DROP VIEW cw_test.vehicledetails;
       cw_test       v       ykwais    false    227    233    231    233    233    235    231    235    231    229    229    227    233    227    233    10            �            1259    23399    vehiclephotos_id_seq    SEQUENCE     �   ALTER TABLE cw_test.vehiclephotos ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME cw_test.vehiclephotos_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);
            cw_test               ykwais    false    10    235            �            1259    23376    vehicles_id_seq    SEQUENCE     �   ALTER TABLE cw_test.vehicles ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME cw_test.vehicles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);
            cw_test               ykwais    false    233    10            �            1259    23354    vehicletypes_id_seq    SEQUENCE     �   ALTER TABLE cw_test.vehicletypes ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME cw_test.vehicletypes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);
            cw_test               ykwais    false    10    227            �          0    23346 
   actionlogs 
   TABLE DATA           Z   COPY cw_test.actionlogs (id, user_id, action_type, action_timestamp, details) FROM stdin;
    cw_test               ykwais    false    225   �m       �          0    23413    bookings 
   TABLE DATA           R   COPY cw_test.bookings (id, user_id, vehicle_id, date_begin, date_end) FROM stdin;
    cw_test               ykwais    false    237   nn       �          0    23362    brands 
   TABLE DATA           1   COPY cw_test.brands (id, brand_name) FROM stdin;
    cw_test               ykwais    false    229   �n       �          0    23370    models 
   TABLE DATA           :   COPY cw_test.models (id, name, price_per_day) FROM stdin;
    cw_test               ykwais    false    231   �n       �          0    23333    users 
   TABLE DATA           a   COPY cw_test.users (id, username, password_hash, email, real_name, role, created_at) FROM stdin;
    cw_test               ykwais    false    223   o       �          0    23400    vehiclephotos 
   TABLE DATA           C   COPY cw_test.vehiclephotos (id, vehicle_id, photo_url) FROM stdin;
    cw_test               ykwais    false    235   �o       �          0    23377    vehicles 
   TABLE DATA           S   COPY cw_test.vehicles (id, type_id, model_id, brand_id, color, status) FROM stdin;
    cw_test               ykwais    false    233   �p       �          0    23355    vehicletypes 
   TABLE DATA           E   COPY cw_test.vehicletypes (id, type_name, price_per_day) FROM stdin;
    cw_test               ykwais    false    227   �p       �           0    0    actionlogs_id_seq    SEQUENCE SET     @   SELECT pg_catalog.setval('cw_test.actionlogs_id_seq', 5, true);
          cw_test               ykwais    false    224            �           0    0    bookings_id_seq    SEQUENCE SET     >   SELECT pg_catalog.setval('cw_test.bookings_id_seq', 1, true);
          cw_test               ykwais    false    236            �           0    0    brands_id_seq    SEQUENCE SET     <   SELECT pg_catalog.setval('cw_test.brands_id_seq', 4, true);
          cw_test               ykwais    false    228            �           0    0    models_id_seq    SEQUENCE SET     <   SELECT pg_catalog.setval('cw_test.models_id_seq', 4, true);
          cw_test               ykwais    false    230            �           0    0    users_id_seq    SEQUENCE SET     ;   SELECT pg_catalog.setval('cw_test.users_id_seq', 3, true);
          cw_test               ykwais    false    222            �           0    0    vehiclephotos_id_seq    SEQUENCE SET     D   SELECT pg_catalog.setval('cw_test.vehiclephotos_id_seq', 24, true);
          cw_test               ykwais    false    234            �           0    0    vehicles_id_seq    SEQUENCE SET     >   SELECT pg_catalog.setval('cw_test.vehicles_id_seq', 4, true);
          cw_test               ykwais    false    232            �           0    0    vehicletypes_id_seq    SEQUENCE SET     B   SELECT pg_catalog.setval('cw_test.vehicletypes_id_seq', 4, true);
          cw_test               ykwais    false    226            �           2606    23353    actionlogs actionlogs_pkey 
   CONSTRAINT     Y   ALTER TABLE ONLY cw_test.actionlogs
    ADD CONSTRAINT actionlogs_pkey PRIMARY KEY (id);
 E   ALTER TABLE ONLY cw_test.actionlogs DROP CONSTRAINT actionlogs_pkey;
       cw_test                 ykwais    false    225            �           2606    23417    bookings bookings_pkey 
   CONSTRAINT     U   ALTER TABLE ONLY cw_test.bookings
    ADD CONSTRAINT bookings_pkey PRIMARY KEY (id);
 A   ALTER TABLE ONLY cw_test.bookings DROP CONSTRAINT bookings_pkey;
       cw_test                 ykwais    false    237            �           2606    23368    brands brands_brand_name_key 
   CONSTRAINT     ^   ALTER TABLE ONLY cw_test.brands
    ADD CONSTRAINT brands_brand_name_key UNIQUE (brand_name);
 G   ALTER TABLE ONLY cw_test.brands DROP CONSTRAINT brands_brand_name_key;
       cw_test                 ykwais    false    229            �           2606    23366    brands brands_pkey 
   CONSTRAINT     Q   ALTER TABLE ONLY cw_test.brands
    ADD CONSTRAINT brands_pkey PRIMARY KEY (id);
 =   ALTER TABLE ONLY cw_test.brands DROP CONSTRAINT brands_pkey;
       cw_test                 ykwais    false    229            �           2606    23375    models models_pkey 
   CONSTRAINT     Q   ALTER TABLE ONLY cw_test.models
    ADD CONSTRAINT models_pkey PRIMARY KEY (id);
 =   ALTER TABLE ONLY cw_test.models DROP CONSTRAINT models_pkey;
       cw_test                 ykwais    false    231            �           2606    23342    users users_pkey 
   CONSTRAINT     O   ALTER TABLE ONLY cw_test.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);
 ;   ALTER TABLE ONLY cw_test.users DROP CONSTRAINT users_pkey;
       cw_test                 ykwais    false    223            �           2606    23344    users users_username_key 
   CONSTRAINT     X   ALTER TABLE ONLY cw_test.users
    ADD CONSTRAINT users_username_key UNIQUE (username);
 C   ALTER TABLE ONLY cw_test.users DROP CONSTRAINT users_username_key;
       cw_test                 ykwais    false    223            �           2606    23406     vehiclephotos vehiclephotos_pkey 
   CONSTRAINT     _   ALTER TABLE ONLY cw_test.vehiclephotos
    ADD CONSTRAINT vehiclephotos_pkey PRIMARY KEY (id);
 K   ALTER TABLE ONLY cw_test.vehiclephotos DROP CONSTRAINT vehiclephotos_pkey;
       cw_test                 ykwais    false    235            �           2606    23383    vehicles vehicles_pkey 
   CONSTRAINT     U   ALTER TABLE ONLY cw_test.vehicles
    ADD CONSTRAINT vehicles_pkey PRIMARY KEY (id);
 A   ALTER TABLE ONLY cw_test.vehicles DROP CONSTRAINT vehicles_pkey;
       cw_test                 ykwais    false    233            �           2606    23360    vehicletypes vehicletypes_pkey 
   CONSTRAINT     ]   ALTER TABLE ONLY cw_test.vehicletypes
    ADD CONSTRAINT vehicletypes_pkey PRIMARY KEY (id);
 I   ALTER TABLE ONLY cw_test.vehicletypes DROP CONSTRAINT vehicletypes_pkey;
       cw_test                 ykwais    false    227            �           2620    23429    bookings after_booking_insert    TRIGGER     ~   CREATE TRIGGER after_booking_insert AFTER INSERT ON cw_test.bookings FOR EACH ROW EXECUTE FUNCTION cw_test.log_user_action();
 7   DROP TRIGGER after_booking_insert ON cw_test.bookings;
       cw_test               ykwais    false    261    237            �           2620    23428    users after_user_insert    TRIGGER     x   CREATE TRIGGER after_user_insert AFTER INSERT ON cw_test.users FOR EACH ROW EXECUTE FUNCTION cw_test.log_user_action();
 1   DROP TRIGGER after_user_insert ON cw_test.users;
       cw_test               ykwais    false    261    223            �           2606    23418    bookings bookings_user_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY cw_test.bookings
    ADD CONSTRAINT bookings_user_id_fkey FOREIGN KEY (user_id) REFERENCES cw_test.users(id) ON DELETE CASCADE;
 I   ALTER TABLE ONLY cw_test.bookings DROP CONSTRAINT bookings_user_id_fkey;
       cw_test               ykwais    false    237    223    3282            �           2606    23423 !   bookings bookings_vehicle_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY cw_test.bookings
    ADD CONSTRAINT bookings_vehicle_id_fkey FOREIGN KEY (vehicle_id) REFERENCES cw_test.vehicles(id) ON DELETE CASCADE;
 L   ALTER TABLE ONLY cw_test.bookings DROP CONSTRAINT bookings_vehicle_id_fkey;
       cw_test               ykwais    false    237    3296    233            �           2606    23407 +   vehiclephotos vehiclephotos_vehicle_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY cw_test.vehiclephotos
    ADD CONSTRAINT vehiclephotos_vehicle_id_fkey FOREIGN KEY (vehicle_id) REFERENCES cw_test.vehicles(id) ON DELETE CASCADE;
 V   ALTER TABLE ONLY cw_test.vehiclephotos DROP CONSTRAINT vehiclephotos_vehicle_id_fkey;
       cw_test               ykwais    false    233    235    3296            �           2606    23394    vehicles vehicles_brand_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY cw_test.vehicles
    ADD CONSTRAINT vehicles_brand_id_fkey FOREIGN KEY (brand_id) REFERENCES cw_test.brands(id) ON DELETE CASCADE;
 J   ALTER TABLE ONLY cw_test.vehicles DROP CONSTRAINT vehicles_brand_id_fkey;
       cw_test               ykwais    false    233    3292    229            �           2606    23389    vehicles vehicles_model_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY cw_test.vehicles
    ADD CONSTRAINT vehicles_model_id_fkey FOREIGN KEY (model_id) REFERENCES cw_test.models(id) ON DELETE CASCADE;
 J   ALTER TABLE ONLY cw_test.vehicles DROP CONSTRAINT vehicles_model_id_fkey;
       cw_test               ykwais    false    231    233    3294            �           2606    23384    vehicles vehicles_type_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY cw_test.vehicles
    ADD CONSTRAINT vehicles_type_id_fkey FOREIGN KEY (type_id) REFERENCES cw_test.vehicletypes(id) ON DELETE CASCADE;
 I   ALTER TABLE ONLY cw_test.vehicles DROP CONSTRAINT vehicles_type_id_fkey;
       cw_test               ykwais    false    233    227    3288            �   �   x�}�A��0E��)|"���!�1�� ����Ҍ����0KK�=�oT��ti�e_S�ɭ�V� R����Fm����R���'M��s�!��]��U�Yq����l��2�/��qX���S������nֶ2�K+yٙ��j�4������.��P{I��&�Y��K]$;?�l���i��Z���f�      �      x������ � �      �   /   x�3�ɯ�/I�2�t�/J�2�t�H-+��I-�2�t������ ��
\      �   @   x�3�t�/���I�T15�30�2�t�545�T14�9#L���	�Jeqqj%��9X F��� �O�      �   �   x�m��
�@����U�p���9i��4Ґ��JhS��dX1j��'T������DÃ�fH+(�깸?�z���k���E�2��ۼ�>�V�n"_+�l�����N*�6�ȼn��u@������ѱ�E�'�u=/C_ݣz+T��.Ir~f�<z/q��ؓU7S�5�!���{����w��I6r���{F)}A      �   �   x�m�K
�0F�q��4���v/BQ���4���H&���8�H�n}�=��,��k}hFą"W�#��P("��zō�i�)�Z8 ����E&�DȤ��D�S�Y���oG�$�L2) ��Y����Y��ovy�o�ﶀ��	�D2�Fv�$�= )<��PR�#      �   @   x�3�4C���T�Ĳ�̜Ĥ�T.cN4
'&g#��p��1gxFf	�zCNJMA����� ��      �   >   x�3�NMI��T15�30�2��T� s�9C�J��9U��|N�̼�2�js�@� �_6     