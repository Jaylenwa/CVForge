package taxonomy

type JobCategoryDTO struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	ParentID         *uint  `json:"parentId"`
	OrderNum         int    `json:"orderNum"`
	IsActive         bool   `json:"isActive"`
}

type JobRoleDTO struct {
	ID         uint   `json:"id"`
	CategoryID uint   `json:"categoryId"`
	Name       string `json:"name"`
	OrderNum   int    `json:"orderNum"`
	IsActive   bool   `json:"isActive"`
}
