package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/ipfs/go-ipfs-api"
	"github.com/scryinfo/iscap_demo/src/sdk"
	"github.com/scryinfo/iscap_demo/src/sdk/core/chainoperations"
	"github.com/scryinfo/iscap_demo/src/sdk/core/ethereum/events"
	"github.com/scryinfo/iscap_demo/src/sdk/scryclient"
	cif "github.com/scryinfo/iscap_demo/src/sdk/scryclient/chaininterfacewrapper"
	"github.com/scryinfo/iscap_demo/src/sdk/util/accounts"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type M map[string]string
type Obj map[string]interface{}

// -------------------- 服务接口 -------------------- //

var (
	accounts_info = [5] string {"0x1eaf7b4bcb0d87c405da62536b3a3b385e1712e0",
		"0xeab13565ab4c9a45f5fd105969468ea932f9d039",
		"0xea31ba64798804034417a2e55419dbbdef17ffb7",
		"0x3777fc8e99190980a72feb85d7bbb741e3e7f1c4",
		"0x045346db8a80519ad793067321cbfabd7c0f5bd0"}
	password = "888888"
	seller = scryclient.NewScryClient(accounts_info[0])
	buyer = scryclient.NewScryClient(accounts_info[4])
	publish_des_list = []Obj{}
	publish_proof_map = Obj{}
	protocolContractAddr = "0x3c4d26e916d79fc3fc925027a79612012462f691"
	tokenContractAddr = "0x5c29f42d640ee25f080cdc648641e8e358459ddc"
	transactionId = big.NewInt(0)
	transactionId_temp = int64(0)
	metaDataIdWithSeller = ""
	downloadInfo = M{}
	ipfsNodeAddr = "/ip4/172.16.0.203/tcp/5001"
)



/**
 * 1. 初始化
 */
func initApp()  {
	wd, _ := os.Getwd()
	err := sdk.Init(
		"http://172.16.0.201:8545/",
		"172.16.0.202:48080",
		protocolContractAddr,
		tokenContractAddr,
		0,
		"/ip4/172.16.0.203/tcp/5001",
		wd + "./testconsole.log",
		"scryapp_goodbay")

	if err != nil {
		fmt.Println("sdk init error")
	}

	// seller和buyer监听事件
	seller.SubscribeEvent("DataPublish",onPublish)
	seller.SubscribeEvent("Buy",onPurchase)
	seller.SubscribeEvent("TransactionCreate",onSellerTransactionCreate)
	seller.SubscribeEvent("TransactionClose",onClose)

	buyer.SubscribeEvent("Approval",onApprovalBuyerTransfer)
	buyer.SubscribeEvent("TransactionCreate",onBuyerTransactionCreate)
	buyer.SubscribeEvent("ReadyForDownload",onReadyForDownload)
	buyer.SubscribeEvent("TransactionClose",onClose)
}

/**
 * 2. seller和buyer事件监听
 */
func onPublish(event events.Event) bool  {
	fmt.Println("seller: onPublish: ", event)
	return true
}

// 买家确定购买数据，买家被通知
func onPurchase(event events.Event) bool  {
	metaDataIdEncWithSeller := event.Data.Get("metaDataIdEncSeller").([]byte)

	metaDataIdEncWithSeller,_ = accounts.GetAMInstance().Decrypt(metaDataIdEncWithSeller,seller.Account.Address,password)
	metaDataIdWithSeller = string(metaDataIdEncWithSeller[:])

	fmt.Println("seller: onPurchase: ", event)
	return true
}

func onSellerTransactionCreate(event events.Event) bool  {
	fmt.Println("seller: TransactionCreate: ", event)
	return true
}

func onClose(event events.Event) bool  {
	fmt.Println("seller: TransactionClose: ", event)
	return true
}

func onApprovalBuyerTransfer(event events.Event) bool  {
	fmt.Println("buyer: onApprovalBuyerTransfer: ", event)
	return true
}

