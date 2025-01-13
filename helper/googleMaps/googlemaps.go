package helper

import (
	"app/internal/modules/log"
	"context"
	"errors"
	"os"

	"googlemaps.github.io/maps"
)

func GetClient() (*maps.Client, error) {
	mapsClient, err := NewMaps()
	if err != nil {
		log.Info("Error creating maps client: %v", err)
		return nil, err
	}

	return mapsClient, nil
}

func NewMaps() (*maps.Client, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		return nil, errors.New("GOOGLE_MAPS_API_KEY is not set")
	}

	log.Info("Google Maps API Key: %s", apiKey)

	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		log.Info("Error creating maps client: %v", err)
		return nil, err
	}

	return client, nil
}

func Geocode(address string) (float64, float64, error) {
	req := &maps.GeocodingRequest{
		Address: address,
	}

	client, err := GetClient()
	if err != nil {
		log.Info("Error getting maps client: %v", err)
		return 0, 0, err
	}

	resp, err := client.Geocode(context.Background(), req)
	if err != nil {
		log.Info("Error geocoding address: %v", err)
		return 0, 0, err
	}

	if len(resp) == 0 {
		return 0, 0, errors.New("no results found")
	}

	return resp[0].Geometry.Location.Lat, resp[0].Geometry.Location.Lng, nil
}
