package waypoint

import (
	"encoding/json"
	"fmt"
	"image/color"
)

type ErrSyntax struct {
	LineNo int
	Line   string
}

func (e ErrSyntax) Error() string {
	return fmt.Sprintf("syntax error:%d: %v", e.LineNo, e.Line)
}

type T struct {
	Id          string
	Name        string
	Description string
	Latitude    float64
	Longitude   float64
	Altitude    float64
	Radius      float64
	Color       color.Color
}

type Collection []*T

func (w *T) MarshalJSON() ([]byte, error) {
	o := map[string]interface{}{
		"id": w.Id,
		"geometry": map[string]interface{}{
			"type":        "Point",
			"coordinates": []float64{w.Latitude, w.Longitude, w.Altitude},
		},
	}
	properties := make(map[string]interface{})
	if w.Color != nil {
		r, g, b, _ := w.Color.RGBA()
		properties["color"] = fmt.Sprintf("#%02x%02x%02x", r/0x101, g/0x101, b/0x101)
	}
	if w.Description != "" {
		properties["description"] = w.Description
	}
	if w.Name != "" {
		properties["name"] = w.Name
	}
	if w.Radius > 0 {
		properties["radius"] = w.Radius
	}
	if len(properties) > 0 {
		o["properties"] = properties
	}
	return json.Marshal(o)
}

func (wc Collection) MarshalJSON() ([]byte, error) {
	o := map[string]interface{}{
		"type":     "FeatureCollection",
		"features": []*T(wc),
	}
	return json.Marshal(o)
}
