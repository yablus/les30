package models

var IDs int

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Friends []int  `json:"friends"`
}

type ReqCreate struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Friends []int  `json:"friends"`
}

type ReqMakeFriends struct {
	Source_id int `json:"source_id"`
	Target_id int `json:"target_id"`
}

type ReqDelete struct {
	Target_id int `json:"target_id"`
}

type ReqUpdate struct {
	NewAge int `json:"new age"`
}
