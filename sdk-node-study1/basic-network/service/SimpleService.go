package service

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)
func (t *ServiceSetup) GetTicket(name string) (string, error){

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "getTicket", Args: [][]byte{[]byte(name)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return "", err
	}

	return string(respone.Payload), nil
}
func (t *ServiceSetup) GetPurchaseRecord(name string) (string, error){

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "getPurchaseRecord", Args: [][]byte{[]byte(name)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return "", err
	}

	return string(respone.Payload), nil
}
func (t *ServiceSetup) GetHistoryForKey(name string) (string, error){

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "getHistoryForKey", Args: [][]byte{[]byte(name)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return "", err
	}

	return string(respone.Payload), nil
}

func (t *ServiceSetup) QueryByName(name string) (string, error){

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryByName", Args: [][]byte{[]byte(name)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return "", err
	}

	return string(respone.Payload), nil
}
func (t *ServiceSetup) AddCompany(args ...string) (string, error) {
                    
        var temArgs [][]byte
       for i:=0;i<len(args);i++{
         temArgs=append(temArgs,[]byte(args[i]))
}
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "addCompany", Args: temArgs}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	return string(respone.Payload), nil
}

func (t *ServiceSetup) AddTicket(args ...string) (string, error) {
                    
        var temArgs [][]byte
       for i:=0;i<len(args);i++{
         temArgs=append(temArgs,[]byte(args[i]))
}
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "addTicket", Args: temArgs}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	return string(respone.Payload), nil
}
func (t *ServiceSetup) AddPurchaseRecord(args ...string) (string, error) {
                    
        var temArgs [][]byte
       for i:=0;i<len(args);i++{
         temArgs=append(temArgs,[]byte(args[i]))
}
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "addPurchaseRecord", Args: temArgs}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	return string(respone.Payload), nil
}