// 买家缴纳押金
func onBuyerTransactionCreate(event events.Event) bool  {
	fmt.Println("buyer: onBuyerTransactionCreate: ", event)
	transactionId = event.Data.Get("transactionId").(*big.Int)
	return true
}

// 买家可以下载数据
func onReadyForDownload(event events.Event) bool  {
	metaDataIdEncWithBuyer := event.Data.Get("metaDataIdEncBuyer").([]byte)
	metaDataId, err := accounts.GetAMInstance().Decrypt(
		metaDataIdEncWithBuyer,
		buyer.Account.Address,
		password)

	data := strings.Split(string(metaDataId[:]),"|")
	downloadInfo["metaDataId"] = data[0]
	downloadInfo["reencryptionKey"] = data[1]

	if err != nil {
		fmt.Println("failed to decrypt meta data id with buyer's private key", err)
	}

	fmt.Println("buyer: onReadyForDownload: ", event)
	return true
}

/**
 * 接口方法
 */


// 消息类型
type PublishMsg struct {
	Data []M
	DesData string
	PriceStragedy M
	Encryptioney string
}


 // API - 卖家发布数据
func PublishHandler(w http.ResponseWriter, r *http.Request)  {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("PublishHandler: parse form error ",err)
	}

	var m PublishMsg
	json.NewDecoder(r.Body).Decode(&m)

	txParam := chainoperations.TransactParams{
		From: common.HexToAddress(seller.Account.Address),
		Password: password,
	}

	var publish_id_list  = make(map[string]string)

	// 根据不同价格策略，发布不同的数据
	for key,properties :=range m.PriceStragedy {
		var publishData  []M

		// 组合不同定价策略的数据
		for _,item := range m.Data{
			var record = make(M)

			for _,property := range strings.Split(properties,",") {
				record[property] = item[property]
			}
			publishData = append(publishData,record)
		}

		pubJson,_ := json.Marshal(Obj{"data":publishData})
		var price,_ = strconv.ParseInt(key,10,64)

		if len(publishData) <3 {
			fmt.Println("发布的数据条数过少")
			return
		}

		var first_record,_ = json.Marshal(publishData[0])
		var second_record,_ = json.Marshal(publishData[1])
		var third_record,_ = json.Marshal(publishData[2])


		publish_id,_ := cif.Publish(
			&txParam,
			big.NewInt(price), // price
			pubJson,
			[][]byte{first_record,second_record,third_record},
			3,
			[]byte(m.DesData))

		publish_id_list[key] = publish_id

		// 数据发布ID、价格、描述信息列表更新
		publish_proof := []M{}
		publish_proof = append(publish_proof,publishData[0])
		publish_proof = append(publish_proof,publishData[1])
		publish_proof = append(publish_proof,publishData[2])

		publish_proof_map[publish_id] = publish_proof
		publish_des_list = append(publish_des_list,Obj{"price":key,"description":m.DesData,"publish_id":publish_id})
	}

	// 返回json字符串给客户端
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin","*")

	//_, err := json.Marshal(publish_id_list)
	json.NewEncoder(w).Encode(Obj{"result":publish_id_list})
}

// API - 数据浏览
func ExploreHandler(w http.ResponseWriter, r *http.Request)  {
	// 返回json字符串给客户端
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin","*")

	json.NewEncoder(w).Encode(Obj{"result":publish_des_list})
}

// API - 准备购买
func PrepareToBuyHandler(w http.ResponseWriter, r *http.Request)  {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("parse form error ",err)
	}
	formData := make(map[string]string)
	json.NewDecoder(r.Body).Decode(&formData)

	publishId := formData["publishId"]
	price,_ := strconv.ParseInt(formData["price"],10,64)

	txParam := chainoperations.TransactParams{
		From: common.HexToAddress(buyer.Account.Address),
		Password:password,
	}

	err = cif.ApproveTransfer(&txParam,common.HexToAddress(protocolContractAddr),big.NewInt(price))
	if err != nil{
		fmt.Println("BuyerApproveTransfer: ",err)
	}else {
		txParam = chainoperations.TransactParams{
			From: common.HexToAddress(buyer.Account.Address),
			Password:password,
		}

		err = cif.PrepareToBuy(&txParam,publishId)
		if err != nil{
			 fmt.Println("fail to prepareToBuy, error: ",err)
		}
	}

	// 返回json字符串给客户端
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin","*")
	json.NewEncoder(w).Encode(Obj{"result":publish_proof_map[publishId]})
}

