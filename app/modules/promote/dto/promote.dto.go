package promotedto

type PromoteDTORequest struct {
	PromoteName string `form:"promote_name" binding:"required"`
	PromoteType string `form:"promote_type" binding:"required"`
	Link        string `form:"link"`
	Status      bool   `form:"status"`
}

type PromoteDTOResponse struct {
	ID          string   `json:"id"`
	PromoteName string   `json:"promote_name"`
	PromoteType string   `json:"promote_type"`
	Status      bool     `json:"status"`
	Link        string   `json:"link"`
	Images      []string `json:"images"`
}

type ImageResponse struct {
	ID   string `json:"id"`
	URL  string `json:"url"`
	Type string `json:"type"`
}

type PromoteResponse struct {
	ID          string          `json:"id"`
	PromoteName string          `json:"promote_name"`
	PromoteType string          `json:"promote_type"`
	Status      bool            `json:"status"`
	Link        string          `json:"link"`
	Images      []ImageResponse `json:"images"`
}

type PromoteResponseAll struct {
	ID          string        `json:"id"`
	PromoteName string        `json:"promote_name"`
	PromoteType string        `json:"promote_type"`
	Link        string        `json:"link"`
	Status      bool          `json:"status"`
	Image       ImageResponse `json:"image"`
}

type PromoteGetAllRequest struct {
	From   int    `form:"from"`
	Size   int    `form:"size"`
	Search string `form:"search"`
	Type   string `form:"type"`
}

type PromoteUpdateStatus struct {
	Status bool `json:"status"`
}

type PromoteUpdateRequest struct {
	PromoteName      string `form:"promote_name"`
	PromoteType      string `form:"promote_type"`
	Link             string `form:"link"`
	RemainingImageID string `form:"remainingImageID"`
}

type PromoteResponseByID struct {
	ID          string          `json:"id"`
	PromoteName string          `json:"promote_name"`
	PromoteType string          `json:"promote_type"`
	Images      []ImageResponse `json:"images"`
	Status      bool            `json:"status"`
	Link        string          `json:"link"`
}
