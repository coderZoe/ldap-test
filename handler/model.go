package handler

type User struct {
	CN              string `json:"cn"`
	SN              string `json:"sn"`
	Mail            string `json:"mail"`
	TelephoneNumber string `json:"telephoneNumber"`
}
