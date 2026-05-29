# movie-reservation project explanation

This document explains the project structure, major components, data flow, and function-level behavior. It also includes diagrams for the main request flows.

## 1. Project overview

This is a Go REST API for movie reservations using Fiber, PostgreSQL (via pgxpool), and JWT auth. The core flow is:

- `main.go` loads environment variables, connects to DB, starts the Fiber app.
- `routes/` defines API endpoints and connects them to `handlers/`.
- `middleware/` enforces JWT auth on protected routes.
- `handlers/` parse requests, validate input, and call `services/`.
- `services/` contain business logic and all DB access via `config.DB`.
- `models/` define request/DB payload structures.

## 2. Entry point and configuration

- `main.go`
  - `config.LoadEnv()` loads `.env` (if present).
  - `config.ConnectDB()` creates a pgx pool using `DB_URL`.
  - `routes.Setup(app)` registers all endpoints.
  - `app.Listen(" :PORT ")` starts the server.

- `config/env.go`
  - `LoadEnv()` loads `.env` with `godotenv`.
  - `GetEnv(key)` reads environment variables.

- `config/db.go`
  - `ConnectDB()` initializes `config.DB` via `pgxpool.New`.

## 3. Routing and middleware

- `routes/routes.go` defines paths under `/api`.
  - Public: `/register`, `/login`.
  - Protected groups:
    - `/movie/*` for movie and booking operations.
    - `/timetable/*` for showtime and admin reporting.
    - `/user/*` for user operations.

- `middleware/auth.go`
  - `Protected` expects `Authorization: Bearer <token>`.
  - Verifies `JWT_SECRET` is set and token is valid.
  - Extracts `user_id` and `role` into `c.Locals`.

## 4. Data models

- `models/User`: `id`, `name`, `email`, `role`, `password`.
- `models/Movie`: `id`, `title`, `description`, `poster_url`, `genre`.
- `models/MovieTimetable`: `id`, `movie_id`, `timings`, `screens`, `show_date`, `normal_price`, `vip_price`.
- `models/Screen`: `id`, `screen_no`, `normal`, `vip`, `type`.
- `models/Bookings` and `models/BookingDetail`: DB booking payloads.

## 5. High-level request flow

```mermaid
flowchart TD
    Client -->|HTTP request| Fiber
    Fiber -->|Route match| Handler
    Handler -->|Parse/validate| Service
    Service -->|DB query/exec| DB[(PostgreSQL)]
    Service -->|Return result| Handler
    Handler -->|JSON response| Client

    Client -->|Protected route| Middleware
    Middleware -->|JWT verify| Handler
```

## 6. Auth and user flows

### 6.1 Register

Handler: `handlers.Register` -> Service: `services.CreateUser`

```mermaid
flowchart TD
    A[POST /api/register] --> B[handlers.Register]
    B --> C[services.CreateUser]
    C --> D[INSERT users]
    D --> B
    B --> E[Response: User created]
```

Key logic:
- `Register` parses `models.User` from body.
- `CreateUser` hashes password with bcrypt and inserts into `users`.

### 6.2 Login

Handler: `handlers.Login` -> Services: `services.GetUserByEmail`, `services.GenerateToken`

```mermaid
flowchart TD
    A[POST /api/login] --> B[handlers.Login]
    B --> C[services.GetUserByEmail]
    C --> D[SELECT users by email]
    D --> B
    B --> E[bcrypt.CompareHashAndPassword]
    E --> F[services.GenerateToken]
    F --> G[JWT signed with JWT_SECRET]
    G --> H[Response: token]
```

Key logic:
- `GetUserByEmail` loads hashed password.
- `Login` verifies password and returns JWT.

### 6.3 Promote user

Handler: `handlers.Promote` -> Service: `services.PromoteToAdmin`

```mermaid
flowchart TD
    A[PATCH /api/user/promote] --> B[handlers.Promote]
    B --> C[services.PromoteToAdmin]
    C --> D[UPDATE users SET role='admin']
    D --> E[Response: User Promoted]
```

## 7. Movie management flows

### 7.1 Add movie (admin)

Handler: `handlers.AddMovie` -> Service: `services.Add_Movie`

