package inner_grpc

import (
	pb "github.com/1821454893q/inner_grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RankGrpcClient struct {
	grpc pb.RankSystemClient
}

func NewRankGrpcClient(grpcAddr string) (*RankGrpcClient, error) {
	// 建立连接
	conn, err := grpc.NewClient(grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	innerClient := pb.NewRankSystemClient(conn)
	return &RankGrpcClient{innerClient}, nil
}

// KickUser 踢人
func (r *RankGrpcClient) KickUser(bid, uid, rankId string, kickType int) error {
	_, err := r.grpc.KickUser(ctx, &pb.KickUserReq{
		Bid:      bid,
		RankId:   rankId,
		Uid:      uid,
		KickType: int32(kickType),
	})
	return err
}

// GetUserRank 获取用户排名
func (r *RankGrpcClient) GetUserRank(bid, uid, rankId string) (*pb.GetUserRankResp, error) {
	resp, err := r.grpc.GetUserRank(ctx, &pb.GetUserRankReq{
		Uid:    uid,
		Bid:    bid,
		RankId: rankId,
	})
	return resp, err
}
