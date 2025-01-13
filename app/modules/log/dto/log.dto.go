package logdto

type LogDTORequest struct {
	MenagerID   string `json:"menager_id"`
	ActionType  string `json:"actiontype"`
	Description string `json:"description"`
	RecordID    string `json:"record_id"`
	TableName   string `json:"table_name"`
	UpdatedBy   string `json:"updated_by"`
}

type LogDTOResponse struct {
	ID          string `json:"id"`
	MenagerID   string `json:"menager_id"`
	ActionType  string `json:"actiontype"`
	Description string `json:"description"`
	RecordID    string `json:"record_id"`
	TableName   string `json:"table_name"`
	UpdatedAt   int64  `json:"updated_at"`
	UpdatedBy   string `json:"updated_by"`
}
