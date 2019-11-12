package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/urfave/cli"

	pb "riggs/pb"
)

var cliGeoCmd = cli.Command{
	Name:   "geo",
	Usage:  "Get IP geodata",
	Action: cliGeo,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "ip",
			Usage: "Ip addr to check. Empty for self IP",
			Value: "",
		},
		cli.BoolFlag{
			Name:  "m, morsify",
			Usage: "Return out in Morse code",
		},
		cli.Float64Flag{
			Name:   "t, timeout",
			Usage:  "Operation timeout",
			EnvVar: "RIGGS_OP_TIMEOUT",
			Value:  5,
		},
	},
}

func cliGeo(c *cli.Context) error {
	grpcConn, err := getGrpcConn(c)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer grpcConn.Close()

	// Make a call
	client := pb.NewIPGetGeoClient(grpcConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Float64("timeout"))*time.Second)
	defer cancel()

	ip := c.String("ip")
	r, err := client.GetGeo(ctx, &pb.IPReq{Ip: ip})
	if err != nil {
		return fmt.Errorf("Could not get GeoIP data: %v", err)
	}
	log.Printf("GeoIP data: %s", r.String())

	if c.Bool("morsify") {
		out, err := morsify(c, grpcConn, r.String())
		if err != nil {
			return err
		}

		log.Printf("Original data: %s", out.GetIn())
		log.Printf("Morsified data: %s", out.GetOut())
	}

	return nil
}
