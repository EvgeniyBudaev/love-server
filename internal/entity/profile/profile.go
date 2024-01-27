package profile

type Profile struct {
	ID          int64  `json:"id"`
	DisplayName string `json:"display_name"`
}

type CreateRequestProfile struct {
	DisplayName string `json:"displayName"`
}
