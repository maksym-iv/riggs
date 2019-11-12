package geo

import (
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	pb "riggs/pb"
)

var (
	tr = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:       30 * time.Second,
			KeepAlive:     30 * time.Second,
			FallbackDelay: -1,
		}).DialContext,
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client = &http.Client{Transport: tr}
)

func Test_getGeo(t *testing.T) {
	testCases := []struct {
		name   string
		ip     string
		geoDb  string
		client *http.Client
		want   error
	}{
		{
			name:   "173.177.164.160",
			ip:     "173.177.164.160",
			geoDb:  "geolocation-db.com",
			client: client,
			want:   nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := getGeo(tc.client, tc.geoDb, tc.ip)
			assert.Equal(t, tc.want, err, "Should be equal")
		})
	}
}

func Test_Execute(t *testing.T) {
	testCases := []struct {
		name string
		geo  *Geo
		want *pb.GeoResp
	}{
		{
			name: "173.177.164.160",
			geo:  New(client, "geolocation-db.com", "173.177.164.160"),
			want: &pb.GeoResp{
				CountryCode: "CA",
				CountryName: "Canada",
				City:        "Montreal",
				Postal:      "H3K",
				Latitude:    "45.4805",
				Longitude:   "-73.5554",
				IPv4:        "173.177.164.160",
				State:       "Quebec",
			},
		},
		{
			name: "127.0.0.1",
			geo:  New(client, "geolocation-db.com", "127.0.0.1"),
			want: &pb.GeoResp{
				CountryCode: "Not found",
				CountryName: "Not found",
				City:        "Not found",
				Postal:      "Not found",
				Latitude:    "Not found",
				Longitude:   "Not found",
				IPv4:        "Not found",
				State:       "Not found",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.geo.Execute()

			if err := tc.geo.Error(); err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.want, tc.geo.Result(), "Should be equal")
		})
	}
}
