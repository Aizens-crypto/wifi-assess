package models

import "time"

// Client is the tracked state for a single wireless station (device),
// built from probe requests and association traffic it sends/receives.
type Client struct {
	MAC string

	// ProbedSSIDs is the set of SSIDs this client has actively probed for.
	// A client probing for many distinct SSIDs — especially ones it isn't
	// currently near — is a common device-fingerprinting/history-leak
	// signal, and feeds Phase 4/5 detection later.
	ProbedSSIDs []string

	// AssociatedBSSID is the AP this client most recently appeared
	// associated with, inferred from association request/response traffic
	// between the two. Empty if no association has been observed yet.
	AssociatedBSSID string

	FirstSeen time.Time
	LastSeen  time.Time

	ProbeRequestCount int
}
