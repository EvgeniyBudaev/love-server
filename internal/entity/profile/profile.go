package profile

import (
	"github.com/EvgeniyBudaev/love-server/internal/entity/pagination"
	"time"
)

type Profile struct {
	ID             uint64                    `json:"id"`
	SessionID      string                    `json:"sessionId"`
	DisplayName    string                    `json:"displayName"`
	Birthday       time.Time                 `json:"birthday"`
	Gender         string                    `json:"gender"`
	Location       string                    `json:"location"`
	Height         uint8                     `json:"height"`
	Weight         uint8                     `json:"weight"`
	Description    string                    `json:"description"`
	IsDeleted      bool                      `json:"isDeleted"`
	IsBlocked      bool                      `json:"isBlocked"`
	IsPremium      bool                      `json:"isPremium"`
	IsShowDistance bool                      `json:"isShowDistance"`
	IsInvisible    bool                      `json:"isInvisible"`
	CreatedAt      time.Time                 `json:"createdAt"`
	UpdatedAt      time.Time                 `json:"updatedAt"`
	LastOnline     time.Time                 `json:"lastOnline"`
	Images         []*ImageProfile           `json:"images"`
	Telegram       *TelegramProfile          `json:"telegram"`
	Navigator      *ResponseNavigatorProfile `json:"navigator"`
	Filter         *FilterProfile            `json:"filters"`
}

type RequestAddProfile struct {
	SessionID        string    `json:"sessionId"`
	UserName         string    `json:"userName"`
	DisplayName      string    `json:"displayName"`
	Birthday         time.Time `json:"birthday"`
	Gender           string    `json:"gender"`
	SearchGender     string    `json:"searchGender"`
	Location         string    `json:"location"`
	Description      string    `json:"description"`
	Height           string    `json:"height"`
	Weight           string    `json:"weight"`
	LookingFor       string    `json:"lookingFor"`
	TelegramID       string    `json:"telegramId"`
	TelegramUserName string    `json:"telegramUserName"`
	Firstname        string    `json:"firstName"`
	Lastname         string    `json:"lastName"`
	LanguageCode     string    `json:"languageCode"`
	AllowsWriteToPm  string    `json:"allowsWriteToPm"`
	QueryID          string    `json:"queryId"`
	Latitude         string    `json:"latitude"`
	Longitude        string    `json:"longitude"`
	AgeFrom          string    `json:"ageFrom"`
	AgeTo            string    `json:"ageTo"`
	Distance         string    `json:"distance"`
	Page             string    `json:"page"`
	Size             string    `json:"size"`
	Image            []byte    `json:"image"`
}

type RequestUpdateProfile struct {
	ID               string    `json:"id"`
	UserName         string    `json:"userName"`
	DisplayName      string    `json:"displayName"`
	Birthday         time.Time `json:"birthday"`
	Gender           string    `json:"gender"`
	SearchGender     string    `json:"searchGender"`
	Location         string    `json:"location"`
	Description      string    `json:"description"`
	Height           string    `json:"height"`
	Weight           string    `json:"weight"`
	LookingFor       string    `json:"lookingFor"`
	TelegramID       string    `json:"telegramId"`
	TelegramUserName string    `json:"telegramUserName"`
	Firstname        string    `json:"firstName"`
	Lastname         string    `json:"lastName"`
	LanguageCode     string    `json:"languageCode"`
	AllowsWriteToPm  string    `json:"allowsWriteToPm"`
	QueryID          string    `json:"queryId"`
	Latitude         string    `json:"latitude"`
	Longitude        string    `json:"longitude"`
	AgeFrom          string    `json:"ageFrom"`
	AgeTo            string    `json:"ageTo"`
	Distance         string    `json:"distance"`
	Page             string    `json:"page"`
	Size             string    `json:"size"`
	Image            []byte    `json:"image"`
}

type RequestDeleteProfile struct {
	ID string `json:"id"`
}

type RequestDeleteProfileImage struct {
	ID string `json:"id"`
}

type ContentListProfile struct {
	ID         uint64                    `json:"id"`
	IsOnline   bool                      `json:"isOnline"`
	LastOnline time.Time                 `json:"lastOnline"`
	Image      *ResponseImageProfile     `json:"image"`
	Navigator  *ResponseNavigatorProfile `json:"navigator"`
}

type ResponseListProfile struct {
	*pagination.Pagination
	Content []*ContentListProfile `json:"content"`
}

type ResponseProfile struct {
	ID        uint64                   `json:"id"`
	SessionID string                   `json:"sessionId"`
	IsDeleted bool                     `json:"isDeleted"`
	IsBlocked bool                     `json:"isBlocked"`
	Image     *ResponseImageProfile    `json:"image"`
	Telegram  *ResponseTelegramProfile `json:"telegram"`
	Filter    *ResponseFilterProfile   `json:"filter"`
}

type ResponseProfileDetail struct {
	ID             uint64                    `json:"id"`
	SessionID      string                    `json:"sessionId"`
	DisplayName    string                    `json:"displayName"`
	Birthday       time.Time                 `json:"birthday"`
	Gender         string                    `json:"gender"`
	Location       string                    `json:"location"`
	Height         uint8                     `json:"height"`
	Weight         uint8                     `json:"weight"`
	Description    string                    `json:"description"`
	IsDeleted      bool                      `json:"isDeleted"`
	IsBlocked      bool                      `json:"isBlocked"`
	IsPremium      bool                      `json:"isPremium"`
	IsShowDistance bool                      `json:"isShowDistance"`
	IsInvisible    bool                      `json:"isInvisible"`
	IsOnline       bool                      `json:"isOnline"`
	CreatedAt      time.Time                 `json:"createdAt"`
	UpdatedAt      time.Time                 `json:"updatedAt"`
	LastOnline     time.Time                 `json:"lastOnline"`
	Images         []*ImageProfile           `json:"images"`
	Telegram       *TelegramProfile          `json:"telegram"`
	Navigator      *ResponseNavigatorProfile `json:"navigator"`
	Filter         *FilterProfile            `json:"filters"`
	Like           *ResponseLikeProfile      `json:"like"`
}

