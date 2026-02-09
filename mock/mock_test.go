package mock

import (
	"testing"

	inner_grpc "github.com/1821454893q/inner_grpc/rpc"
)

func TestAss(t *testing.T) {
	c, err := inner_grpc.NewASSGrpcClient("192.168.10.5:9988")
	if err != nil {
		t.Error(err)
	}

	resp, err := c.Rewards("com.yifan.ass", "123456", 0, map[string]int64{"coin": 100})
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
