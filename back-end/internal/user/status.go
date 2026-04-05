package user

type Status string

const (
	Unknown  Status = ""
	Pending  Status = "pending"
	Active   Status = "active"
	Inactive Status = "inactive"
)
