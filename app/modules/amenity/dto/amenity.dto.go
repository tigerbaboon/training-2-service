package amenitydto

type AmenityDTOResponse struct {
	ID          string `json:"id"`
	AmenityName string `json:"amenity_name"`
	Icons       string `json:"icons"`
}

type AmenityRequest struct {
	AmenityName string `json:"amenity_name"`
	Icons       string `json:"icons"`
}

type AmenityGetAllRequest struct {
	From   int    `form:"from"`
	Size   int    `form:"size"`
	Search string `form:"search"`
}
