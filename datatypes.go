package main

type dtConfirm struct {
	Email string            `json:"email,omitempty"`
	Forms map[string]string `json:"forms,omitempty"`
}

type dtUser struct {
	PublicID string     `json:"publicID"`
	Email    string     `json:"email"`
	Forms    []string   `json:"forms,omitempty"`
	Confirm  *dtConfirm `json:"confirm,omitempty"`
}
