package models

import "time"

type Bookings struct {
	ID          int      `json:"id"`
	UserID      int      `json:"user_id"`
	TimeTableID int      `json:"timetable_id"`
	ScreenID    int      `json:"screen_id"`
	Reservation []string `json:"reservation"`
	DateTime    string   `json:"date_time"`
}

type BookingDetail struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	TimetableID int       `json:"timetable_id"`
	ScreenID    int       `json:"screen_id"`
	Reservation []string  `json:"reservation"`
	DateTime    time.Time `json:"date_time"`
}
