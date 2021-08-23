package repository

import (
	`fmt`
)

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
	RawURI  string `gorm:"raw_uri"`
	
	// optional field to get headers of associated request.
	Header Header `gorm:"request_header" json:"headers"`
	
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
	IP string
	
	// On which domain strixeye agent found this suspicion
	DomainID string
	
	// timestamp of suspicion creation in epoch milliseconds
	CreatedAt int64
	
	// Single associated request to this trip. I could've embedded it,
	// I know. Some changes around HTTP Cookies made us put requests in another table,
	// but for cli it's fine to embed, i guess.
	Request Request
	
	// Static Features of trip
	StaticChecks []StaticCheck
}

type IP struct {
	ID string
	
	// ipv4 or ipv6 address of request owner.
	Ip string
	
	// To whom this ip belongs.
	ProfileID string `gorm:"profile_id;size:36"`
}

// TableName is name of the table in database
func (i IP) TableName() string {
	return "ips"
}

// TableName is name of the table in database
func (Trip) TableName() string {
	return "trips"
}

// TripQueryArgs are arguments you can use to customize your queries. Multiple fields can be used at once,
// also empty query args is not a problem.
type TripQueryArgs struct {
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
	
	// If true, queries nested fields like request headers.
	Verbose bool
}

func (q TripQueryArgs) String() string {
	var query string
	
	query = fmt.Sprintf("%s\nDisplaying maximum %d rows", query, q.Limit)
	
	if q.SuspectIds != nil && len(q.SuspectIds) != 0 {
		query = fmt.Sprintf(
			"%s\nQuerying %d suspects with ids: %s", query, len(q.SuspectIds), q.SuspectIds,
		)
	}
	if q.Endpoints != nil && len(q.Endpoints) != 0 {
		query = fmt.Sprintf("%s\nQuerying only %d endpoints: %s", query, len(q.Endpoints), q.Endpoints)
	}
	
	if q.TripsIds != nil && len(q.TripsIds) != 0 {
		query = fmt.Sprintf(
			"%s\nQuerying only %d trips with ids: %s", query, len(q.TripsIds),
			q.TripsIds,
		)
	}
	
	if q.SinceTime > 0 {
		query = fmt.Sprintf("%s\nQuerying only trips that came after: %d", query, q.SinceTime)
	}
	
	return query
}

// Header keeps request headers of a single request.
type Header struct {
	// Primary Key
	ID string
	
	// Foreign Key
	RequestID string
	
	Method     string `gorm:"method"`
	URI        string `gorm:"uri"`
	UserAgent  string `gorm:"user_agent"`
	Host       string `gorm:"host"`
	RawHeaders string `gorm:"raw_headers"`
}

func (h Header) TableName() string {
	return "request_headers"
}
