package models

type PlanType string
type StatusType string

const (
	FreePlan   PlanType = "free"
	StarPlan   PlanType = "star"
	GalaxyPlan PlanType = "galaxy"
)

const (
	ActiveStatus   StatusType = "active"
	InactiveStatus StatusType = "inactive"
	PendingStatus  StatusType = "pending"
)

type Premium struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Plan      PlanType   `json:"plan"`
	Status    StatusType `json:"status"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}
