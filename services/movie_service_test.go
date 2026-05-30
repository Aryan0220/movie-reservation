package services

import (
	"reflect"
	"testing"

	"booking-system/models"
	"booking-system/testutils"

	"github.com/pashagolub/pgxmock/v2"
)

func TestAddMovie(t *testing.T) {
	mock := testutils.NewMockDB(t)
	movie := models.Movie{Title: "Dune", Description: "Sci-fi", PosterURL: "url", Genre: []string{"Sci-Fi"}}

	mock.ExpectExec("INSERT INTO movies").
		WithArgs(movie.Title, movie.Description, movie.PosterURL, movie.Genre).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	if err := Add_Movie(movie); err != nil {
		t.Fatalf("Add_Movie error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUpdateMovie(t *testing.T) {
	mock := testutils.NewMockDB(t)
	movie := models.Movie{Title: "Dune", Description: "Updated", PosterURL: "url2", Genre: []string{"Sci-Fi"}}

	mock.ExpectExec("UPDATE movies").
		WithArgs(movie.Title, movie.Description, movie.PosterURL, movie.Genre, movie.Title).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	if err := Update_Movie(movie); err != nil {
		t.Fatalf("Update_Movie error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDeleteMovie(t *testing.T) {
	mock := testutils.NewMockDB(t)
	movie := models.Movie{Title: "Dune"}

	mock.ExpectExec("DELETE FROM movies").
		WithArgs(movie.Title).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	if err := Delete_Movie(movie); err != nil {
		t.Fatalf("Delete_Movie error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetMovies(t *testing.T) {
	mock := testutils.NewMockDB(t)

	rows := pgxmock.NewRows([]string{"title"}).
		AddRow("Dune").
		AddRow("Arrival")
	mock.ExpectQuery("SELECT title FROM MOVIES").WillReturnRows(rows)

	movies, err := Get_Movies()
	if err != nil {
		t.Fatalf("Get_Movies error: %v", err)
	}
	if !reflect.DeepEqual(movies, []string{"Dune", "Arrival"}) {
		t.Fatalf("unexpected movies: %+v", movies)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestViewSeats(t *testing.T) {
	mock := testutils.NewMockDB(t)
	request := SeatStatusRequest{ShowTimeID: 1, ScreenID: 2, ShowTime: "15:30"}
	status := map[string]string{"A1": "available", "A2": "booked"}

	rows := pgxmock.NewRows([]string{"seat_status"}).AddRow(status)
	mock.ExpectQuery("SELECT seat_status FROM show_seats").
		WithArgs(request.ShowTimeID, request.ScreenID, request.ShowTime).
		WillReturnRows(rows)

	result := ViewSeats(request)
	if !reflect.DeepEqual(result, status) {
		t.Fatalf("unexpected seat status: %+v", result)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
