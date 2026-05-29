# Movie Reservation API Documentation

This documentation provides a comprehensive guide to the Movie Reservation System API. It is designed to be easily parsed by both humans and LLM agents.

## Base URL
All endpoints are relative to: `http://<host>:<port>/api`

Health check: `GET /health` (outside `/api`)

## Authentication
Most endpoints require a JWT token for authorization.
1. Use the `/register` endpoint to create an account.
2. Use the `/login` endpoint to obtain a JWT token.
3. Include the token in the `Authorization` header of subsequent requests:
   `Authorization: Bearer <your_token>`

The token is valid for 72 hours.

---

## API Endpoints

### 1. Authentication & User Management

#### Register User
- **Method:** `POST`
- **Path:** `/register`
- **Auth Required:** No
- **Request Body:**
  ```json
  {
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepassword",
    "admin": false
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "message": "User created"
  }
  ```

#### Login User
- **Method:** `POST`
- **Path:** `/login`
- **Auth Required:** No
- **Request Body:**
  ```json
  {
    "email": "john@example.com",
    "password": "securepassword"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "token": "eyJhbGci..."
  }
  ```

#### Promote User to Admin
- **Method:** `PATCH`
- **Path:** `/user/promote`
- **Auth Required:** Yes (Admin role)
- **Request Body:**
  ```json
  {
    "email": "user@example.com"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "message": "User Promoted to Admin"
  }
  ```

#### Get Current User Status
- **Method:** `GET`
- **Path:** `/user/me`
- **Auth Required:** Yes
- **Response (200 OK):**
  ```json
  {
    "message": "Protected route"
  }
  ```

---

### 2. Movie Management

#### Add Movie
- **Method:** `POST`
- **Path:** `/movie/add`
- **Auth Required:** Yes (Admin role)
- **Request Body:**
  ```json
  {
    "title": "Inception",
    "description": "A thief who steals corporate secrets through the use of dream-sharing technology.",
    "poster_url": "http://example.com/poster.jpg",
    "genre": ["Sci-Fi", "Action"]
  }
  ```

#### Update Movie
- **Method:** `PATCH`
- **Path:** `/movie/update`
- **Auth Required:** Yes (Admin role)
- **Request Body:** Same as Add Movie. Updates are matched by `title`.

#### Delete Movie
- **Method:** `DELETE`
- **Path:** `/movie/delete`
- **Auth Required:** Yes (Admin role)
- **Request Body:**
  ```json
  {
    "title": "Inception"
  }
  ```

#### Get All Movie Titles
- **Method:** `GET`
- **Path:** `/movie/get`
- **Auth Required:** Yes
- **Response (200 OK):**
  ```json
  ["Inception", "Interstellar"]
  ```

#### Get Movie Timings
- **Method:** `GET`
- **Path:** `/movie/timings`
- **Auth Required:** Yes
- **Request Body:** (Raw String) `"Movie Title"`
- **Response (200 OK):** Array of `MovieTimetable` objects.

---

### 3. Bookings & Reservations

#### Reserve Ticket
- **Method:** `POST`
- **Path:** `/movie/reserve`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "timetable_id": 1,
    "screen_id": 1,
    "seats": ["A1", "A2"],
    "date_time": "2026-05-20 18:00:00"
  }
  ```
- **Note:** `date_time` can be in RFC3339 format or `YYYY-MM-DD HH:MM:SS`.

#### Cancel Reservation
- **Method:** `DELETE`
- **Path:** `/movie/cancel`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "booking_id": 123
  }
  ```

#### Get All Bookings (Admin Only)
- **Method:** `GET`
- **Path:** `/timetable/all/bookings`
- **Auth Required:** Yes (Admin role)
- **Response:** Array of booking details.

---

### 4. Timetable & Capacity (Admin Only)

#### Add Showtime
- **Method:** `POST`
- **Path:** `/timetable/add`
- **Auth Required:** Yes (Admin role)
- **Request Body:**
  ```json
  {
    "movie_id": 1,
    "schedule": [
      {"screen_id": 1, "timings": ["12:00", "15:00"]},
      {"screen_id": 2, "timings": ["18:00"]}
    ],
    "show_date": "2026-05-20T00:00:00Z",
    "normal_price": 10,
    "vip_price": 20
  }
  ```

#### Update Showtime
- **Method:** `PATCH`
- **Path:** `/timetable/update`
- **Auth Required:** Yes (Admin role)
- **Request Body:** Same as Add Showtime, but must include `id`.

#### Get Capacity
- **Method:** `POST`
- **Path:** `/timetable/capacity`
- **Auth Required:** Yes (Admin role)
- **Request Body:**
  ```json
  {
    "timetable_id": 1,
    "screen_id": 1,
    "date_time": "2026-05-20 18:00:00"
  }
  ```
- **Response:**
  ```json
  {
    "total": 100,
    "available": 85
  }
  ```

#### Get Movie Revenue
- **Method:** `POST`
- **Path:** `/timetable/revenue`
- **Auth Required:** Yes (Admin role)
- **Request Body:**
  ```json
  {
    "movie_id": 1
  }
  ```
- **Response:**
  ```json
  {
    "movie_id": 1,
    "revenue": 1500
  }
  ```

---

### 5. Screens & Seat Status

#### View Seat Status
- **Method:** `GET`
- **Path:** `/screen/view_seats`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "showtime_id": 1,
    "screen_id": 1,
    "show_time": "18:00"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "A1": "available",
    "A2": "booked"
  }
  ```

---

## Data Models

### Movie
```json
{
  "id": 1,
  "title": "String",
  "description": "String",
  "poster_url": "String",
  "genre": ["String"]
}
```

### MovieTimetable
```json
{
  "id": 1,
  "movie_id": 1,
  "schedule": [
    {"screen_id": 1, "timings": ["HH:MM", "HH:MM:SS"]}
  ],
  "show_date": "ISO8601 Date",
  "normal_price": 10,
  "vip_price": 20
}
```

### Booking
```json
{
  "id": 1,
  "user_id": 1,
  "timetable_id": 1,
  "screen_id": 1,
  "reservation": ["A1", "A2"],
  "date_time": "String"
}
```

### Screen
```json
{
  "id": 1,
  "auditorium_number": 1,
  "normal_seats": ["A1", "A2"],
  "vip_seats": ["V1", "V2"],
  "type": "IMAX"
}
```

### User
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "admin": false
}
```
