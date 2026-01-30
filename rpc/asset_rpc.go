package inner_grpc

import (
	pb "github.com/1821454893q/inner_grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AssetGrpcClient 客户端结构体
type AssetGrpcClient struct {
	grpc pb.AssetGoeInnerClient
}

func NewAssetGrpcClient(grpcAddr string) (*AssetGrpcClient, error) {
	// 建立连接
	conn, err := grpc.NewClient(grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	innerClient := pb.NewAssetGoeInnerClient(conn)
	return &AssetGrpcClient{innerClient}, nil
}

// GrantRewards 发放奖励
func (a *AssetGrpcClient) GrantRewards(bid, uid string, aid int32, rewards map[string]int64) (map[string]int64, error) {
	resp, err := a.grpc.Rewards(ctx, &pb.RewardGrantRequest{
		Bid:     bid,
		Uid:     uid,
		Aid:     aid,
		Rewards: rewards,
	})
	if err != nil {
		return nil, err
	}
	return resp.Updates, nil
}
