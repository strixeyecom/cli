package suspicion

/*
	Created by aomerk at 5/21/21 for project cli
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// Suspicion is the table we use too keep noteworthy anomalies,
// mostly created by our ai backed security engine or request based static engine.
type Suspicion struct {
	// Suspicions id. Stored as string, uuid v4
	Id string

	// To whom this suspicion belongs. You know it as Suspect and Suspect Id,
	// but agent obviously doesn't think everybody is suspect, there are visitors who arent suspect.
	//
	// Suspects are a subset of visitors.
	ProfileID string `gorm:"profile_id;size:36"`

	// What is the associated trip(request-response pair)'s id.
	TripId string

	// On which domain strixeye agent found this suspicion
	DomainId string

	// timestamp of suspicion creation in epoch milliseconds
	CreatedAt uint64
}

type Ip struct {
	ID string
	Ip string
	// To whom this ip belongs.
	ProfileID string `gorm:"profile_id;size:36"`
}

func (i Ip) TableName() string {
	return "ips"
}

func (Suspicion) TableName() string {
	return "anomalies"
}

type QueryArgs struct {
	// how many results do you want to retrieve
	Limit int

	// list of suspicion ids to return
	SuspicionIds []string

	// List of suspects, that you want the resulting suspicions belong to.
	SuspectIds []string

	// Request-Response Pair ids the suspicions must relate to
	TripsIds []string
}
