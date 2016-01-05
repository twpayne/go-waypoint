package waypoint

import (
	"encoding/json"
	"fmt"
	"io"
)

type GeoJSONFormat struct{}

type GeoJSONWaypoint struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	}
	Properties struct {
		Color       string  `json:"color"`
		Description string  `json:"description`
		Radius      float64 `json:"radius"`
	}
}

type GeoJSONWaypointFeatureCollection struct {
	Type     string            `json:"type"`
	Features []GeoJSONWaypoint `json:"features"`
}

func NewGeoJSONFormat() *GeoJSONFormat {
	return &GeoJSONFormat{}
}

func (*GeoJSONFormat) Extension() string {
	return "json"
}

func (*GeoJSONFormat) Name() string {
	return "geojson"
}

func (*GeoJSONFormat) Read(r io.Reader) (Collection, error) {
	var wfc GeoJSONWaypointFeatureCollection
	if err := json.NewDecoder(r).Decode(&wfc); err != nil {
		return nil, err
	}
	if wfc.Type != "FeatureCollection" {
		return nil, fmt.Errorf("expected FeatureCollection, got %v", wfc.Type)
	}
	var c Collection
	for _, f := range wfc.Features {
		if f.Type != "Feature" {
			return nil, fmt.Errorf("expected Feature, got %v", f.Type)
		}
		if f.Geometry.Type != "Point" {
			return nil, fmt.Errorf("expected Point, got %v", f.Geometry.Type)
		}
		// FIXME check size of f.Geometry.Coordinates
		t := &T{
			Id:          f.Id,
			Description: f.Properties.Description,
			Latitude:    f.Geometry.Coordinates[0],
			Longitude:   f.Geometry.Coordinates[1],
			Altitude:    f.Geometry.Coordinates[2],
			Radius:      f.Properties.Radius,
			//Color:       f.Properties.Color, // FIXME
		}
		c = append(c, t)
	}
	return c, nil
}

func (*GeoJSONFormat) Write(w io.Writer, wc Collection) error {
	if err := json.NewEncoder(w).Encode(wc); err != nil {
		return err
	}
	return nil
}

func (w *T) MarshalJSON() ([]byte, error) {
	o := map[string]interface{}{
		"id": w.Id,
		"geometry": map[string]interface{}{
			"type":        "Point",
			"coordinates": []float64{w.Latitude, w.Longitude, w.Altitude},
		},
		"type": "Feature",
	}
	properties := make(map[string]interface{})
	if w.Color != nil {
		r, g, b, _ := w.Color.RGBA()
		properties["color"] = fmt.Sprintf("#%02x%02x%02x", r/0x101, g/0x101, b/0x101)
	}
	if w.Description != "" {
		properties["description"] = w.Description
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
