package managerdto

type ManagerDTOResponse struct {
	Username    string `json:"username"`
	ManagerName string `json:"manager_name"`
}

type ManagerDTORequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	ManagerName string `json:"manager_name"`
}

type ManagerLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ManagerLoginResponse struct {
	Username    string `json:"username"`
	ManagerName string `json:"manager_name"`
	Token       string `json:"token"`
}

type ManagerUpdateRequest struct {
	ManagerName string `json:"manager_name"`
}
