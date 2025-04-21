package handlers

import (
	"context"
	"time"

	"github.com/farinmark/Backend-trainee-assignment-spring-2025/internal/service"
	pb "github.com/farinmark/Backend-trainee-assignment-spring-2025/proto"
)

type GRPCServer struct {
	pb.UnimplementedPVZServiceServer
	svc *service.Service
}

func NewGRPCServer(svc *service.Service) *GRPCServer { return &GRPCServer{svc: svc} }

func (g *GRPCServer) ListPVZ(ctx context.Context, _ *pb.Empty) (*pb.PVZList, error) {
	pvs, err := g.svc.ListPVZ(ctx, time.Time{}, time.Now(), 0, 0)
	if err != nil {
		return nil, err
	}
	resp := &pb.PVZList{}
	for _, p := range pvs {
		resp.Pvzs = append(resp.Pvzs, &pb.PVZ{Id: p.ID, City: p.City, RegisteredAt: p.RegisteredAt})
	}
	return resp, nil
}
