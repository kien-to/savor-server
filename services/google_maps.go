package services

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"googlemaps.github.io/maps"
)

type GoogleMapsService struct {
	client *maps.Client
}

type DistanceResult struct {
	Distance string `json:"distance"`
	Duration string `json:"duration"`
	Meters   int    `json:"meters"`
	Seconds  int    `json:"seconds"`
}

type DirectionsResult struct {
	GoogleMapsURL string `json:"googleMapsUrl"`
	Distance      string `json:"distance"`
	Duration      string `json:"duration"`
}

func NewGoogleMapsService() (*GoogleMapsService, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_MAPS_API_KEY environment variable is required")
	}

	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Google Maps client: %v", err)
	}

	return &GoogleMapsService{client: client}, nil
}

func (g *GoogleMapsService) CalculateDistance(originLat, originLng, destLat, destLng float64) (*DistanceResult, error) {
	origin := fmt.Sprintf("%f,%f", originLat, originLng)
	destination := fmt.Sprintf("%f,%f", destLat, destLng)

	req := &maps.DistanceMatrixRequest{
		Origins:      []string{origin},
		Destinations: []string{destination},
		Mode:         maps.TravelModeDriving,
		Units:        maps.UnitsMetric,
	}

	log.Printf("Sending Distance Matrix request to Google Maps: %+v", req)

	resp, err := g.client.DistanceMatrix(context.Background(), req)
	if err != nil {
		log.Printf("ERROR: Google Maps Distance Matrix API call failed: %v", err)
		return nil, fmt.Errorf("failed to calculate distance: %v", err)
	}

	if len(resp.Rows) == 0 || len(resp.Rows[0].Elements) == 0 {
		log.Printf("ERROR: No distance data returned from Google Maps. Response: %+v", resp)
		return nil, fmt.Errorf("no distance data returned")
	}

	element := resp.Rows[0].Elements[0]
	if element.Status != "OK" {
		log.Printf("ERROR: Google Maps Distance Matrix element status not OK: %s. Response: %+v", element.Status, resp)
		return nil, fmt.Errorf("distance calculation failed: %s", element.Status)
	}

	log.Printf("Successfully received distance data from Google Maps: %s, %s", element.Distance.HumanReadable, element.Duration.String())

	return &DistanceResult{
		Distance: element.Distance.HumanReadable,
		Duration: element.Duration.String(),
		Meters:   element.Distance.Meters,
		Seconds:  int(element.Duration.Seconds()),
	}, nil
}

func (g *GoogleMapsService) GetDirectionsURL(originLat, originLng, destLat, destLng float64) string {
	baseURL := "https://www.google.com/maps/dir/"
	origin := fmt.Sprintf("%f,%f", originLat, originLng)
	destination := fmt.Sprintf("%f,%f", destLat, destLng)

	params := url.Values{}
	params.Add("api", "1")
	params.Add("origin", origin)
	params.Add("destination", destination)
	params.Add("travelmode", "driving")

	return baseURL + "?" + params.Encode()
}

func (g *GoogleMapsService) GetDirections(originLat, originLng, destLat, destLng float64) (*DirectionsResult, error) {
	origin := fmt.Sprintf("%f,%f", originLat, originLng)
	destination := fmt.Sprintf("%f,%f", destLat, destLng)

	req := &maps.DirectionsRequest{
		Origin:      origin,
		Destination: destination,
		Mode:        maps.TravelModeDriving,
	}

	log.Printf("Sending Directions request to Google Maps: %+v", req)

	routes, _, err := g.client.Directions(context.Background(), req)
	if err != nil {
		log.Printf("ERROR: Google Maps Directions API call failed: %v", err)
		return nil, fmt.Errorf("failed to get directions: %v", err)
	}

	if len(routes) == 0 {
		log.Printf("ERROR: No routes found from Google Maps.")
		return nil, fmt.Errorf("no routes found")
	}

	route := routes[0]
	leg := route.Legs[0]

	log.Printf("Successfully received directions from Google Maps: %s, %s", leg.Distance.HumanReadable, leg.Duration.String())

	return &DirectionsResult{
		GoogleMapsURL: g.GetDirectionsURL(originLat, originLng, destLat, destLng),
		Distance:      leg.Distance.HumanReadable,
		Duration:      leg.Duration.String(),
	}, nil
}

// Global instance
var GoogleMaps *GoogleMapsService

func InitializeGoogleMaps() error {
	var err error
	GoogleMaps, err = NewGoogleMapsService()
	if err != nil {
		log.Printf("Warning: Google Maps service not initialized: %v", err)
		return err
	}
	log.Println("Google Maps service initialized successfully")
	return nil
}
