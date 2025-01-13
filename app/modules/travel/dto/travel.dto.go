package traveldto

type TravelResponse struct {
	ID                string   `json:"id"`
	TravelTitle       string   `json:"travel_title"`
	TravelDetail      string   `json:"travel_detail"`
	ImagesMain        string   `json:"images_main"`
	Images            []string `json:"images"`
	Status            bool     `json:"status"`
	Address           string   `json:"address"`
	LocationLatitute  float64  `json:"location_latitute"`
	LocationLongitute float64  `json:"location_longitute"`
}

type ImageResponse struct {
	ID   string `json:"id"`
	URL  string `json:"url"`
	Type string `json:"type"`
}

type TravelResponseByID struct {
	ID                string          `json:"id"`
	TravelTitle       string          `json:"travel_title"`
	TravelDetail      string          `json:"travel_detail"`
	ImageMain         ImageResponse   `json:"image_main"`
	Images            []ImageResponse `json:"images"`
	Status            bool            `json:"status"`
	Address           string          `json:"address"`
	LocationLatitute  float64         `json:"location_latitute"`
	LocationLongitute float64         `json:"location_longitute"`
}

type TravelResponseByGetAll struct {
	ID          string        `json:"id"`
	TravelTitle string        `json:"travel_title"`
	ImageMain   ImageResponse `json:"image_main"`
	Status      bool          `json:"status"`
}

type TravelRequest struct {
	TravelTitle       string  `form:"travel_title"`
	TravelDetail      string  `form:"travel_detail"`
	Address           string  `form:"address"`
	LocationLatitute  float64 `form:"location_latitute"`
	LocationLongitute float64 `form:"location_longitute"`
}

type TravelGetAllRequest struct {
	From   int    `form:"from"`
	Size   int    `form:"size"`
	Search string `form:"search"`
}

type TravelUpdateRequest struct {
	TravelTitle          string   `form:"travel_title"`
	TravelDetail         string   `form:"travel_detail"`
	Address              string   `form:"address"`
	LocationLatitute     float64  `form:"location_latitute"`
	LocationLongitute    float64  `form:"location_longitute"`
	RemainingImageIDs    []string `form:"remainingImageIDs"`
	RemainingImageMainID string   `form:"remainingImageMainID"`
}

type TravelUpdateStatus struct {
	Status bool `json:"status"`
}

type TravelStatusResponse struct {
	ID          string `json:"id"`
	TravelTitle string `json:"travel_title"`
	Status      bool   `json:"status"`
}
