package inner_grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"strconv"
)

var (
	ctx = context.Background()
)

type ASSGrpcClient struct {
	grpc ArchiveInnerClient
}

func NewASSGrpcClient(grpcAddr string) (*ASSGrpcClient, error) {
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

// RemoveBlackList 云存档移除黑名单接口
func (a *ASSGrpcClient) RemoveBlackList(uid, bid string, aid, action int) error {
	_, err := a.grpc.DisAbnormal(ctx, &DisAbnormalReq{
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
	_, err = a.grpc.WebDestroy(ctx, &WebDestroyReq{
		Uid:  uid,
		Bid:  bid,
		Aid:  int32(tAid),
		Code: verifyCode,
	})
	if err != nil {
		switch Code(status.Convert(err).Code()) {
		case Code_WebDestroyCodeExpiredErr:
			return 1
		case Code_WebDestroyCodeErr:
			return 2
		case Code_WebDestroyCodeNotFoundErr:
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
	if list == nil || len(list) == 0 {
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
		err = client.Send(&BatchAbnuserReq{
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

		if resp.Code == Code_Success {
			res = append(res, "200")
		} else if resp.Code == Code_UserNotExistErr {
			res = append(res, "-2")
		} else if resp.Code == Code_ParamsErr {
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
	userInfo, err := a.grpc.GetUserInfo(ctx, &GetUserInfoReq{Uid: uid, Bid: bid})
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
	userInfo, err := a.grpc.GetUserInfoDid(ctx, &GetUserInfoDidReq{Did: did, Bid: bid})
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
	_, err := a.grpc.ModifyArchive(ctx, &ModifyArchiveReq{
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
	resp, err := a.grpc.GetUserInfoByKey(ctx, &GetUserInfoByKeyReq{Key: key})
	if err != nil {
		return nil, err
	}
	return &UserKey{
		Uid: resp.Uid,
		Bid: resp.Bid,
		Aid: int(resp.Aid),
	}, nil
}

func (a *ASSGrpcClient) ArchiveIpForbid(bid, uid string) (*QueryIPForbidUserResp, error) {
	resp, err := a.grpc.QueryIPForbidUser(ctx, &QueryIPForbidUserReq{
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
	_, err := a.grpc.ClearArchive(ctx, &ClearArchiveReq{
		UserInfo: &BaseInfo{
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
	_, err := a.grpc.DeleteArchive(ctx, &DeleteArchiveReq{
		UserInfo: &BaseInfo{
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
