package types

// SportEventUpdateTicket is the ticket type for sport event update
type SportEventUpdateTicket struct {
	Events []SportEvent `json:"events"`
}
