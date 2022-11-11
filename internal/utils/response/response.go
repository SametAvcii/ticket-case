package response

type ResponseData struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
}

func Response(Status int, Data interface{}) ResponseData {
	return ResponseData{
		Data:   Data,
		Status: Status,
	}
}
