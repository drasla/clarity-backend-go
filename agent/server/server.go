package server

import (
	"agent/services"
	pb "clarity/generated/proto"
	"context"
	"log"
)

type AgentServer struct {
	pb.UnimplementedAgentServiceServer
}

func NewAgentServer() *AgentServer {
	return &AgentServer{}
}

func (s *AgentServer) CreateHosting(_ context.Context, req *pb.CreateHostingRequest) (*pb.CommonResponse, error) {
	log.Printf("[요청] 호스팅 생성 - ID: %s, Domain: %s", req.FtpId, req.Domain)

	err := services.CreateHosting(req.FtpId, req.FtpPw, req.StorageQuotaMb)
	if err != nil {
		log.Printf("[실패] %v", err)
		return &pb.CommonResponse{
			Result:  false,
			Message: err.Error(),
		}, nil
	}

	log.Printf("[성공] 호스팅 생성 완료 - ID: %s", req.FtpId)
	return &pb.CommonResponse{
		Result:  true,
		Message: "호스팅 계정이 성공적으로 생성되었습니다.",
	}, nil
}
func (s *AgentServer) Ping(_ context.Context, _ *pb.Empty) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		Message: "pong",
		Version: "1.0.0",
	}, nil
}
