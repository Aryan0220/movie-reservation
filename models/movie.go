package models

type Movie struct {
	ID	int				`json:"id"`
	Title	string 		`json:"title"`
	Description	string  `json:"description"`
	PosterURL	string	`json:"poster_url"`
	Genre	[]string	`json:"genre"`
}

type MovieTimetable struct {
	ID	int		`json:"id"`
	MovieID	int 	`json:"movie_id"`
	Timings	[]string	`json:"timings"`
	Screens []int		`json:"screens"`
	ShowDate	string		`json:"show_date"`
	NormalPrice	int		`json:"normal_price"`
	VipPrice	int		`json:"vip_price"`
}


