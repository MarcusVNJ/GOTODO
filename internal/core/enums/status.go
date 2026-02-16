package enums

type Status string

const (
	Pending Status= "PENDING"
	InProcess Status= "IN_PROCESS"
	Completed Status= "COMPLETED"
	Cancelled Status= "CANCELLED"
)

func (status Status) String() string {
	return string(status)
}