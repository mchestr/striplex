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
