package models

type Claims struct {
	Audience interface{} `json:"aud"`
	Azp      string      `json:"azp"`
}
