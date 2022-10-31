package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	pb "github.com/tonx22/gocloudcamp/pb"
	Models "github.com/tonx22/gocloudcamp/pkg/models"
	"github.com/tonx22/gocloudcamp/pkg/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type server struct {
	pb.UnimplementedConfigSvcServer
	service service.ConfigService
}

func (s *server) SetConfig(ctx context.Context, in *pb.ConfigRequest) (*pb.ConfigRequest, error) {
	rsp, err := s.processGRPCRequest(ctx, in, "setConfig")
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (s *server) GetConfig(ctx context.Context, in *pb.ConfigRequest) (*pb.ConfigRequest, error) {
	rsp, err := s.processGRPCRequest(ctx, in, "getConfig")
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (s *server) UpdConfig(ctx context.Context, in *pb.ConfigRequest) (*pb.ConfigRequest, error) {
	rsp, err := s.processGRPCRequest(ctx, in, "updConfig")
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (s *server) DelConfig(ctx context.Context, in *pb.ConfigRequest) (*pb.ConfigRequest, error) {
	rsp, err := s.processGRPCRequest(ctx, in, "delConfig")
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (s *server) processGRPCRequest(ctx context.Context, in *pb.ConfigRequest, method string) (*pb.ConfigRequest, error) {
	req, err := decodeGRPCRequest(ctx, in)
	if err != nil {
		return nil, err
	}

	var resp *Models.ConfigRequest
	svc := s.service

	switch method {
	case "setConfig":
		resp, err = svc.SetConfig(ctx, req)
	case "getConfig":
		resp, err = svc.GetConfig(ctx, req)
	case "updConfig":
		resp, err = svc.UpdConfig(ctx, req)
	case "delConfig":
		resp, err = svc.DelConfig(ctx, req)
	default:
		return nil, errors.New("unknown method")
	}
	if err != nil {
		return nil, err
	}

	rsp, err := encodeGRPCResponse(ctx, resp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func decodeGRPCRequest(_ context.Context, grpcReq interface{}) (*Models.ConfigRequest, error) {
	r := grpcReq.(*pb.ConfigRequest)
	req := Models.ConfigRequest{Service: r.Service, Version: int(r.Version), Used: r.Used}
	err := json.Unmarshal(r.Data, &req.Data)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func encodeGRPCResponse(_ context.Context, response interface{}) (*pb.ConfigRequest, error) {
	r := response.(*Models.ConfigRequest)
	resp := pb.ConfigRequest{Service: r.Service, Version: int32(r.Version), Used: r.Used}
	data, err := json.Marshal(r.Data)
	if err != nil {
		return nil, err
	}
	resp.Data = data
	return &resp, nil
}

func StartNewGRPCServer(s interface{}, grpcPort int) error {
	svc := s.(service.ConfigService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterConfigSvcServer(grpcServer, &server{service: svc})

	ch := make(chan error)
	go func() {
		ch <- grpcServer.Serve(lis)
	}()

	var e error
	select {
	case e = <-ch:
		return e
	case <-time.After(time.Second * 1):
	}
	log.Printf("GRPC server listening at %v", lis.Addr())
	return nil
}
