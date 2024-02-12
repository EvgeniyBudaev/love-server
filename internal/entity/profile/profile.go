package profile

import (
	"github.com/EvgeniyBudaev/love-server/internal/entity/pagination"
	"time"
)

type Profile struct {
	ID             uint64              `json:"id"`
	DisplayName    string              `json:"displayName"`
	Birthday       time.Time           `json:"birthday"`
	Gender         string              `json:"gender"`
	SearchGender   string              `json:"searchGender"`
	Location       string              `json:"location"`
	Description    string              `json:"description"`
	Height         uint8               `json:"height"`
	Weight         uint8               `json:"weight"`
	LookingFor     string              `json:"lookingFor"`
	IsDeleted      bool                `json:"isDeleted"`
	IsBlocked      bool                `json:"isBlocked"`
	IsPremium      bool                `json:"isPremium"`
	IsShowDistance bool                `json:"isShowDistance"`
	IsInvisible    bool                `json:"isInvisible"`
	CreatedAt      time.Time           `json:"createdAt"`
	UpdatedAt      time.Time           `json:"updatedAt"`
	LastOnline     time.Time           `json:"lastOnline"`
	Images         []*ImageProfile     `json:"images"`
	Complaints     []*ComplaintProfile `json:"complaints"`
	Telegram       *TelegramProfile    `json:"telegram"`
	Navigator      *NavigatorProfile   `json:"navigator"`
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
	Gender          string    `json:"gender"`
	SearchGender    string    `json:"searchGender"`
	Location        string    `json:"location"`
	Description     string    `json:"description"`
	Height          string    `json:"height"`
	Weight          string    `json:"weight"`
	LookingFor      string    `json:"lookingFor"`
	Latitude        string    `json:"latitude"`
	Longitude       string    `json:"longitude"`
	Image           []byte    `json:"image"`
}

type RequestUpdateProfile struct {
	TelegramID      string    `json:"telegramId"`
	UserName        string    `json:"username"`
	Firstname       string    `json:"firstName"`
	Lastname        string    `json:"lastName"`
	LanguageCode    string    `json:"languageCode"`
	AllowsWriteToPm string    `json:"allowsWriteToPm"`
	QueryID         string    `json:"queryId"`
	ID              string    `json:"id"`
	DisplayName     string    `json:"displayName"`
	Birthday        time.Time `json:"birthday"`
	Gender          string    `json:"gender"`
	SearchGender    string    `json:"searchGender"`
	Location        string    `json:"location"`
	Description     string    `json:"description"`
	Height          string    `json:"height"`
	Weight          string    `json:"weight"`
	LookingFor      string    `json:"lookingFor"`
	Latitude        string    `json:"latitude"`
	Longitude       string    `json:"longitude"`
	Image           []byte    `json:"image"`
}

type RequestDeleteProfile struct {
	ID string `json:"id"`
}

type RequestDeleteProfileImage struct {
	ID string `json:"id"`
}

type ContentListProfile struct {
	ID         uint64                    `json:"id"`
	LastOnline time.Time                 `json:"lastOnline"`
	Image      *ResponseImageProfile     `json:"image"`
	Navigator  *ResponseNavigatorProfile `json:"navigator"`
}

type ResponseListProfile struct {
	*pagination.Pagination
	Content []*ContentListProfile `json:"content"`
}

type ResponseProfile struct {
	ID           uint64                   `json:"id"`
	SearchGender string                   `json:"searchGender"`
	Image        *ResponseImageProfile    `json:"image"`
	Telegram     *ResponseTelegramProfile `json:"telegram"`
}

type ComplaintProfile struct {
	ID        uint64 `json:"id"`
	ProfileID uint64 `json:"profileId"`
	Reason    string `json:"reason"`
}

type QueryParamsProfileList struct {
	pagination.Pagination
	ProfileID    string `json:"profileId"`
	AgeFrom      string `json:"ageFrom"`
	AgeTo        string `json:"ageTo"`
	SearchGender string `json:"searchGender"`
}

type QueryParamsGetProfileByTelegramID struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
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

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type NavigatorProfile struct {
	ID        uint64 `json:"id"`
	ProfileID uint64 `json:"profileId"`
	Location  *Point `json:"location"`
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

type ResponseTelegramProfile struct {
	TelegramID uint64 `json:"telegramId"`
}

type ResponseNavigatorProfile struct {
	Distance string `json:"distance"`
}
