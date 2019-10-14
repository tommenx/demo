package types

type GetDataRequest struct {
}

type Data struct {
	JobName string   `json:"job_name"`
	Values  []string `json:"values"`
}

type GetDataResponse struct {
	Data []Data `json:"data"`
	//Data   map[string][]string `json:"data"`
	Status string `json:"status"`
	Code   int32  `json:"code"`
}
