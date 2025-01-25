package tablet

import (
	"context"
	ipb "final/proto/internal-api"
	"log"
)

func (s *TabletServiceServer) RegisterMyself() error {
	request := &ipb.RegisterTabletRequest{
		TabletAddress: s.TabletAddress,
		MaxTableCnt:   int32(s.MaxTableCnt),
	}

	_, err := (*s.MasterClient).RegisterTablet(context.Background(), request)
	if err != nil {
		return err
	}

	return nil
}

func (s *TabletServiceServer) UnRegisterMyself() error {
	request := &ipb.UnregisterTabletRequest{
		TabletAddress: s.TabletAddress,
	}

	_, err := (*s.MasterClient).UnregisterTablet(context.Background(), request)
	if err != nil {
		log.Fatal("fail to unregister")
		return err
	}

	return nil
}

func (s *TabletServiceServer) Heartbeat(ctx context.Context, request *ipb.HeartbeatRequest) (*ipb.HeartbeatResponse, error) {
	return &ipb.HeartbeatResponse{
		Success: true,
	}, nil
}
