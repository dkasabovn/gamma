package core

type ApiResponse struct {
	Data  map[string]interface{} `json:"data"`
	Error int                    `json:"error_code"`
}

// ApiConverter converts a map to a response object
func ApiConverter(data map[string]interface{}, errorCode int) *ApiResponse {
	return &ApiResponse{
		Data:  data,
		Error: errorCode,
	}
}
