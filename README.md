# 项目统一grpc调用封装

```
    // 初始化客户端
    grpc := ass_grpc.NewASSGrpcClient("192.168.10.5:443")
    
    // 执行API 移除黑名单
    err := grpc.RemoveBlackList("uid", "bid" , 0, 0)
    if err != nil {
        // 异常处理
    }
```