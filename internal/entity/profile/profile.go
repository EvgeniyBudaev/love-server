package profile

import (
	"github.com/EvgeniyBudaev/love-server/internal/entity/pagination"
	"time"
)

type Gender string

const (
	Male   Gender = "man"
	Female Gender = "woman"
)

type Profile struct {
	ID          uint64              `json:"id"`
	DisplayName string              `json:"displayName"`
	Birthday    time.Time           `json:"birthday"`
	Gender      Gender              `json:"gender"`
	Location    string              `json:"location"`
	Description string              `json:"description"`
	IsDeleted   bool                `json:"isDeleted"`
	IsBlocked   bool                `json:"isBlocked"`
	IsPremium   bool                `json:"isPremium"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	LastOnline  time.Time           `json:"lastOnline"`
	Images      []*ImageProfile     `json:"images"`
	Complaints  []*ComplaintProfile `json:"complaints"`
	Telegram    *TelegramProfile    `json:"telegram"`
}

type RequestAddProfile struct {
	TelegramID      string    `json:"telegramId"`
	UserName        string    `json:"username"`
	Firstname       string    `json:"firstName"`
	Lastname        string    `json:"lastName"`
	LanguageCode    string    `json:"languageCode"`
	AllowsWriteToPm string    `json:"allowsWriteToPm"`
	QueryID         string    `json:"queryId"`
	DisplayName     string    `json:"displayName"`
	Birthday        time.Time `json:"birthday"`
	Gender          Gender    `json:"gender"`
	Location        string    `json:"location"`
	Description     string    `json:"description"`
	Image           []byte    `json:"image"`
}

type ContentListProfile struct {
	ID         uint64                `json:"id"`
	LastOnline time.Time             `json:"lastOnline"`
	Image      *ResponseImageProfile `json:"image"`
}

type ResponseListProfile struct {
	*pagination.Pagination
	Content []*ContentListProfile `json:"content"`
}

type ComplaintProfile struct {
	ID        uint64 `json:"id"`
	ProfileID uint64 `json:"profileId"`
	Reason    string `json:"reason"`
}

type QueryParamsProfileList struct {
	pagination.Pagination
}

type TelegramProfile struct {
	ID              uint64 `json:"id"`
	ProfileID       uint64 `json:"profileId"`
	TelegramID      uint64 `json:"telegramId"`
	UserName        string `json:"username"`
	Firstname       string `json:"firstName"`
	Lastname        string `json:"lastName"`
	LanguageCode    string `json:"languageCode"`
	AllowsWriteToPm bool   `json:"allowsWriteToPm"`
	QueryID         string `json:"queryId"`
}

type ImageProfile struct {
	ID        uint64    `json:"id"`
	ProfileID uint64    `json:"profileId"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	IsDeleted bool      `json:"isDeleted"`
	IsBlocked bool      `json:"isBlocked"`
	IsPrimary bool      `json:"isPrimary"`
	IsPrivate bool      `json:"isPrivate"`
}

type ResponseImageProfile struct {
	Url string `json:"url"`
}
