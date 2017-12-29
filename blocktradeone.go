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
type ContractHold struct{
	ObjectType string `json:"docType"`
	ContractN string `json:"contractN"`
	ContractCode string `json:"contractCode"`
	ContractStatus string `json:"contractStatus"`//"0"代表等待开始，"1"代表开始合约，"2"代表已经开始,
	                    //"3"代表行权，"4"代表已终结,"5"代表互操作
	ContractFunctionName string `json:"contractFunctionName"`
    AccIdA string `json:"accIdA"`//合约关联方A
    AccIdB string `json:"accIdB"`
    AccIdC string `json:"accIdC"`
    AccIdD string `json:"accIdD"`
    AccIdE string `json:"accIdE"`
    CorContractNA string `json:"corContractNA"`//关联合约
    CorContractNB string `json:"corContractNB"`
    CorContractNC string `json:"corContractNC"`
    actionType string `json:"actionType"`
    LastSwapTime string `json:"lastSwapTime"`
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
	if args[0]=="blocktradeone"{
		return t.blocktradeone(stub,args)
	}
	return shim.Error("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'")
}
//设计合约
//bigdealone,contractN,actionType,accId,pw,accIdN,channelId,chaincodeToCall,
func (t *SimpleChaincode) blocktradeone(stub shim.ChaincodeStubInterface,args []string)pb.Response{
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
	invokeArgs:=util.ToChaincodeArgs("invoke","getContractHold",contractNStr)
	response := stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
	if response.Status!=shim.OK{
		errStr := fmt.Sprintf("getContractHold Failed to invoke chaincode. Got error: %s", string(response.Payload))
		fmt.Printf(errStr)
		return shim.Error(errStr)
	}
	contractHoldBytes:=response.Payload
    var contractHold ContractHold
    err:=json.Unmarshal(contractHoldBytes,&contractHold)
    if err!=nil{
    	return shim.Error("contractHoldBytes Unmarshal failed")
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
	result:=getPermission(stub,accIdStr,pwStr,chaincodeToCall,channelId)
    if result!=success{
        return shim.Error(result)
    }
	//账户流水
	var accFlowAAsset,accFlowBAsset,accFlowAHold,accFlowBHold AccFlow
	//获得最新的账户流水编号
	accFlowIdStr,result:=getAccFlowId(stub,chaincodeToCall,channelId)
    if result!=success{
        return shim.Error(result)
    }
	accFlowId,err:=strconv.Atoi(accFlowIdStr)
	if err!=nil{
		return shim.Error("accFlowId Atoi failed")
	}
	var acctAssetA,acctAssetB AcctAsset
	var sseHoldA,sseHoldB SSEHold
	if contractHold.ContractStatus=="0"&&actionType=="1"{
		if accIdStr!=contractHold.AccIdA{
			return shim.Error("accId is wrong")
		}
		//获得A的资金账户
		acctAssetA,result=getAcctAsset(stub,contractHold.AccIdA,chaincodeToCall,channelId)
		if result!=success{
			return shim.Error(result)
		}
		acctAssetB,result=getAcctAsset(stub,contractHold.AccIdB,chaincodeToCall,channelId)
		if result!=success{
			return shim.Error(result)
		}
	    //判断资金是否够锁钱
	    if acctAssetA.AvaMoney<65000{
	    	return shim.Error("accIdA avamoney is not enough")
	    }
	    //获得账户B的持仓
	    sseHoldB,result=getSSEHoldByAccAndProduct(stub,contractHold.AccIdB,"sh0003",chaincodeToCall,channelId)
	    if result!=success{
	    	return shim.Error(result)
	    }
	    sseHoldA,result=getSSEHoldByAccAndProduct(stub,contractHold.AccIdA,"sh0003",chaincodeToCall,channelId)
	    if result!=success{
	    	return shim.Error(result)
	    }
        //判断是否有足够的持仓数据
        if sseHoldB.HoldNum<1000{
        	return shim.Error("accountB holdNum is not enough")
        }
        //A给B转钱
        acctAssetA.AvaMoney-=65000
        acctAssetB.AvaMoney+=65000
	    //B给A转券
	    sseHoldA.HoldNum+=1000
	    sseHoldB.HoldNum-=1000
	    //账户流水
        timeStr:=getCurrTime(stub,chaincodeToCall,channelId)
        accFlowAAsset=AccFlow{
        	"accFlow",
        	accFlowIdStr,
        	contractHold.AccIdA,
        	"money",
        	"65000",
        	"1",
        	contractNStr,
        	timeStr,
        }
        accFlowId=accFlowId+1
        accFlowIdStr=strconv.Itoa(accFlowId)
        accFlowBAsset=AccFlow{
        	"accFlow",
        	accFlowIdStr,
        	contractHold.AccIdB,
        	"money",
        	"65000",
        	"0",
        	contractNStr,
        	timeStr,
        }
	    //账户流水
	    accFlowId=accFlowId+1
        accFlowIdStr=strconv.Itoa(accFlowId)
        accFlowAHold=AccFlow{
        	"accFlow",
        	accFlowIdStr,
        	contractHold.AccIdA,
        	"sh0003",
        	"1000",
        	"0",
        	contractNStr,
        	timeStr,
        }
        accFlowId=accFlowId+1
        accFlowIdStr=strconv.Itoa(accFlowId)
        accFlowBHold=AccFlow{
        	"accFlow",
        	accFlowIdStr,
        	contractHold.AccIdB,
        	"sh0003",
        	"1000",
        	"1",
        	contractNStr,
        	timeStr,
        }
	    contractHold.ContractStatus="2"
	}else{
		return shim.Error("the contract has started")
	}
	//保存最新的状态
	result=saveContractHold(stub,contractHold,contractNStr,chaincodeToCall,channelId)
    if result!=success{
            return shim.Error(result)
    }
    //保存资金账户
    result=saveAcctAsset(stub,acctAssetA,contractHold.AccIdA,chaincodeToCall,channelId)
    if result!=success{
    	return shim.Error(result)
    }
    result=saveAcctAsset(stub,acctAssetB,contractHold.AccIdB,chaincodeToCall,channelId)
    if result!=success{
    	return shim.Error(result)
    }
    //保存持仓账户
    result=saveSSEHoldByAccAndProduct(stub,sseHoldA,contractHold.AccIdA,"sh0003",chaincodeToCall,channelId)
    if result!=success{
    	return shim.Error(result)
    }
    result=saveSSEHoldByAccAndProduct(stub,sseHoldB,contractHold.AccIdB,"sh0003",chaincodeToCall,channelId)
    if result!=success{
    	return shim.Error(result)
    }
    //保存最新的账户流水号
    result=saveAccFlowId(stub,accFlowIdStr,chaincodeToCall,channelId)
    if result!=success{
        return shim.Error(result)
    }
    //保存账户流水
    result=saveAccFlow(stub,accFlowAAsset,accFlowAAsset.AccFlowId,chaincodeToCall,channelId)
    if result!=success{
    	return shim.Error(result)
    }
    result=saveAccFlow(stub,accFlowBAsset,accFlowBAsset.AccFlowId,chaincodeToCall,channelId)
    if result!=success{
    	return shim.Error(result)
    }
    result=saveAccFlow(stub,accFlowAHold,accFlowAHold.AccFlowId,chaincodeToCall,channelId)
    if result!=success{
    	return shim.Error(result)
    }
    result=saveAccFlow(stub,accFlowBHold,accFlowBHold.AccFlowId,chaincodeToCall,channelId)
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
//获得持仓账户
func getSSEHoldByAccAndProduct(stub shim.ChaincodeStubInterface,accIdStr,productCode,chaincodeToCall,channelId string)(SSEHold,string){
    var sseHold SSEHold 
    invokeArgs:=util.ToChaincodeArgs("invoke","getSSEHoldByAccAndProduct",accIdStr,productCode)
    response := stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
    if response.Status!=shim.OK{
        errStr := fmt.Sprintf("A B lockMoney Failed to invoke chaincode. Got error: %s", string(response.Payload))
        fmt.Printf(errStr)
        return sseHold,errStr
    }
    sseHoldBytes:=response.Payload
    err:=json.Unmarshal(sseHoldBytes,&sseHold) 
    if err!=nil{
        return sseHold,"sseHoldBytes Unmarshal failed"
    }
    return sseHold,success
}
//保存持仓账户信息
func saveSSEHoldByAccAndProduct(stub shim.ChaincodeStubInterface,sseHold SSEHold,accIdStr,productCode,chaincodeToCall,channelId string)string{
    sseHoldBytes,err:=json.Marshal(sseHold)
    if err!=nil{
        return "sseHoldBytes marshal failed"
    }
    invokeArgs:=util.ToChaincodeArgs("invoke","saveSSEHoldByAccAndProduct",accIdStr,productCode,string(sseHoldBytes))
    response:= stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
    if response.Status!=shim.OK{
        errStr := fmt.Sprintf("A B lockMoney Failed to invoke chaincode. Got error: %s", string(response.Payload))
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
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
