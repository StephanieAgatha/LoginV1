package response

type ResponseApi struct {
	StatusCode int         `json:"status_code"`
	Messages   string      `json:"messages"`
	Payload    interface{} `json:"payload,omitempty"`
}
