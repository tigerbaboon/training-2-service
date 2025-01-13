package contactdto

type ContactRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	LineID      string `json:"line_id"`
	HouseID     string `json:"house_id"`
}

type ContactResponse struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	LineID      string `json:"line_id"`
	HouseID     string `json:"house_id"`
}
