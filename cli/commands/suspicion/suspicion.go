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
	ProfileId string

	// What is the associated trip(request-response pair)'s id.
	TripId string

	// On which domain strixeye agent found this suspicion
	DomainId string

	// timestamp of suspicion creation in epoch milliseconds
	CreatedAt uint64
}

func (Suspicion) TableName() string {
	return "anomalies"
}
