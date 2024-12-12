package tablet

//
//func (s *TabletServiceServer) RegisterMyself(conn *grpc.ClientConn) (*pb.RegisterTabletResponse, error) {
//	request := &pb.RegisterTabletRequest{
//		TabletAddress: s.TabletAddress,
//	}
//
//	response, err := conn.RegisterTablet(context.Background(), request)
//	if err != nil {
//		log.Fatalf("Could not register tablet server: %v", err)
//	}
//	log.Printf("Registration response: Status Code: %d, Message: %s", response.StatusCode, response.Message)
//}
