package tweats_reader

import (
	"context"
	"errors"
	"strings"

	"github.com/AishwaryaRK/tweats/datamodel"

	"google.golang.org/api/option"
	sheets "google.golang.org/api/sheets/v4"
)

const API_KEY = ""
const SPREADSHEET_ID = ""

func Read() ([]datamodel.Tweep, error) {
	ctx := context.Background()
	sheetsService, err := sheets.NewService(ctx, option.WithAPIKey(API_KEY))
	if err != nil {
		return nil, err
	}
	readRange := "tweats!A2:J"
	resp, err := sheetsService.Spreadsheets.Values.Get(SPREADSHEET_ID, readRange).Do()
	if err != nil {
		return nil, err
	}

	var tweeps []datamodel.Tweep
	if len(resp.Values) == 0 {
		return nil, errors.New("No data found.")
	} else {
		for _, row := range resp.Values {
			//TO-DO: refactor to remove index hardcoding
			tweep := datamodel.Tweep{
				LDAP:                          row[1].(string),
				OfficeLocation:                row[2].(string),
				Interests:                     strings.Split(row[3].(string), ", "),
				AllergiesNDieteryRestrictions: row[4].(string),
				TweatLocation:                 row[5].(string),
				Availabilities: []datamodel.Availability{
					{
						Weekday: 1,
						TimeSlots: []datamodel.TimeSlot{
							{
								Start: 12,
								End:   13,
							},
						},
					},
				},
			}
			tweeps = append(tweeps, tweep)
		}
	}

	return tweeps, nil
}
