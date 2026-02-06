package preset

type ContentPreset struct {
	ID       uint
	ExternalID string
	Name     string
	Language string
	RoleID   uint
	DataJSON string
	IsActive bool
}
