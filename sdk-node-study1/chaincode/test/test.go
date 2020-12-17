package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"math/rand"
	"strconv"
	"time"
)

const Admin = "Admin"
const Pwd="123456"
//定义智能合约结构体
type SmartContract struct{
}
const AdminName="skyhuihui"
const TokenKey = "Token"
type TransactionRecord struct {
	From        string  `json:"From"`
	To          string  `json:"To"`
	TokenSymbol string `json:"TokenSymbol"`
	Amount      float64 `json:"Amount"`
	TxId        string  `json:"TxId"`
}
type Currency struct {

	TokenName   string              `json:"TokenName"`
	TokenSymbol string              `json:"TokenSymbol"`
	TotalSupply float64             `json:"TotalSupply"`
	User        map[string]float64  `json:"User"` //某代币下各个用户持有的数量
	Record      []TransactionRecord `json:"Record"`
}
type Token struct {
	Currency map[string]Currency    `json:"Currency"`
}

type Member struct{
	MemberID int `json:"member_id"` //会员编号
	MemberName string `json:"member_name"`//会员姓名
	MemberPwd string `"json:member_pwd"`  //会员密码
	MemberClass string `json:"member_class"`//会员类别：ABC三类
	MemberLevel int `json:"member_level"`//会员级别
	DeviceID string `json:"device_id"` //登录绑定设备号
	SafeCode string `"json:safe_code"` //安全码
	BalanceOf map[string]float64            `json:"BalanceOf"`//对应币的数量
	Frozen    bool               `json:"Frozen"`//账户是否冻结
}


