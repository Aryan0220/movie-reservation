package models

type Screen struct {
	ID 		int 		`json:"id"`
	ScreenNo int		`json:"screen_no"`
	Normal	[]string	`json:"normal"`
	Vip		[]string	`json:"vip"`
	Type	string		`json:"type"`
}
