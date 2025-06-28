package models

type Answer struct {
	Text    string `json:"text"`
	Correct bool   `json:"correct"`
	Reason  string `json:"reason"`
}

type Question struct {
	Question string   `json:"question"`
	Answers  []Answer `json:"answers"`
}
