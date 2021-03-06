package maps

import (
	"errors"
	"strings"

	"github.com/photoprism/photoprism/internal/maps/osm"
	"github.com/photoprism/photoprism/internal/maps/places"
)

/* TODO

(SELECT pl.loc_label as album_name, pl.loc_country, YEAR(ph.taken_at) as taken_year, round(count(ph.id)) as photo_count FROM photos ph
        JOIN places pl ON ph.place_id = pl.id AND pl.id <> 1
        GROUP BY album_name, taken_year HAVING photo_count > 5) UNION (
            SELECT c.country_name AS album_name, pl.loc_country, YEAR(ph.taken_at) as taken_year, round(count(ph.id)) as photo_count FROM photos ph
        JOIN places pl ON ph.place_id = pl.id AND pl.id <> 1
            JOIN countries c ON c.id = pl.loc_country
        GROUP BY album_name, taken_year
        HAVING photo_count > 10)
ORDER BY loc_country, album_name, taken_year;

*/

// Photo location
type Location struct {
	ID          string
	LocName     string
	LocCategory string
	LocSuburb   string
	LocLabel    string
	LocCity     string
	LocState    string
	LocCountry  string
	LocSource   string
}

type LocationSource interface {
	CellID() string
	CountryCode() string
	Category() string
	Name() string
	City() string
	Suburb() string
	State() string
	Source() string
}

func NewLocation(id string) *Location {
	result := &Location{
		ID: id,
	}

	return result
}

func (l *Location) Query(api string) error {
	switch api {
	case "osm":
		return l.QueryOSM()
	case "places":
		return l.QueryPlaces()
	}

	return errors.New("maps: reverse lookup disabled")
}

func (l *Location) QueryPlaces() error {
	s, err := places.FindLocation(l.ID)

	if err != nil {
		return err
	}

	l.LocSource = s.Source()
	l.LocName = s.Name()
	l.LocCity = s.City()
	l.LocSuburb = s.Suburb()
	l.LocState = s.State()
	l.LocCountry = s.CountryCode()
	l.LocCategory = s.Category()
	l.LocLabel = s.Label()

	return nil
}

func (l *Location) QueryOSM() error {
	s, err := osm.FindLocation(l.ID)

	if err != nil {
		return err
	}

	return l.Assign(s)
}

func (l *Location) Assign(s LocationSource) error {
	l.LocSource = s.Source()

	l.ID = s.CellID()

	if l.Unknown() {
		l.LocCategory = "unknown"
		return errors.New("maps: unknown location")
	}

	l.LocName = s.Name()
	l.LocCity = s.City()
	l.LocSuburb = s.Suburb()
	l.LocState = s.State()
	l.LocCountry = s.CountryCode()
	l.LocCategory = s.Category()
	l.LocLabel = l.label()

	return nil
}

func (l *Location) Unknown() bool {
	return l.ID == ""
}

func (l *Location) label() string {
	if l.Unknown() {
		return "Unknown"
	}

	var countryName = l.CountryName()
	var loc []string

	shortCountry := len([]rune(countryName)) <= 20

	if l.LocCity != "" {
		loc = append(loc, l.LocCity)
	}

	if shortCountry && l.LocState != "" && l.LocCity != l.LocState {
		loc = append(loc, l.LocState)
	}

	if countryName != "" {
		loc = append(loc, countryName)
	}

	return strings.Join(loc[:], ", ")
}

func (l Location) Name() string {
	return l.LocName
}

func (l Location) Category() string {
	return l.LocCategory
}

func (l Location) Suburb() string {
	return l.LocSuburb
}

func (l Location) Label() string {
	return l.LocLabel
}

func (l Location) City() string {
	return l.LocCity
}

func (l Location) State() string {
	return l.LocState
}

func (l Location) CountryCode() string {
	return l.LocCountry
}

func (l Location) CountryName() string {
	return CountryNames[l.LocCountry]
}

func (l Location) Source() string {
	return l.LocSource
}
