package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	pb "github.com/tonx22/gocloudcamp/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

type ConfigService interface {
	SetConfig(ctx context.Context, r ConfigRequest) (*ConfigRequest, error)
	GetConfig(ctx context.Context, r ConfigRequest) (*ConfigRequest, error)
	UpdConfig(ctx context.Context, r ConfigRequest) (*ConfigRequest, error)
	DelConfig(ctx context.Context, r ConfigRequest) (*ConfigRequest, error)
}

type configService struct {
	GRPCClient pb.ConfigSvcClient
}

func NewGRPCClient() (*configService, error) {
	defaultHost, ok := os.LookupEnv("GRPC_HOST")
	if !ok {
		defaultHost = "localhost"
	}
	defaultPort, ok := os.LookupEnv("GRPC_PORT")
	if !ok {
		defaultPort = "50051"
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	serverAddr := fmt.Sprintf("%s:%s", defaultHost, defaultPort)
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("fail to dial: %s", err))
	}
	svc := configService{GRPCClient: pb.NewConfigSvcClient(conn)}
	return &svc, nil
}

func (svc configService) SetConfig(ctx context.Context, r ConfigRequest) (*ConfigRequest, error) {
	res, err := svc.processGRPCRequest(ctx, r, "setConfig")
	if err != nil {
		return nil, err
	}
	return res, err
}

func (svc configService) GetConfig(ctx context.Context, r ConfigRequest) (*ConfigRequest, error) {
	res, err := svc.processGRPCRequest(ctx, r, "getConfig")
	if err != nil {
		return nil, err
	}
	return res, err
}

func (svc configService) UpdConfig(ctx context.Context, r ConfigRequest) (*ConfigRequest, error) {
	res, err := svc.processGRPCRequest(ctx, r, "updConfig")
	if err != nil {
		return nil, err
	}
	return res, err
}

func (svc configService) DelConfig(ctx context.Context, r ConfigRequest) (*ConfigRequest, error) {
	res, err := svc.processGRPCRequest(ctx, r, "delConfig")
	if err != nil {
		return nil, err
	}
	return res, err
}

func (svc configService) processGRPCRequest(ctx context.Context, r ConfigRequest, method string) (*ConfigRequest, error) {
	req, err := encodeGRPCRequest(ctx, r)
	if err != nil {
		return nil, err
	}

	var resp *pb.ConfigRequest
	switch method {
	case "setConfig":
		resp, err = svc.GRPCClient.SetConfig(context.Background(), req)
	case "getConfig":
		resp, err = svc.GRPCClient.GetConfig(context.Background(), req)
	case "updConfig":
		resp, err = svc.GRPCClient.UpdConfig(context.Background(), req)
	case "delConfig":
		resp, err = svc.GRPCClient.DelConfig(context.Background(), req)
	default:
		return nil, errors.New("unknown method")
	}
	if err != nil {
		return nil, err
	}

	res, err := decodeGRPCResponse(context.Background(), resp)
	if err != nil {
		return nil, err
	}
	return res, err
}

func encodeGRPCRequest(_ context.Context, request interface{}) (*pb.ConfigRequest, error) {
	r := request.(ConfigRequest)
	req := pb.ConfigRequest{Service: r.Service, Version: r.Version, Used: r.Used}
	req.Data, _ = json.Marshal(r.Data)
	return &req, nil
}

func decodeGRPCResponse(_ context.Context, grpcResp interface{}) (*ConfigRequest, error) {
	r := grpcResp.(*pb.ConfigRequest)
	resp := ConfigRequest{Service: r.Service, Version: r.Version, Used: r.Used}
	err := json.Unmarshal(r.Data, &resp.Data)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ConfigRequest struct {
	Service string
	Data    map[string]interface{}
	Version int32
	Used    bool
}
