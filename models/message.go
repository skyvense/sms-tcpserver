package models

// Message represents the structure of the incoming JSON data
type Message struct {
	Txt   string `json:"txt"`
	Num   string `json:"num"`
	Cmd   string `json:"cmd"`
	Metas struct {
		Tz     int `json:"tz"`
		Min    int `json:"min"`
		SeqNum int `json:"seqNum"`
		RefNum int `json:"refNum"`
		Year   int `json:"year"`
		Sec    int `json:"sec"`
		MaxNum int `json:"maxNum"`
		Mon    int `json:"mon"`
		Hour   int `json:"hour"`
		Day    int `json:"day"`
	} `json:"metas"`
}
