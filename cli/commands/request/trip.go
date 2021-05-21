package trip

/*
	Created by aomerk at 5/21/21 for project cli
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// Request is an HTTP request that has been parsed by strixeye sensors and extracted relevant information.
// A request and response combined together is called a Trip.
type Request struct {
	// request's own id, in uuid v4 kept as string
	ID string

	// self explanatory fields
	RawBody string `gorm:"raw_body"`
	RawUri  string `gorm:"raw_uri"`

	// foreign key to associated trip.
	TripID string
}

// TableName is name of the table in database
func (Request) TableName() string {
	return "requests"
}

// Trip is the table where we log incoming request/response pairs. If no response exists in your agent,
// don't worry.
//
// Response part of trip is experimental only.
type Trip struct {
	// Suspicions id. Stored as string, uuid v4
	ID string

	// To whom this trip belongs.
	ProfileID string `gorm:"profile_id;size:36"`

	// What is the associated client's ip address
	Ip string

	// On which domain strixeye agent found this suspicion
	DomainId string

	// timestamp of suspicion creation in epoch milliseconds
	CreatedAt uint64

	// Single associated request to this trip. I could've embedded it,
	// I know. Some changes around HTTP Cookies made us put requests in another table,
	// but for cli it's fine to embed, i guess.
	Request Request
}

type Ip struct {
	ID string

	// ipv4 or ipv6 address of request owner.
	Ip string

	// To whom this ip belongs.
	ProfileID string `gorm:"profile_id;size:36"`
}

// TableName is name of the table in database
func (i Ip) TableName() string {
	return "ips"
}

// TableName is name of the table in database
func (Trip) TableName() string {
	return "trips"
}

// QueryArgs are arguments you can use to customize your queries. Multiple fields can be used at once,
// also empty query args is not a problem.
type QueryArgs struct {
	// how many results do you want to retrieve
	Limit int

	// get only profiles who has detected since given epoch "millisecond" timestamp
	SinceTime int64

	// Request-Response Pair ids the suspicions must relate to
	TripsIds []string

	// List of suspects, that you want the resulting suspicions belong to.
	SuspectIds []string

	// Get trip that are only requested to given endpoints
	Endpoints []string
}