```mermaid
flowchart TD
    A[POST /api/movie/add] --> B[handlers.AddMovie]
    B -->|Role check| C{Role == admin}
    C -->|yes| D[services.Add_Movie]
    D --> E[INSERT movies]
    E --> F[Response: Movie added]
    C -->|no| G[403]
```

### 7.2 Update movie (admin)

Handler: `handlers.UpdateMovie` -> Service: `services.Update_Movie`

```mermaid
flowchart TD
    A[PATCH /api/movie/update] --> B[handlers.UpdateMovie]
    B --> C{Role == admin}
    C -->|yes| D[services.Update_Movie]
    D --> E[UPDATE movies WHERE title]
    E --> F[Response: Movie Updated]
    C -->|no| G[403]
```

### 7.3 Delete movie (admin)

Handler: `handlers.DeleteMovie` -> Service: `services.Delete_Movie`

```mermaid
flowchart TD
    A[DELETE /api/movie/delete] --> B[handlers.DeleteMovie]
    B --> C{Role == admin}
    C -->|yes| D[services.Delete_Movie]
    D --> E[DELETE FROM movies WHERE title]
    E --> F[Response: Movie Deleted]
    C -->|no| G[403]
```

### 7.4 Get movies

Handler: `handlers.GetMovies` -> Service: `services.Get_Movies`

```mermaid
flowchart TD
    A[GET /api/movie/get] --> B[handlers.GetMovies]
    B --> C[services.Get_Movies]
    C --> D[SELECT title FROM movies]
    D --> E[Response: list of titles]
```

### 7.5 Get movie timings

Handler: `handlers.GetMovieTimings` -> Service: `services.GetMovieTimings`

```mermaid
flowchart TD
    A[GET /api/movie/timings] --> B[handlers.GetMovieTimings]
    B --> C[services.GetMovieTimings]
    C --> D[SELECT timetable by movie title]
    D --> E[Response: timetable]
```

## 8. Timetable management flows

### 8.1 Add showtime (admin)

Handler: `handlers.AddShowTime` -> Service: `services.AddShowTime`

```mermaid
flowchart TD
    A[POST /api/timetable/add] --> B[handlers.AddShowTime]
    B --> C{Role == admin}
    C -->|yes| D[validateTimetableInput]
    D --> E[services.AddShowTime]
    E --> F[INSERT movie_timetables]
    F --> G[Response: Showtime added]
    C -->|no| H[403]
```

Validation checks in `validateTimetableInput`:
- Required IDs, timings, screens, and prices.
- Show date cannot be in the past.
- No duplicate timings or screens.

### 8.2 Update showtime (admin)

Handler: `handlers.UpdateShowTime` -> Service: `services.UpdateShowTime`

```mermaid
flowchart TD
    A[PATCH /api/timetable/update] --> B[handlers.UpdateShowTime]
    B --> C{Role == admin}
    C -->|yes| D[validateTimetableInput]
    D --> E[services.UpdateShowTime]
    E --> F[UPDATE movie_timetables]
    F --> G[Response: Showtime updated]
    C -->|no| H[403]
```

## 9. Booking and capacity flows

### 9.1 Reserve movie

Handler: `handlers.ReserveMovie` -> Service: `services.ReserveTicket`

```mermaid
flowchart TD
    A[POST /api/movie/reserve] --> B[handlers.ReserveMovie]
    B --> C[requireUserID]
    C --> D[parseDateTime]
    D --> E[services.ReserveTicket]

    E --> F[getTimetableByID]
    F --> G[Validate date and timing]
    G --> H[getScreenByID]
    H --> I[normalizeSeats]
    I --> J[validateSeatsExist]
    J --> K[getBookedSeats]
    K --> L{Any seat already booked?}
    L -->|yes| M[ErrSeatUnavailable]
    L -->|no| N[INSERT bookings]
    N --> O[Response: Reservation created]
```

Key checks:
- Reservation time must be in the future.
- Timetable date and time must match requested show.
- Seats must exist in the selected screen.
- Duplicate seats in the request are rejected.
- Already booked seats return conflict.

### 9.2 Cancel reservation

Handler: `handlers.CancelReservation` -> Service: `services.CancelReservation`

