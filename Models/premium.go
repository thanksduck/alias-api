package models

import "time"

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
	ID             uint32     `json:"-"`
	UserID         uint32     `json:"-"`
	SubscriptionID string     `json:"suid"`
	Username       string     `json:"username"`
	Mobile         string     `json:"mobile"`
	Plan           PlanType   `json:"plan"`
	Status         StatusType `json:"status"`
	Gateway        string     `json:"gateway"`
	TransactionID  string     `json:"txnid"`
	MerchentUserID string     `json:"muid"`
	CreatedAt      time.Time  `json:"-"`
	UpdatedAt      time.Time  `json:"-"`
}
