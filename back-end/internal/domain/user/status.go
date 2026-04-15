package user

type Status string

const (
	Unknown  Status = ""
	Active   Status = "active"
	Inactive Status = "inactive"
)
