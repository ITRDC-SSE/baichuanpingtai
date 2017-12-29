package main


import (
       "encoding/json"
       "math"
       "time"
       "strings"
       "bytes"
	   "fmt"
	   "strconv" 
	   "github.com/hyperledger/fabric/core/chaincode/shim"
	   pb "github.com/hyperledger/fabric/protos/peer"
)
const success="success"
// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
type AcctAssetTime struct{
	AcctId string `json:"acctId"`
	AvaMoney int `json:"avaMoney"`
	Time string `json:"time"`
}
type AcctAsset struct{
	ObjectType string `json:"docType"`
	AcctId string `json:"acctId"`
	AvaMoney int `json:"avaMoney"`
}
type AcctMoneyFrozen struct{
	ObjectType string `json:"docType"`
	AcctId string `json:"acctId"`
	ContractN string `json:"contractN"`
	FrozenMoney int `json:"frozenMoney"`
}
type PriceTable struct{
	ObjectType string `json:"docType"`
	Time string `json:"time"`
	ProductCode string `json:"productCode"`
	Price float64 `json:"price"`
}
type SSEHold struct{
	ObjectType string `json:"docType"`
	AccId string `json:"accId"`
	ProductCode string `json:"productCode"`
	HoldNum int `json:"holdNum"`
	FrozenSecNum int `json:"frozenSecNum"`
}
type SSEHoldFrozen struct{
	ObjectType string `json:"docType"`
	AccId string `json:"accId"`
	ProductCode string `json:"productCode"`
	ContractN string `json:"contractN"`
	FrozenSecNum int `json:"frozenSecNum"`
}
type ContractHold struct{
	ObjectType string `json:"docType"`
	ContractN string `json:"contractN"`
	ContractCode string `json:"contractCode"`
	ContractStatus string `json:"contractStatus"`//"0"代表等待开始，"1"代表开始合约，"2"代表已经开始,
//"3"代表行权，"4"代表已终结,"5"代表互操作,"6"代表待联通,"7"代表联通,"8"代表解质押
	ContractCCID string `json:"contractCCID"`//合约名
	ContractFunctionName string `json:"contractFunctionName"`
    AccIdA string `json:"accIdA"`//合约关联方A
    AccIdB string `json:"accIdB"`
    AccIdC string `json:"accIdC"`
    AccIdD string `json:"accIdD"`
    AccIdE string `json:"accIdE"`
    CorContractNA string `json:"corContractNA"`//关联合约
    CorContractNB string `json:"corContractNB"`
    CorContractNC string `json:"corContractNC"`
    TransType string `json:"transType"`//"0"代表不可转让，"1"代表可转让
    LastSwapTime string `json:"lastSwapTime"`
}
//账户流水表
type AccFlow struct{
	ObjectType string `json:"docType"`
	AccFlowId string `json:"accFlowId"`//账户流水编号
	AccId string `json:"accId"`
	AssetId string `json:"assetId"`
	AssetNum string `json:"assetNum"`
	SType string `json:"sType"`//"0"代表增加,"1"代表减少,"2"代表锁定,"3"代表解锁
	ContractN string `json:"contractN"`
	Time string `json:"time"`
} 
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response  {
    fmt.Println("########### example_cc Init ###########")
    
	return shim.Success(nil)
}
//获得时间戳
func getTimeStamp() string{
	//获取时间戳
    timestamp:= time.Now().Unix()
    //格式化为字符串,tm为Time类型
    tm := time.Unix(timestamp, 0)
    Time:=tm.Format("20060102")
    return Time
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
		return shim.Error("Unknown supported call")
}
// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
        fmt.Println("########### example_cc Invoke ###########")
	function, args := stub.GetFunctionAndParameters()
	
	if function != "invoke" {
                return shim.Error("Unknown function call")
	}
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting at least 2")
	}
	if args[0] == "query" {
		// queries an entity state
		return t.query(stub, args)
	}
	//获得当前时间
	if args[0]=="getCurrTime"{
        return t.getCurrTime(stub,args)
	}
	//获得资金账户和钱
	if args[0]=="getAcctAssetAndTime"{
        return t.getAcctAssetAndTime(stub,args)
	}
	//获得股票价格
	if args[0]=="getStockPrice"{
		return t.getStockPrice(stub,args)
	}
	//查询账户余额
	if args[0]=="getAccMoney"{
		return t.getAccMoney(stub,args)
	}
	//锁定证券
	if args[0]=="lockStock"{
		return t.lockStock(stub,args)
	}
	//解锁证券
	if args[0]=="unlockStock"{
		return t.unlockStock(stub,args)
	}
	//划转证券
	if args[0]=="moveStock"{
		return t.moveStock(stub,args)
	}
	//锁定资金
	if args[0]=="lockMoney"{
		return t.lockMoney(stub,args)
	}
	//解锁资金
	if args[0]=="unlockMoney"{
		return t.unlockMoney(stub,args)
	}
	//划转资金
	if args[0]=="moveMoney"{
		return t.moveMoney(stub,args)
	}
	//获得合约持仓
	if args[0]=="getContractHold"{
		return t.getContractHold(stub,args)
	}
	//保存合约持仓
	if args[0]=="saveContractHold"{
		return t.saveContractHold(stub,args)
	}
	//设立合约
	if args[0]=="newContract"{
		return t.newContract(stub,args)
	}
	//启动合约
	if args[0]=="startContract"{
		return t.startContract(stub,args)
	}
	//获取权限
	if args[0]=="getAccPermissionOut"{
		return t.getAccPermissionOut(stub,args)
	}
	//出金入金
	if args[0]=="initAccMoney"{
		return t.initAccMoney(stub,args)
	}
	//添加股票价格
	if args[0]=="addStockPrice"{
		return t.addStockPrice(stub,args)
	}
	//添加SSE股票持仓
	if args[0]=="addStockHold"{
		return t.addStockHold(stub,args)
	}
	//设置时间
	if args[0]=="setTime"{
		return t.setTime(stub,args)
	}
	//获得某个账户的所有持仓数据
	if args[0]=="getSSEHoldByAccId"{
		return t.getSSEHoldByAccId(stub,args)
	}
	//获得所有的合约持仓
	if args[0]=="getContractHoldAll"{
		return t.getContractHoldAll(stub,args)
	}
	//获得某个账户的合约持仓
	if args[0]=="getContractHoldByAccId"{
		return t.getContractHoldByAccId(stub,args)
	}
	//获得某一个合约的状态
	if args[0]=="getContractHoldState"{
		return t.getContractHoldState(stub,args)
	}
	//查询合约内容
	if args[0]=="getContractContent"{
		return t.getContractContent(stub,args)
	}
	//查询账户流水
	if args[0]=="getAccFlow"{
		return t.getAccFlow(stub,args)
	}
	//获得指定时间的价格均值
	if args[0]=="getAveragePrice"{
		return t.getAveragePrice(stub,args)
	}
	//获得指定时间的价格标准差
	if args[0]=="getStandDeviationPrice"{
		return t.getStandDeviationPrice(stub,args)
	}
	//获得账户资产
	if args[0]=="getAcctAsset"{
		return t.getAcctAsset(stub,args)
	}
	//保存账户资产
	if args[0]=="saveAcctAsset"{
		return t.saveAcctAsset(stub,args)
	}
	//获得某个账户某个产品的持仓数据
	if args[0]=="getSSEHoldByAccAndProduct"{
		return t.getSSEHoldByAccAndProduct(stub,args)
	}
	//保存某个账户某个产品的持仓数据
	if args[0]=="saveSSEHoldByAccAndProduct"{
		return t.saveSSEHoldByAccAndProduct(stub,args)
	}
    //保存账户流水
    if args[0]=="saveAccFlowOut"{
		return t.saveAccFlowOut(stub,args)
	}
	//获得当前账户流水号
	if args[0]=="getAccFlowId"{
		return t.getAccFlowId(stub,args)
	}
	//保存当前的账户流水号
	if args[0]=="saveAccFlowId"{
		return t.saveAccFlowId(stub,args)
	}
	//保存股票价格测试数据
	if args[0]=="saveStockPriceTest"{
		return t.saveStockPriceTest(stub,args)
	}
	//保存账户资金
	if args[0]=="saveAcctAssetTest"{
		return t.saveAcctAssetTest(stub,args)
	}
	//保存账户密码
	if args[0]=="saveAcctPwTest"{
		return t.saveAcctPwTest(stub,args)
	}
	//保存合约持仓
	if args[0]=="saveContractHoldTest"{
		return t.saveContractHoldTest(stub,args)
	}
	//保存股票持仓数据
	if args[0]=="saveSSEHoldTest"{
		return t.saveSSEHoldTest(stub,args)
	}
	if args[0]=="saveDataTest"{
		return t.saveDataTest(stub,args)
	}
	//测试用GetStateByRange
	if args[0]=="getStateByRangeTest"{
		return t.getStateByRangeTest(stub,args)
	}
	//根据合约编号修改字段
	if args[0]=="updateNameByContractN"{
		return t.updateNameByContractN(stub,args)
	}
	//根据合约编号修改字段
	if args[0]=="updateFunctionNameByContractN"{
		return t.updateFunctionNameByContractN(stub,args)
	}
	//修改合约状态，为了方便测试
	if args[0]=="updateStatusByContractN"{
		return t.updateStatusByContractN(stub,args)
	}
	return shim.Error("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'")
}
//获得资金账户和时间
//getAcctAssetAndTime,acctId,pw
func(t *SimpleChaincode)getAcctAssetAndTime(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=3{
		return shim.Error("args length is wrong")
	}
	acctId:=args[1]
	pw:=args[2]
	result:=getAccPermission(stub,acctId,pw)
	if result!=success{
		return shim.Error("the pw is wrong")
	}
	acctAssetBytes,err:=stub.GetState(acctId)
	if err!=nil{
		return shim.Error("acctAssetBytes GetState failed")
	}
	var acctAsset AcctAsset
	err=json.Unmarshal(acctAssetBytes,&acctAsset)
	if err!=nil{
		return shim.Error("acctAssetBytes Unmarshal failed")
	}
	acctAssetTime:=AcctAssetTime{
		acctAsset.AcctId,
		acctAsset.AvaMoney,
        "/",
	}
	acctAssetTimeBytes,err:=json.Marshal(acctAssetTime)
	if err!=nil{
		return shim.Error("acctAssetTimeBytes Marshal failed")
	}
	return shim.Success(acctAssetTimeBytes)
}
//修改合约状态，为了方便测试
//updateStatusByContractN,contractN,newStatus
func (t *SimpleChaincode) updateStatusByContractN(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=3{
		return shim.Error("args length is wrong")
	}
	contractN:=args[1]
	newStatus:=args[2]
	contractHoldBytes,err:=stub.GetState(args[1])
	var contractHold ContractHold
	err=json.Unmarshal(contractHoldBytes,&contractHold)
	if err!=nil{
		return shim.Error("contractHoldBytes failed")
	}
	contractHold.ContractStatus=newStatus
    contractHoldBytes,err=json.Marshal(contractHold)
	if err!=nil{
		return shim.Error("contractHoldBytes Marshal failed")
	}
	err=stub.PutState(contractN,contractHoldBytes)
	if err!=nil{
		return shim.Error("contractHoldBytes PutState failed")
	}
	return shim.Success([]byte(success))
}
//根据合约编号修改字段
//updateNameByContractN,contractfunctionName,newName
func(t *SimpleChaincode) updateFunctionNameByContractN(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	contractHoldBytes,err:=stub.GetState(args[1])
	var contractHold ContractHold
	err=json.Unmarshal(contractHoldBytes,&contractHold)
	if err!=nil{
		return shim.Error("contractHoldBytes failed")
	}
	contractHold.ContractFunctionName=args[2]
	contractHoldBytes,err=json.Marshal(contractHold)
	if err!=nil{
		return shim.Error("contractHoldBytes Marshal failed")
	}
	err=stub.PutState(args[1],contractHoldBytes)
	if err!=nil{
		return shim.Error("contractHoldBytes PutState failed")
	}
	return shim.Success([]byte(success))
}
//根据合约编号修改字段
//updateNameByContractN,contractfunctionName,newName
func(t *SimpleChaincode) updateNameByContractN(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	contractHoldBytes,err:=stub.GetState(args[1])
	var contractHold ContractHold
	err=json.Unmarshal(contractHoldBytes,&contractHold)
	if err!=nil{
		return shim.Error("contractHoldBytes failed")
	}
	contractHold.ContractCCID=args[2]
	contractHoldBytes,err=json.Marshal(contractHold)
	if err!=nil{
		return shim.Error("contractHoldBytes Marshal failed")
	}
	err=stub.PutState(args[1],contractHoldBytes)
	if err!=nil{
		return shim.Error("contractHoldBytes PutState failed")
	}
	return shim.Success([]byte(success))
}
//保存股票持仓数据
func (t *SimpleChaincode) saveSSEHoldTest(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	for i:=1;i<len(args);i++{
		str:=strings.Split(args[i],",")
		// contractN:=str[0]
		// contractCode:=str[1]
		// contractStatus:="0"
		holdNum,err:=strconv.Atoi(str[2])
		if err!=nil{
			return shim.Error("holdNum Atoi failed")
		}
		var frozenNum int
		if str[3]!=""{
			frozenNum,err=strconv.Atoi(str[3])
			if err!=nil{
				return shim.Error("frozenNum Atoi failed")
			}
		}else{
			frozenNum=0
		}
		sseHold:=SSEHold{
			"sseHold",
			str[0],
			str[1],
			holdNum,
			frozenNum,
		}
		sseHoldBytes,err:=json.Marshal(sseHold)
        if err!=nil{
        	return shim.Error("sseHoldBytes Marshal failed")
        }
        err=stub.PutState(str[0]+str[1],sseHoldBytes)
        if err!=nil{
        	return shim.Error("acctPw PutState failed")
        }
    }
	return shim.Success([]byte(success))
}
func (t *SimpleChaincode) saveDataTest(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	fmt.Println("########### example_cc Init ###########")
	for i:=1;i<len(args);i++{
		str:=strings.Split(args[i],",")
		contractHold:=ContractHold{
			"contractHold",
			str[0],
			str[1],
			str[2],
			str[3],
			str[4],
			str[5],
			str[6],
			str[7],
			str[8],
			str[9],
			str[10],
			str[11],
			str[12],
			str[13],
			str[14],
		}
		contractHoldBytes,err:=json.Marshal(contractHold)
        if err!=nil{
        	return shim.Error("pTableBytes Marshal failed")
        }
        err=stub.PutState(str[0],contractHoldBytes)
        if err!=nil{
        	return shim.Error("acctPw PutState failed")
        }
    }
	return shim.Success([]byte(success))
}
//保存合约持仓
func (t *SimpleChaincode) saveContractHoldTest(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)<2{
		return shim.Error("args number is wrong")
	}
	for i:=1;i<len(args);i++{
		str:=strings.Split(args[i],",")
		if len(str)!=15{
			var ss string
			for i,aa:=range str{
               ss=strconv.Itoa(i)+aa
			}
			return shim.Error("length is wrong"+ss+"xxxx"+strconv.Itoa(len(str)))
		}
		contractHold:=ContractHold{
			"contractHold",
			str[0],
			str[1],
			str[2],
			str[3],
			str[4],
			str[5],
			str[6],
			str[7],
			str[8],
			str[9],
			str[10],
			str[11],
			str[12],
			str[13],
			str[14],
		}
		contractHoldBytes,err:=json.Marshal(contractHold)
        if err!=nil{
        	return shim.Error("pTableBytes Marshal failed")
        }
        err=stub.PutState(str[0],contractHoldBytes)
        if err!=nil{
        	return shim.Error("acctPw PutState failed")
        }
    }
	return shim.Success([]byte(success))
}
//保存账户密码
func(t *SimpleChaincode)saveAcctPwTest(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	for i:=1;i<len(args);i++{
		str:=strings.Split(args[i],",")
        err:=stub.PutState(str[0]+"pw",[]byte(str[1]))
        if err!=nil{
        	return shim.Error("acctPw PutState failed")
        }
    }
	return shim.Success([]byte(success))
}
//saveAcctAssetTest,
func(t *SimpleChaincode) saveAcctAssetTest(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	for i:=1;i<len(args);i++{
    	str:=strings.Split(args[i],",")
    	avaMoney,err:=strconv.Atoi(str[1])
        if err!=nil{
        	return shim.Error("price ParseFloat failed")
        }
    	acctAsset:=AcctAsset{
        	"acctAsset",
        	str[0],
        	avaMoney,
        }
        acctAssetBytes,err:=json.Marshal(acctAsset)
        if err!=nil{
        	return shim.Error("pTableBytes Marshal failed")
        }
        err=stub.PutState(str[0],acctAssetBytes)
        if err!=nil{
        	return shim.Error("acctAssetBytes PutState failed")
        }
    }
	return shim.Success([]byte(success))
}
//getStateByRangeTest,startDate,endDate,productCode
func(t *SimpleChaincode) getStateByRangeTest(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=4{
		return shim.Error("args length is wrong")
	}
	startDate:=args[1]
	endDate:=args[2]
	productCode:=args[3]
	resultsIterator, err := stub.GetStateByRange(startDate+productCode, endDate+productCode)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
    // buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResultKey, queryResultValue, err := resultsIterator.Next()
		if !strings.Contains(queryResultKey,productCode){
			continue
		}
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResultKey)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResultValue))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("- getMarblesByRange queryResult:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}
//保存股票价格测试数据
//saveStockPriceTest
func(t *SimpleChaincode) saveStockPriceTest(stub shim.ChaincodeStubInterface,args []string)pb.Response{
    for i:=1;i<len(args);i++{
    	str:=strings.Split(args[i],",")
    	price,err:=strconv.ParseFloat(str[2],64)
        if err!=nil{
        	return shim.Error("price ParseFloat failed")
        }
    	pTable:=PriceTable{
        	"priceTable",
        	str[1],
        	str[0],
        	price,
        }
        pTableBytes,err:=json.Marshal(pTable)
        if err!=nil{
        	return shim.Error("pTableBytes Marshal failed")
        }
        err=stub.PutState(str[1]+str[0],pTableBytes)
        if err!=nil{
        	return shim.Error("pTableBytes PutState failed")
        }
    }
	return shim.Success([]byte(success))
}
//保存某个账户某个产品的持仓数据
//saveSSEHoldByAccAndProduct,accId,product,sseHold
func(t *SimpleChaincode) saveSSEHoldByAccAndProduct(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=4{
		return shim.Error("args length is wrong")
	}
	accIdStr:=args[1]
	productCode:=args[2]
	sseHoldStr:=args[3]
	err:=stub.PutState(accIdStr+productCode,[]byte(sseHoldStr))
	if err!=nil{
		return shim.Error("sseHoldStr PutState failed")
	}
	return shim.Success([]byte(success))
}
//保存当前的账户流水号
//saveAccFlowId,accFlowId
func (t *SimpleChaincode)saveAccFlowId(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=2{
		return shim.Error("args length is wrong")
	}
	if args[1]==""{
		return shim.Error("accFlowId can not be nil")
	}
	accFlowId:=args[1]
	err:=stub.PutState("accFlow",[]byte(accFlowId))
	if err!=nil{
		return shim.Error("accFlowId putstate failed")
	}
	return shim.Success([]byte(accFlowId))
}
//获得当前账户流水号
//getAccFlowId
func (t *SimpleChaincode) getAccFlowId(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	accFlowId,result:=getId(stub,"accFlow")
	if result!=success{
		return shim.Error(result)
	}
	accFlowIdStr:=strconv.Itoa(accFlowId)
	return shim.Success([]byte(accFlowIdStr))
}
 //保存账户流水
 //saveAccFlowOut,accFlowIdStr,accFlow
func (t *SimpleChaincode) saveAccFlowOut(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=3{
		return shim.Error("args length is wrong")
	}
	accFlow:=args[2]
	accFlowIdStr:=args[1]
	err:=stub.PutState(accFlowIdStr,[]byte(accFlow))
	if err!=nil{
		return shim.Error("accFlow putstate failed")
	}
	return shim.Success([]byte(success))
}
//保存账户资产
//saveAcctAsset,accId,acctAsset,
func (t *SimpleChaincode) saveAcctAsset(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=3{
		return shim.Error("args length is wrong")
	}
	if args[1]==""||args[2]==""{
		return shim.Error("args is nil")
	}
	accIdStr:=args[1]
	acctAssetStr:=args[2]
	err:=stub.PutState(accIdStr,[]byte(acctAssetStr))
	if err!=nil{
		return shim.Error("acctAssetStr putstate failed")
	}
	return shim.Success([]byte(success))
}
//获得某个账户某个产品的持仓数据
//getSSEHoldByAccAndProduct,accId,productCode
func (t *SimpleChaincode) getSSEHoldByAccAndProduct(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=3{
		return shim.Error("args length is wrong")
	}
	if args[1]==""||args[2]==""{
		return shim.Error("args can not be nil")
	}
	accIdStr:=args[1]
	productCode:=args[2]
	sseHoldBytes,err:=stub.GetState(accIdStr+productCode)
	if err!=nil{
		return shim.Error("sseHoldBytes GetState failed")
	}
	if sseHoldBytes==nil{
		sseHold:=SSEHold{
			"sseHold",
			accIdStr,
			productCode,
			0,
			0,
		}
		sseHoldBytes,err=json.Marshal(sseHold)
		if err!=nil{
			return shim.Error("sseHoldBytes Marshal failed")
		}
	}
	return shim.Success(sseHoldBytes)
}
//获得账户资产
//getAcctAsset,accIdStr
func (t *SimpleChaincode) getAcctAsset(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=2{
		return shim.Error("args length is wrong")
	}
	if args[1]==""{
		return shim.Error("accId can not be nil")
	}
	accIdStr:=args[1]
	acctAssetBytes,err:=stub.GetState(accIdStr)
	if err!=nil{
		return shim.Error("acctAssetBytes getstate failed")
	}
	if acctAssetBytes==nil{
		acctAsset:=AcctAsset{
			"acctAsset",
			accIdStr,
			0,
		}
		acctAssetBytes,err=json.Marshal(acctAsset)
		if err!=nil{
			return shim.Error("acctAssetBytes Marshal failed")
		}
	}
    return shim.Success(acctAssetBytes)
}
//获得指定时间的价格标准差
//getStandDeviationPrice,startDate,endDate,productCode
func (t *SimpleChaincode) getStandDeviationPrice(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=4{
		return shim.Error("args length is wrong")
	}
	startDate:=args[1]
	endDate:=args[2]
	productCode:=args[3]
	averageValue,result:=getAverageValue(stub,startDate,endDate,productCode)
	if result!=success{
		return shim.Error(result)
	}
	resultsIterator, err := stub.GetStateByRange(startDate+productCode, endDate+productCode)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	var total,i float64
    for resultsIterator.HasNext() {
		_, queryResultValue, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		var pTable PriceTable
		err=json.Unmarshal(queryResultValue,&pTable)
		if err != nil {
			return shim.Error(err.Error())
		}
		if pTable.ProductCode==productCode{
			total+=(pTable.Price-averageValue)*(pTable.Price-averageValue)
            i++
		}
	}
	if i==0{
		return shim.Error("i is 0")
	}
	averageValue=total/i
	deviationV:=math.Sqrt(averageValue)
	deviationStr:=strconv.FormatFloat(deviationV,'g',2,64)
	return shim.Success([]byte(deviationStr))
}
//获得指定时间的价格均值
//getAveragePrice,startDate,endDate,productCode
func (t *SimpleChaincode) getAveragePrice(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=4{
		return shim.Error("args length is wrong")
	}
	startDate:=args[1]
	endDate:=args[2]
	productCode:=args[3]
    averageValue,result:=getAverageValue(stub,startDate,endDate,productCode)
	if result!=success{
		return shim.Error(result)
	}
	averagePriceStr:=strconv.FormatFloat(averageValue,'g',2,64)
	return shim.Success([]byte(averagePriceStr))
}
func getAverageValue(stub shim.ChaincodeStubInterface,startDate,endDate,productCode string)(float64,string){
	resultsIterator, err := stub.GetStateByRange(startDate+productCode, endDate+productCode)
	if err != nil {
		return 0,err.Error()
	}
	defer resultsIterator.Close()
	var total,i float64
    for resultsIterator.HasNext() {
		_, queryResultValue, err := resultsIterator.Next()
		if err != nil {
			return 0,err.Error()
		}
		var pTable PriceTable
		err=json.Unmarshal(queryResultValue,&pTable)
		if err != nil {
			return 0,err.Error()
		}
		if pTable.ProductCode==productCode{
			total+=pTable.Price
            i++
		}
	}
	if i==0{
		return 0,"i is 0"
	}
	return total/i,success
}
//查询账户流水
//getAccFlow,accId,pw,flag,"0"代表该账户的所有流水,"1"代表所有账户的流水
func (t *SimpleChaincode) getAccFlow(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=4{
		return shim.Error("args length is wrong")
	}
	if args[1]==""||args[2]==""||args[3]==""{
		return shim.Error("args can not be nil")
	}
	accIdStr:=args[1]
	pwStr:=args[2]
	flag:=args[3]
	// if accIdStr!="admin"{
	// 	return shim.Success([]byte("only admin can invoke"))
	// }
	result:=getAccPermission(stub,accIdStr,pwStr)
	if result!=success{
		return shim.Error(result)
	}
	var queryString string
	if flag=="0"{
		queryString=fmt.Sprintf("{\"selector\":{\"docType\":\"accFlow\",\"accId\":\"%s\"}}",accIdStr)		
	}else{
		queryString= fmt.Sprintf("{\"selector\":{\"docType\":\"accFlow\"}}")
	}
	 
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}
//获得合约内容
//getContractContent,contractN,accId,pw
func (t *SimpleChaincode) getContractContent(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=4{
		return shim.Error("args length is wrong")
	}
	if args[1]==""||args[2]==""||args[3]==""{
		return shim.Error("args can not be nil")
	}
    contractNStr:=args[1]
    accIdStr:=args[2]
    pwStr:=args[3]
    result:=getAccPermission(stub,accIdStr,pwStr)
    if result!=success{
    	return shim.Error(result)
    }
    contractContentBytes,err:=stub.GetState(contractNStr+"content")
    if err!=nil{
    	return shim.Error("contractContentBytes getstate failed")
    }
    if contractContentBytes==nil{
    	return shim.Error("there is not the contract")
    }
    return shim.Success(contractContentBytes)
}
//获得一个合约的状态
//getContractHoldState,contractN,accId,pw
func (t *SimpleChaincode) getContractHoldState(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=4{
		return shim.Error("args length is wrong")
	}
	if args[1]==""||args[2]==""||args[3]==""{
		return shim.Error("args can not be nil")
	}
	contractNStr:=args[1]
	accIdStr:=args[2]
	pw:=args[3]
	result:=getAccPermission(stub,accIdStr,pw)
	if result!=success{
		return shim.Error(result)
	}
	contractHoldBytes,err:=stub.GetState(contractNStr)
	if err!=nil{
		return shim.Error("contractHoldBytes getstate failed")
	}
	var contractHold ContractHold
	err=json.Unmarshal(contractHoldBytes,&contractHold)
	if err!=nil{
		return shim.Error("contractHoldBytes Unmarshal failed")
	}
	return shim.Success([]byte(contractHold.ContractStatus))
}
//获得某个账户的所有合约持仓
//getContractHoldByAccId,accId,pw
func (t *SimpleChaincode) getContractHoldByAccId(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=3{
		return shim.Success([]byte("args length is wrong"))
	}
	if args[1]==""||args[2]==""{

	}
	accIdStr:=args[1]
	pwStr:=args[2]
	result:=getAccPermission(stub,accIdStr,pwStr)
	if result!=success{
		return shim.Error("the pw is wrong")
	}
	queryResults, err := getQueryResultForQueryString2(stub, accIdStr)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}
//获得所有的合约持仓,只有超级管理员才可以调用
//getContractHoldAll,accId,pw
func (t *SimpleChaincode) getContractHoldAll(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=3{
		return shim.Success([]byte("args length is wrong"))
	}
	if args[1]==""||args[2]==""{
		return shim.Success([]byte("args can not be nil"))
	}
	accIdStr:=args[1]
	pwStr:=args[2]
	if accIdStr!="admin"{
		return shim.Error("no permission")
	}
	result:=getAccPermission(stub,accIdStr,pwStr)
	if result!=success{
		return shim.Error(result)
	}
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"contractHold\"}}")
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}
//获得某个账户的持仓数据
//getSSEHoldByAccId,accId,pw
func (t *SimpleChaincode) getSSEHoldByAccId(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=3{
		return shim.Success([]byte("args length is not 2"))
	}
	if args[1]==""||args[2]==""{
		return shim.Error("args can not be nil")
	}
	accIdStr:=args[1]
	pwStr:=args[2]
	result:=getAccPermission(stub,accIdStr,pwStr)
	if result!=success{
		return shim.Error("the pw is wrong")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"sseHold\",\"accId\":\"%s\"}}", accIdStr)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}
