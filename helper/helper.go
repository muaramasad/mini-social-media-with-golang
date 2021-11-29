package helper

type ResponseData struct {
	Status int `json: "status"`
	Message string `json: "message"`
	Data interface{} `json: "data"`
}