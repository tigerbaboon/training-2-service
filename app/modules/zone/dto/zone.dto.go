package zonedto

type ZoneRequest struct {
	ZoneName string  `form:"zone_name"`
	Lat      float64 `form:"lat"`
	Long     float64 `form:"long"`
}

type ZoneResponse struct {
	ID       string        `json:"id"`
	ZoneName string        `json:"zone_name"`
	Lat      float64       `json:"lat"`
	Long     float64       `json:"long"`
	Images   ImageResponse `json:"image"`
}

type ZoneResponses struct {
	ID       string        `json:"id"`
	ZoneName string        `json:"zone_name"`
	Lat      float64       `json:"lat"`
	Long     float64       `json:"long"`
	Images   ImageResponse `json:"images"`
}

type ZoneForHouse struct {
	ID       string `json:"id"`
	ZoneName string `json:"zone_name"`
}

type ImageResponse struct {
	ID   string `json:"id"`
	URL  string `json:"url"`
	Type string `json:"type"`
}

type ZoneDTOResponse struct {
	ID       string   `json:"id"`
	ZoneName string   `json:"zone_name"`
	Lat      float64  `json:"lat"`
	Long     float64  `json:"long"`
	Images   []string `json:"image"`
}

type ZoneGetAllRequest struct {
	From   int    `form:"from"`
	Size   int    `form:"size"`
	Search string `form:"search"`
}

type ZoneResponseByID struct {
	ID       string          `json:"id"`
	ZoneName string          `json:"zone_name"`
	Lat      float64         `json:"lat"`
	Long     float64         `json:"long"`
	Images   []ImageResponse `json:"image"`
}

type ZoneUpdateRequest struct {
	ZoneName          string   `form:"zone_name"`
	Lat               float64  `form:"lat"`
	Long              float64  `form:"long"`
	RemainingImageIDs []string `form:"remainingImageIDs"`
}
