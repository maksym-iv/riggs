package serve

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/urfave/cli"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/peer"

	geo "riggs/geo"
	morse "riggs/morse"
	pb "riggs/pb"
)

var Command = cli.Command{
	Name:   "serve",
	Usage:  "start Riggs gRPC server",
	Action: serve,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "l, listen",
			Usage:  "Server listen `ADDRESS:PORT`",
			EnvVar: "RIGGS_LISTEN",
			Value:  "127.0.0.1:8080",
		},
		cli.StringFlag{
			Name:   "geo-db",
			Usage:  "Geo DB `FQDN`",
			EnvVar: "RIGGS_GEO_DB",
			Value:  "geolocation-db.com",
		},
	},
}

// IPGetGeo service
type geoServer struct {
	pb.UnimplementedIPGetGeoServer
	geoDb  string
	client *http.Client
}

func (s *geoServer) GetGeo(ctx context.Context, in *pb.IPReq) (*pb.GeoResp, error) {
	log.Printf("Received: %v", in.GetIp())

	ip := in.GetIp()
	if ip == "" {
		p, ok := peer.FromContext(ctx)
		if !ok {
			return nil, fmt.Errorf("Could not get client address.")
		}

		ip = p.Addr.(*net.TCPAddr).IP.String()
	}

	g := geo.New(s.client, s.geoDb, ip)
	g.Execute()

	if g.GetError() != nil {
		return nil, g.GetError()
	}

	return g.Result(), nil
}

// Healthcheck service
type morsifyServer struct {
	pb.UnimplementedMorsifyServer
}

func (s *morsifyServer) MorsifyText(ctx context.Context, in *pb.MorseCode) (*pb.MorseCode, error) {
	log.Printf("Received: %v", in.GetIn())
	m, err := morse.ToMorse(in.GetIn())
	in.Out = m
	return in, err
}

// Healthcheck service
type healthServer struct {
	healthpb.UnimplementedHealthServer
}

func (s *healthServer) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	log.Printf("Received health probe: %v", in.GetService())
	return &healthpb.HealthCheckResponse{Status: 1}, nil
}

func newHttpClient() *http.Client {
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:       10 * time.Second,
			KeepAlive:     10 * time.Second,
			FallbackDelay: -1,
		}).DialContext,
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	return &http.Client{Transport: tr}
}

func serve(c *cli.Context) error {
	lis, err := net.Listen("tcp", c.String("listen"))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	httpClient := newHttpClient()

	s := grpc.NewServer()

	pb.RegisterIPGetGeoServer(s, &geoServer{
		geoDb:  c.String("geo-db"),
		client: httpClient,
	})
	pb.RegisterMorsifyServer(s, &morsifyServer{})
	healthpb.RegisterHealthServer(s, &healthServer{})

	log.Printf("Riggs server starting on %s", c.String("listen"))

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}
