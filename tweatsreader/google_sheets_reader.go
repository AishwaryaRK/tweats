package tweatsreader

import (
	"context"
	"errors"
	"github.com/AishwaryaRK/tweats/config"
	"github.com/AishwaryaRK/tweats/datamodel"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"strings"
)

// Read reads google sheets and return the tweeps
func Read() ([]datamodel.Tweep, error) {
	ctx := context.Background()
	sheetsService, err := sheets.NewService(ctx, option.WithAPIKey(config.GOOGLE_SHEET_API_KEY))
	if err != nil {
		return nil, err
	}
	readRange := "tweats!A2:K"
	resp, err := sheetsService.Spreadsheets.Values.Get(config.GOOGLE_SHEET_SPREAD_SHEET_ID, readRange).Do()
	if err != nil {
		return nil, err
	}

	var tweeps []datamodel.Tweep
	if len(resp.Values) == 0 {
		return nil, errors.New("no data found")
	}

	for _, row := range resp.Values {
		var availabilites []datamodel.Availability
		days := make(map[int][]datamodel.TimeSlot)
		if len(row) > 7 {
			days = calculate_availability(row[7].(string), days, 11, 12)
			if len(row) > 8 {
				days = calculate_availability(row[8].(string), days, 12, 13)
				if len(row) > 9 {
					days = calculate_availability(row[9].(string), days, 13, 14)
					if len(row) > 10 {
						days = calculate_availability(row[10].(string), days, 14, 15)
					}
				}
			}
		}

		for k, v := range days {
			availabilites = append(availabilites,
				datamodel.Availability{
					Weekday:   k,
					TimeSlots: v,
				})
		}

		tweep := datamodel.Tweep{
			LDAP:                          row[1].(string),
			Name:                          row[2].(string),
			OfficeLocation:                row[3].(string),
			Interests:                     strings.Split(row[4].(string), ", "),
			AllergiesNDieteryRestrictions: row[5].(string),
			TweatLocation:                 row[6].(string),
			Availabilities:                availabilites,
		}
		tweeps = append(tweeps, tweep)
	}

	return tweeps, nil
}

func calculate_availability(avail_string string, days map[int][]datamodel.TimeSlot, start int, end int) map[int][]datamodel.TimeSlot {
	for _, day := range strings.Split(avail_string, ", ") {
		dayIndex := -1
		switch day {
		case "MONDAY":
			dayIndex = 1
		case "TUESDAY":
			dayIndex = 2
		case "WEDNESDAY":
			dayIndex = 3
		case "THURSDAY":
			dayIndex = 4
		case "FRIDAY":
			dayIndex = 5
		default:
			return days
		}

		if val, ok := days[dayIndex]; ok {
			val = append(val, datamodel.TimeSlot{
				Start: start,
				End:   end,
			})
			days[dayIndex] = val
		} else {
			days[dayIndex] = []datamodel.TimeSlot{
				{
					Start: start,
					End:   end,
				},
			}
		}
	}
	return days
}
