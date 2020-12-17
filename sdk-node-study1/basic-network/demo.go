package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const Admin = "Admin"
const Pwd="123456"
//定义智能合约结构体
type SmartContract struct{

}

/*
定义会员结构体
A、B、C三类
*/
type Member struct{
	MemberID int `json:"member_id"` //会员编号
	MemberName string `json:"member_name"`//会员姓名
	MemberPwd string `"json:member_pwd"`  //会员密码
	MemberClass string `json:"member_class"`//会员类别：ABC三类
	MemberLevel int `json:"member_level"`//会员级别
	DeviceID string `json:"device_id"` //登录绑定设备号
	SafeCode string `"json:safe_code"` //安全码
}

//链码初始化
func (t *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("member manager init ...")
	return shim.Success(nil)
}

//Invoke函数
func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function,args:=stub.GetFunctionAndParameters()
	if function=="addMember"{
		return t.addMember(stub,args)
	}else if function=="delMember"{
		return t.delMember(stub,args)
	}else if function=="initLedger"{
		return t.initLedger(stub,args)
	}else if function=="queryByName"{
		return t.queryByName(stub,args)
	}
	return shim.Error("Invalid function name，input correct funciton name.")
}
func delmenber(stub shim.ChaincodeStubInterface, member Member)bool{
	//msg:=fmt.Sprintf("会员删除失败，%s",memberName)
	err:=stub.DelState(member.MemberName)
	if err!=nil{
		return false
	}
	return true
}

//保存member
func Putmember(stub shim.ChaincodeStubInterface, member Member) ([]byte, bool) {
	b, err := json.Marshal(member)
	if err != nil {
		return nil, false
	}
	// 保存member状态
	err = stub.PutState(string(member.MemberID), b)
	if err != nil {
		return nil, false
	}
	err = stub.PutState(member.MemberName, b)
	if err != nil {
		return nil, false
	}
	return b, true
}
// 根据会员ID查询信息状态
// args: MemberID
func GetIDInfo(stub shim.ChaincodeStubInterface, MemberID string) (Member, bool) {
	var member Member
	// 根据会员ID查询信息状态
	b, err := stub.GetState(MemberID)
	if err != nil {
		return member, false
	}
	if b == nil {
		return member, false
	}
	// 对查询到的状态进行反序列化
	err = json.Unmarshal(b, &member)
	if err != nil {
		return member, false
	}
	// 返回结果
	return member, true
}
//根据会员名字查询信息
func GetNameInfo(stub shim.ChaincodeStubInterface, MemberName string) (Member, bool) {
	var member Member
	// 根据会员ID查询信息状态
	b, err := stub.GetState(MemberName)
	if err != nil {
		return member, false
	}
	if b == nil {
		return member, false
	}
	// 对查询到的状态进行反序列化
	err = json.Unmarshal(b, &member)
	if err != nil {
		return member, false
	}
	// 返回结果
	return member, true
}
func (t *SmartContract) queryByName(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=1{
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	name:=args[0]
	// 根据会员ID查询信息状态
	b, err := stub.GetState(name)
	if err != nil {
		return shim.Error("查找失败")
	}
	if b == nil {
	 return shim.Error("找不到该账户")
	}

	return shim.Success(b)

}
func (t *SmartContract) initLedger(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	name := "skyhuihui"
	member := Member{
		MemberID:0,
		MemberName:name,
		MemberPwd:Pwd,
		MemberClass:Admin,
	}
	re,ok:=Putmember(stub,member)
	if !ok{
		return shim.Error("管理员账号创建失败")
	}
	err := stub.SetEvent("InitLedger", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
    ans:="管理员创建成功"+string(re)
	return shim.Success([]byte(ans))
}
//删除会员
func (t *SmartContract)delMember(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=2{
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	Opmember:=args[0]
	a,ok:=GetNameInfo(stub,Opmember)
	if !ok {
		return shim.Error("Cannot find operator")
	}else if a.MemberClass!=Admin{
		return shim.Error("You should enough privilege to do")
	}
	Delmembername:=args[1]
	var member Member
	member,exist:=GetNameInfo(stub,Delmembername)
	if !exist{
        return shim.Error("failed to delete")
	}
	tag:=delmenber(stub,member)
	if  !tag{
		return shim.Error("failed to delete")

	}
	return shim.Success(nil)
}
//添加会员
func (t *SmartContract) addMember(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	//当前操作者
	if len(args)!=5{
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	//查询当前操作者的身份信息、并确认待添加会员姓名是否冲突

	currOperatorName := args[0]
	//操作者的信息
	currOperator,ok := GetNameInfo(stub,currOperatorName)
	if !ok{
		return shim.Error("Nnvalid operator name")
	}

	level,err:=strconv.Atoi(args[4])
	if err!=nil{
		return shim.Error(err.Error())
	}
	rand.Seed(time.Now().Unix())
	id:=rand.Intn(10000000)+1
     member:=Member{
     	MemberID:id,
     	MemberName:args[1],
     	MemberPwd:args[2],
     	MemberClass:args[3],
     	MemberLevel:level,
     }
	fmt.Println("member:", member)
	//返回会员名字的信息
	_,exist := GetNameInfo(stub,member.MemberName)
	if exist{
		return shim.Error("要添加的会员名字已经存在")
	}
	//权限判断
	isAuthorized := false
	if member.MemberClass == "A" {
		if member.MemberLevel == 1 {
			if currOperator.MemberClass == "Admin"{
				isAuthorized = true
			}
		}
		if member.MemberLevel == 2 {
			if currOperator.MemberClass == "A" &&  currOperator.MemberLevel == 1{
				isAuthorized = true
			}
		}
		if member.MemberLevel == 3 {
			if currOperator.MemberClass == "A" &&  currOperator.MemberLevel == 2{
				isAuthorized = true
			}
		}
	} else if member.MemberClass == "B" {
		if member.MemberLevel == 1 {
			if currOperator.MemberClass == "Admin"{
				isAuthorized = true
			}
		}
		if member.MemberLevel == 2 {
			if currOperator.MemberClass == "B" &&  currOperator.MemberLevel == 1{
				isAuthorized = true
			}
		}
		if member.MemberLevel == 3 {
			if currOperator.MemberClass == "B" &&  currOperator.MemberLevel == 2{
				isAuthorized = true
			}
		}
		if member.MemberLevel == 4 {
			if currOperator.MemberClass == "B" &&  currOperator.MemberLevel == 3{
				isAuthorized = true
			}
		}
	} else if member.MemberClass == "C" {
		isAuthorized = true
	}
	//如果权限满足
	if isAuthorized {
		_,b1:=Putmember(stub,member)
			if !b1{
				shim.Error("b保存信息失败")
			}
	} else {
		fmt.Println("添加会员失败！权限不满足，当前会员等级" + currOperator.MemberClass + string(currOperator.MemberLevel) + "，待添加会员等级" + member.MemberClass + string(member.MemberLevel))
		return shim.Error("Error, Authorized Fail !")
	}
	return shim.Success([]byte("会员创建成功"))
}
//Go语言的入口Main函数
func main(){
	err := shim.Start(new(SmartContract))
	if err != nil{
		fmt.Printf("Error creating new Smart Contract: %", err)
	}
}