//Invoke函数
func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function,args:=stub.GetFunctionAndParameters()
	//peer chaincode invoke -C mychannel -n mycc -c '{"function":"addMember","Args":["skyhuihui","lion","123456","A","1"]}'
	if function=="addMember"{
		return t.addMember(stub,args)
	}else if function=="delMember"{
		return t.delMember(stub,args)//删除成员（没有考虑删除成员之后货币的总量如何变化）peer chaincode invoke -C mychannel -n mycc -c '{"function":"delMember","Args":["skyhuihui","lion"]}' -o orderer.example.com:7050
	}else if function=="initLedger"{
		return t.initLedger(stub,args)//生成管理员账号peer chaincode invoke -C mychannel -n mycc -c '{"function":"initLedger","Args":[]}' -o orderer.example.com:7050
	}else if function=="queryByName"{
		return t.queryByName(stub,args)//根据名字返回成员信息 peer chaincode invoke -C mychannel -n mycc -c '{"function":"queryByName","Args":["skyhuihui"]}'
	}else if function=="initCurrency"{
		return t.initCurrency(stub,args)//创建代币peer chaincode invoke -C mychannel -n mycc -c '{"function":"initCurrency","Args":["Netkiller Token","NKC","1000000","skyhuihui"]}'
	}else if function=="showToken"{
		return t.showToken(stub,args)//展示所有的代币信息 peer chaincode invoke -C mychannel -n mycc -c '{"function":"showToken","Args":[]}'
	}else if function=="showTokenUser"{
		return t.showTokenUser(stub,args)//展示某个代币下的信息peer chaincode invoke -C mychannel -n mycc -c '{"function":"showTokenUser","Args":["NKC"]}'
	}else if function=="transferToken"{
		return t.transferToken(stub,args)//转账 peer chaincode invoke -C mychannel -n mycc -c '{"function":"transferToken","Args":["skyhuihui","lion","NKC","12.584"]}'
	}else if function=="frozenAccount"{
		return t.frozenAccount(stub,args)// 冻结某个账户,该账户冻结之后无法转账peer chaincode invoke -C mychannel -n mycc -c '{"function":"frozenAccount","Args":["lion","true","skyhuihui"]}'
	}else if function=="tokenHistory"{
		return  t.tokenHistory(stub,args)// 查询某个币的交易记录 peer chaincode invoke -C mychannel -n mycc -c '{"function":"tokenHistory","Args":["NKC"]}'
	}else if function=="userTokenHistory"{
		return t.userTokenHistory(stub,args)//查询某个币某个用户的记录 peer chaincode invoke -C mychannel -n mycc -c '{"function":"userTokenHistory","Args":["NKC","lion"]}'
	}else if function=="getHistoryForKey" {
		return t.getHistoryForKey(stub, args)//查询某个键的所有交易记录 peer chaincode invoke -C mychannel -n mycc -c '{"function":"getHistoryForKey","Args":["lion"]}'
	}else if function=="burnToken"{
		return t.burnToken(stub,args)//回收某个账户的特定代币数量 peer chaincode invoke -C mychannel -n mycc -c '{"function":"burnToken","Args":["NKC","5000","123","skyhuihui"]}'
	}else if function=="mintToken" {
		return t.mintToken(stub, args)//增加某个币的数量，peer chaincode invoke -C mychannel -n mycc -c '{"function":"mintToken","Args":["NKC","5000","skyhuihui"]}'
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
	//err = stub.PutState(string(member.MemberID), b)
	//if err != nil {
	//	return nil, false
	//}
	err = stub.PutState(member.MemberName, b)
	if err != nil {
		return nil, false
	}
	return b, true
}
// 根据会员ID查询信息状态
// args: MemberID
//func GetIDInfo(stub shim.ChaincodeStubInterface, MemberID string) (Member, bool) {
//	var member Member
//	// 根据会员ID查询信息状态
//	b, err := stub.GetState(MemberID)
//	if err != nil {
//		return member, false
//	}
//	if b == nil {
//		return member, false
//	}
//	// 对查询到的状态进行反序列化
//	err = json.Unmarshal(b, &member)
//	if err != nil {
//		return member, false
//	}
//	// 返回结果
//	return member, true
//}
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
	// 根据会员名字查询信息状态
	b, err := stub.GetState(name)
	if err != nil {
		return shim.Error("failed to find")
	}
	if b == nil {
		return shim.Error("The account is empty")
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
		Frozen:false,
		BalanceOf: map[string]float64{},
	}
	re,ok:=Putmember(stub,member)
	if !ok{
		return shim.Error("管理员账号创建失败")
	}
	err := stub.SetEvent("InitLedger", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	ans:="initLedger success,"+string(re)
	return shim.Success([]byte(ans))
}
//删除会员
func (t *SmartContract)delMember(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=2{
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	Operator:=args[0]
	a,ok:=GetNameInfo(stub,Operator)
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
		Frozen:false,
		BalanceOf: map[string]float64{},
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
		memberbytes,b1:=Putmember(stub,member)
		if !b1{
			return shim.Error("b保存信息失败")
		}
		return shim.Success(memberbytes)
	} else {
		fmt.Println("添加会员失败！权限不满足，当前会员等级" + currOperator.MemberClass + string(currOperator.MemberLevel) + "，待添加会员等级" + member.MemberClass + string(member.MemberLevel))
		return shim.Error("Error, Authorized Fail !")
	}
	return shim.Success(nil)
}
//判断会员是否具有该货币
func isCurrency(member Member,_currency string) bool {

	if _,ok:=member.BalanceOf[_currency];!ok{
		return false
	}
	return true
}
//转账
func Transfer(stub shim.ChaincodeStubInterface,_from Member, _to Member, _currency string, _value float64) ([]byte,bool){

	var rev []byte

	if _from.Frozen {
		msg := "From 账号冻结"
		rev, _ = json.Marshal(msg)
		return rev,false
	}
	if _to.Frozen {
		msg := "To 账号冻结"
		rev, _ = json.Marshal(msg)
		return rev,false

	}
	//这里特别定义，只要转账者携带该货币，即使被转账者没有携带也可以转账，直接让被转账者生成新货币
	if !isCurrency(_from,_currency) {
		msg := "货币符号不存在"
		rev, _ = json.Marshal(msg)
		return rev,false
	}
	tokenAsbytes,err:=stub.GetState(TokenKey)

	if err!=nil{
		msg := "获取货币信息失败"
		rev, _ = json.Marshal(msg)
		return rev,false
	}
	token :=Token{}
	err= json.Unmarshal(tokenAsbytes,&token)
	if err!=nil {
		msg := "反序列失败"
		rev, _ = json.Marshal(msg)
		return rev, false
	}
	if _from.BalanceOf[_currency] >= _value {
		_from.BalanceOf[_currency] -= _value
		_to.BalanceOf[_currency] +=_value
		//更新转账和被转账者的信息
		_,ok:=Putmember(stub,_from)
		if !ok{
			msg := "更新转账者信息失败"
			rev, _ = json.Marshal(msg)
			return rev,false
		}
		_,tag:=Putmember(stub,_to)
		if !tag{
			msg := "更新接受者信息失败"
			rev, _ = json.Marshal(msg)
			return rev,false
		}
		//更新token的用户记录
		token.Currency[_currency].User[_from.MemberName] = _from.BalanceOf[_currency]
		token.Currency[_currency].User[_to.MemberName] = _to.BalanceOf[_currency]
		
		//将交易记录纳入代币当中，存放在区块链上
		TransferRecord:=TransactionRecord{_from.MemberName,_to.MemberName,_currency,_value,stub.GetTxID()}
		recordList := make([]TransactionRecord, 0)
		recordList = append(token.Currency[_currency].Record,TransferRecord)
		var cur Currency
		cur=Currency{token.Currency[_currency].TokenName,token.Currency[_currency].TokenSymbol,token.Currency[_currency].TotalSupply,token.Currency[_currency].User,recordList}
		token.Currency[_currency]=cur
		tokenAsBytes2,err:= json.Marshal(token)
		if err!=nil{
			msg := "序列化代币信息失败"
			rev, _ = json.Marshal(msg)
			return rev,false
		}
		err= stub.PutState(TokenKey,tokenAsBytes2)
		if err!=nil{
			msg := "更新代币信息失败"
			rev, _ = json.Marshal(msg)
			return rev,false
		}

		msg :="success to transfer"
		rev, _ = json.Marshal(msg)
		return rev,true
	} else {
		msg := "账户余额不足"
		rev, _ = json.Marshal(msg)
		return rev,false
	}
}

