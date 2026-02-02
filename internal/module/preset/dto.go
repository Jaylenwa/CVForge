package preset

type ContentPresetDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Language string `json:"language"`
	RoleID   uint   `json:"roleId"`
	DataJSON string `json:"dataJson"`
	IsActive bool   `json:"isActive"`
}
