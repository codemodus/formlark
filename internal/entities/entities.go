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
	ID        uint64    `json:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updateAt,omitempty"`
	DeletedAt NullTime  `json:"deletedAt,omitempty"`
	BlockedAt NullTime  `json:"blockedAt,omitempty"`

	Email string `json:"email,omitempty"`

	ConfirmedAt NullTime `json:"confirmedAt,omitempty"`
}

// UserRecord ...
type UserRecord struct {
	Email string `json:"email,omitempty"`
}

// UserRequiz ...
type UserRequiz struct {
	URL        string     `json:"url,omitempty"`
	UserRecord UserRecord `json:"user,omitempty"`
}

// UserReferral ...
type UserReferral struct {
	Email string `json:"email,omitempty"`
}

// Message ...
type Message struct {
	ID        uint64    `json:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`

	UserID uint64 `json:"userID,omitempty"`

	ReplyTo string            `json:"replyTo,omitempty"`
	Subject string            `json:"subject,omitempty"`
	Form    map[string]string `json:"form,omitempty"`
}

// MessageRecord ...
type MessageRecord struct {
	Form map[string]string `json:"form,omitempty"`
}

// MessageByUserIDRecord ...
type MessageByUserIDRecord struct {
	UserID  uint64        `json:"userID,omitempty"`
	Message MessageRecord `json:"message,omitempty"`
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
