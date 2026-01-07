package inner_grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RankGrpcClient struct {
	grpc RankSystemClient
}

func NewRankGrpcClient(grpcAddr string) (*ASSGrpcClient, error) {
	// 建立连接
	conn, err := grpc.NewClient(grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	innerClient := NewArchiveInnerClient(conn)
	return &ASSGrpcClient{innerClient}, nil
}

// KickUser 踢人
func (r *RankGrpcClient) KickUser(bid, uid, rankId string, kickType int) error {
	_, err := r.grpc.KickUser(ctx, &KickUserReq{
		Bid:      bid,
		RankId:   rankId,
		Uid:      uid,
		KickType: int32(kickType),
	})
	return err
}

// GetUserRank 获取用户排名
func (r *RankGrpcClient) GetUserRank(bid, uid, rankId string) (*GetUserRankResp, error) {
	resp, err := r.grpc.GetUserRank(ctx, &GetUserRankReq{
		Uid:    uid,
		Bid:    bid,
		RankId: rankId,
	})
	return resp, err
}
