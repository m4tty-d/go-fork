package types

type Client struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type Server struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
