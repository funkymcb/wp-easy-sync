package models

// Member stores all necessarry data for wordpress account creation
type WVCMember struct {
	LoginName string
	FirstName string `json:"firstName"`
	LastName  string `json:"familyName"`
	Email     string `json:"privateEmail,omitempty"`
}
