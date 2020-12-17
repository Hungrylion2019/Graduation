/**
  author: hanxiaodong
 */
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"time"
)
type Company struct {
	CompanyName  string                `json:"company_name"`
	Location    string                 `json:"location"`
	RestaurantName string             `json:"restaurant_name"`
	Record   map[string]Ticket          `json:"record"`
	LeagalPerson string                  `json:"leagal_person"`
	Linkman    map[string]int            `json:"linkman"`
}
type Ticket struct {
      PurchaseRecord  []string           `json:"purchase_record"`
      TicketRecord  []string                     `json:"ticket_record"`
}

type SimpleChaincode struct {
} 

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response{

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "addCompany" {
		return t.addCompany(stub, args)
	} else if function == "addPurchaseRecord" {
		return t.addPurchaseRecord(stub, args)
	} else if function == "queryByName" {
		return t.queryByName(stub, args) 
	}  else if function == "addTicket" {
		return t.addTicket(stub, args)
	} else if function == "getHistoryForKey" {
		return t.getHistoryForKey(stub, args)
	} else if function == "getTicket" {
		return t.getTicket(stub, args)
	} else if function == "getPurchaseRecord" {
		return t.getPurchaseRecord(stub, args)
	}else if function == "delCompany" {
		return t.delCompany(stub, args)
	}
	return shim.Error("Invalid function name，input correct funciton name.")
}

//
func (t *SimpleChaincode) delCompany(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	Delcompanyname := args[0]
	err := stub.DelState(Delcompanyname)
	if err!=nil {
		return shim.Error(err.Error())

	}
	return shim.Success([]byte("success to DelCompany"))
}
//peer chaincode invoke -C mychannel -n company -c '{"function":"addCompany","Args":["Huanong","hubei","fandian","lion","lzc"]}'
func(t *SimpleChaincode)addCompany(stub shim.ChaincodeStubInterface, args []string)pb.Response{

	if len(args) != 6{
	   return shim.Error("Incorrect number of arguments. Expecting 6")
	}
        
        _,exist:=GetNameInfo(stub,args[0])
       if exist{ 
          return shim.Error("要创建的公司名字已经存在")
       }
    
         var company Company
	companyname:=args[0]
	location:=args[1]
	restaurant:=args[2]
	legalperson:=args[3]
	linkname:=args[4]
	tel,_:=strconv.Atoi(args[5])
	linkman:=map[string]int{linkname:tel}
	company=Company{
		CompanyName:companyname,
		Location:location,
		RestaurantName:restaurant,
		Record:map[string]Ticket{companyname:{}},
		LeagalPerson:legalperson,
		Linkman:linkman,
	}
	companyAsBytes,ok:=Putcompany(stub,company)
	if !ok{
		return shim.Error("failed to create CompanyAccount")
	}
	err := stub.SetEvent("Addmember", []byte{})
	if err != nil {
	return shim.Error(err.Error())
	}
	return shim.Success(companyAsBytes)
}
func(t *SimpleChaincode)addTicket(stub shim.ChaincodeStubInterface, args []string)pb.Response{

	if len(args) != 2{
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	name:=args[0]
	company,ok:=GetNameInfo(stub,name)
	if !ok{
		return shim.Error("the company is empty")
	}
	ticketrecord:=[]string{}
	ticketrecord=append(company.Record[name].TicketRecord,args[1])
	ticket:=Ticket{

		company.Record[name].PurchaseRecord,
		ticketrecord,
	}
	company.Record[name]=ticket
	a,tag:=Putcompany(stub,company)
	if !tag{
		return shim.Error("failed to upload Ticket")
	}
	err := stub.SetEvent("AddTicket", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(a)
}

func(t *SimpleChaincode)addPurchaseRecord(stub shim.ChaincodeStubInterface, args []string)pb.Response{

	if len(args) != 2{
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	name:=args[0]
	company,ok:=GetNameInfo(stub,name)
	if !ok{
		return shim.Error("the company is empty")
	}

	Purchaserecord:=[]string{}
	Purchaserecord=append(company.Record[name].PurchaseRecord,args[1])
	purchase:=Ticket{
		TicketRecord:company.Record[name].TicketRecord,
		PurchaseRecord:Purchaserecord,
	}
	company.Record[name]=purchase
	a,tag:=Putcompany(stub,company)
	if !tag{
		return shim.Error("failed to upload PurchaseRecord")
	}
	err := stub.SetEvent("AddPurchaseRecord", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(a)
}
func(t *SimpleChaincode)queryByName(stub shim.ChaincodeStubInterface, args []string)pb.Response{

	if len(args) != 1{
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	name := args[0]
	// 根据公司名字查询信息状态
	b, err := stub.GetState(name)
	if err != nil {
		return shim.Error("failed to find")
	}
	if b == nil {
		return shim.Error("The companyname is empty")
	}
	return shim.Success(b)
}
func(t *SimpleChaincode)getTicket(stub shim.ChaincodeStubInterface, args []string)pb.Response{

	if len(args) != 1{
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	name := args[0]
	// 根据公司名字查询信息状态
   company,ok:= GetNameInfo(stub,name)
   if !ok{
   	return shim.Error("the Company name is empty")
   }
   ticket:=company.Record[name].TicketRecord
   b,err:=json.Marshal(ticket)
   if err!=nil{
   	return shim.Error("failed to marshall")
   }
	return shim.Success(b)
}
func(t *SimpleChaincode)getPurchaseRecord(stub shim.ChaincodeStubInterface, args []string)pb.Response{

	if len(args) != 1{
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	name := args[0]
	// 根据公司名字查询信息状态
	company,ok:= GetNameInfo(stub,name)
	if !ok{
		return shim.Error("the Company name is empty")
	}
	purchase:=company.Record[name].PurchaseRecord
	b,err:=json.Marshal(purchase)
	if err!=nil{
		return shim.Error("failed to marshall")
	}
	return shim.Success(b)
}
func (t *SimpleChaincode) getHistoryForKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := args[0]

	// 返回某个键的所有历史值
	resultsIterator, err := stub.GetHistoryForKey(key)
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
//保存company
func Putcompany(stub shim.ChaincodeStubInterface, company Company) ([]byte, bool) {
	b, err := json.Marshal(company)
	if err != nil {
		return nil, false
	}
	err = stub.PutState(company.CompanyName, b)
	if err != nil {
		return nil, false
	}
	return b, true
}
func GetNameInfo(stub shim.ChaincodeStubInterface, CompanyName string) (Company, bool) {
	var company Company
	// 根据会员名字查询信息状态
	b, err := stub.GetState(CompanyName)
	if err != nil {
		return company, false
	}
	if b == nil {
		return company, false
	}
	// 对查询到的状态进行反序列化
	err = json.Unmarshal(b, &company)
	if err != nil {
		return company, false
	}
	// 返回结果
	return company, true
}

func main(){
	err := shim.Start(new(SimpleChaincode))
	if err != nil{
		fmt.Printf("启动SimpleChaincode时发生错误: %s", err)
	}
}
