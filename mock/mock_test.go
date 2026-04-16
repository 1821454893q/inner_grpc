package mock

import (
	"context"
	"fmt"
	"testing"
	"time"

	inner_grpc "github.com/1821454893q/inner_grpc/rpc"
)

func TestAss(t *testing.T) {
	c, err := inner_grpc.NewASSGrpcClient("192.168.10.5:9988")
	if err != nil {
		t.Error(err)
	}

	resp, err := c.Rewards("com.yifan.ass", "123456", 0, map[string][]byte{"coin": []byte("100")})
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)

	// resp, err = c.Rewards("com.yifan.ass", "123456", 0, map[string]float32{"coin": 100})
	// if err != nil {
	// 	t.Error(err)
	// }
	// t.Log(resp)
}

func TestFacebook(t *testing.T) {
	c, err := inner_grpc.NewFacebookGrpcClient("localhost:46125")
	if err != nil {
		t.Error(err)
	}

	resp, err := c.ListRankInfo("com.yifan.ass", []string{}, 0, 20)
	if err != nil {
		t.Error(err)
	}

	// values: key -> json 内容
	fmt.Println("=== Values ===")
	for k, v := range resp.Values {
		fmt.Printf("  [%s] => %s\n", k, v)
	}

	// player: key -> 玩家列表
	fmt.Println("\n=== Players ===")
	for k, playerList := range resp.Player {
		fmt.Printf("  [%s]\n", k)
		for _, p := range playerList.List {
			fmt.Printf("    uid=%-20s name=%-20s score=%-8d updated=%s\n",
				p.Uid, p.Name, p.Score,
				time.Unix(p.UpdateTimestamp, 0).Format("2006-01-02 15:04:05"),
			)
		}
	}

	// create_timestamp: key -> 创建时间
	fmt.Println("\n=== CreateTimestamp ===")
	for k, ts := range resp.CreateTimestamp {
		fmt.Printf("  [%s] => %s\n", k, time.Unix(ts, 0).Format("2006-01-02 15:04:05"))
	}
}

func TestFacebookScore(t *testing.T) {
	c, err := inner_grpc.NewFacebookGrpcClient("localhost:46125")
	if err != nil {
		t.Error(err)
	}

	err = c.UpdatePlayerScore("com.yifan.ass", "25361181180210862", "24959949823707476", 3200, "", "icon1")
	if err != nil {
		t.Error(err)
	}

}

func TestRankGuildRankList(t *testing.T) {
	c, err := inner_grpc.NewRankGuildGrpcClient("192.168.10.96:51415")
	if err != nil {
		t.Error(err)
	}

	resp, err := c.GuildRankList(context.Background(), "com.yifan.ass", "guild_test_geo", 2, 5, "")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("total=%d\n", resp.Total)
	for _, item := range resp.List {
		fmt.Printf("  guild_id=%s rank=%d score=%d ranking_key=%s\n", item.GuildId, item.Rank, item.Score, item.RankingKey)
	}
}

func TestRankGuildMemberList(t *testing.T) {
	c, err := inner_grpc.NewRankGuildGrpcClient("192.168.10.96:51415")
	if err != nil {
		t.Error(err)
		return
	}

	resp, err := c.GuildMemberRankList(context.Background(), "com.yifan.ass", "guild_test_geo", "5970", "", 1, 2)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("total=%d\n", resp.Total)
	for _, item := range resp.List {
		fmt.Printf("  user_id=%s score=%d rank=%s ranking_key=%s\n", item.UserId, item.Score, item.Rank, item.RankingKey)
	}
}
