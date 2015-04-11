package waypoint

import (
	"errors"
	"fmt"
	"image/color"
	"io"
)

var ErrUnknownFormat = errors.New("waypoint: unknown format")

type ErrSyntax struct {
	LineNo int
	Line   string
}

func (e ErrSyntax) Error() string {
	return fmt.Sprintf("syntax error:%d: %v", e.LineNo, e.Line)
}

type T struct {
	Id          string
	Description string
	Latitude    float64
	Longitude   float64
	Altitude    float64
	Radius      float64
	Color       color.Color
}

func Equal(t1, t2 *T) bool {
	if t1.Id != t2.Id {
		return false
	}
	if t1.Description != t2.Description {
		return false
	}
	if t1.Latitude != t2.Latitude {
		return false
	}
	if t1.Longitude != t2.Longitude {
		return false
	}
	if t1.Altitude != t2.Altitude {
		return false
	}
	if t1.Radius != t2.Radius {
		return false
	}
	if t1.Color == nil {
		if t2.Color != nil {
			return false
		}
	} else if t2.Color == nil {
		return false
	} else {
		r1, g1, b1, a1 := t1.Color.RGBA()
		r2, g2, b2, a2 := t2.Color.RGBA()
		if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
			return false
		}
	}
	return true
}

type Collection []*T

type Format interface {
	Extension() string
	Name() string
	Read(io.Reader) (Collection, error)
	Write(io.Writer, Collection) error
}

func Read(rs io.ReadSeeker) (Collection, Format, error) {
	var formats = []Format{
		NewCompeGPSFormat(),
		NewFormatGeoFormat(),
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

func Write(w io.Writer, c Collection, format string) error {
	var f Format
	switch format {
	case "compegps":
		f = NewCompeGPSFormat()
	case "formatgeo":
		f = NewFormatGeoFormat()
	case "oziexplorer":
		f = NewOziExplorerFormat()
	case "seeyou":
		f = NewSeeYouFormat()
	default:
		return ErrUnknownFormat
	}
	return f.Write(w, c)
}
