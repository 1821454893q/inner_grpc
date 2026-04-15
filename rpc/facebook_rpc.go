package inner_grpc

import (
	pb "github.com/1821454893q/inner_grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FacebookGrpcClient struct {
	grpc pb.FacebookSystemClient
}

func NewFacebookGrpcClient(grpcAddr string) (*FacebookGrpcClient, error) {
	// 建立连接
	conn, err := grpc.NewClient(grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	innerClient := pb.NewFacebookSystemClient(conn)
	return &FacebookGrpcClient{innerClient}, nil
}

// ListRankInfo 获取排位赛数据
func (f *FacebookGrpcClient) ListRankInfo(appKey string, keys []string, pageNumber, pageSize int) (*pb.ListRankInfoResp, error) {
	resp, err := f.grpc.ListRankInfo(ctx, &pb.ListRankInfoReq{
		AppKey:     appKey,
		Keys:       keys,
		PageNumber: int32(pageNumber),
		PageSize:   int32(pageSize),
	})
	return resp, err
}

// UpdatePlayerScore 修改用户分数
func (f *FacebookGrpcClient) UpdatePlayerScore(appKey, key, uid string, score int64, name, avatar string) error {
	_, err := f.grpc.UpdatePlayerScore(ctx, &pb.UpdatePlayerScoreReq{
		AppKey: appKey,
		Key:    key,
		Uid:    uid,
		Score:  score,
		Name:   name,
		Avatar: avatar,
	})
	return err
}

// DeleteRank 删除排位赛
func (f *FacebookGrpcClient) DeleteRank(appKey, key string) error {
	_, err := f.grpc.DeleteRank(ctx, &pb.DeleteRankReq{
		AppKey: appKey,
		Key:    key,
	})
	return err
}
