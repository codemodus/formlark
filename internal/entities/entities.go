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
	DeletedAt NullTime
	BlockedAt NullTime

	Email string

	ConfirmedAt NullTime
	Token       string
}

// UserRequiz ...
type UserRequiz struct {
	URL  string
	User User
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

	UserID uint64

	ReplyTo string
	Subject string
	Form    map[string]string
}

// MessageRecord ...
type MessageRecord struct {
	UserID uint64
	Form   map[string]string
}
