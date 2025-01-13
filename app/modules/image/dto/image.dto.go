package imagedto

type ImageRequest struct {
	StationID string `json:"s_id"`
	ImageURL  string `json:"image_url"`
}

type ImageResponse struct {
	ImageName string `json:"image_name"`
	ImageURL  string `json:"image_url"`
}
