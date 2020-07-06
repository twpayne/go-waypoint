package waypoint

import (
	"encoding/json"
	"fmt"
	"io"
)

// A GeoJSONFormat is a GeoJSON format.
type GeoJSONFormat struct{}

// A GeoJSONWaypoint is a GeoJSON waypoint.
type GeoJSONWaypoint struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	}
	Properties struct {
		Color       string  `json:"color"`
		Description string  `json:"description"`
		Radius      float64 `json:"radius"`
	}
}

// A GeoJSONWaypointFeatureCollection is a GeoJSON FeatureCollection of GeoJSON
// waypoints.
type GeoJSONWaypointFeatureCollection struct {
	Type     string            `json:"type"`
	Features []GeoJSONWaypoint `json:"features"`
}

// NewGeoJSONFormat returns a new GeoJSONFormat.
func NewGeoJSONFormat() *GeoJSONFormat {
	return &GeoJSONFormat{}
}

// Extension returns f's extension.
func (f *GeoJSONFormat) Extension() string {
	return "json"
}

// Name returns f's name.
func (f *GeoJSONFormat) Name() string {
	return "geojson"
}

// Read reads a Collection from r.
func (f *GeoJSONFormat) Read(r io.Reader) (Collection, error) {
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
			ID:          f.ID,
			Description: f.Properties.Description,
			Latitude:    f.Geometry.Coordinates[0],
			Longitude:   f.Geometry.Coordinates[1],
			Altitude:    f.Geometry.Coordinates[2],
			Radius:      f.Properties.Radius,
			// Color:       f.Properties.Color, // FIXME
		}
		c = append(c, t)
	}
	return c, nil
}

// Write writes c to w.
func (f *GeoJSONFormat) Write(w io.Writer, wc Collection) error {
	return json.NewEncoder(w).Encode(wc)
}

// MarshalJSON implements encoding/json.Marshaler.
func (w *T) MarshalJSON() ([]byte, error) {
	o := map[string]interface{}{
		"id": w.ID,
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

// MarshalJSON implements encoding/json.Marshaler.
func (wc Collection) MarshalJSON() ([]byte, error) {
	o := map[string]interface{}{
		"type":     "FeatureCollection",
		"features": []*T(wc),
	}
	return json.Marshal(o)
}
