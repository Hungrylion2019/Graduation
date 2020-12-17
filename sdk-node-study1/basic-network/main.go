package main

import (
	"basic-network/sdkInit"
	"basic-network/service"
        "basic-network/web"
        "basic-network/web/controller"
	"fmt"
	"os"
)

const (
	configFile  = "config.yaml"
	initialized = false
	SimpleCC    = "companytest"
)

func main() {

	initInfo := &sdkInit.InitInfo{

		ChannelID:     "lion",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/hyperledger/sdk-node-study1/basic-network/channel-artifacts/channel.tx",

		OrgAdmin:        "Admin",
		OrgName:         "Org1",
		OrdererOrgName:  "orderer.example.com",
		ChaincodeID:     SimpleCC,
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/hyperledger/sdk-node-study1/chaincode/food",
		UserName:        "User1",
	}

	sdk, err := sdkInit.SetupSDK(configFile, initialized)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	defer sdk.Close()

	err = sdkInit.CreateChannel(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	channelClient, err := sdkInit.InstallAndInstantiateCC(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(channelClient)
	//====================================
	serviceSetup := service.ServiceSetup{
		ChaincodeID: SimpleCC,
		Client:      channelClient,
	}
        fmt.Println("创建第一家公司信息")
        add0:=[]string{"Huanong","hubei","fandian","lion","lzc","1827131234"}
	msg, err := serviceSetup.AddCompany(add0...)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg)
	}
        fmt.Println("开始查询Huanong公司的信息")
	msg, err = serviceSetup.QueryByName("Huanong")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg)
	}
        fmt.Println("添加票据信息")
        add1:=[]string{"Huanong","psdasdaioju"}
	msg, err = serviceSetup.AddTicket(add1...)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg)
        }
        fmt.Println("开始查询Huanong公司的信息")
	msg, err = serviceSetup.QueryByName("Huanong")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg)
	}
        fmt.Println("添加购买信息")
        add2:=[]string{"Huanong","paioasdasdju"}
	msg, err = serviceSetup.AddPurchaseRecord(add2...)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg)
        }
        fmt.Println("开始查询Huanong公司的信息")
	msg, err = serviceSetup.QueryByName("Huanong")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg)
	}
        fmt.Println("查询Huanong公司的票据信息")
	msg, err = serviceSetup.GetTicket("Huanong")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg)
	}
        fmt.Println("查询Huanong公司的购买信息")
	msg, err = serviceSetup.GetPurchaseRecord("Huanong")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg)
	}
        fmt.Println("开始查询Huanong公司的记录")
	msg, err = serviceSetup.GetHistoryForKey("Huanong")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg)
	}
       	//===========================================//


	app := controller.Application{

		Fabric: &serviceSetup,

	}

	web.WebStart(&app)
        
}
