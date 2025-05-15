package models

type BaseResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type CheckServerAccessResponse struct {
	BaseResponse
	HasAccess bool `json:"has_access"`
}

type GetSubscriptionsResponse struct {
	BaseResponse
	Subscriptions []SubscriptionSummary `json:"subscriptions"`
}

type GetCurrentUserResponse struct {
	BaseResponse
	User *UserInfo `json:"user"`
}

// SetUserNotesRequest represents the request to set notes for a user
type SetUserNotesRequest struct {
	UserID int    `json:"user_id"`
	Notes  string `json:"notes"`
}

// SetUserNotesResponse represents the response to setting notes for a user
type SetUserNotesResponse struct {
	BaseResponse
}
