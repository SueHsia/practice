package main

import (
	"context"
	"database/sql"
	"rpc"

	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/apache/thrift/lib/go/thrift"
)

type info struct {
	Goodsid   int64
	Goodsname string
	Address   string
	Phone     string
	Des       string
}

type rpcServer struct{}

func (serv *rpcServer) DoQuery(ctx context.Context, GoodsId int64, UserName string, Token string) (*rpc.DataRes, error) {
	dict := info{}
	temGoods := &rpc.LostGoodInfo{
		GoodsId:   dict.Goodsid,
		GoodsName: dict.Goodsname,
		Address:   dict.Address,
		Phone:     dict.Phone,
		Des:       dict.Des,
	}
	db, _ := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/lost_and_found?charset=utf8")
	db.QueryRow("select id,name,address,phone,des from lost_goods where id = ?", int(GoodsId)).Scan(&temGoods.GoodsId, &temGoods.GoodsName, &temGoods.Address, &temGoods.Phone, &temGoods.Des)
	res := &rpc.DataRes{
		Code: 1,
		Msg:  "success",
		Gift: temGoods,
	}
	return res, nil
}

func (serv *rpcServer) GoodsList(ctx context.Context, UserId int64, UserName string, Token string) ([]*rpc.LostGoodInfo, error) {
	db, _ := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/lost_and_found?charset=utf8")
	dict := info{}
	temGoods := &rpc.LostGoodInfo{
		GoodsId:   dict.Goodsid,
		GoodsName: dict.Goodsname,
		Address:   dict.Address,
		Phone:     dict.Phone,
		Des:       dict.Des,
	}

	var count int
	db.QueryRow("select count(*) from lost_goods where user_id=?", UserId).Scan(&count)
	rData := make([]*rpc.LostGoodInfo, count)
	rows, _ := db.Query("select id,name,address,phone,des from lost_goods where user_id=?", UserId)
	i := 0
	for rows.Next() {
		rows.Scan(&temGoods.GoodsId, &temGoods.GoodsName, &temGoods.Address, &temGoods.Phone, &temGoods.Des)
		rData[i] = temGoods
		i += 1
	}
	return rData, nil
}

func main() {
	transport, err := thrift.NewTServerSocket(":9090")
	if err != nil {
		panic(err)
	}
	// var serv rpc.LuckyService = &rpcServer{}
	serv := &rpcServer{}
	processor := rpc.NewLuckyServiceProcessor(serv)
	transportFactory := thrift.NewTBufferedTransportFactory(8192)
	protocolFactory := thrift.NewTCompactProtocolFactory()
	server := thrift.NewTSimpleServer4(
		processor,
		transport,
		transportFactory,
		protocolFactory,
	)
	if err = server.Serve(); err != nil {
		panic(err)
	}
}
