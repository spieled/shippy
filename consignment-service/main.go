package main

import (
    "log"
    "net"
    pb "github.com/spieled/shippy/consignment-service/proto/consignment"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

const (
    port = ":50051"
)

type IRepository interface {
    Create(*pb.Consignment) (*pb.Consignment, error)
    GetAll() ([]*pb.Consignment, error)
}

type Repository struct {
    consignments []*pb.Consignment
}

func (r *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
    r.consignments = append(r.consignments, consignment)
    return consignment, nil
}

func (r *Repository) GetAll() ([]*pb.Consignment, error) {
    return r.consignments, nil
}

type service struct {
    repo IRepository
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
    consignment, err := s.repo.Create(req)
    if err != nil {
        return nil, err
    }
    return &pb.Response{Created: true, Consignment: consignment}, nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
    consignments, err := s.repo.GetAll()
    if err != nil {
        return nil, err 
    }
    return &pb.Response{Consignments: consignments}, nil
}

func main() {
    repo := &Repository{}
    lsr, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen port %s. %v", port, err)
    }
    log.Printf("listen at port %s", port)
    server := grpc.NewServer()
    pb.RegisterShippingServiceServer(server, &service{repo})
    reflection.Register(server)
    if err := server.Serve(lsr); err != nil {
        log.Fatalf("failed to serve. %v", err)
    }
}
