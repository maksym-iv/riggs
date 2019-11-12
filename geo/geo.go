package geo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	pb "riggs/pb"
)

type Geo struct {
	client  *http.Client
	ip      string
	geoDb   string
	geoResp *pb.GeoResp
	err     error
}

func getGeo(client *http.Client, geoDb string, ip string) (*http.Response, error) {
	url := fmt.Sprintf("http://%s/json/%s", geoDb, ip)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	return resp, err
}

func New(client *http.Client, geoDb string, ip string) *Geo {
	return &Geo{
		client:  client,
		ip:      ip,
		geoDb:   geoDb,
		geoResp: &pb.GeoResp{},
	}
}

func (g *Geo) Execute() {
	r, err := getGeo(g.client, g.geoDb, g.ip)
	if err != nil {
		g.err = err
		return
	}

	if r.StatusCode != 200 {
		g.err = fmt.Errorf("Got non 200 status code. %s. Exiting.", r.Status)
		return
	}
	defer r.Body.Close()

	// d := json.NewDecoder(r.Body)
	// g.err = d.Decode(g.geoResp)
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r.Body); err != nil {
		g.err = err
		return
	}

	g.err = g.UnmarshalJSON(buf.Bytes())
}

// UnmarshalJSON - dummy unmarshal. It is a hack. It is here to workaround not static field type in geoloc API
func (g *Geo) UnmarshalJSON(b []byte) error {
	var v map[string]interface{}

	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	for k, val := range v {
		switch val.(type) {
		case float64:
			v[k] = fmt.Sprintf("%v", val)
		}
	}

	jsonbody, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonbody, g.geoResp)
}

func (g *Geo) GetError() error {
	return g.err
}

func (g *Geo) Error() string {
	return g.err.Error()
}

func (g *Geo) Result() *pb.GeoResp {
	return g.geoResp
}
