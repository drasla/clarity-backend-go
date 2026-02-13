package agent

import (
	pb "clarity/generated/proto"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var AgentClient pb.AgentServiceClient

func InitAgentClient(agentHost string) {
	conn, err := grpc.NewClient(agentHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Agent 연결 실패: %v", err)
	}

	AgentClient = pb.NewAgentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := AgentClient.Ping(ctx, &pb.Empty{})
	if err != nil {
		log.Printf("⚠️ Agent 응답 없음 (나중에 재시도 됩니다): %v", err)
	} else {
		log.Printf("✅ Agent 연결 성공! 버전: %s", resp.Version)
	}
}

func RequestHosting(ftpId, ftpPw, domain string, quotaMB int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req := &pb.CreateHostingRequest{
		FtpId:          ftpId,
		FtpPw:          ftpPw,
		Domain:         domain,
		StorageQuotaMb: quotaMB,
	}

	resp, err := AgentClient.CreateHosting(ctx, req)
	if err != nil {
		return err
	}

	if !resp.Result {
		return fmt.Errorf("agent 처리 실패: %s", resp.Message)
	}

	return nil
}
