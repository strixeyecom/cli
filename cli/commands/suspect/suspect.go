package suspect

import (
	`fmt`
	
	`github.com/usestrix/cli/cli/commands/suspicion`
)

/*
	Created by aomerk at 5/21/21 for project cli
*/

/*
	Model declaration for suspect information
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// Suspect is the main table in strixeye agent. However, it is called a profile,
// because strixeye's underlying technology treats security threats on profile(visitor) base,
// and not all of them are actually "suspect"s
type Suspect struct {
	ID         string
	Suspicions []*suspicion.Suspicion `gorm:"anomalies;foreignKey:profile_id"`
	Ips        []*suspicion.Ip        `gorm:"ips;foreignKey:profile_id"`
	Score      int64
}

// TableName as I explained in type definition, agent knows suspects as profiles.
func (suspect Suspect) TableName() string {
	return "profiles"
}

// QueryArgs are arguments you can use to customize your queries. Multiple fields can be used at once,
// also empty query args is not a problem.
type QueryArgs struct {
	Limit      int
	SuspectIds []string
	
	// Minimum risk score of queried suspects. Higher means they are more likely to attack.
	MinScore int64
	
	// get only profiles who has detected since given epoch "millisecond" timestamp
	SinceTime int64
	
	// 	most fields are kept in different tables, bound via foreign keys and have nested relations
	// 	to get which fields you want to load other than the default, set it via fields argument
	Fields  []string
	Verbose bool
	
}

func (q QueryArgs) String() string {
	var query string
	
	query = fmt.Sprintf("%s\nDisplaying maximum %d rows", query, q.Limit)
	
	if q.SuspectIds != nil && len(q.SuspectIds) != 0 {
		query = fmt.Sprintf(
			"%s\nQuerying %d suspects with ids: %s", query, len(q.SuspectIds), q.SuspectIds,
		)
	}
	if q.SinceTime > 0 {
		query = fmt.Sprintf("%s\nQuerying only suspects that came after: %d", query, q.SinceTime)
	}
	
	if q.MinScore > 0 {
		query = fmt.Sprintf("%s\nQuerying only suspects with score higher than: %d", query, q.MinScore)
	}
	
	return query
}
