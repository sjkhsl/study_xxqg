package lib

type T struct {
	Question    string      `json:"question"`
	Answer      string      `json:"answer"`
	WrongAnswer string      `json:"wrongAnswer"`
	Option      string      `json:"option"`
	Datetime    interface{} `json:"datetime"`
}
