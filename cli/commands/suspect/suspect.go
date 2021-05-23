package suspect

import "github.com/usestrix/cli/cli/commands/suspicion"

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
	Score      uint64
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
	Score      float64

	// get only profiles who has detected since given epoch "millisecond" timestamp
	SinceTime int64

	// 	most fields are kept in different tables, bound via foreign keys and have nested relations
	// 	to get which fields you want to load other than the default, set it via fields argument
	Fields []string
}
