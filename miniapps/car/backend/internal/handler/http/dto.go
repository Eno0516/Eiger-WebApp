package httphandler

import (
	"time"

	"github.com/enoro/eiger-webapp/miniapps/car/backend/internal/domain/entity"
)

type carJSON struct {
	ID          string `json:"id"`
	CircleID    string `json:"circleId"`
	Name        string `json:"name"`
	PlateNumber string `json:"plateNumber"`
	Available   bool   `json:"available"`
}

type reservationJSON struct {
	ID         string     `json:"id"`
	CarID      string     `json:"carId"`
	CircleID   string     `json:"circleId"`
	MemberID   string     `json:"memberId"`
	Start      time.Time  `json:"start"`
	End        time.Time  `json:"end"`
	Status     string     `json:"status"`
	DistanceKm int        `json:"distanceKm"`
	ReturnNote string     `json:"returnNote,omitempty"`
	ReturnedAt *time.Time `json:"returnedAt,omitempty"`
}

type locationJSON struct {
	CarID         string    `json:"carId"`
	ReservationID string    `json:"reservationId"`
	MemberID      string    `json:"memberId"`
	Link          string    `json:"link"`
	Note          string    `json:"note"`
	ReportedAt    time.Time `json:"reportedAt"`
}

func toCarJSON(c entity.Car) carJSON {
	return carJSON{ID: c.ID, CircleID: c.CircleID, Name: c.Name, PlateNumber: c.PlateNumber, Available: c.Available}
}

func toReservationJSON(r entity.Reservation) reservationJSON {
	out := reservationJSON{
		ID: r.ID, CarID: r.CarID, CircleID: r.CircleID, MemberID: r.MemberID,
		Start: r.Period.Start(), End: r.Period.End(), Status: string(r.Status),
		DistanceKm: r.DistanceKm, ReturnNote: r.ReturnNote,
	}
	if !r.ReturnedAt.IsZero() {
		returnedAt := r.ReturnedAt
		out.ReturnedAt = &returnedAt
	}
	return out
}

func toLocationJSON(l entity.LocationReport) locationJSON {
	return locationJSON{
		CarID: l.CarID, ReservationID: l.ReservationID, MemberID: l.MemberID,
		Link: l.Link.String(), Note: l.Note, ReportedAt: l.ReportedAt,
	}
}

func toCarsJSON(cars []entity.Car) []carJSON {
	out := make([]carJSON, 0, len(cars))
	for _, c := range cars {
		out = append(out, toCarJSON(c))
	}
	return out
}

func toReservationsJSON(rs []entity.Reservation) []reservationJSON {
	out := make([]reservationJSON, 0, len(rs))
	for _, r := range rs {
		out = append(out, toReservationJSON(r))
	}
	return out
}