//Token键，值是所有以token为名，值为currency的键值对
func (t *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {

	token := &Token{Currency: map[string]Currency{}}

	tokenAsBytes, err := json.Marshal(token)
	err = stub.PutState(TokenKey, tokenAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Init Token %s \n", string(tokenAsBytes))
	}
	err = stub.SetEvent("tokenInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}
func (t *SmartContract) showToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	tokenAsBytes, err := stub.GetState(TokenKey)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("GetState(%s)) %s \n", TokenKey, string(tokenAsBytes))
	}
	return shim.Success(tokenAsBytes)
}
func (t *SmartContract) showTokenUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	_token := args[0]
	token := Token{}
	existAsBytes, err := stub.GetState(TokenKey)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("GetState(%s)) %s \n", TokenKey, string(existAsBytes))
	}
	json.Unmarshal(existAsBytes, &token)
	if _, ok := token.Currency[_token]; !ok {
		return shim.Error("The Token doesn't exist")
	}//1.3版本增加的部分，判断查询代币是否存在
	reToekn, err := json.Marshal(token.Currency[_token])
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Account balance %s \n", string(reToekn))
	}
	return shim.Success(reToekn)
}

func (t *SmartContract) initCurrency(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	operator:=args[3]
	member, exist := GetNameInfo(stub,operator)
	if !exist{
		return shim.Error("The administrator account is empty")
	}else if member.MemberClass!=Admin {
		return shim.Error("You should enough privilege to do")
	}

	_name := args[0]
	_symbol := args[1]
	_supply, _ := strconv.ParseFloat(args[2], 64)
	_account := args[3]
	token := Token{}
	existAsBytes, err := stub.GetState(TokenKey)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("GetState(%s)) %s \n", TokenKey, string(existAsBytes))
	}
	json.Unmarshal(existAsBytes, &token)
	if _, ok := token.Currency[_symbol]; ok {
		return shim.Error("Token has been created")
	}
	user := make(map[string]float64)
	user[_account] = _supply
	token.Currency[_symbol] = Currency{TokenName: _name, TokenSymbol: _symbol, TotalSupply: _supply,User: user,Record:[]TransactionRecord{}}
	tokenAsBytes, _ := json.Marshal(token)
	err = stub.PutState(TokenKey, tokenAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Init Token %s \n", string(tokenAsBytes))
	}
	member.BalanceOf[_symbol]=_supply
	_,ok:=Putmember(stub,member)
	if !ok{
		return shim.Error("failed to put member")
	}


	err = stub.SetEvent("tokenInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(tokenAsBytes)
}
//货币交易
func (t *SmartContract) transferToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	_from := args[0]
	_to := args[1]
	_currency := args[2]
	_amount, _ := strconv.ParseFloat(args[3], 64)

	if _amount <= 0 {
		return shim.Error("Incorrect number of amount")
	}
	memberfrom,exist:=GetNameInfo(stub,_from)
	if !exist{
		return shim.Error("Invalid user to transfer")
	}
	memberto,ok:=GetNameInfo(stub,_to)
	if !ok{
		return shim.Error("Invaild user to receive")
	}
	//转账是否正确
	isAuthorized := false
	if memberfrom.MemberClass=="A"{
		if memberto.MemberClass=="B"|| memberto.MemberClass=="C"{
			isAuthorized=true
		}
	}else if memberfrom.MemberClass=="B"{
		if memberto.MemberClass=="C"||memberto.MemberClass=="A"{
			isAuthorized=true
		}
	}else if memberfrom.MemberClass=="C"{
		if memberto.MemberClass=="B" {
			isAuthorized = true
		}
	}else if memberfrom.MemberClass==Admin{
		if memberto.MemberClass=="A" {
			isAuthorized = true
		}
	}
	if isAuthorized {

	result,ok := Transfer(stub, memberfrom, memberto, _currency, _amount)
		if !ok {
			return shim.Error(string(result))
		}
		err := stub.SetEvent("tokenInvoke", []byte{})
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(result)

	}else{
		return shim.Error("Error, Authorized Fail !")
	}

       return shim.Success(nil)
}
//冻结账户
func (t *SmartContract) frozenAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	_account := args[0]
	_status := args[1]
	Operator:=args[2]
	a,ok:=GetNameInfo(stub,Operator)
	if !ok {
		return shim.Error("Cannot find operator")
	}else if a.MemberClass!=Admin{
		return shim.Error("You should get enough privilege to do")
	}

	member,exist:= GetNameInfo(stub,_account)
	if !exist{
		return shim.Error("Cannot find member")
	}

	var status bool
	if _status == "true" {
		status = true
	} else {
		status = false
	}

	member.Frozen = status
	b,tag:=Putmember(stub,member)
	if !tag{
		return shim.Error("Failed to change frozen ")
	}else{
		fmt.Printf("frozenAccount - end %s \n",b)
	}

	err := stub.SetEvent("tokenInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//获取代币交易记录
func (t *SmartContract) tokenHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	_currency := args[0]

	tokenAsBytes, err := stub.GetState(TokenKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	token := Token{}
	json.Unmarshal(tokenAsBytes, &token)
	resultAsBytes, _ := json.Marshal(token.Currency[_currency].Record)

	fmt.Printf("Token Record %s \n", string(resultAsBytes))
	return shim.Success(resultAsBytes)
}
//获取某个用户某个代币交易记录
func (t *SmartContract) userTokenHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	_currency := args[0]
	_account := args[1]

	tokenAsBytes, err := stub.GetState(TokenKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	token := Token{}
	json.Unmarshal(tokenAsBytes, &token)
	var userRecord []TransactionRecord
	index := 0
	for k, v := range token.Currency[_currency].Record {
		if token.Currency[_currency].Record[k].From == _account || token.Currency[_currency].Record[k].To == _account {
			userRecord = append(userRecord, v)
			index++
		}
	}

	resultAsBytes, _ := json.Marshal(userRecord)
	fmt.Printf("Token Record nums %d \n", index)
	fmt.Printf("Token Record  %s \n", string(resultAsBytes))
	return shim.Success(resultAsBytes)
}

func (t *SmartContract) getHistoryForKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	marbleId := args[0]

	// 返回某个键的所有历史值
	resultsIterator, err := stub.GetHistoryForKey(marbleId)
	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResult.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		//时间戳格式化
		txtimestamp := queryResult.Timestamp
		tm := time.Unix(txtimestamp.Seconds, 0)
		timeString := tm.Format("2006-01-02 03:04:05 PM")
		buffer.WriteString(timeString)
		buffer.WriteString("\"")

		buffer.WriteString("{\"Value\":")
		buffer.WriteString("\"")
		buffer.WriteString(string(queryResult.Value))
		buffer.WriteString("\"")

		buffer.WriteString("{\"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(queryResult.IsDelete))
		buffer.WriteString("\"")

		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("- getMarblesByRange queryResult:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func (t *SmartContract) burnToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	_currency := args[0]
	_amount, _ := strconv.ParseFloat(args[1], 64)
	_account := args[2]
	operator:=args[3]
	member, exist := GetNameInfo(stub,operator)
	if !exist{
		return shim.Error("The administrator account is empty")
	}else if member.MemberClass!=Admin {
		return shim.Error("You should enough privilege to do")
	}

	burnmember,ok:=GetNameInfo(stub,_account)
	if !ok{
		return shim.Error("the destroyed token account is empty")
	}
	tag:=isCurrency(burnmember,_currency)
	if !tag{
		return shim.Error("the currency doesn't exist")
	}

	tokenAsBytes, err := stub.GetState(TokenKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Token before %s \n", string(tokenAsBytes))

	token := Token{}

	json.Unmarshal(tokenAsBytes, &token)
	//member的代币要减少，总token的代币也要减少
	if burnmember.BalanceOf[_currency] >= _amount {
		cur := token.Currency[_currency]
		cur.TotalSupply -= _amount  //货币销毁，所以这里货币总量要减少
		cur.User[burnmember.MemberName] -= _amount
		token.Currency[_currency] = cur
		burnmember.BalanceOf[_currency] -= _amount
		fmt.Println("success to recycle")
		//更新burnmenber
		a,ok:=Putmember(stub,burnmember)
		if !ok{
			return shim.Error("failed to update member ")
		}else{
			fmt.Println("Admin after:",a)
		}
		//更新代币
		tokenAsBytes, err = json.Marshal(token)
		if err != nil {
			return shim.Error(err.Error())
		}
		err = stub.PutState(TokenKey, tokenAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Printf("Token after %s \n", string(tokenAsBytes))
		fmt.Printf("burnToken %s \n", string(tokenAsBytes))
	} else {
		return shim.Error("burnmember's token is not enough to decrease")
	}
	err = stub.SetEvent("tokenInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SmartContract) mintToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	_currency := args[0]
	_amount, _ := strconv.ParseFloat(args[1], 64)
	_account := args[2]
	member, exist := GetNameInfo(stub,_account)
	if !exist{
		return shim.Error("The administrator account is empty")
	}else if member.MemberClass!=Admin {
		return shim.Error("You should enough privilege to do")
	}
	ok:=isCurrency(member,_currency)
	if !ok{
		return shim.Error("the destroyed token account is empty")
	}

	tokenAsBytes, err := stub.GetState(TokenKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Token before %s \n", string(tokenAsBytes))
	token := Token{}
	json.Unmarshal(tokenAsBytes, &token)
	//代币更新，管理员代币更新
	cur := token.Currency[_currency]
	cur.TotalSupply += _amount    //增发导致总量增加
	cur.User[member.MemberName] += _amount
	token.Currency[_currency] = cur
	member.BalanceOf[_currency] += _amount //更新管理员代币
	fmt.Println("success to increase")
	//更新menber
	a,ok:=Putmember(stub,member)
	if !ok{
		return shim.Error("failed to update member ")
	}else{
		fmt.Println("Admin after : ",a)
	}
	//更新代币
	tokenAsBytes, err = json.Marshal(token)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(TokenKey, tokenAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Token after %s \n", string(tokenAsBytes))
	fmt.Printf("minToken %s \n", string(tokenAsBytes))

	err = stub.SetEvent("tokenInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}
//Go语言的入口Main函数
func main(){
	err := shim.Start(new(SmartContract))
	if err != nil{
		fmt.Printf("Error creating new Smart Contract: %", err)
	}
}











