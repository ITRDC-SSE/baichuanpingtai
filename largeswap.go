package main
import (
        "encoding/json"
        "github.com/hyperledger/fabric/common/util"
		"fmt"
		"strconv" 
		"github.com/hyperledger/fabric/core/chaincode/shim"
		pb "github.com/hyperledger/fabric/protos/peer"
)
const success="success"
// channelId:="mychannel"
// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
type PriceTable struct{
	ObjectType string `json:"docType"`
	Time string `json:"time"`
	ProductCode string `json:"productCode"`
	Price int `json:"price"`
}
//持仓表
type SSEHold struct{
	ObjectType string `json:"docType"`
	AccId string `json:"accId"`
	ProductCode string `json:"productCode"`
	HoldNum int `json:"holdNum"`
	FrozenSecNum int `json:"frozenSecNum"`
}
type AcctAsset struct{
	ObjectType string `json:"docType"`
	AcctId string `json:"acctId"`
	AvaMoney int `json:"avaMoney"`
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

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting at least 2")
	}
	//设计合约
	if args[0]=="bigswap"{
		return t.bigswap(stub,args)
	}
	return shim.Error("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'")
}
//设计合约
//bigSwp,contractN,actionType,accId,pw,accIdN,channelId,chaincodeToCall,
func (t *SimpleChaincode) bigswap(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=8{
		return shim.Error("args numbers is wrong")
	}
	if args[1]==""||args[2]==""||args[3]==""||args[4]==""||args[5]==""||args[6]==""||args[7]==""{
		return shim.Error("args can not be nil")
	}
	contractNStr:=args[1]
	actionType:=args[2]
	accIdStr:=args[3]
	pwStr:=args[4]
	accIdNStr:=args[5]
	channelId:=args[6]
	chaincodeToCall:=args[7]
	contractHold,result:=getContractHold(stub,contractNStr,chaincodeToCall,channelId)
    if result!=success{
        return shim.Error(result)
    }
    var accIdN string
    if accIdNStr=="1"{
    	accIdN=contractHold.AccIdA
    }else if accIdNStr=="2"{
    	accIdN=contractHold.AccIdB
    }else if accIdNStr=="3"{
    	accIdN=contractHold.AccIdC
	}else if accIdNStr=="4"{
		accIdN=contractHold.AccIdD
	}else if accIdNStr=="5"{
		accIdN=contractHold.AccIdE
	}
	if accIdStr!=accIdN{
		return shim.Error("accId is wrong")
	}
	//获得权限
	result=getPermission(stub,accIdStr,pwStr,chaincodeToCall,channelId)
    if result!=success{
        return shim.Error(result)
    }
    //获得时间
    ti:=getCurrTime(stub,chaincodeToCall,channelId)
	if contractHold.ContractStatus=="0"&&actionType=="1"{
		if accIdStr!=contractHold.AccIdA{
			return shim.Error("accId is wrong")
		}
        contractHold.LastSwapTime=getCurrTime(stub,chaincodeToCall,channelId)
        contractHold.ContractStatus="1"//开始合约
	}else if contractHold.ContractStatus=="1"&&actionType=="5"{//互操作
		if accIdStr!=contractHold.AccIdA&&accIdStr!=contractHold.AccIdB{
			return shim.Error("accId is wrong")
		}
		//获得股票价格表
		currStockPrice,result:=getStockPrice(stub,ti,"sh0003",chaincodeToCall,channelId)
		if result!=success{
			return shim.Error(result)
		}
		if contractHold.LastSwapTime==""{
			return shim.Error("contractHold.LastSwapTime is nil")
		}
		lastStockPrice,result:=getStockPrice(stub,contractHold.LastSwapTime,"sh0003",chaincodeToCall,channelId)
		if result!=success{
			return shim.Error(result)
		}
        amtA:=(currStockPrice.Price-lastStockPrice.Price)*1000
        currTime,err:=strconv.Atoi(currStockPrice.Time)
        if err!=nil{
        	return shim.Error("currTime Atoi failed")
        }
        lastTime,err:=strconv.Atoi(lastStockPrice.Time)
        if err!=nil{
        	return shim.Error("lastTime Atoi failed")
        }
        amtB:=(currTime-lastTime)*100
        //获得A账户的资金
        acctAssetA,result:=getAcctAsset(stub,contractHold.AccIdA,chaincodeToCall,channelId)
	    if result!=success{
	        return shim.Error(result)
	    }
	    //看A是否有足够的钱
	    if acctAssetA.AvaMoney<amtB{
	    	return shim.Error("accIdA money is not enough")
	    }
	    //获得B账户的资金
        acctAssetB,result:=getAcctAsset(stub,contractHold.AccIdB,chaincodeToCall,channelId)
	    if result!=success{
	        return shim.Error(result)
	    }
	    //看B是否有足够的钱
	    if acctAssetB.AvaMoney<amtA{
	    	return shim.Error("accIdB money is not enough")
	    }
        acctAssetA.AvaMoney=acctAssetA.AvaMoney-amtB+amtA
        acctAssetB.AvaMoney=acctAssetB.AvaMoney+amtB-amtA
        //保存账户资金
        result=saveAcctAsset(stub,acctAssetA,contractHold.AccIdA,chaincodeToCall,channelId)
	    if result!=success{
	        return shim.Error(result)
	    }
	    result=saveAcctAsset(stub,acctAssetB,contractHold.AccIdB,chaincodeToCall,channelId)
	    if result!=success{
	        return shim.Error(result)
	    }
	    //保存账户流水
	    //获得当前时间
	    timeStr:=getCurrTime(stub,chaincodeToCall,channelId)
	    //获得当前流水编号
	    currAccFlowIdStr,result:=getAccFlowId(stub,chaincodeToCall,channelId)
        if result!=success{
            return shim.Error(result)
        }
        currAccFlowId,err:=strconv.Atoi(currAccFlowIdStr)
        if err!=nil{
        	return shim.Error("currAccFlowId Atoi failed")
        }
        //A给B转钱,A的账户流水
    	accFlowAAsset:=AccFlow{
	    	"accFlow",
	    	currAccFlowIdStr,
	    	contractHold.AccIdA,
	    	"money",
	    	strconv.Itoa(amtA),
	    	"1",
	    	contractNStr,
	    	timeStr,
        }
        //保存账户流水
        result=saveAccFlow(stub,accFlowAAsset,accFlowAAsset.AccFlowId,chaincodeToCall,channelId)
	    if result!=success{
	        return shim.Error(result)
	    }
        //A给B转钱,B的账户流水
        currAccFlowId=currAccFlowId+1
        currAccFlowIdStr=strconv.Itoa(currAccFlowId)
        accFlowBAsset:=AccFlow{
	    	"accFlow",
	    	currAccFlowIdStr,
	    	contractHold.AccIdB,
	    	"money",
	    	strconv.Itoa(amtA),
	    	"0",
	    	contractNStr,
	    	timeStr,
        }
        //保存账户流水
        result=saveAccFlow(stub,accFlowBAsset,accFlowBAsset.AccFlowId,chaincodeToCall,channelId)
	    if result!=success{
	        return shim.Error(result)
	    }
	    //B给A转钱,A的账户流水
	    currAccFlowId=currAccFlowId+1
	    currAccFlowIdStr=strconv.Itoa(currAccFlowId)
	    accFlowAAsset=AccFlow{
	    	"accFlow",
	    	currAccFlowIdStr,
	    	contractHold.AccIdA,
	    	"money",
	    	strconv.Itoa(amtA),
	    	"0",
	    	contractNStr,
	    	timeStr,
        }
        //保存账户流水
        result=saveAccFlow(stub,accFlowAAsset,accFlowAAsset.AccFlowId,chaincodeToCall,channelId)
	    if result!=success{
	        return shim.Error(result)
	    }
	    //B给A转钱，B的账户流水
        currAccFlowId=currAccFlowId+1
        currAccFlowIdStr=strconv.Itoa(currAccFlowId)
        accFlowBAsset=AccFlow{
	    	"accFlow",
	    	currAccFlowIdStr,
	    	contractHold.AccIdB,
	    	"money",
	    	strconv.Itoa(amtA),
	    	"1",
	    	contractNStr,
	    	timeStr,
        }
        //保存账户流水
        result=saveAccFlow(stub,accFlowBAsset,accFlowBAsset.AccFlowId,chaincodeToCall,channelId)
	    if result!=success{
	        return shim.Error(result)
	    }
        //保存最新的账户流水号
        result=saveAccFlowId(stub,currAccFlowIdStr,chaincodeToCall,channelId)
        if result!=success{
            return shim.Error(result)
        }
        //不改变状态
	    // contractHold.ContractStatus="4"
	}else if contractHold.LastSwapTime==ti&&actionType=="4"{
        contractHold.ContractStatus="4"
    }else if contractHold.ContractStatus=="4"{
        return shim.Error("contract has finished")
    }else{
		return shim.Error("unvalid Status")
	}
	//保存最新的状态
	result=saveContractHold(stub,contractHold,contractNStr,chaincodeToCall,channelId)
    if result!=success{
            return shim.Error(result)
    }
	return shim.Success([]byte(success))
}
//获得系统时间
func getCurrTime(stub shim.ChaincodeStubInterface,chaincodeToCall,channelId string)string{
    invokeArgs:=util.ToChaincodeArgs("invoke","getCurrTime")
    response:= stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
    if response.Status!=shim.OK{
        errStr := fmt.Sprintf("getCurrTime Failed to invoke chaincode. Got error: %s", string(response.Payload))
        fmt.Printf(errStr)
        return errStr
    }
    return string(response.Payload)
}
//保存合约最新的状态
func saveContractHold(stub shim.ChaincodeStubInterface,contractHold ContractHold,contractNStr,chaincodeToCall,channelId string)string{
    contractHoldBytes,err:=json.Marshal(contractHold)
    if err!=nil{
        return "contractHold marshal failed"
    }
    fmt.Println("bigdealtwo"+";"+string(contractHoldBytes))
    invokeArgs:=util.ToChaincodeArgs("invoke","saveContractHold",contractNStr,string(contractHoldBytes))
    response := stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
    if response.Status!=shim.OK{
        errStr := fmt.Sprintf("7 Failed to invoke chaincode. Got error: %s", string(response.Payload))
        fmt.Printf(errStr)
        return errStr
    }
    return success
}
//保存最新的账户流水编号
func saveAccFlowId(stub shim.ChaincodeStubInterface,accFlowIdStr,chaincodeToCall,channelId string)string{
    invokeArgs:=util.ToChaincodeArgs("invoke","saveAccFlowId",accFlowIdStr)
    response:=stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
    if response.Status!=shim.OK{
        errStr := fmt.Sprintf("B saveAcctAsset Failed to invoke chaincode. Got error: %s", string(response.Payload))
        fmt.Printf(errStr)
        return errStr
    }
    return success
}
//获得最新的账户流水编号
func getAccFlowId(stub shim.ChaincodeStubInterface,chaincodeToCall,channelId string)(string,string){
    invokeArgs:=util.ToChaincodeArgs("invoke","getAccFlowId")
    response := stub.InvokeChaincode(chaincodeToCall,invokeArgs,channelId)
    if response.Status!=shim.OK{
        errStr := fmt.Sprintf("2 getContractHold Failed to invoke chaincode. Got error: %s", string(response.Payload))
        fmt.Printf(errStr)
        return "0",errStr
    }
    return string(response.Payload),success
}
//获得股票价格表
func getStockPrice(stub shim.ChaincodeStubInterface,ti,productCode,chaincodeToCall,channelId string)(PriceTable,string){
	var priceTable PriceTable
	invokeArgs:=util.ToChaincodeArgs("invoke","getStockPrice",ti,productCode)
	response := stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
	if response.Status!=shim.OK{
    	errStr := fmt.Sprintf("A B moveMoney Failed to invoke chaincode. Got error: %s", string(response.Payload))
		fmt.Printf(errStr)
		return priceTable,errStr
    }
    currStockPriceBytes:=response.Payload
    err:=json.Unmarshal(currStockPriceBytes,&priceTable)
    if err!=nil{
    	return priceTable,"currStockPriceBytes Unmarshal failed"
    }
    return priceTable,success
}
//获得权限
func getPermission(stub shim.ChaincodeStubInterface,accIdStr,pwStr,chaincodeToCall,channelId string)string{
    invokeArgs:=util.ToChaincodeArgs("invoke","getAccPermissionOut",accIdStr,pwStr)
    response := stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
    if response.Status!=shim.OK{
        errStr := fmt.Sprintf(" 1 getContractHold Failed to invoke chaincode. Got error: %s", string(response.Payload))
        fmt.Printf(errStr)
        return errStr+"getPermission"
    }
    return success
}
//获得合约持仓
func getContractHold(stub shim.ChaincodeStubInterface,contractNStr,chaincodeToCall,channelId string)(ContractHold,string){
    var contractHold ContractHold
    invokeArgs:=util.ToChaincodeArgs("invoke","getContractHold",contractNStr)
    response := stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
    if response.Status!=shim.OK{
        errStr := fmt.Sprintf("5 getContractHold Failed to invoke chaincode. Got error: %s", string(response.Payload))
        fmt.Printf(errStr)
        return contractHold,errStr
    }
    contractHoldBytes:=response.Payload
    err:=json.Unmarshal(contractHoldBytes,&contractHold)
    if err!=nil{
        return contractHold,"contractHoldBytes Unmarshal failed"
    }
    return contractHold,success
}
//保存账户流水
func saveAccFlow(stub shim.ChaincodeStubInterface,accFlow AccFlow,accIdStr,chaincodeToCall,channelId string)string{
    accFlowBytes,err:=json.Marshal(accFlow)
    if err!=nil{
        return "accFlowBytes Marshal failed"
    }
    invokeArgs:=util.ToChaincodeArgs("invoke","saveAccFlowOut",accIdStr,string(accFlowBytes))
    response := stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
    if response.Status!=shim.OK{
        errStr := fmt.Sprintf("getContractHold Failed to invoke chaincode. Got error: %s", string(response.Payload))
        fmt.Printf(errStr)
        return errStr
    }
    return success
}
//获得资金账户
func getAcctAsset(stub shim.ChaincodeStubInterface,accIdStr,chaincodeToCall,channelId string)(AcctAsset,string){
    var acctAsset AcctAsset
    invokeArgs:=util.ToChaincodeArgs("invoke","getAcctAsset",accIdStr)
    response:= stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
    if response.Status!=shim.OK{
        errStr := fmt.Sprintf("getContractHold Failed to invoke chaincode. Got error: %s", string(response.Payload))
        fmt.Printf(errStr)
        return acctAsset,errStr
    }
    acctAssetBytes:=response.Payload
    err:=json.Unmarshal(acctAssetBytes,&acctAsset)
    if err!=nil{
        return acctAsset,"acctAssetABytes Unmarshal failed"
    }
    return acctAsset,success
}
//保存资金账户
func saveAcctAsset(stub shim.ChaincodeStubInterface,acctAsset AcctAsset,accIdStr,chaincodeToCall,channelId string)string{
    acctAssetBytes,err:=json.Marshal(acctAsset)
    if err!=nil{
        return "acctAssetBytes Marshal failed"
    }
    invokeArgs:=util.ToChaincodeArgs("invoke","saveAcctAsset",accIdStr,string(acctAssetBytes))
    response:= stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
    if response.Status!=shim.OK{
        errStr := fmt.Sprintf("getContractHold Failed to invoke chaincode. Got error: %s", string(response.Payload))
        fmt.Printf(errStr)
        return errStr
    }
    return success
}
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
