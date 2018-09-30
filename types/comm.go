package types

type Request struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type Response struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
