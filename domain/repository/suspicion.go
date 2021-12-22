package repository

import (
	"fmt"
	"time"
)

/*
	Created by aomerk at 5/21/21 for project cli
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

type StaticCheck struct {
	Id       string
	TripId   string   ` protobuf:"bytes,2,opt,name=TripId"              json:"TripId,omitempty"    gorm:"trip_id;size:36"`
	ModSecID string   ` protobuf:"bytes,3,opt,name=ModSecID"            json:"ModSecID,omitempty"  gorm:"mod_sec_id"`
	Group    string   ` protobuf:"bytes,4,opt,name=Group"               json:"Group,omitempty"     gorm:"group"`
	Message  string   ` protobuf:"bytes,5,opt,name=Message"             json:"Message,omitempty"   gorm:"message"`
	Tags     []string ` protobuf:"bytes,6,rep,name=Tags"                json:"Tags,omitempty"      gorm:"-"`
	PL       int32    ` protobuf:"varint,7,opt,name=PL"                 json:"PL,omitempty"        gorm:"pl"`
}

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
	Domain   Domain `gorm:"-"`
	// timestamp of suspicion creation in epoch milliseconds
	CreatedAt int64
	Trip      Trip `gorm:"-"`
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

// QueryArgs are arguments you can use to customize your queries. Multiple fields can be used at once,
// also empty query args is not a problem.
type SuspicionQueryArgs struct {
	// how many results do you want to retrieve
	Limit int

	// list of suspicion ids to return
	SuspicionIds []string

	// List of suspects, that you want the resulting suspicions belong to.
	SuspectIds []string

	// Request-Response Pair ids the suspicions must relate to
	TripsIds []string

	// get only suspicions who has detected since given epoch "millisecond" timestamp
	SinceTime int64
}

func (q SuspicionQueryArgs) String() string {
	var query string

	query = fmt.Sprintf("%s\nDisplaying maximum %d rows", query, q.Limit)

	if q.SuspectIds != nil && len(q.SuspectIds) != 0 {
		query = fmt.Sprintf(
			"%s\nQuerying %d suspects with ids: %s", query, len(q.SuspectIds), q.SuspectIds,
		)
	}
	if q.SuspicionIds != nil && len(q.SuspicionIds) != 0 {
		query = fmt.Sprintf("%s\nQuerying only %d suspicions: %s", query, len(q.SuspicionIds), q.SuspicionIds)
	}

	if q.TripsIds != nil && len(q.TripsIds) != 0 {
		query = fmt.Sprintf("%s\nQuerying only %d trips with ids: %s", query, len(q.TripsIds), q.TripsIds)
	}

	if q.SinceTime > 0 {
		time.Unix(q.SinceTime, 0).Format(time.RFC822Z)
		query = fmt.Sprintf("%s\nQuerying only suspicions that came after: %d", query, q.SinceTime)
	}

	return query
}

// Domain keeps information about a single domain
type Domain struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
}

type DomainMessage struct {
	Data   Domain `json:"data"`
	Status string `json:"status"`
}
