package eventapi

// User is typical user structure, that comes from rubykube Event API.
type User struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Level int    `json:"level"`
	Otp   bool   `json:"otp_enabled"`
	State string `json:"state"`
}

// Record is required payload of event.
// Each event should have record, otherwice it is not valid!
type Record struct {
	User     *User  `json:"user"`
	Language string `json:"language"`
}

// Event represents basic event structure.
// * "record" - object payload (always).
// * "changes" - list of fields, which was changed (on update).
type Event struct {
	Record  map[string]interface{} `json:"record"`
	Changes map[string]interface{} `json:"changes"`
}
