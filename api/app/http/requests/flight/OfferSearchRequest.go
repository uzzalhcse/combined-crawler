package flightrequests

import (
	"combined-crawler/api/app/http/requests"
)

type FlightSearchRequest struct {
	OriginLocationCode      string `json:"origin_location_code" form:"origin_location_code" validate:"required"`
	DestinationLocationCode string `json:"destination_location_code" form:"destination_location_code" validate:"required"`
	DepartureDate           string `json:"departure_date" form:"departure_date" validate:"required"`
	ReturnDate              string `json:"return_date" form:"return_date" validate:"required"`
	Adult                   string `json:"adult" form:"adult" validate:"required"`
	IncludedAirlineCodes    string `json:"included_airline_codes" form:"included_airline_codes"`
	Max                     string `json:"max" form:"max"`
	*requests.Request
}
