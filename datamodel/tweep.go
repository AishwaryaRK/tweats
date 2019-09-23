package datamodel

type Interest int

const (
	Sports Interest = iota
	MusicAndEnterainment 
	CurrentAffairs
	HighTech
	TravelAndLifestyle
	NatureLovers
	Anything
)

type TimeSlot struct {
	Start int
	End int
}

type Availability struct {
	Weekday int
	TimeSlots []TimeSlot
}

var (
	InterestDescMapping = map[Interest]string {
		Sports: "Sports",
		MusicAndEnterainment: "Music And Enterainment",
		CurrentAffairs: "Current Affairs",
		HighTech: "High-Tech",
		TravelAndLifestyle: "Travel and Lifestyle",
		NatureLovers: "Nature Lovers",
		Anything: "Anything Under The Sun",
	}
)

type Tweep struct {
	LDAP string
	OfficeLocation string
	Interest []Interest
	AllergiesNDieteryRestrictions string
	EatAtOffice bool
	Availabilities []Availability
}