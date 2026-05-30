SELECT 'CREATE DATABASE moviedb'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'moviedb')\gexec

\c moviedb

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    admin BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    poster_url TEXT NOT NULL,
    genre TEXT[] NOT NULL
);
    
CREATE TABLE screens (
    id SERIAL PRIMARY KEY,
    auditorium_number INTEGER NOT NULL,
    normal_seats TEXT[] NOT NULL,
    vip_seats TEXT[] NOT NULL,
    type TEXT NOT NULL
);

CREATE TABLE showtimes (
    id SERIAL PRIMARY KEY,
    movie_id INTEGER NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    schedule JSONB NOT NULL,
    show_date DATE NOT NULL,
    normal_price INTEGER NOT NULL,
    vip_price INTEGER NOT NULL
);

CREATE TABLE show_seats (
    id SERIAL PRIMARY KEY,
    showtime_id INTEGER NOT NULL REFERENCES showtimes(id) ON DELETE CASCADE,
    screen_id INTEGER NOT NULL REFERENCES screens(id) ON DELETE CASCADE,
    show_time TEXT NOT NULL,
    seat_status JSONB NOT NULL,
    UNIQUE(showtime_id, screen_id, show_time)
);

CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    showtime_id INTEGER NOT NULL REFERENCES showtimes(id) ON DELETE CASCADE,
    screen_id INTEGER NOT NULL REFERENCES screens(id) ON DELETE CASCADE,
    seats TEXT[] NOT NULL,
    booking_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO movies (title, description, poster_url, genre) VALUES
(
    'Inception',
    'A skilled thief enters dreams to steal secrets from the subconscious.',
    'https://image.tmdb.org/t/p/original/xlaY2zyzMfkhk0HSC5VUwzoZPU1.jpg',
    ARRAY['Sci-Fi', 'Thriller']
),
(
    'Avengers: Endgame',
    'The Avengers assemble once more to reverse the damage caused by Thanos.',
    'https://image.tmdb.org/t/p/original/ulzhLuWrPK07P1YkdWQLZnQh1JL.jpg',
    ARRAY['Action', 'Adventure', 'Superhero']
),
(
    'Interstellar',
    'A team of astronauts travel through a wormhole in search of a new home for humanity.',
    'https://image.tmdb.org/t/p/original/gEU2QniE6E77NI6lCU6MxlNBvIx.jpg',
    ARRAY['Sci-Fi', 'Drama']
),
(
    'Joker',
    'A failed comedian descends into madness and becomes Gotham''s infamous villain.',
    'https://image.tmdb.org/t/p/original/mZuAPY4ETMQPHhCXIcJ90kd6RaS.jpg',
    ARRAY['Crime', 'Drama', 'Thriller']
),
(
    'The Dark Knight',
    'Batman faces the Joker in a battle for Gotham City.',
    'https://image.tmdb.org/t/p/original/qJ2tW6WMUDux911r6m7haRef0WH.jpg',
    ARRAY['Action', 'Crime', 'Drama']
);

INSERT INTO screens (
    auditorium_number,
    normal_seats,
    vip_seats,
    type
) VALUES

(
    1,
    ARRAY[
        'A1','A2','A3','A4','A5','A6','A7','A8',
        'B1','B2','B3','B4','B5','B6','B7','B8',
        'C1','C2','C3','C4','C5','C6','C7','C8'
    ],
    ARRAY[
        'V1','V2','V3','V4','V5','V6'
    ],
    '2D'
),

(
    2,
    ARRAY[
        'A1','A2','A3','A4','A5','A6',
        'B1','B2','B3','B4','B5','B6',
        'C1','C2','C3','C4','C5','C6',
        'D1','D2','D3','D4','D5','D6'
    ],
    ARRAY[
        'VIP1','VIP2','VIP3','VIP4'
    ],
    '3D'
),

(
    3,
    ARRAY[
        'A1','A2','A3','A4','A5','A6','A7','A8','A9','A10',
        'B1','B2','B3','B4','B5','B6','B7','B8','B9','B10',
        'C1','C2','C3','C4','C5','C6','C7','C8','C9','C10'
    ],
    ARRAY[
        'V1','V2','V3','V4','V5','V6','V7','V8'
    ],
    'IMAX'
),

(
    4,
    ARRAY[
        'A1','A2','A3','A4','A5',
        'B1','B2','B3','B4','B5',
        'C1','C2','C3','C4','C5',
        'D1','D2','D3','D4','D5'
    ],
    ARRAY[
        'VIP1','VIP2','VIP3','VIP4','VIP5'
    ],
    '4DX'
),

(
    5,
    ARRAY[
        'A1','A2','A3','A4','A5','A6','A7',
        'B1','B2','B3','B4','B5','B6','B7',
        'C1','C2','C3','C4','C5','C6','C7',
        'D1','D2','D3','D4','D5','D6','D7'
    ],
    ARRAY[
        'R1','R2','R3','R4','R5','R6'
    ],
    'Dolby Atmos'
);

INSERT INTO showtimes (
    movie_id,
    schedule,
    show_date,
    normal_price,
    vip_price
) VALUES

(
    1,
    '[{"screen_id": 1, "timings": ["10:00", "14:00"]}]'::jsonb,
    '2027-06-01',
    200,
    350
),
(
    1,
    '[{"screen_id": 3, "timings": ["18:00", "21:30"]}]'::jsonb,
    '2027-06-02',
    250,
    450
),

(
    2,
    '[{"screen_id": 2, "timings": ["09:30", "13:30"]}]'::jsonb,
    '2027-06-01',
    220,
    400
),
(
    2,
    '[{"screen_id": 5, "timings": ["17:00", "20:30"]}]'::jsonb,
    '2027-06-02',
    250,
    450
),

(
    3,
    '[{"screen_id": 3, "timings": ["11:00", "15:00"]}]'::jsonb,
    '2027-06-01',
    250,
    450
),
(
    3,
    '[{"screen_id": 1, "timings": ["19:00", "22:00"]}]'::jsonb,
    '2027-06-03',
    200,
    350
),

(
    4,
    '[{"screen_id": 4, "timings": ["12:00", "16:00"]}]'::jsonb,
    '2027-06-01',
    230,
    420
),
(
    4,
    '[{"screen_id": 2, "timings": ["18:30", "21:00"]}]'::jsonb,
    '2027-06-03',
    220,
    400
),

(
    5,
    '[{"screen_id": 5, "timings": ["10:30", "14:30"]}]'::jsonb,
    '2027-06-02',
    250,
    450
),
(
    5,
    '[{"screen_id": 3, "timings": ["19:30", "22:30"]}]'::jsonb,
    '2027-06-03',
    250,
    450
);

INSERT INTO show_seats (
    showtime_id,
    screen_id,
    show_time,
    seat_status
)
SELECT
    st.id,
    (schedule_item->>'screen_id')::INTEGER,
    timing,
    (
        SELECT jsonb_object_agg(seat, 'available')
        FROM (
            SELECT unnest(s.normal_seats || s.vip_seats) AS seat
        ) seats
    )
FROM showtimes st
CROSS JOIN LATERAL jsonb_array_elements(st.schedule) schedule_item
CROSS JOIN LATERAL jsonb_array_elements_text(schedule_item->'timings') timing
JOIN screens s
    ON s.id = (schedule_item->>'screen_id')::INTEGER;