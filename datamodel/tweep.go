package datamodel

//24hr format
type TimeSlot struct {
	Start int
	End   int
}

type Availability struct {
	//Monday = 1, Friday = 5
	Weekday   int
	TimeSlots []TimeSlot
}

type Tweep struct {
	LDAP                          string
	OfficeLocation                string
	Interest                      []string
	AllergiesNDieteryRestrictions string
	TweatLocation                 string
	Availabilities                []Availability
}
