package main

import (
	"context"
	"fmt"

	"github.com/mloncode/sonic"
	"github.com/kelseyhightower/envconfig"
	"github.com/rakyll/portmidi"
	"github.com/src-d/lookout"
	"gopkg.in/src-d/lookout-sdk.v0/pb"
	"github.com/src-d/lookout/util/grpchelper"
	"google.golang.org/grpc"
	log "gopkg.in/src-d/go-log.v1"
)

type config struct {
	Host           string `envconfig:"HOST" default:"0.0.0.0"`
	Port           int    `envconfig:"PORT" default:"2001"`
	DataServiceURL string `envconfig:"DATA_SERVICE_URL" default:"ipv4://localhost:10301"`
}

func main() {
	var conf config
	envconfig.MustProcess("SONIC", &conf)
	log.Infof("Starting...")

	grpcAddr, err := pb.ToGoGrpcAddress(conf.DataServiceURL)
	if err != nil {
		log.Errorf(err, "failed to parse DataService addres %s", conf.DataServiceURL)
		return
	}

	conn, err := grpchelper.DialContext(
		context.Background(),
		grpcAddr,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.FailFast(false)),
	)
	if err != nil {
		log.Errorf(err, "cannot create connection to DataService %s", grpcAddr)
		return
	}

	if err := portmidi.Initialize(); err != nil {
		log.Errorf(err, "can't initializer portmidi")
		return
	}
	defer portmidi.Terminate()

	if portmidi.CountDevices() == 0 {
		log.Errorf(nil, "no midi devices")
		return
	}

	analyzer := &sonic.Analyzer{
		DataClient: lookout.NewDataClient(conn),
		DeviceID:   portmidi.DefaultOutputDeviceID(),
	}

	server := grpchelper.NewServer()
	lookout.RegisterAnalyzerServer(server, analyzer)

	analyzerURL := fmt.Sprintf("ipv4://%s:%d", conf.Host, conf.Port)
	lis, err := pb.Listen(analyzerURL)
	if err != nil {
		log.Errorf(err, "failed to start analyzer gRPC server on %s", analyzerURL)
		return
	}

	log.Infof("server has started on '%s'", analyzerURL)
	err = server.Serve(lis)
	if err != nil {
		log.Errorf(err, "gRPC server failed listening on %v", lis)
	}
	return
}