type QueryParamsProfileList struct {
	pagination.Pagination
	SessionID    string `json:"sessionId"`
	AgeFrom      string `json:"ageFrom"`
	AgeTo        string `json:"ageTo"`
	SearchGender string `json:"searchGender"`
	LookingFor   string `json:"lookingFor"`
	Distance     string `json:"distance"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
}

type QueryParamsGetProfileByTelegramID struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type QueryParamsGetProfileByUserID struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type QueryParamsGetProfileByID struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type QueryParamsGetProfileDetail struct {
	ViewerID  string `json:"viewerId"`
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

type FilterProfile struct {
	ID           uint64 `json:"id"`
	ProfileID    uint64 `json:"profileId"`
	SearchGender string `json:"searchGender"`
	LookingFor   string `json:"lookingFor"`
	AgeFrom      uint8  `json:"ageFrom"`
	AgeTo        uint8  `json:"ageTo"`
	Distance     uint64 `json:"distance"`
	Page         uint64 `json:"page"`
	Size         uint64 `json:"size"`
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

type ResponseFilterProfile struct {
	ID           uint64 `json:"id"`
	SearchGender string `json:"searchGender"`
	LookingFor   string `json:"lookingFor"`
	AgeFrom      uint8  `json:"ageFrom"`
	AgeTo        uint8  `json:"ageTo"`
	Distance     uint64 `json:"distance"`
	Page         uint64 `json:"page"`
	Size         uint64 `json:"size"`
}

type ResponseNavigatorProfile struct {
	Distance float64 `json:"distance"`
}

type ReviewProfile struct {
	ID         uint64    `json:"id"`
	ProfileID  uint64    `json:"profileId"`
	Message    string    `json:"message"`
	Rating     float32   `json:"rating"`
	HasDeleted bool      `json:"hasDeleted"`
	HasEdited  bool      `json:"hasEdited"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type ResponseReviewProfile struct {
	ID         uint64    `json:"id"`
	ProfileID  uint64    `json:"profileId"`
	Message    string    `json:"message"`
	Rating     float32   `json:"rating"`
	HasDeleted bool      `json:"hasDeleted"`
	HasEdited  bool      `json:"hasEdited"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	SessionID  string    `json:"sessionId"`
}

type QueryParamsReviewList struct {
	pagination.Pagination
	ProfileID string `json:"profileId"`
}

type ContentReviewProfile struct {
	ID          uint64    `json:"id"`
	ProfileID   uint64    `json:"profileId"`
	Message     string    `json:"message"`
	Rating      float32   `json:"rating"`
	HasDeleted  bool      `json:"hasDeleted"`
	HasEdited   bool      `json:"hasEdited"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DisplayName string    `json:"displayName"`
	SessionID   string    `json:"sessionId"`
}

type ResponseListReview struct {
	*pagination.Pagination
	RatingAverage            float32                 `json:"ratingAverage"`
	CountItemsTodayByProfile uint                    `json:"countItemsTodayByProfile"`
	Content                  []*ContentReviewProfile `json:"content"`
}

type RequestAddReview struct {
	ProfileID string `json:"profileId"`
	Message   string `json:"message"`
	Rating    string `json:"rating"`
}

type RequestUpdateReview struct {
	ID        string `json:"id"`
	ProfileID string `json:"profileId"`
	Message   string `json:"message"`
	Rating    string `json:"rating"`
}

type RequestDeleteReview struct {
	ID string `json:"id"`
}

type LikeProfile struct {
	ID        uint64    `json:"id"`
	ProfileID uint64    `json:"profileId"`
	HumanID   uint64    `json:"humanId"`
	IsLiked   bool      `json:"isLiked"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RequestAddLike struct {
	SessionID string `json:"sessionId"`
	HumanID   string `json:"humanId"`
}

type RequestUpdateLike struct {
	ID string `json:"id"`
}

type RequestDeleteLike struct {
	ID string `json:"id"`
}

type ResponseLikeProfile struct {
	ID        *uint64    `json:"id"`
	IsLiked   bool       `json:"isLiked"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type BlockedProfile struct {
	ID            uint64    `json:"id"`
	ProfileID     uint64    `json:"profileId"`
	BlockedUserID uint64    `json:"blockedUserId"`
	IsBlocked     bool      `json:"isLiked"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type RequestAddBlock struct {
	SessionID     string `json:"sessionId"`
	BlockedUserID string `json:"blockedUserId"`
}

type RequestUpdateBlock struct {
	ID string `json:"id"`
}

type ComplaintProfile struct {
	ID              uint64    `json:"id"`
	ProfileID       uint64    `json:"profileId"`
	ComplaintUserID uint64    `json:"complaintUserId"`
	Reason          string    `json:"reason"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type RequestAddComplaint struct {
	SessionID       string `json:"sessionId"`
	ComplaintUserID string `json:"complaintUserId"`
	Reason          string `json:"reason"`
}
