# movie-reservation

API documentation with examples for each endpoint.

## Base URL

`http://localhost:<PORT>/api`

`<PORT>` comes from the `PORT` environment variable.

## Authentication

Protected routes require a JWT in the `Authorization` header:

```
Authorization: Bearer <token>
```

The token is returned by `POST /api/login` and is valid for 72 hours.

## Date and time format

`date_time` accepts either:

- RFC3339, e.g. `2026-05-17T19:30:00Z`
- Local time layout `YYYY-MM-DD HH:MM:SS`, e.g. `2026-05-17 19:30:00`

## Endpoints

### Auth

#### POST /api/register

Creates a user.

Request body:

```json
{
	"name": "Jane Doe",
	"email": "jane@example.com",
	"password": "secret123",
	"role": "user"
}
```

Response:

```json
{
	"message": "User created"
}
```

#### POST /api/login

Authenticates a user and returns a token.

Request body:

```json
{
	"email": "jane@example.com",
	"password": "secret123"
}
```

Response:

```json
{
	"token": "<jwt>"
}
```

### Movies (protected)

All routes in this section require `Authorization: Bearer <token>`.

#### POST /api/movie/add (admin)

Adds a movie.

Request body:

```json
{
	"title": "Inception",
	"description": "Dream heist thriller",
	"poster_url": "https://example.com/poster.jpg",
	"genre": ["Sci-Fi", "Thriller"]
}
```

Response:

```json
{
	"message": "Movie added successfully"
}
```

#### PATCH /api/movie/update (admin)

Updates a movie (matched by `title`).

Request body:

```json
{
	"title": "Inception",
	"description": "Updated description",
	"poster_url": "https://example.com/poster2.jpg",
	"genre": ["Sci-Fi", "Thriller"]
}
```

Response:

```json
{
	"message": "Movie Updated"
}
```

#### DELETE /api/movie/delete (admin)

Deletes a movie (matched by `title`).

Request body:

```json
{
	"title": "Inception"
}
```

Response:

```json
{
	"message": "Movie Deleted"
}
```

#### GET /api/movie/get

Returns a list of movie titles.

Response:

```json
[
	"Inception",
	"Interstellar"
]
```

#### GET /api/movie/timings

Returns the timetable for a movie title. Note: this endpoint expects a JSON string as the request body.

Request body:

```json
"Inception"
```

Response:

```json
{
	"id": 3,
	"movie_id": 12,
	"timings": ["10:00", "14:00", "19:30"],
	"screens": [1, 2],
	"show_date": "2026-05-17T00:00:00Z",
	"normal_price": 10,
	"vip_price": 15
}
```

#### POST /api/movie/reserve

Reserves seats for a showtime.

Request body:

```json
{
	"timetable_id": 3,
	"screen_id": 1,
	"seats": ["A1", "A2"],
	"date_time": "2026-05-17T19:30:00Z"
}
```

Response:

```json
{
	"message": "Reservation created"
}
```

#### DELETE /api/movie/cancel

Cancels an existing reservation (must be owner).

Request body:

```json
{
	"booking_id": 42
}
```

Response:

```json
{
	"message": "Reservation canceled"
}
```

### Timetable (protected)

All routes in this section require `Authorization: Bearer <token>`.

#### POST /api/timetable/add (admin)

Adds showtimes for a movie.

Request body:

```json
{
	"movie_id": 12,
	"timings": ["10:00", "14:00", "19:30"],
	"screens": [1, 2],
	"show_date": "2026-05-17T00:00:00Z",
	"normal_price": 10,
	"vip_price": 15
}
```

Response:

```json
{
	"message": "Showtime added successfully"
}
```

#### PATCH /api/timetable/update (admin)

Updates an existing showtime.

Request body:

```json
{
	"id": 3,
	"movie_id": 12,
	"timings": ["11:00", "15:00"],
	"screens": [1, 3],
	"show_date": "2026-05-17T00:00:00Z",
	"normal_price": 12,
	"vip_price": 18
}
```

Response:

```json
{
	"message": "Showtime updated successfully"
}
```

#### POST /api/timetable/capacity (admin)

Returns total and available seats for a showtime.

Request body:

```json
{
	"timetable_id": 3,
	"screen_id": 1,
	"date_time": "2026-05-17 19:30:00"
}
```

Response:

```json
{
	"total": 120,
	"available": 98
}
```

#### POST /api/timetable/revenue (admin)

Returns total revenue for a movie up to now.

Request body:

```json
{
	"movie_id": 12
}
```

Response:

```json
{
	"movie_id": 12,
	"revenue": 1840
}
```

#### GET /api/timetable/all/bookings (admin)

Returns all bookings (latest first).

Response:

```json
[
	{
		"id": 42,
		"user_id": 7,
		"timetable_id": 3,
		"screen_id": 1,
		"reservation": ["A1", "A2"],
		"date_time": "2026-05-17T19:30:00Z"
	}
]
```

### User (protected)

All routes in this section require `Authorization: Bearer <token>`.

#### PATCH /api/user/promote

Promotes a user to admin.

Request body:

```json
{
	"email": "jane@example.com"
}
```

Response:

```json
{
	"message": "User Promoted to Admin"
}
```

#### GET /api/user/me

Simple auth check.

Response:

```json
{
	"message": "Protected route"
}
```

## Example curl usage

```bash
curl -X POST http://localhost:8080/api/login \
	-H "Content-Type: application/json" \
	-d '{"email":"jane@example.com","password":"secret123"}'
```



docker run -it --name movie-database -p 5432:5432 -v postgres_data:/var/lib/postgresql/data -e POSTGRES_DB=moviedb -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=zoro movie-postgres

docker build -t movie-postgres -f ./db/Dockerfile db 

docker volume rm postgres_data    
docker run -it --name movie-database -p 5432:5432 -v postgres_data:/var/lib/postgresql -e POSTGRES_DB=moviedb -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=zoro movie-postgres
CMD ["psql", "-U", "postgres", "-d", "moviedb"]

- name: Login to Docker Hub
  uses: docker/login-action@v3
  with:
    username: ${{ secrets.DOCKER_USERNAME }}
    password: ${{ secrets.DOCKER_PASSWORD }}