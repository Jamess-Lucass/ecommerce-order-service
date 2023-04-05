package requests

type Email struct {
	To      []string `json:"to"`
	From    string   `json:"from"`
	CC      []string `json:"cc,omitempty"`
	BCC     []string `json:"bcc,omitempty"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}
