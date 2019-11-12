package client

import (
	"context"
	"fmt"
	"time"

	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

var Command = cli.Command{
	Name:  "cli",
	Usage: "cli operations",
	Subcommands: []cli.Command{
		cliGeoCmd,
		cliMorseCmd,
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "riggs-addr",
			Usage:  "Riggs server `ADDRESS:PORT`",
			EnvVar: "RIGGS_ADDR",
			Value:  "127.0.0.1:8080",
		},
		cli.IntFlag{
			Name:   "connTimeout",
			Usage:  "Riggs server connection timeout",
			EnvVar: "RIGGS_CONN_TIMEOUT",
			Value:  5,
			Hidden: true,
		},
	},
}

func getGrpcConn(c *cli.Context) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.GlobalInt("connTimeout"))*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, c.GlobalString("riggs-addr"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		err := fmt.Errorf("did not connect to grpc server for %d seconds: %v", c.GlobalInt("connTimeout"), err)
		return nil, err
	}
	return conn, nil
}
