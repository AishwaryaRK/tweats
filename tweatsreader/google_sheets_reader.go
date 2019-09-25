package tweatsreader

import (
	"context"
	"errors"
	"github.com/AishwaryaRK/tweats/datamodel"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Constants
const (
	APIKey        = "******"
	SpreadSheetID = "******"
)

// Read reads google sheets and return the tweeps
func Read() ([]datamodel.Tweep, error) {
	ctx := context.Background()
	sheetsService, err := sheets.NewService(ctx, option.WithAPIKey(APIKey))
	if err != nil {
		return nil, err
	}
	readRange := "tweats!A2:J"
	resp, err := sheetsService.Spreadsheets.Values.Get(SpreadSheetID, readRange).Do()
	if err != nil {
		return nil, err
	}

	var tweeps []datamodel.Tweep
	if len(resp.Values) == 0 {
		return nil, errors.New("no data found")
	}
	for _, row := range resp.Values {
		//TO-DO: refactor to remove index hardcoding
		tweep := datamodel.Tweep{
			LDAP:                          row[1].(string),
			Name:                          row[2].(string),
			OfficeLocation:                row[3].(string),
			Interests:                     strings.Split(row[4].(string), ", "),
			AllergiesNDieteryRestrictions: row[5].(string),
			TweatLocation:                 row[6].(string),
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

	return tweeps, nil
}
