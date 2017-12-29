package main
import (
        "encoding/json"
        "github.com/hyperledger/fabric/common/util"
        "fmt"
        "math"
        "strconv" 
        "github.com/hyperledger/fabric/core/chaincode/shim"
        pb "github.com/hyperledger/fabric/protos/peer"
)
const (
       success="success"
       contractN="c3003"
      )
// channelId:="mychannel"
// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
type PriceTable struct{
    ObjectType string `json:"docType"`
    Time string `json:"time"`
    ProductCode string `json:"productCode"`
    Price float64 `json:"price"`
}
type ContractHold struct{
    ObjectType string `json:"docType"`
    ContractN string `json:"contractN"`
    ContractCode string `json:"contractCode"`
    ContractStatus string `json:"contractStatus"`//"0"代表等待开始，"1"代表开始合约，"2"代表已经开始,
//"3"代表行权，"4"代表已终结,"5"代表互操作,"6"代表待联通,"7"代表联通,"8"代表解质押,"9"已完成,"10"执行合约
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
    //设计合约
    if args[0]=="transprofittwo"{
        return t.transprofittwo(stub,args)
    }
    return shim.Error("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'")
}
//设计合约
//transprofitone,contractN,actionType,accId,pw,accIdN,channelId,chaincodeToCall,
func (t *SimpleChaincode) transprofittwo(stub shim.ChaincodeStubInterface,args []string)pb.Response{
    if len(args)!=8{
        return shim.Error("args numbers is wrong")
    }
    if args[1]==""||args[2]==""||args[3]==""||args[4]==""||args[5]==""{
        return shim.Error("args can not be nil")
    }
    contractNStr:=args[1]
    actionType:=args[2]
    accIdStr:=args[3]
    pwStr:=args[4]
    accIdNStr:=args[5]
    channelId:=args[6]
    chaincodeToCall:=args[7]
    if contractNStr!=contractN{
        return shim.Error("contractN must be c3003")
    }
    //获得合约持仓
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
    if contractHold.ContractStatus=="9"&&actionType=="10"{
        netAccountStr,result:=getNetAmountByAcctIdAndContractN(stub,"c3002",contractHold.AccIdB,chaincodeToCall,channelId)
        if result!=success{
            return shim.Error(result)
        }
        kk,err:=strconv.ParseFloat(netAccountStr,64)
        if err!=nil{
            return shim.Error("netAccountStr Atoi failed")
        }
        var zz float64
        zzBytes,err:=stub.GetState("zz")
        if err!=nil{
            return shim.Error("zzBytes getstate failed")
        }
        if zzBytes!=nil{
            zz,err=strconv.ParseFloat(string(zzBytes),64)
            if err!=nil{
                return shim.Error("zzBytes Atoi failed")
            }
        }
        gg:=kk-zz
        time,result:=getCurrTime(stub,chaincodeToCall,channelId)
        if result!=success{
            return shim.Error(result)
        }
        price,result:=getStockPrice(stub,time,"sh0002",chaincodeToCall,channelId)
        if result!=success{
            return shim.Error(result)
        }
        transMoney:=(price/40)*gg
        transMoney=math.Abs(transMoney)
        if gg>0{
            //B向A转钱
            result=moveMoney(stub,contractHold.AccIdB,contractHold.AccIdA,strconv.Itoa(int(transMoney)),contractNStr,chaincodeToCall,channelId)
            if result!=success{
                return shim.Error(result)
            }
        }else if gg<0{
            //A向B转钱
            result=moveMoney(stub,contractHold.AccIdA,contractHold.AccIdB,strconv.Itoa(int(transMoney)),contractNStr,chaincodeToCall,channelId)
            if result!=success{
                return shim.Error(result)
            }
        }
        zz+=gg
        err=stub.PutState("zz",[]byte(strconv.FormatFloat(zz,'g',2,64)))
        if err!=nil{
            return shim.Error("zz PutState failed")
        }
    }else if contractHold.ContractStatus=="2"&&actionType=="9"{
        contractHold.ContractStatus="9"
    }else{
        return shim.Error("actionType or contractHold.ContractStatus is wrong")
    }
    //保存最新的状态
    result=saveContractHold(stub,contractHold,contractNStr,chaincodeToCall,channelId)
    if result!=success{
        return shim.Error(result)
    }
    return shim.Success([]byte(success))
}
//获得系统时间
func getCurrTime(stub shim.ChaincodeStubInterface,chaincodeToCall,channelId string)(string,string){
    invokeArgs:=util.ToChaincodeArgs("invoke","getCurrTime")
    response:= stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
    if response.Status!=shim.OK{
        errStr := fmt.Sprintf("getCurrTime Failed to invoke chaincode. Got error: %s", string(response.Payload))
        fmt.Printf(errStr)
        return "",errStr
    }
    return string(response.Payload),success
}
//获得产品价格
func getStockPrice(stub shim.ChaincodeStubInterface,time,productCode,chaincodeToCall,channelId string)(float64,string){
    invokeArgs:=util.ToChaincodeArgs("invoke","getStockPrice",time,productCode)
    response:= stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
    if response.Status!=shim.OK{
        errStr := fmt.Sprintf("getCurrTime Failed to invoke chaincode. Got error: %s", string(response.Payload))
        fmt.Printf(errStr)
        return 0,errStr
    }
    var price PriceTable
    err:=json.Unmarshal(response.Payload,&price)
    if err!=nil{
        return 0,"price Unmarshal failed"
    }
    return price.Price,success
}
//账户之间转钱
func moveMoney(stub shim.ChaincodeStubInterface,accIdFrom,accIdTo,money,contractNStr,chaincodeToCall,channelId string)string{
    invokeArgs:=util.ToChaincodeArgs("invoke","moveMoney",accIdFrom,accIdTo,money,contractN)
    response:= stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
    if response.Status!=shim.OK{
        errStr := fmt.Sprintf("getCurrTime Failed to invoke chaincode. Got error: %s", string(response.Payload))
        fmt.Printf(errStr)
        return errStr
    }
    return success
}
//获得净额
func getNetAmountByAcctIdAndContractN(stub shim.ChaincodeStubInterface,contractNStr,acctId,chaincodeToCall,channelId string)(string,string){
    invokeArgs:=util.ToChaincodeArgs("invoke","getNetAmountByAcctIdAndContractN",contractNStr,acctId)
    response:= stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
    if response.Status!=shim.OK{
        errStr := fmt.Sprintf("getCurrTime Failed to invoke chaincode. Got error: %s", string(response.Payload))
        fmt.Printf(errStr)
        return "",errStr
    }
    return string(response.Payload),success
}
//保存合约最新的状态
func saveContractHold(stub shim.ChaincodeStubInterface,contractHold ContractHold,contractNStr,chaincodeToCall,channelId string)string{
    contractHoldBytes,err:=json.Marshal(contractHold)
    if err!=nil{
        return "contractHold marshal failed"
    }
    fmt.Println("saveContractHold"+string(contractHoldBytes))
    invokeArgs:=util.ToChaincodeArgs("invoke","saveContractHold",contractNStr,string(contractHoldBytes))
    response := stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
    if response.Status!=shim.OK{
        errStr := fmt.Sprintf("7 Failed to invoke chaincode. Got error: %s", string(response.Payload))
        fmt.Printf(errStr)
        return errStr
    }
    return success
}
//获得合约持仓
func getContractHold(stub shim.ChaincodeStubInterface,contractNStr,chaincodeToCall,channelId string)(ContractHold,string){
    var contractHold ContractHold
    invokeArgs:=util.ToChaincodeArgs("invoke","getContractHold",contractNStr)
    response := stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
    if response.Status!=shim.OK{
        errStr := fmt.Sprintf("getContractHold Failed to invoke chaincode. Got error: %s", string(response.Payload))
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
