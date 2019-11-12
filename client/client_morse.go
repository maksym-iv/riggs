package client

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/urfave/cli"
	"google.golang.org/grpc"

	pb "riggs/pb"
)

var cliMorseCmd = cli.Command{
	Name:   "morsify",
	Usage:  "Transform test to morse code",
	Action: cliMorse,
	Flags: []cli.Flag{
		cli.Float64Flag{
			Name:   "t, timeout",
			Usage:  "Operation timeout",
			EnvVar: "RIGGS_OP_TIMEOUT",
			Value:  5,
		},
	},
}

func cliMorse(c *cli.Context) error {
	grpcConn, err := getGrpcConn(c)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer grpcConn.Close()

	// Make a call
	text := strings.Join(c.Args(), " ")
	out, err := morsify(c, grpcConn, text)
	if err != nil {
		return err
	}

	log.Printf("Original data: %s", out.GetIn())
	log.Printf("Morsified data: %s", out.GetOut())
	return nil
}

func morsify(c *cli.Context, conn *grpc.ClientConn, text string) (*pb.MorseCode, error) {
	client := pb.NewMorsifyClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Float64("timeout"))*time.Second)
	defer cancel()

	r, err := client.MorsifyText(ctx, &pb.MorseCode{In: text})
	if err != nil {
		return nil, fmt.Errorf("Could not encode to Morse: %v", err)
	}

	return r, nil
}