func getQueryResultForQueryString2(stub shim.ChaincodeStubInterface, accIdStr string) ([]byte, error) {
	// fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)
    accIdArray:=[]string{"accIdA","accIdB","accIdC","accIdD","accIdE"}
    var buffer bytes.Buffer
	buffer.WriteString("[")
	var queryString string
	contractN:=""
	bArrayMemberAlreadyWritten2:=false
    for _,accId:=range accIdArray{
    	queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"contractHold\",\""+accId+"\":\"%s\"}}",accIdStr)
    	resultsIterator, err := stub.GetQueryResult(queryString)
		if err != nil {
			return nil, err
		}
		defer resultsIterator.Close()
		// buffer is a JSON array containing QueryRecords
		bArrayMemberAlreadyWritten := false
		for resultsIterator.HasNext() {
			_, queryResultRecord, err := resultsIterator.Next()
			if err != nil {
				return nil, err
			}
			// Add a comma before array members, suppress it for the first array member
			
			//buffer.WriteString("{")
			// buffer.WriteString("{\"Key\":")
			// buffer.WriteString("\"")
			// buffer.WriteString(queryResultKey)
			// buffer.WriteString("\"")

			// buffer.WriteString(", \"Record\":")
			// Record is a JSON object, so we write as-is
			var contractHold ContractHold
			err=json.Unmarshal(queryResultRecord,&contractHold)
			if err != nil {
				return nil, err
			}
			if contractN!=""&&strings.Contains(contractN,contractHold.ContractN){
				continue
			}
			if bArrayMemberAlreadyWritten == true ||bArrayMemberAlreadyWritten2==true{
				buffer.WriteString(",")
			}
			contractN+=contractHold.ContractN
			buffer.WriteString(string(queryResultRecord))
			// if index!=(len(accIdArray)-1){
			// 	buffer.WriteString("},")
			// }else{
			// 	buffer.WriteString("}")
			// }
			bArrayMemberAlreadyWritten = true
			bArrayMemberAlreadyWritten2=true
		}
		
    }
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		_, queryResultRecord, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		// buffer.WriteString("{")
		// buffer.WriteString("{\"Key\":")
		// buffer.WriteString("\"")
		// buffer.WriteString(queryResultKey)
		// buffer.WriteString("\"")

		// buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResultRecord))
		// buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}
