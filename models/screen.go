package models

type SeatStatus map[string]string

type Screen struct {
	ID 		int 		`json:"id"`
	ScreenNo int		`json:"auditorium_number"`
	Normal	[]string	`json:"normal_seats"`
	Vip		[]string	`json:"vip_seats"`
	Type	string		`json:"type"`
}

type ShowSeat struct {
	ID int 	  `json:"id"`
	ShowtimeID int `json:"showtime_id"`
	ScreenID int `json:"screen_id"`
	ShowTime string `json:"show_time"`
	SeatStatus SeatStatus `json:"seat_status"`
}