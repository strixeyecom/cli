package config

import `encoding/json`

/*
	Created by aomerk at 5/20/21 for project strixeye
*/

/*
 
 */

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// Config is base interface to implement for strixeye configuration structs.
type Config interface {
	json.Marshaler
	json.Unmarshaler
	
	// Since most of the config is crucial, validation process is highly encouraged.
	Validate() error
	
	// Configs are mostly kept as files.
	Save(filePath string) error
}
