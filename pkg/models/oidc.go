package models

var Claims struct {
	Audience interface{} `json:"aud"`
	Azp      string      `json:"azp"`
}
