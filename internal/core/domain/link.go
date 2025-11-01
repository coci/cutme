package domain

type Link struct {
	Code      string `json:"code,omitempty"`
	Link      string `json:"link"`
	CreatedAt int64  `json:"created_at,omitempty"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
}
