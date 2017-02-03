package entities

import "time"

// Empty ...
type Empty struct{}

// NullTime ...
type NullTime struct {
	Valid bool
	Time  time.Time
}

// User ...
type User struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt NullTime `json:"omitempty"`
	BlockedAt NullTime `json:"omitempty"`

	Email string

	ConfirmedAt NullTime `json:"omitempty"`
	Token       string   `json:"omitempty"`
}

// UserRecord ...
type UserRecord struct {
	Email string
}

// UserRequiz ...
type UserRequiz struct {
	URL  string
	User UserRecord
}

// UserReferral ...
type UserReferral struct {
	Email string
	Token string
}

// Message ...
type Message struct {
	ID        uint64
	CreatedAt time.Time

	UserID uint64 `json:"omitempty"`

	ReplyTo string            `json:"omitempty"`
	Subject string            `json:"omitempty"`
	Form    map[string]string `json:"omitempty"`
}

// MessageRecord ...
type MessageRecord struct {
	Form map[string]string
}

// MessageByUserIDRecord ...
type MessageByUserIDRecord struct {
	UserID  uint64
	Message MessageRecord
}

// MarshalJSON ...
func (n *NullTime) MarshalJSON() ([]byte, error) {
	return []byte(n.Time.Format(`"2006-01-02T15:04:05Z"`)), nil
}

// UnmarshalJSON ...
func (n *NullTime) UnmarshalJSON(data []byte) (err error) {
	t, err := time.Parse("\"2006-01-02T15:04:05Z\"", string(data))
	if err != nil {
		return err
	}

	n.Time = t
	n.Valid = true

	return nil
}

// IsZero ...
func (n *NullTime) IsZero() bool {
	return !n.Valid
}
