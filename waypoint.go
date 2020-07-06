package waypoint

import (
	"errors"
	"fmt"
	"image/color"
	"io"
)

// errUnknownFormat is returned when the format is unknown.
var errUnknownFormat = errors.New("waypoint: unknown format")

// An errSyntax is a syntax error.
type errSyntax struct {
	LineNo int
	Line   string
}

func (e errSyntax) Error() string {
	return fmt.Sprintf("syntax error:%d: %v", e.LineNo, e.Line)
}

// A T is a waypoint.
type T struct {
	ID          string
	Description string
	Latitude    float64
	Longitude   float64
	Altitude    float64
	Radius      float64
	Color       color.Color
}

// A Collection is a collection of Ts.
type Collection []*T

// A Format is a waypoint format, with metadata and methods to read and write.
type Format interface {
	Extension() string
	Name() string
	Read(io.Reader) (Collection, error)
	Write(io.Writer, Collection) error
}

// New returns a new Format. format must be a known format.
func New(format string) (Format, error) {
	switch format {
	case "compegps":
		return NewCompeGPSFormat(), nil
	case "formatgeo":
		return NewFormatGeoFormat(), nil
	case "geojson":
		return NewGeoJSONFormat(), nil
	case "oziexplorer":
		return NewOziExplorerFormat(), nil
	case "seeyou":
		return NewSeeYouFormat(), nil
	default:
		return nil, errUnknownFormat
	}
}

// Read tries to read waypoints from rs using all known formats. When
// successful, it returns the waypoints and the original format.
func Read(rs io.ReadSeeker) (Collection, Format, error) {
	formats := []Format{
		NewCompeGPSFormat(),
		NewFormatGeoFormat(),
		NewGeoJSONFormat(),
		NewOziExplorerFormat(),
		NewSeeYouFormat(),
	}
	offset, err := rs.Seek(0, 1)
	if err != nil {
		return nil, nil, err
	}
	for _, format := range formats {
		if c, err := format.Read(rs); err == nil {
			return c, format, nil
		}
		if _, err := rs.Seek(offset, 0); err != nil {
			return nil, nil, err
		}
	}
	return nil, nil, errUnknownFormat
}

// Write writes wc to w in format. Format must be a known format.
func Write(w io.Writer, wc Collection, format string) error {
	f, err := New(format)
	if err != nil {
		return err
	}
	return f.Write(w, wc)
}
