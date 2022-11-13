package response

type ResponseData struct {
	Ticket interface{} `json:"ticket"`
	Status int         `json:"status"`
}

type ResponseOnlyStatus struct {
	Status int `json:"status"`
}

func Response(Status int, Data interface{}) ResponseData {
	return ResponseData{
		Ticket: Data,
		Status: Status,
	}
}

func ResponseStatus(Status int) ResponseOnlyStatus {
	return ResponseOnlyStatus{
		Status: Status,
	}
}
