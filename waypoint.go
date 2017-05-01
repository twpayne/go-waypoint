package waypoint

import (
	"errors"
	"fmt"
	"image/color"
	"io"
)

var (
	// ErrUnknownFormat is returned when the format is unknown.
	ErrUnknownFormat = errors.New("waypoint: unknown format")
)

// An ErrSyntax is a syntax error.
type ErrSyntax struct {
	LineNo int
	Line   string
}

func (e ErrSyntax) Error() string {
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

func equal(t1, t2 *T) error {
	if t1.ID != t2.ID {
		return fmt.Errorf("ID mismatch: want %v, got %v", t1.ID, t2.ID)
	}
	if t1.Description != t2.Description {
		return fmt.Errorf("Description mismatch: want %v, got %v", t1.Description, t2.Description)
	}
	if t1.Latitude != t2.Latitude {
		return fmt.Errorf("Latitude mismatch: want %v, got %v", t1.Latitude, t2.Latitude)
	}
	if t1.Longitude != t2.Longitude {
		return fmt.Errorf("Longitude mismatch: want %v, got %v", t1.Longitude, t2.Longitude)
	}
	if t1.Altitude != t2.Altitude {
		return fmt.Errorf("Altitude mismatch: want %v, got %v", t1.Altitude, t2.Altitude)
	}
	if t1.Radius != t2.Radius {
		return fmt.Errorf("Radius mismatch: want %v, got %v", t1.Radius, t2.Radius)
	}
	if t1.Color == nil {
		if t2.Color != nil {
			return fmt.Errorf("Color mismatch: want %#v, got %#v", t1.Color, t2.Color)
		}
	} else if t2.Color == nil {
		return fmt.Errorf("Color mismatch: want %#v, got %#v", t1.Color, t2.Color)
	} else {
		r1, g1, b1, a1 := t1.Color.RGBA()
		r2, g2, b2, a2 := t2.Color.RGBA()
		if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
			return fmt.Errorf("Color mismatch: want %#v, got %#v", t1.Color, t2.Color)
		}
	}
	return nil
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
		return nil, ErrUnknownFormat
	}
}

// Read tries to read waypoints from rs using all known formats. When
// successful, it returns the waypoints and the original format.
func Read(rs io.ReadSeeker) (Collection, Format, error) {
	var formats = []Format{
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
	return nil, nil, ErrUnknownFormat
}

// Write writes c to w in format. Format must be a known format.
func Write(w io.Writer, c Collection, format string) error {
	f, err := New(format)
	if err != nil {
		return err
	}
	return f.Write(w, c)
}