// API - 购买
func BuyHandler(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin","*")
	transactionId_temp = transactionId.Int64()

	if transactionId.Int64() != 0 {
		txParam := chainoperations.TransactParams{
			From: common.HexToAddress(buyer.Account.Address),
			Password:password,
		}

		err := cif.BuyData(&txParam,transactionId)

		if err != nil{
			fmt.Println("fail to buy data, error: ",err)
		}

		transactionId = big.NewInt(0)
		json.NewEncoder(w).Encode(Obj{"result":transactionId_temp})
	} else {
		json.NewEncoder(w).Encode(Obj{"result":"none"})
	}
}

// API - 查看数据被购买状态
func DataPurchaseStateHandler(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin","*")

	if metaDataIdWithSeller != "" {
		json.NewEncoder(w).Encode(Obj{"result":M{"metaDataId":metaDataIdWithSeller,"buyer":buyer.Account.Address}})
	}else {
		json.NewEncoder(w).Encode(Obj{"result":"none"})
	}
}

// API - 数据被购买后，卖家给买家发送代理重加密密钥
func ProxyReEncryptionHandler(w http.ResponseWriter, r *http.Request)  {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("parse form error ",err)
	}

	formData := make(map[string]string)
	json.NewDecoder(r.Body).Decode(&formData)
	metaDataIdWithSeller = metaDataIdWithSeller +"|"+ formData["reencryptKey"]
	metaDataIdEncWithSeller,_ := accounts.GetAMInstance().Encrypt([]byte(metaDataIdWithSeller),seller.Account.Address)

	metaDataIdEncWithBuyer, err := accounts.GetAMInstance().ReEncrypt(
		metaDataIdEncWithSeller,
		seller.Account.Address,
		buyer.Account.Address,
		password,
	)

	if err != nil {
		fmt.Println("failed to ReEncrypt meta data id with buyer's public key")
	}

	txParam := chainoperations.TransactParams{
		From: common.HexToAddress(seller.Account.Address),
		Password:password,
	}

	err = cif.SubmitMetaDataIdEncWithBuyer(
		&txParam,
		big.NewInt(transactionId_temp),
		metaDataIdEncWithBuyer)

	if err != nil {
		fmt.Println("failed to SubmitMetaDataIdEncWithBuyer, error: ", err)
	}

	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin","*")
	json.NewEncoder(w).Encode(Obj{"result":"ok"})
}


// API - 买家查看数据下载信息
func DataDownloadHandler(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin","*")

	metaData := downloadInfo["metaDataId"]

	sh  := shell.NewShell(ipfsNodeAddr)
	err := sh.Get(metaData,"./ipfs")

	if err != nil {
		fmt.Println("error get ipfs data from metadataId")
	}

	if len(downloadInfo) >0 {
		json.NewEncoder(w).Encode(Obj{"result":M{"metaDataId":downloadInfo["metaDataId"],"reencryptionKey":downloadInfo["reencryptionKey"]}})
	}else {
		json.NewEncoder(w).Encode(Obj{"result":"none"})
	}
}


//--------------------------------------------------------------------------------------//


func main() {
	initApp()
	r := mux.NewRouter()
	r.HandleFunc("/publish", PublishHandler)
	r.HandleFunc("/explore", ExploreHandler)
	r.HandleFunc("/prepare2buy", PrepareToBuyHandler)
	r.HandleFunc("/buy", BuyHandler)
	r.HandleFunc("/getPurchaseState", DataPurchaseStateHandler)
	r.HandleFunc("/reencrypt", ProxyReEncryptionHandler)
	r.HandleFunc("/getDownloadInfo", DataDownloadHandler)


	http.Handle("/", r)
	http.ListenAndServe(":3001", nil)
}