```mermaid
flowchart TD
    A[DELETE /api/movie/cancel] --> B[handlers.CancelReservation]
    B --> C[requireUserID]
    C --> D[services.CancelReservation]
    D --> E[SELECT booking by id]
    E --> F{Owner matches user?}
    F -->|no| G[ErrNotOwner]
    F -->|yes| H{Showtime in future?}
    H -->|no| I[ErrPastReservation]
    H -->|yes| J[DELETE booking]
    J --> K[Response: Reservation canceled]
```

### 9.3 Get capacity (admin)

Handler: `handlers.GetCapacity` -> Service: `services.GetCapacity`

```mermaid
flowchart TD
    A[POST /api/timetable/capacity] --> B[handlers.GetCapacity]
    B --> C{Role == admin}
    C -->|yes| D[parseDateTime]
    D --> E[services.GetCapacity]
    E --> F[getTimetableByID]
    F --> G[getScreenByID]
    G --> H[getBookedSeats]
    H --> I[Compute total vs available]
    I --> J[Response: total, available]
    C -->|no| K[403]
```

### 9.4 Get revenue (admin)

Handler: `handlers.GetRevenue` -> Service: `services.GetMovieRevenue`

```mermaid
flowchart TD
    A[POST /api/timetable/revenue] --> B[handlers.GetRevenue]
    B --> C{Role == admin}
    C -->|yes| D[services.GetMovieRevenue]
    D --> E[Query bookings + timetables]
    E --> F[getScreenByID (cached)]
    F --> G[Count normal/vip seats]
    G --> H[Sum revenue]
    H --> I[Response: revenue]
    C -->|no| J[403]
```

### 9.5 Get all bookings (admin)

Handler: `handlers.GetAllReservations` -> Service: `services.GetAllBookings`

```mermaid
flowchart TD
    A[GET /api/timetable/all/bookings] --> B[handlers.GetAllReservations]
    B --> C{Role == admin}
    C -->|yes| D[services.GetAllBookings]
    D --> E[SELECT bookings ORDER BY date_time]
    E --> F[Response: booking list]
    C -->|no| G[403]
```

## 10. Service helpers and error flow

Key errors from `services/booking_service.go`:
- `ErrSeatUnavailable` -> HTTP 409
- `ErrInvalidShowtime` or `ErrPastReservation` -> HTTP 400
- `ErrNotFound` -> HTTP 404
- `ErrNotOwner` -> HTTP 403

The handlers map these error values to status codes in a consistent way.

Helper functions used by booking services:
- `getTimetableByID`, `getScreenByID` (DB lookups)
- `getBookedSeats` (aggregation of reserved seats)
- `normalizeSeats`, `validateSeatsExist`
- `matchesTiming` (compares showtime to timetable slots)
- `isSameDate` (date comparison)

## 11. Database shape (inferred)

From SQL in services, the DB likely has tables:

- `users(id, name, email, role, password)`
- `movies(id, title, description, poster_url, genre)`
- `movie_timetables(id, movie_id, timings, screens, show_date, normal_price, vip_price)`
- `screens(id, screen_no, normal, vip, type)`
- `bookings(id, user_id, timetable_id, screen_id, reservation, date_and_time)`

Note: `genre`, `timings`, `screens`, `normal`, `vip`, and `reservation` are stored as array columns (Go slices in models).

## 12. Test coverage overview

- `main_test.go` and `config/env_test.go` validate `.env` loading behavior.
- `middleware/auth_test.go` tests JWT error cases and valid token pass-through.
- `services/services_test.go` checks JWT generation.

## 13. Outstanding or stubbed areas

- `services/screen_service.go` contains placeholders and is not used.
- No explicit schema migrations or DB setup scripts are included.

## 14. End-to-end flow summary

```mermaid
flowchart TD
    Start[Client] --> Login[POST /api/login]
    Login --> Token[JWT]
    Token --> Protected[Protected endpoints]

    Protected --> MovieAdd[Add/Update/Delete movie]
    Protected --> Timetable[Add/Update showtime]
    Protected --> Reserve[Reserve seats]
    Protected --> Cancel[Cancel reservation]
    Protected --> Capacity[Get capacity]
    Protected --> Revenue[Get revenue]

    MovieAdd --> DB[(DB)]
    Timetable --> DB
    Reserve --> DB
    Cancel --> DB
    Capacity --> DB
    Revenue --> DB
```