//设置时间
//setTime,time
func (t *SimpleChaincode) setTime(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=2{
		return shim.Error("args length is wrong")
	}
	if args[1]==""{
		shim.Error("the time can not be nil")
	}
	err:=stub.PutState("currTime",[]byte(args[1]))
	if err!=nil{
		return shim.Error("currTime putstate failed")
	}
	return shim.Success([]byte(success))
}
//添加SSE股票持仓
//addStockHold,accIdStr,productCode,holdNum,frozenSecNum
func (t *SimpleChaincode) addStockHold(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=5{
		return shim.Error("args length is wrong")
	}
	if args[1]==""||args[2]==""{
		return shim.Error("args can not be nil")
	}
	accIdStr:=args[1]
	productCode:=args[2]
	holdNumStr:=args[3]
	frozenSecNumStr:=args[4]
	var err error
	var holdNum,frozenSecNum int
	if holdNumStr==""{
		holdNum=0
	}else{
		holdNum,err=strconv.Atoi(holdNumStr)
		if err!=nil{
			return shim.Error("holdNum atoi failed")
		}
	}
	if frozenSecNumStr==""{
		frozenSecNum=0
	}else{
		frozenSecNum,err=strconv.Atoi(frozenSecNumStr)
		if err!=nil{
			return shim.Error("frozenSecNum atoi failed")
		}
	}
	sseHold:=SSEHold{
		"sseHold",
		accIdStr,
		productCode,
		holdNum,
		frozenSecNum,
	}
	sseHoldBytes,err:=json.Marshal(sseHold)
	if err!=nil{
		return shim.Error("sseHoldBytes Marshal failed")
	}
	err=stub.PutState(accIdStr+productCode,sseHoldBytes)
	if err!=nil{
		return shim.Error("sseHoldBytes putstate failed")
	}
	accFlowId,result:=getId(stub,"accFlow")
	if result!=success{
		return shim.Error(result)
	}
	accFlowIdStr:=strconv.Itoa(accFlowId)
	result=saveAccFlow(stub,accFlowIdStr,accIdStr,productCode,holdNumStr,"0","c0000")
	if result!=success{
		return shim.Error(result)
	}
	return shim.Success([]byte(success))
}
//添加股票价格
//addStockPrice,time,productCode,price
func (t *SimpleChaincode) addStockPrice(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=4{
		return shim.Error("args length is wrong")
	}
	time:=args[1]
	productCode:=args[2]
	priceStr:=args[3]
	price,err := strconv.ParseFloat(priceStr,64)
    if err!=nil{
    	return shim.Error("priceStr ParseFloat failed")
    }
	priceTable:=PriceTable{
		"priceTable",
		time,
		productCode,
        price,
	}
	priceBytes,err:=json.Marshal(priceTable)
	if err!=nil{
		return shim.Error("priceBytes Marshal failed")
	}
	err=stub.PutState(time+productCode,priceBytes)
	if err!=nil{
		return shim.Error("price putstate failed")
	}
	return shim.Success([]byte(success))
}
//出金入金
//initAccMoney,accIdStr,money,fname
func (t *SimpleChaincode) initAccMoney(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=4{
		return shim.Error("args length is wrong")
	}
	if args[1]==""||args[2]==""||args[3]==""{
		return shim.Error("args can not be nil")
	}
	accIdStr:=args[1]
	moneyStr:=args[2]
	fname:=args[3]
	money,err:=strconv.Atoi(moneyStr)
	if err!=nil{
		return shim.Error("money atoi failed")
	}
	acctAssetBytes,err:=stub.GetState(accIdStr)
	if err!=nil{
		return shim.Error("accMoneyBytes getstate failed")
	}
	var acctAsset AcctAsset
	if acctAssetBytes==nil{
		acctAsset=AcctAsset{
			"acctAsset",
		    accIdStr,
		    0,
		}
	}else{
        err:=json.Unmarshal(acctAssetBytes,&acctAsset)
        if err!=nil{
        	return shim.Error("acctAssetBytes Unmarshal failed")
        }
	}
    var sType string
	if fname=="in"{
		sType="0"
		acctAsset.AvaMoney+=money
	}else if fname=="out"{
		sType="1"
		if acctAsset.AvaMoney<money{
			return shim.Error("money is more than accmoney")
		}else{
			acctAsset.AvaMoney-=money
		}
	}else{
		return shim.Error("fname must be in or out")
	}
	acctAssetBytes,err=json.Marshal(acctAsset)
	if err!=nil{
		return shim.Error("acctAssetBytes Marshal failed")
	}
	err=stub.PutState(accIdStr,acctAssetBytes)
	if err!=nil{
		return shim.Error("accMoneyStr putstate failed")
	}
	accFlowId,result:=getId(stub,"accFlow")
	if result!=success{
		return shim.Error(result)
	}
	accFlowIdStr:=strconv.Itoa(accFlowId)
	result=saveAccFlow(stub,accFlowIdStr,accIdStr,"money",moneyStr,sType,"c9999")
	if result!=success{
		return shim.Error(result)
	}
	return shim.Success([]byte(success))
}
//获取权限
//getAccPermissionOut,accIdStr,pw
func (t *SimpleChaincode) getAccPermissionOut(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=3{
		return shim.Error("args length is wrong")
	}
	if args[1]==""||args[2]==""{
		return shim.Error("args can not be nil")
	} 
	accIdStr:=args[1]
	pwStr:=args[2]
	pwBytes,err:=stub.GetState(accIdStr+"pw")
	if err!=nil{
		return shim.Error("pwBytes GetState failed")
	}
	if pwBytes==nil{
		return shim.Error("there is not the accountId")
	}else if pwStr!=string(pwBytes) {
		return shim.Error("password is wrong")
	}
	return shim.Success([]byte(success))
}
//启动合约
//startContract,accIdBStr,pwStr,contractNStr
func (t *SimpleChaincode) startContract(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=4{
		return shim.Error("args length is wrong")
	}
	if args[1]==""||args[2]==""||args[3]==""{
		return shim.Error("args can not be nil")
	}
	//获得权限
	accIdBStr:=args[1]
	pwStr:=args[2]
	contractNStr:=args[3]
	result:=getAccPermission(stub,accIdBStr,pwStr)
	if result!=success{
		return shim.Error(result)
	}
	contractHoldBytes,err:=stub.GetState(contractNStr)
	if err!=nil{
		return shim.Error("contractNStr getstate failed")
	}
	if contractHoldBytes==nil{
		return shim.Error("there is not the contractN")
	}
	var contractHold ContractHold
	err=json.Unmarshal(contractHoldBytes,&contractHold)
	if err!=nil{
		return shim.Error("contractHoldBytes Unmarshal failed")
	}
	if contractHold.ContractStatus!="0"{
		return shim.Error("the contractstate is not 0")
	}
	contractHold.ContractStatus="2"
	contractHoldBytes,err=json.Marshal(contractHold)
	if err!=nil{
		return shim.Error("contractHoldbytes marshal failed")
	}
	err=stub.PutState(contractNStr,contractHoldBytes)
	if err!=nil{
		return shim.Error("contractHoldbytes PutState failed")
	}
	return shim.Success([]byte(success))
}
//设立合约
//newContract,contractN,contractCode,contractFunctionName,accIdAStr,accIdBStr
//accIdCStr,accIdDStr,accIdEStr,corContractNA,corContractNB,corContractNC
func (t *SimpleChaincode) newContract(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=14{
		return shim.Error("args number is wrong")
	}
	// contractN:=getId(stub,"contractN")
	contractNStr:=args[1]
	contractCode:=args[2]
	contractCCID:=args[3]
	contractFunctionName:=args[4]
	accIdAStr:=args[5]
	accIdBStr:=args[6]
	accIdCStr:=args[7]
	accIdDStr:=args[8]
	accIdEStr:=args[8]
	corContractNA:=args[10]
	corContractNB:=args[11]
	corContractNC:=args[12]
	isTransType:=args[13]
	contractBytes,err:=stub.GetState(contractNStr)
	if err!=nil{
		return shim.Error("contractBytes getstate failed")
	}
	if contractBytes!=nil{
		return shim.Error("the contractN has existed")
	}
	contractHold:=ContractHold{
		"contractHold",
		contractNStr,
		contractCode,
        "0",//合约状态
        contractCCID,
        contractFunctionName,
        accIdAStr,
        accIdBStr,
        accIdCStr,
        accIdDStr,
        accIdEStr,
        corContractNA,
        corContractNB,
        corContractNC,
        isTransType,//"0"代表true
        "0",//上次转换时间
	}
	//
    contractHoldBytes,err:=json.Marshal(contractHold)
    if err!=nil{
    	return shim.Error("contractHoldBytes Marshal failed")
    }
    err=stub.PutState(contractNStr,contractHoldBytes)
    if err!=nil{
    	return shim.Error("contractHoldBytes Marshal failed")
    }
    return shim.Success([]byte(success))
}
//getAccPermission,accId,pw,
func getAccPermission(stub shim.ChaincodeStubInterface,accIdStr,pw string)string{
	pwBytes,err:=stub.GetState(accIdStr+"pw")
	if err!=nil{
		return "pwBytes getstate failed"
	}
	if pwBytes==nil{
		return "there is not the pw"
	}
	if pw!=string(pwBytes){
		return "no permission"
	}
	return success
}
//保存合约持仓
//saveContractHold,contractNStr,contractHold
func (t *SimpleChaincode) saveContractHold(stub shim.ChaincodeStubInterface,args []string)pb.Response{

	if len(args)!=3{
		return shim.Error("args number is wrong")
	}
	contractNStr:=args[1]
	contractHoldStr:=args[2]
	fmt.Println("dbcc;"+contractNStr+";"+contractHoldStr)
	err:=stub.PutState(contractNStr,[]byte(contractHoldStr))
	if err!=nil{
		return shim.Error("contractHoldStr putstate failed")
	}
	return shim.Success([]byte(success))
}
//获得合约持仓
//getContractHold,contractNStr
func (t *SimpleChaincode) getContractHold(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=2{
		return shim.Error("getContractHold failed")
	}
	if args[1]==""{
		return shim.Error("the second args can not be nil")
	}
	contractHoldbytes,err:=stub.GetState(args[1])
	if err!=nil{
		return shim.Error("contractHoldbytes getstate failed")
	}
	if contractHoldbytes==nil{
		return shim.Error("there is not the contract")
	}
	return shim.Success(contractHoldbytes)
}
//划转资金
//moveMoney,accIdFrom,accIdTo,money,contractN,accFlowIdStr
func (t *SimpleChaincode) moveMoney(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=6{
		return shim.Error("args length is wrong")
	}
	if args[1]==""||args[2]==""||args[3]==""||args[4]==""{
		return shim.Error("args can not be nil")
	}
	accIdFromStr:=args[1]
	accIdToStr:=args[2]
	moneyStr:=args[3]
	contractNStr:=args[4]
	accFlowIdStr:=args[5]
	moveMoney,err:=strconv.Atoi(args[3])
	if err!=nil{
		return shim.Error("money atoi failed")
	}
	acctAssetFromBytes,err:=stub.GetState(accIdFromStr)
	if err!=nil{
		return shim.Error("avaMoneyFromBytes getstate failed")
	}
	var acctAssetFrom,acctAssetTo AcctAsset
	err=json.Unmarshal(acctAssetFromBytes,&acctAssetFrom)
	if err!=nil{
		return shim.Error("avaMoneyFrom atoi failed")
	}
	if acctAssetFrom.AvaMoney<moveMoney{
		return shim.Error("moveMoney is more than avaMoneyFrom")
	}
	acctAssetFrom.AvaMoney-=moveMoney
	acctAssetToBytes,err:=stub.GetState(accIdToStr)
	if err!=nil{
		return shim.Error("avaMoneyToBytes getstate failed")
	}
	err=json.Unmarshal(acctAssetToBytes,&acctAssetTo)
	if err!=nil{
		return shim.Error("acctAssetToBytes unmarshal failed")
	}
	acctAssetTo.AvaMoney+=moveMoney
	acctAssetFromBytes,err=json.Marshal(acctAssetFrom)
	if err!=nil{
		return shim.Error("acctAssetFromBytes Marshal failed")
	}
	acctAssetToBytes,err=json.Marshal(acctAssetTo)
	if err!=nil{
		return shim.Error("acctAssetToBytes Marshal failed")
	}
	//保存
	err=stub.PutState(accIdFromStr,acctAssetFromBytes)
	if err!=nil{
		return shim.Error("avaMoneyFromStr pustate failed")
	}
	err=stub.PutState(accIdToStr,acctAssetToBytes)
	if err!=nil{
		return shim.Error("avaMoneyToStr putstate failed")
	}
	var accFlowId int
	var result string
	if accFlowIdStr==""{
		accFlowId,result=getId(stub,"accFlow")
		if result!=success{
			return shim.Error(result)
		}
		accFlowIdStr=strconv.Itoa(accFlowId)
	}else{
		accFlowId,err=strconv.Atoi(accFlowIdStr)
		if err!=nil{
			return shim.Error("accFlowIdStr Atoi failed")
		}
	}
	result=saveAccFlow(stub,accFlowIdStr,accIdFromStr,"money",moneyStr,"1",contractNStr)
	if result!=success{
		return shim.Error(result)
	}
	accFlowIdStr=strconv.Itoa(accFlowId+1)
	result=saveAccFlow(stub,accFlowIdStr,accIdToStr,"money",moneyStr,"0",contractNStr)
	if result!=success{
		return shim.Error(result)
	}
	//保存最新的id
    err=stub.PutState("accFlow",[]byte(accFlowIdStr))
    if err!=nil{
    	return shim.Error("currIdStr putstate failed")
    }
	return shim.Success([]byte("done"))
}
//保存账户流水
func saveAccFlow(stub shim.ChaincodeStubInterface,accFlowIdStr,accIdStr,assetIdStr,assetNumStr,sType,contractN string)string{
	//获得当前时间
	var timeStr string
	timeBytes,err:=stub.GetState("currTime")
	if timeBytes!=nil{
		timeStr=string(timeBytes)
	}else{
		timeStr=time.Now().Format("20060102")
	}
	accFlow:=AccFlow{
		"accFlow",
		accFlowIdStr,
		accIdStr,
		assetIdStr,
		assetNumStr,
		sType,
		contractN,
		timeStr,
	}
	accFlowBytes,err:=json.Marshal(accFlow)
	if err!=nil{
		return "accFlowBytes marshal failed"
	}
	err=stub.PutState(accFlowIdStr,accFlowBytes)
	if err!=nil{
		return "accFlowBytes putstate failed"
	}
	return success
}
//将Id的初始化统一
func getId(stub shim.ChaincodeStubInterface,f string)(int,string){
	currIdByte,err:=stub.GetState(f)
	if err!=nil{
		return 0,"currIdByte getstate failed"
	}
	var currId int
	if currIdByte==nil{
		if f=="accFlow"{
			currId=100000
		}
	}else{
		currId,err=strconv.Atoi(string(currIdByte))
        if err!=nil{
    	    return 0,"currId Atoi failed;"+string(currIdByte)
        }
        currId+=1
    }
    currIdStr:=strconv.Itoa(currId)
    //保存最新的id
    err=stub.PutState(f,[]byte(currIdStr))
    if err!=nil{
    	return 0,"currIdStr putstate failed"
    }
	return currId,success
}
//解锁资金
//unlockMoney,accId,contractN,unlockMoney,accFlowIdStr
func (t *SimpleChaincode) unlockMoney(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=5{
		return shim.Error("args length is wrong")
	}
	if args[1]==""||args[2]==""||args[3]==""{
		return shim.Error("args can not be nil")
	}
	accIdStr:=args[1]
	contractNStr:=args[2]
	moneyStr:=args[3]
	accFlowIdStr:=args[4]
	unlockMoney,err:=strconv.Atoi(moneyStr)
    acctMoneyFrozenBytes,err:=stub.GetState(accIdStr+contractNStr)
    if err!=nil{
    	return shim.Error("frozenMoney getstate failed")
    }
    if acctMoneyFrozenBytes==nil{
    	return shim.Error("the accid and contractN doesn't have the frozenMoney")
    }
    var acctMoneyFrozen AcctMoneyFrozen
    err=json.Unmarshal(acctMoneyFrozenBytes,&acctMoneyFrozen)
    if err!=nil{
    	return shim.Error("acctMoneyFrozenBytes Unmarshal failed")
    }
    if acctMoneyFrozen.FrozenMoney<unlockMoney{
    	return shim.Error("unlockMoney is more than frozenMoney")
    }
    acctMoneyFrozen.FrozenMoney-=unlockMoney
    acctMoneyFrozenBytes,err=json.Marshal(acctMoneyFrozen)
    if err!=nil{
    	return shim.Error("acctMoneyFrozenBytes Marshal failed")
    }
    acctAsset,result:=getAcctAsset(stub,accIdStr)
    if result!=success{
    	return shim.Error(result)
    }
    acctAsset.AvaMoney+=unlockMoney
    acctAssetBytes,err:=json.Marshal(acctAsset)
    if err!=nil{
    	return shim.Error("acctAssetBytes Marshal failed")
    }
    //更新state
    err=stub.PutState(accIdStr+contractNStr,acctMoneyFrozenBytes)
    if err!=nil{
    	return shim.Error("frozenMoneyStr putstate failed")
    }
    err=stub.PutState(accIdStr,acctAssetBytes)
    if err!=nil{
    	return shim.Error("avaMoneyStr putstate failed")
    }
    var accFlowId int
    if accFlowIdStr==""{
    	accFlowId,result=getId(stub,"accFlow")
		if result!=success{
			return shim.Error(result)
		}
		accFlowIdStr=strconv.Itoa(accFlowId)
    }else{
    	accFlowId,err=strconv.Atoi(accFlowIdStr)
    	if err!=nil{
    		return shim.Error("accFlowIdStr Atoi failed")
    	}
    }
    result=saveAccFlow(stub,accFlowIdStr,accIdStr,"money",moneyStr,"3",contractNStr)
	if result!=success{
		return shim.Error(result)
	}
	err=stub.PutState("accFlow",[]byte(accFlowIdStr))
    return shim.Success([]byte("done"))
}
func getAcctAsset (stub shim.ChaincodeStubInterface,accIdStr string)(AcctAsset,string){
	var acctAsset AcctAsset
	acctAssetBytes,err:=stub.GetState(accIdStr)
    if err!=nil{
    	return acctAsset,"acctAssetBytes getstate failed"
    }
    
    if acctAssetBytes==nil{
    	acctAsset=AcctAsset{
    		"acctAsset",
    		accIdStr,
    		0,
    	}
    }else{
    	err=json.Unmarshal(acctAssetBytes,&acctAsset)
    	if err!=nil{
    		return acctAsset,"acctAssetBytes Unmarshal failed"
    	}
    }
    return acctAsset,success
}
//锁定资金
//lockMoney,accId,contractN,lockMoney,accFlowIdStr
func (t *SimpleChaincode) lockMoney(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=5{
		return shim.Error("args length is wrong")
	}
	if args[1]==""||args[2]==""{
		return shim.Error("args can not be nil")
	}
	accIdStr:=args[1]
	contractNStr:=args[2]
	moneyStr:=args[3]
	accFlowIdStr:=args[4]
	lockMoney,err:=strconv.Atoi(args[3])
	if err!=nil{
		return shim.Error("lockMoney atoi failed")
	}
	acctAsset,result:=getAcctAsset(stub,accIdStr)
	if result!=success{
		return shim.Error(result)
	}
	if acctAsset.AvaMoney<lockMoney{
		return shim.Error("lockMoney is more than avaMoney")
	}
	acctAsset.AvaMoney-=lockMoney
	acctAssetBytes,err:=json.Marshal(acctAsset)
	if err!=nil{
		return shim.Error("acctAssetBytes Marshal failed")
	}
	err=stub.PutState(accIdStr,acctAssetBytes)
	if err!=nil{
		return shim.Error("acctAssetBytes putstate failed")
	}
	//更新资金冻结表
	acctMoneyFrozen,result:=getAcctMoneyFrozen(stub,accIdStr,contractNStr)
    if result!=success{
    	return shim.Error(result)
    }
    acctMoneyFrozen.FrozenMoney+=lockMoney
    acctMoneyFrozenBytes,err:=json.Marshal(acctMoneyFrozen)
    if err!=nil{
    	return shim.Error("acctMoneyFrozenBytes Marshal failed")
    }
    err=stub.PutState(accIdStr+contractNStr,acctMoneyFrozenBytes)
    if err!=nil{
    	return shim.Error("acctMoneyFrozenBytes putstate failed")
    }
    var accFlowId int
    if accFlowIdStr==""{
    	accFlowId,result=getId(stub,"accFlow")
		if result!=success{
			return shim.Error(result)
		}
	    accFlowIdStr=strconv.Itoa(accFlowId)
    }
    result=saveAccFlow(stub,accFlowIdStr,accIdStr,"money",moneyStr,"2",contractNStr)
	if result!=success{
		return shim.Error(result)
	}
	err=stub.PutState("accFlow",[]byte(accFlowIdStr))
	if err!=nil{
		return shim.Error("accFlowIdStr PutState failed")
	}
	return shim.Success([]byte(string(acctAssetBytes)+","+string(acctMoneyFrozenBytes)))
}
func getAcctMoneyFrozen(stub shim.ChaincodeStubInterface,accIdStr,contractNStr string)(AcctMoneyFrozen,string){
    var acctMoneyFrozen AcctMoneyFrozen
	acctMoneyFrozenBytes,err:=stub.GetState(accIdStr+contractNStr)
    if err!=nil{
    	return acctMoneyFrozen,"acctMoneyFrozenBytes getstate failed"
    }
    if acctMoneyFrozenBytes==nil{
    	acctMoneyFrozen=AcctMoneyFrozen{
    		"acctMoneyFrozen",
    		accIdStr,
    		contractNStr,
    		0,
    	}
    }else{
    	err:=json.Unmarshal(acctMoneyFrozenBytes,&acctMoneyFrozen)
        if err!=nil{
        	return acctMoneyFrozen,"acctMoneyFrozenBytes Unmarshal failed"
        }
    }
    return acctMoneyFrozen,success
}
//划转证券moveStock,accIdFrom,accIdTo,AHoldStr,productCode,moveNum,contractNStr,accFlowIdStr
func (t *SimpleChaincode) moveStock(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=8{
		return shim.Error("args length is wrong")
	}
	if args[2]==""||args[3]==""||args[5]==""||args[6]==""{
		return shim.Error("args can not be nil")
	}
	accIdFromStr:=args[1]
	accIdToStr:=args[2]
	FromHoldStr:=args[3]
	productCodeStr:=args[4]
	moveNumStr:=args[5]
	contractNStr:=args[6]
	accFlowIdStr:=args[7]
	moveNum,err:=strconv.Atoi(moveNumStr)
	if err!=nil{
		return shim.Error("moveNum atoi failed")
	}
	var sseHoldFrom SSEHold
	if FromHoldStr!=""{
		err=json.Unmarshal([]byte(FromHoldStr),&sseHoldFrom)
		if err!=nil{
			return shim.Error("sseHoldBytes unmarshal failed")
		}
	}else{
		sseHoldFromBytes,err:=stub.GetState(accIdFromStr+productCodeStr)
		if err!=nil{
			return shim.Error("sseHoldFromBytes getstate failed")
		}
		if sseHoldFromBytes==nil{
			return shim.Error("there is not the hold")
		}
		err=json.Unmarshal(sseHoldFromBytes,&sseHoldFrom)
		if err!=nil{
			return shim.Error("sseHoldBytes unmarshal failed")
		}
	}
	sseHoldToBytes,err:=stub.GetState(accIdToStr+productCodeStr)
	if err!=nil{
		return shim.Error("sseHoldToBytes getstate failed")
	}
	var sseHoldTo SSEHold
	if sseHoldToBytes==nil{
		sseHoldTo=SSEHold{
			"sseHold",
			accIdToStr,
			productCodeStr,
			0,
			0,
		}
	}else{
		err=json.Unmarshal(sseHoldToBytes,&sseHoldTo)
		if err!=nil{
			return shim.Error("sseHoldBytes unmarshal failed")
		}
	}
	if sseHoldFrom.HoldNum<moveNum{
		return shim.Error("moveNum is more than sseHoldFrom holdnum")
	}
	sseHoldFrom.HoldNum-=moveNum
	sseHoldTo.HoldNum+=moveNum
	sseHoldFromBytes,err:=json.Marshal(sseHoldFrom)
	if err!=nil{
		return shim.Error("sseHoldFromBytes marshal failed")
	}
	sseHoldToBytes,err=json.Marshal(sseHoldTo)
	if err!=nil{
		return shim.Error("sseHoldToBytes marshal failed")
	}
	err=stub.PutState(accIdFromStr+productCodeStr,sseHoldFromBytes)
	if err!=nil{
		return shim.Error("sseHoldFromBytes putstate failed")
	}
	err=stub.PutState(accIdToStr+productCodeStr,sseHoldToBytes)
	if err!=nil{
		return shim.Error("sseHoldToBytes putstate failed")
	}
	var accFlowId int
	var result string
	if accFlowIdStr==""{
		accFlowId,result=getId(stub,"accFlow")
		if result!=success{
			return shim.Error(result)
		}
		accFlowIdStr=strconv.Itoa(accFlowId)
	}else{
		accFlowId,err=strconv.Atoi(accFlowIdStr)
		if err!=nil{
			return shim.Error("accFlowIdStr Atoi failed")
		}
	}
    result=saveAccFlow(stub,accFlowIdStr,accIdFromStr,productCodeStr,moveNumStr,"1",contractNStr)
	if result!=success{
		return shim.Error(result)
	}
	accFlowIdStr=strconv.Itoa(accFlowId+1)
	result=saveAccFlow(stub,accFlowIdStr,accIdToStr,productCodeStr,moveNumStr,"0",contractNStr)
	if result!=success{
		return shim.Error(result)
	}
	//保存最新的id
    err=stub.PutState("accFlow",[]byte(accFlowIdStr))
    if err!=nil{
    	return shim.Error("currIdStr putstate failed")
    }
	return shim.Success([]byte("done"))
}
//解锁证券
//unlockStock,accId,contractN,productCode,unlockNum,accFlowIdStr
func (t *SimpleChaincode) unlockStock(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=6{
		return shim.Error("args length is wrong")
	}
	if args[1]==""||args[2]==""||args[3]==""||args[4]==""{
		return shim.Error("args can not be nil")
	}
	accIdStr:=args[1]
	contractNStr:=args[2]
	productCodeStr:=args[3]
	unlockNumStr:=args[4]
	accFlowIdStr:=args[5]
    sseHoldFrozen,result:=getSSEHoldFrozen(stub,accIdStr,productCodeStr,contractNStr)
	if result!=success{
		return shim.Error(result)
	}
    unlockNum,err:=strconv.Atoi(unlockNumStr)
    if err!=nil{
    	return shim.Error("lockNum atoi failed"+unlockNumStr)
    }
    if sseHoldFrozen.FrozenSecNum<unlockNum{
    	return shim.Error("unlockNum is more than the lockNum")
    }
    sseHoldFrozen.FrozenSecNum-=unlockNum
    sseHoldFrozenBytes,err:=json.Marshal(sseHoldFrozen)
    if err!=nil{
    	return shim.Error("sseHoldFrozen Marshal sseHoldFrozenBytes failed")
    }
    err=stub.PutState(accIdStr+productCodeStr+contractNStr,sseHoldFrozenBytes)
    if err!=nil{
    	return shim.Error("lockNumStr putstate failed")
    }
    sseHoldBytes,err:=stub.GetState(accIdStr+productCodeStr)
    if err!=nil{
    	return shim.Error("sseHoldBytes getstate failed")
    }
    var sseHold SSEHold
    err=json.Unmarshal(sseHoldBytes,&sseHold)
    if err!=nil{
    	return shim.Error("sseHoldBytes unmarshal failed")
    }
    if sseHold.FrozenSecNum<unlockNum{
    	return shim.Error("unlockNum is more than ssehold frozenSecNum")
    }
    sseHold.HoldNum+=unlockNum
    sseHold.FrozenSecNum-=unlockNum
    sseHoldBytes,err=json.Marshal(sseHold)
    if err!=nil{
    	return shim.Error("sseHoldBytes marshal failed")
    }
    err=stub.PutState(accIdStr+productCodeStr,sseHoldBytes)
    if err!=nil{
    	return shim.Error("sseHoldBytes putstate failed")
    }
    var accFlowId int
    if accFlowIdStr==""{
    	accFlowId,result=getId(stub,"accFlow")
		if result!=success{
			return shim.Error(result)
		}
		accFlowIdStr=strconv.Itoa(accFlowId)
    }
    result=saveAccFlow(stub,accFlowIdStr,accIdStr,productCodeStr,unlockNumStr,"3",contractNStr)
	if result!=success{
		return shim.Error(result)
	}
	err=stub.PutState("accFlow",[]byte(accFlowIdStr))
	if err!=nil{
		return shim.Error("accFlowIdStr PutState failed")
	}
    return shim.Success(sseHoldBytes)
}
//锁定证券
//lockStock,accIdStr,contractN,productCode,lockNum,accFlowIdStr
func (t *SimpleChaincode) lockStock(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=6{
		return shim.Error("args length is wrong")
	}
	if args[1]==""||args[2]==""||args[3]==""||args[4]==""{
    	return shim.Error("args can not be nil")
    }
	accIdStr:=args[1]
	contractNStr:=args[2]
	productCodeStr:=args[3]
	lockNumStr:=args[4]
	accFlowIdStr:=args[5]
	sseHoldBytes,err:=stub.GetState(accIdStr+productCodeStr)
	if err!=nil{
		return shim.Error("sseHoldBytes getstate failed")
	}
	if sseHoldBytes==nil{
		return shim.Error("there is not the ssehold")
	}
	var ssehold SSEHold
	err=json.Unmarshal(sseHoldBytes,&ssehold)
	lockNum,err:=strconv.Atoi(lockNumStr)
	if err!=nil{
		return shim.Error("lockNum atoi failed")
	}
	if ssehold.HoldNum<lockNum{
		return shim.Error("ssehold is not enough")
	}
	ssehold.HoldNum-=lockNum
	ssehold.FrozenSecNum+=lockNum
	sseHoldBytes,err=json.Marshal(ssehold)
	if err!=nil{
		return shim.Error("sseHoldBytes marshal failed")
	}
	err=stub.PutState(accIdStr+productCodeStr,sseHoldBytes)
	if err!=nil{
		return shim.Error("sseHoldBytes putstate failed")
	}
	sseHoldFrozen,result:=getSSEHoldFrozen(stub,accIdStr,productCodeStr,contractNStr)
	if result!=success{
		return shim.Error(result)
	}
	sseHoldFrozen.FrozenSecNum+=lockNum
	sseHoldFrozenBytes,err:=json.Marshal(sseHoldFrozen)
	if err!=nil{
		return shim.Error("sseHoldFrozenBytes Marshal failed")
	}
	err=stub.PutState(accIdStr+productCodeStr+contractNStr,sseHoldFrozenBytes)
	if err!=nil{
		return shim.Error("sseHoldFrozenBytes putstate failed")
	}
	if accFlowIdStr==""{
		accFlowId,result:=getId(stub,"accFlow")
		if result!=success{
			return shim.Error(result)
		}
		accFlowIdStr=strconv.Itoa(accFlowId)
	}
    result=saveAccFlow(stub,accFlowIdStr,accIdStr,productCodeStr,lockNumStr,"2",contractNStr)
	if result!=success{
		return shim.Error(result)
	}
	err=stub.PutState("accFlow",[]byte(accFlowIdStr))
	if err!=nil{
		return shim.Error("accFlowIdStr PutState failed")
	}
	return shim.Success([]byte(string(sseHoldBytes)+","+string(sseHoldFrozenBytes)))
}
func getSSEHoldFrozen(stub shim.ChaincodeStubInterface,accIdStr,productCode,contractNStr string)(SSEHoldFrozen,string){
	var sseHoldFrozen SSEHoldFrozen
	sseHoldFrozenBytes,err:=stub.GetState(accIdStr+productCode+contractNStr)
	if err!=nil{
		return sseHoldFrozen,"sseHoldFrozenBytes GetState failed"
	}
	if sseHoldFrozenBytes==nil{
		sseHoldFrozen=SSEHoldFrozen{
			"sseHoldFrozen",
			accIdStr,
			productCode,
			contractNStr,
			0,
		}
	}else{
		err=json.Unmarshal(sseHoldFrozenBytes,&sseHoldFrozen)
		if err!=nil{
			return sseHoldFrozen,"sseHoldFrozenBytes Unmarshal failed"
		}
	}
	return sseHoldFrozen,success
}
//查询可用余额
//getAccMoney,accId,pw,moneyType
//moneyType=1代表可用余额+冻结余额，2代表可用余额，3代表冻结余额
func (t *SimpleChaincode) getAccMoney(stub shim.ChaincodeStubInterface,args []string)pb.Response{
    if len(args)!=4{
    	return shim.Error("args length is wrong")
    }
    if args[1]==""||args[2]==""{
    	return shim.Error("args can not be nil")
    }
    accIdStr:=args[1]
    pw:=args[2]
    moneyType:=args[3]
    result:=getAccPermission(stub,accIdStr,pw)
    if result!=success{
    	return shim.Error(result)
    }
    if moneyType=="1"{
    	avaMoney,err:=getMoney(stub,accIdStr,"2")
    	if err!=success{
    		return shim.Error(err)
    	}
    	freezeMoney,err:=getMoney(stub,accIdStr,"3")
    	if err!=success{
    		return shim.Error(err)
    	}
    	allMoney:=avaMoney+freezeMoney
    	allMoneyStr:=strconv.Itoa(allMoney)
    	return shim.Success([]byte(allMoneyStr))
    }else if moneyType=="2"{
    	avaMoney,err:=getMoney(stub,accIdStr,"2")
    	if err!=success{
    		return shim.Error(err)
    	}
    	avaMoneyStr:=strconv.Itoa(avaMoney)
    	return shim.Success([]byte(avaMoneyStr))
    }else if moneyType=="3"{
    	freezeMoney,err:=getMoney(stub,accIdStr,"3")
    	if err!=success{
    		return shim.Error(err)
    	}
    	freezeMoneyStr:=strconv.Itoa(freezeMoney)
    	return shim.Success([]byte(freezeMoneyStr))
    }
    return shim.Error("moneytype must be 1 or 2 or 3")
}
//获得余额 2代表可用余额，3代表冻结余额
func getMoney(stub shim.ChaincodeStubInterface,accIdStr,moneyType string)(int,string){
	if moneyType=="2"{
		acctAsset,result:=getAcctAsset(stub,accIdStr)
		if result!=success{
			return 0,result
		}
		return acctAsset.AvaMoney,success
	}else if moneyType =="3"{
		//获得所有合约编码
        contractNBytes,err:=stub.GetState(accIdStr+"contractN")
        if err!=nil{
        	return 0,"contractNBytes getstate failed"
        }
        contractNArr:=strings.Fields(string(contractNBytes))
        //根据合约编码和账户获得冻结资金
        var allFreezeMoney int
        for _,contractN:=range contractNArr{
        	acctMoneyFrozenBytes,err:=stub.GetState(accIdStr+contractN)
            if err!=nil{
            	return 0,"acctMoneyFrozenBytes getstate failed"
            }
            if acctMoneyFrozenBytes==nil{
            	return 0,"acctMoneyFrozenBytes is nil"
            }
            var acctMoneyFrozen AcctMoneyFrozen
            err=json.Unmarshal(acctMoneyFrozenBytes,&acctMoneyFrozen)
            if err!=nil{
            	return 0,"acctMoneyFrozenBytes Unmarshal failed"
            }
            allFreezeMoney+=acctMoneyFrozen.FrozenMoney
        }
        return allFreezeMoney,success
	}
	return 0,"error"
	
}
//返回价格时间表
//getStockPrice,time,productCode
func (t *SimpleChaincode) getStockPrice(stub shim.ChaincodeStubInterface,args []string)pb.Response{
    if len(args)!=3{
    	return shim.Error("args length is wrong")
    }
    if args[1]==""||args[2]==""{
    	return shim.Error("args can not be nil")
    }
    time:=args[1]
    productCode:=args[2]
	priceBytes,err:=stub.GetState(time+productCode)
	if err!=nil{
		return shim.Error("priceBytes getstate failed")
	}
	if priceBytes==nil{
		return shim.Error("there is not the price")
	}
	return shim.Success(priceBytes)
}
//getCurrTime
func (t *SimpleChaincode) getCurrTime(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	timeBytes,err:=stub.GetState("currTime")
	if err!=nil{
		return shim.Error("time getstate failed")
	}
	if timeBytes!=nil{
		return shim.Success(timeBytes)
	}
	currTime:=time.Now().Unix()
	currTm:=time.Unix(currTime,0)
	return shim.Success([]byte(currTm.Format("20060102")))
}
// Query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var A string // Entities
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}
	A = args[1]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
