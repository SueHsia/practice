namespace go rpc
namespace php rpc
namespace py rpc

# 失物详情

struct LostGoodInfo{
    1:i64    GoodsId    =0
    2:string GoodsName  =""
    3:string Address    =""
    4:string Phone      =""
    5:string Des        =""
}

# 返回值
struct DataRes{
    1:i64           Code
    2:string        Msg
    3:LostGoodInfo  Gift
}

# 服务接口
service LuckyService {
    # 抽奖的方法
    DataRes DoQuery(1:i64 GoodsId, 2:string UserName, 3:string Token),
    list<LostGoodInfo> GoodsList(1:i64 UserId, 2:string UserName, 3:string Token),
}
