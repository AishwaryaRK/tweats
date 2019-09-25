package datamodel

// TimeSlot defines the start and end time, in 24hr format
type TimeSlot struct {
	Start int
	End   int
}

// Availability defines the availability
type Availability struct {
	//Monday = 1, Friday = 5
	Weekday   int
	TimeSlots []TimeSlot
}

// Tweep defines a tweep
type Tweep struct {
	Name                          string
	LDAP                          string
	Name                          string
	OfficeLocation                string
	Interests                     []string
	AllergiesNDieteryRestrictions string
	TweatLocation                 string
	Availabilities                []Availability
}
