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
	if args[0] == "query" {
		// queries an entity state
		return t.query(stub, args)
	}
	//设计合约
	if args[0]=="europtiondemo"{
		return t.europtiondemo(stub,args)
	}
	return shim.Error("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'")
}
//设计合约
//europtiondemo,contractN,type,acctId,pw,channelId,chaincodeToCall,
func (t *SimpleChaincode) europtiondemo(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=7{
		return shim.Error("args numbers is wrong")
	}
	if args[1]==""||args[2]==""{
		return shim.Error("args can not be nil")
	}
	contractNStr:=args[1]
	transType:=args[2]
	acctIdStr:=args[3]
	pw:=args[4]
	channelId:=args[5]
	chaincodeToCall:=args[6]
	result:=executeConrtact(stub,contractNStr,transType,acctIdStr,pw,channelId,chaincodeToCall)
	if result!=success{
		return shim.Error(result)
	}
	return shim.Success([]byte(success))
}
func executeConrtact(stub shim.ChaincodeStubInterface,contractNStr,transType,acctIdStr,pw,channelId,chaincodeToCall string)string{
	f:="getContractHold"
	invokeArgs:=util.ToChaincodeArgs("invoke",f,contractNStr)
	response := stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
	if response.Status!=shim.OK{
		errStr := fmt.Sprintf("11Failed to invoke chaincode. Got error: %s", string(response.Payload))
		fmt.Printf(errStr)
		return errStr
	}
	contractHoldBytes:=response.Payload
    var contractHold ContractHold
    err:=json.Unmarshal(contractHoldBytes,&contractHold)
    if err!=nil{
    	return "contractHoldBytes Unmarshal failed"
    }
    accIdAStr:=contractHold.AccIdA
	accIdBStr:=contractHold.AccIdB
	if acctIdStr!=contractHold.AccIdA&&acctIdStr!=contractHold.AccIdB{
		return "accId is wrong"
	}
	//获得最新的账户流水编号
    accFlowIdStr,result:=getAccFlowId(stub,chaincodeToCall,channelId)
    if result!=success{
        return result
    }
    accFlowId,err:=strconv.Atoi(accFlowIdStr)
    if err!=nil{
        return "accFlowId Atoi failed"
    }
	if contractHold.ContractStatus=="0"&&transType=="1"{
		currTimeStr:=getCurrTime(stub,chaincodeToCall,channelId)
        currTime,err:=strconv.Atoi(currTimeStr)
        if err!=nil{
        	return "currTimeStr atoi failed"
        }
		if currTime<20170101&&currTime>=20170110{
			return "time is wrong"
		}
		f:="lockStock"
		invokeArgs:=util.ToChaincodeArgs("invoke",f,accIdAStr,contractNStr,"sh0001","1000",accFlowIdStr)
	    response = stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
	    if response.Status!=shim.OK{
	    	errStr := fmt.Sprintf("8Failed to invoke chaincode. Got error: %s", string(response.Payload))
			fmt.Printf(errStr)
			return errStr+accFlowIdStr
	    }
	    accFlowId+=1
		invokeArgs=util.ToChaincodeArgs("invoke","moveMoney",accIdBStr,accIdAStr,"500",contractNStr,strconv.Itoa(accFlowId))
		response = stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
		if response.Status!=shim.OK{
	    	errStr := fmt.Sprintf("9Failed to invoke chaincode. Got error: %s", string(response.Payload))
			fmt.Printf(errStr)
			invokeArgs:=util.ToChaincodeArgs("invoke","unlockStock",accIdAStr,contractNStr,"sh0001","1000")
			response=stub.InvokeChaincode(chaincodeToCall,invokeArgs,channelId)
			if response.Status!=shim.OK{
				errStr := fmt.Sprintf("12Failed to invoke chaincode. Got error: %s", string(response.Payload))
				fmt.Printf(errStr)
				return errStr
			}
			return "B has not enough Money"
	    }
		contractHold.ContractStatus="2"
	}else if contractHold.ContractStatus=="2"&&transType=="3"{
		invokeArgs:=util.ToChaincodeArgs("invoke","getCurrTime","")
	    response = stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
	    if response.Status!=shim.OK{
	    	errStr := fmt.Sprintf("1Failed to invoke chaincode. Got error: %s", string(response.Payload))
			fmt.Printf(errStr)
			return errStr
	    }
	    exerciseDate,err:=strconv.Atoi(string(response.Payload))
	    if err!=nil{
	    	return "exerciseDate atoi failed"
	    }
	    if exerciseDate<20170720||exerciseDate>20170805{
	    	return "the exercise date must between 20170720 and 20170805"
	    }
	    invokeArgs=util.ToChaincodeArgs("invoke","unlockStock",accIdAStr,contractNStr,"sh0001","1000",accFlowIdStr)
	    response = stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
	    if response.Status!=shim.OK{
	    	errStr := fmt.Sprintf("2Failed to invoke chaincode. Got error: %s", string(response.Payload))
			fmt.Printf(errStr)
			return errStr
	    }
	    AHoldBytes:=response.Payload
	    accFlowId+=1
	    invokeArgs=util.ToChaincodeArgs("invoke","moveMoney",accIdBStr,accIdAStr,"6000",contractNStr,strconv.Itoa(accFlowId))
	    response = stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
	    if response.Status!=shim.OK{
			invokeArgs:=util.ToChaincodeArgs("invoke","lockStock",accIdAStr,contractNStr,"sh0001","1000")
		    response = stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
		    if response.Status!=shim.OK{
		    	errStr := fmt.Sprintf("3Failed to invoke chaincode. Got error: %s", string(response.Payload))
				fmt.Printf(errStr)
				return errStr
		    }
			return "B has not enough money"
	    }
	    accFlowId+=3
	    invokeArgs=util.ToChaincodeArgs("invoke","moveStock",accIdAStr,accIdBStr,string(AHoldBytes),"sh0001","1000",contractNStr,strconv.Itoa(accFlowId))
	    response = stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
	    if response.Status!=shim.OK{
	    	errStr := fmt.Sprintf("4Failed to invoke chaincode. Got error: %s", string(response.Payload))
			fmt.Printf(errStr)
			return errStr
	    }
	    contractHold.ContractStatus="4"
	}else if contractHold.ContractStatus=="2"&&transType=="5"{
		invokeArgs:=util.ToChaincodeArgs("invoke","getCurrTime","")
	    response = stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
	    if response.Status!=shim.OK{
	    	errStr := fmt.Sprintf("5Failed to invoke chaincode. Got error: %s", string(response.Payload))
			fmt.Printf(errStr)
			return errStr
	    }
	    exerciseDate,err:=strconv.Atoi(string(response.Payload))
	    if err!=nil{
	    	return "exerciseDate atoi failed"
	    }
	    if exerciseDate>20170805{
	    	invokeArgs:=util.ToChaincodeArgs("invoke","unlockStock",accIdAStr,contractNStr,"sh0001","1000",accFlowIdStr)
		    response = stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
		    if response.Status!=shim.OK{
		    	errStr := fmt.Sprintf("6Failed to invoke chaincode. Got error: %s", string(response.Payload))
				fmt.Printf(errStr)
				return errStr
		    }
		    contractHold.ContractStatus="4"
	    }
	}else{
		return "status is unvalid,please check the status  of the contract"
	}
	//保存最新的状态
	contractHoldBytes,err=json.Marshal(contractHold)
	if err!=nil{
		return "contractHold marshal failed"
	}
	invokeArgs=util.ToChaincodeArgs("invoke","saveContractHold",contractNStr,string(contractHoldBytes))
    response = stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
    if response.Status!=shim.OK{
    	errStr := fmt.Sprintf("7Failed to invoke chaincode. Got error: %s", string(response.Payload))
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
// Query callback representing the query of a chaincode
//orgName+username+money 查询资金账户信息
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface,args []string) pb.Response {

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
