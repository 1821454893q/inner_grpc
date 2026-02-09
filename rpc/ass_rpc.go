package inner_grpc

import (
	"context"
	"errors"
	"strconv"

	pb "github.com/1821454893q/inner_grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	ctx = context.Background()
)

type ASSGrpcClient struct {
	grpc pb.ArchiveInnerClient
}

func NewASSGrpcClient(grpcAddr string) (*ASSGrpcClient, error) {
	// 建立连接
	conn, err := grpc.NewClient(grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	innerClient := pb.NewArchiveInnerClient(conn)
	return &ASSGrpcClient{innerClient}, nil
}

// RemoveBlackList 云存档移除黑名单接口
func (a *ASSGrpcClient) RemoveBlackList(uid, bid string, aid, action int) error {
	_, err := a.grpc.DisAbnormal(ctx, &pb.DisAbnormalReq{
		Uid:    uid,
		Bid:    bid,
		Aid:    int32(aid),
		Action: int32(action),
	})
	return err
}

// DestroyUserByCode 根据授权码删除用户的某个存档
func (a *ASSGrpcClient) DestroyUserByCode(uid, aid, bid, verifyCode string) int {
	tAid, err := strconv.Atoi(aid)
	if err != nil {
		return -1
	}
	_, err = a.grpc.WebDestroy(ctx, &pb.WebDestroyReq{
		Uid:  uid,
		Bid:  bid,
		Aid:  int32(tAid),
		Code: verifyCode,
	})
	if err != nil {
		switch pb.Code(status.Convert(err).Code()) {
		case pb.Code_WebDestroyCodeExpiredErr:
			return 1
		case pb.Code_WebDestroyCodeErr:
			return 2
		case pb.Code_WebDestroyCodeNotFoundErr:
			return 0
		}
		return -1
	}

	return 200
}

type AbnormalUser struct {
	Uid    string
	Bid    string
	Aid    string
	Reason int
}

func (a *ASSGrpcClient) BatchAbnormal(list []*AbnormalUser) ([]string, error) {
	if len(list) == 0 {
		return nil, errors.New("param list is nil")
	}

	client, err := a.grpc.BatchAbnuser(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(list))

	for _, user := range list {
		tAid, err := strconv.Atoi(user.Aid)
		if err != nil {
			res = append(res, "-1")
			continue
		}
		err = client.Send(&pb.BatchAbnuserReq{
			Uid:    user.Uid,
			Bid:    user.Bid,
			Aid:    int32(tAid),
			Reason: int32(user.Reason),
		})
		if err != nil {
			res = append(res, "500")
			continue
		}

		resp, err := client.Recv()
		if err != nil {
			res = append(res, "500")
			continue
		}

		if resp.Code == pb.Code_Success {
			res = append(res, "200")
		} else if resp.Code == pb.Code_UserNotExistErr {
			res = append(res, "-2")
		} else if resp.Code == pb.Code_ParamsErr {
			res = append(res, "-1")
		} else {
			res = append(res, "500")
		}
	}

	return res, nil
}

type OldUserInfo struct {
	Did        string
	CreateTime int
	Aid        []int
	Sid        []string
}

type UserInfo struct {
	Did     string       `json:"did"`
	Uid     string       `json:"uid"`
	RegTime int          `json:"regtime"`
	Sis     []SocialInfo `json:"sis"`
}

type SocialInfo struct {
	Sid string `json:"sid"`
	Aid int    `json:"aid"`
}

// ListAidByUid 通过uid bid获取用户所有的aid
func (a *ASSGrpcClient) ListAidByUid(uid, bid string) (*UserInfo, error) {
	userInfo, err := a.grpc.GetUserInfo(ctx, &pb.GetUserInfoReq{Uid: uid, Bid: bid})
	if err != nil {
		return nil, err
	}
	res := &UserInfo{
		Did:     userInfo.Did,
		Uid:     userInfo.Uid,
		RegTime: int(userInfo.RegTime),
	}
	for _, info := range userInfo.List {
		res.Sis = append(res.Sis, SocialInfo{Sid: info.Sid, Aid: int(info.Aid)})
	}
	return res, nil
}

// ListUserAssByDId 通过uid bid获取用户所有的aid
func (a *ASSGrpcClient) ListUserAssByDId(did, bid string) (*UserInfo, error) {
	userInfo, err := a.grpc.GetUserInfoDid(ctx, &pb.GetUserInfoDidReq{Did: did, Bid: bid})
	if err != nil {
		return nil, err
	}
	res := &UserInfo{
		Did:     userInfo.Did,
		Uid:     userInfo.Uid,
		RegTime: int(userInfo.RegTime),
	}
	for _, info := range userInfo.List {
		res.Sis = append(res.Sis, SocialInfo{Sid: info.Sid, Aid: int(info.Aid)})
	}
	return res, nil
}

// UpdateArchive 修改用户云存档
func (a *ASSGrpcClient) UpdateArchive(uid, bid, key, value string, aid int) error {
	_, err := a.grpc.ModifyArchive(ctx, &pb.ModifyArchiveReq{
		Uid:   uid,
		Bid:   bid,
		Aid:   int32(aid),
		Key:   key,
		Value: value,
	})
	if err != nil {
		return err
	}
	return nil
}

type UserKey struct {
	Uid string `json:"uid"`
	Bid string `json:"bid"`
	Aid int    `json:"aid"`
}

// ArchiveUserInfoByKey 获取存档userKey
func (a *ASSGrpcClient) ArchiveUserInfoByKey(key string) (*UserKey, error) {
	resp, err := a.grpc.GetUserInfoByKey(ctx, &pb.GetUserInfoByKeyReq{Key: key})
	if err != nil {
		return nil, err
	}
	return &UserKey{
		Uid: resp.Uid,
		Bid: resp.Bid,
		Aid: int(resp.Aid),
	}, nil
}

func (a *ASSGrpcClient) ArchiveIpForbid(bid, uid string) (*pb.QueryIPForbidUserResp, error) {
	resp, err := a.grpc.QueryIPForbidUser(ctx, &pb.QueryIPForbidUserReq{
		Bid: bid,
		Uid: uid,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *ASSGrpcClient) ClearArchive(bid, uid string, aid int) error {
	if bid == "" || uid == "" || aid < 0 || aid > 9 {
		return errors.New("param errors")
	}
	_, err := a.grpc.ClearArchive(ctx, &pb.ClearArchiveReq{
		UserInfo: &pb.BaseInfo{
			Bid: bid,
			Uid: uid,
			Aid: int32(aid),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *ASSGrpcClient) DeleteArchive(bid, uid string, aid int) error {
	if bid == "" || uid == "" || aid < 0 || aid > 9 {
		return errors.New("param errors")
	}
	_, err := a.grpc.DeleteArchive(ctx, &pb.DeleteArchiveReq{
		UserInfo: &pb.BaseInfo{
			Bid: bid,
			Uid: uid,
			Aid: int32(aid),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *ASSGrpcClient) Rewards(bid, uid string, aid int, rewards map[string]int64) (map[string]int64, error) {
	resp, err := a.grpc.Rewards(ctx, &pb.RewardGrantRequest{
		Bid:     bid,
		Uid:     uid,
		Aid:     int32(aid),
		Rewards: rewards,
	})
	if err != nil {
		return nil, err
	}
	return resp.Updates, nil
}
