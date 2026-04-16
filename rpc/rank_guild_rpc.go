package inner_grpc

import (
	"context"

	pb "github.com/1821454893q/inner_grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RankGuildGrpcClient struct {
	grpc pb.RankGuildSystemClient
}

func NewRankGuildGrpcClient(grpcAddr string) (*RankGuildGrpcClient, error) {
	conn, err := grpc.NewClient(grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	innerClient := pb.NewRankGuildSystemClient(conn)
	return &RankGuildGrpcClient{innerClient}, nil
}

func (r *RankGuildGrpcClient) GuildRankList(ctx context.Context, appKey, rankingId string, page, pageSize int64, guildId string) (*pb.GuildRankListResp, error) {
	resp, err := r.grpc.GuildRankList(ctx, &pb.GuildRankListReq{
		AppKey:    appKey,
		RankingId: rankingId,
		Page:      page,
		PageSize:  pageSize,
		GuildId:   guildId,
	})
	return resp, err
}

func (r *RankGuildGrpcClient) GuildMemberRankList(ctx context.Context, appKey, rankingId, guildId, userId string, page, pageSize int64) (*pb.GuildMemberRankListResp, error) {
	resp, err := r.grpc.GuildMemberRankList(ctx, &pb.GuildMemberRankListReq{
		AppKey:    appKey,
		RankingId: rankingId,
		GuildId:   guildId,
		UserId:    userId,
		Page:      page,
		PageSize:  pageSize,
	})
	return resp, err
}
