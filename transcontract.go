package main
import (
        
        "encoding/json"
         "github.com/hyperledger/fabric/common/util"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)
const (
	success="success"
	)
	// channelId:="mychannel"
// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
type ContractHold struct{
	ObjectType string `json:"docType"`
	ContractN string `json:"contractN"`
	ContractCode string `json:"contractCode"`
	ContractStatus string `json:"contractStatus"`//"0"代表等待开始，"1"代表开始合约，"2"代表已经开始,
	                    //"3"代表行权，"4"代表已终结
	ContractFunctionName string `json:"contractFunctionName"`
    AccIdA string `json:"accIdA"`//合约义务方
    AccIdB string `json:"accIdB"`//合约权利方
    AccIdC string `json:"accIdC"`
    AccIdD string `json:"accIdD"`
    AccIdE string `json:"accIdE"`
    CorContractNA string `json:"corContractNA"`//关联合约
    CorContractNB string `json:"corContractNB"`
    CorContractNC string `json:"corContractNC"`
    TransType string `json:"transType"`
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
	//转让合约
	if args[0]=="moveContract"{
		return t.moveContract(stub,args)
	}
	return shim.Error("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'")
}
//转让合约
//moveContract,accIdFrom,accIdTo,contractN,amt,accIdN,channelId,chaincodeToCall
func (t *SimpleChaincode) moveContract(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	if len(args)!=8{
		return shim.Error("args length is wrong")
	}
	if args[1]==""||args[2]==""||args[3]==""||args[4]==""{
		return shim.Error("args can not be nil")
	}
	accIdFrom:=args[1]
	accIdTo:=args[2]
	contractNStr:=args[3]
	amtStr:=args[4]
	accIdN:=args[5]
	channelId:=args[6]
	chaincodeToCall:=args[7]
	invokeArgs:=util.ToChaincodeArgs("invoke","getContractHold",contractNStr)
	response := stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
	if response.Status!=shim.OK{
		errStr := fmt.Sprintf("Failed to invoke chaincode. Got error: %s", string(response.Payload))
		fmt.Printf(errStr)
		return shim.Error(errStr)
	}
	contractHoldbytes:=response.Payload
	var contractHold ContractHold
	err:=json.Unmarshal(contractHoldbytes,&contractHold)
	if err!=nil{
		return shim.Error("contractHoldbytes Unmarshal failed")
	}
	if accIdN=="1"{
		if contractHold.AccIdA!=accIdFrom{
			return shim.Error("the accIdFrom is wrong")
		}
		contractHold.AccIdA=accIdTo
	}else if accIdN=="2"{
		if contractHold.AccIdB!=accIdFrom{
			return shim.Error("the accIdFrom is wrong")
		}
		contractHold.AccIdB=accIdTo
	}else if accIdN=="3"{
		if contractHold.AccIdC!=accIdFrom{
			return shim.Error("the accIdFrom is wrong")
		}
		contractHold.AccIdC=accIdTo
	}else if accIdN=="4"{
		if contractHold.AccIdD!=accIdFrom{
			return shim.Error("the accIdFrom is wrong")
		}
		contractHold.AccIdD=accIdTo
	}else if accIdN=="5"{
		if contractHold.AccIdE!=accIdFrom{
			return shim.Error("the accIdFrom is wrong")
		}
		contractHold.AccIdE=accIdTo
	}else{
		return shim.Error("accIdN is wrong")
	}
	if contractHold.TransType!="1"{
		return shim.Error("the contract is not available")
	}
	//获得最新的账户流水编号
    accFlowIdStr,result:=getAccFlowId(stub,chaincodeToCall,channelId)
    if result!=success{
        return shim.Error(result)
    }
	//转让钱
	invokeArgs=util.ToChaincodeArgs("invoke","moveMoney",accIdTo,accIdFrom,amtStr,contractNStr,accFlowIdStr)
	response = stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
	if response.Status!=shim.OK{
		errStr := fmt.Sprintf("Failed to invoke chaincode. Got error: %s", string(response.Payload))
		fmt.Printf(errStr)
		return shim.Error(errStr)
	}
	contractHoldBytes,err:=json.Marshal(contractHold)
	if err!=nil{
		return shim.Error("contractHoldBytes marshal failed")
	}
	//保存最新的合约状态
	invokeArgs=util.ToChaincodeArgs("invoke","saveContractHold",contractNStr,string(contractHoldBytes))
	response = stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
	if response.Status!=shim.OK{
		errStr := fmt.Sprintf("Failed to invoke chaincode. Got error: %s", string(response.Payload))
		fmt.Printf(errStr)
		return shim.Error(errStr)
	}
	return shim.Success([]byte(success))
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
