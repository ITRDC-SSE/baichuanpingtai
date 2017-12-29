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
const success="success"
	// channelId:="mychannel"
// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
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
	                    //"3"代表行权，"4"代表已终结
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
	if args[0]=="complexoptiondemo"{
		return t.complexoptiondemo(stub,args)
	}
	return shim.Error("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'")
}
//设计合约
//complexOptionDemo,contractN,type,accId,pw,accIdN,channelId,chaincodeToCall,
func (t *SimpleChaincode) complexoptiondemo(stub shim.ChaincodeStubInterface,args []string)pb.Response{
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
    //获得账户A的资金账户
	acctAssetA,result:=getAcctAsset(stub,contractHold.AccIdA,chaincodeToCall,channelId)
	if result!=success{
		return shim.Error(result)
	}
	acctAssetB,result:=getAcctAsset(stub,contractHold.AccIdB,chaincodeToCall,channelId)
	if result!=success{
		return shim.Error(result)
	}
	acctAssetC,result:=getAcctAsset(stub,contractHold.AccIdC,chaincodeToCall,channelId)
	if result!=success{
		return shim.Error(result)
	}
	acctAssetD,result:=getAcctAsset(stub,contractHold.AccIdD,chaincodeToCall,channelId)
	if result!=success{
		return shim.Error(result)
	}
	acctAssetE,result:=getAcctAsset(stub,contractHold.AccIdE,chaincodeToCall,channelId)
	if result!=success{
		return shim.Error(result)
	}
	//账户流水
	var accFlowAAsset,accFlowBAsset,accFlowCAsset,accFlowDAsset,accFlowEAsset AccFlow
	//获得最新的账户流水编号
	invokeArgs:=util.ToChaincodeArgs("invoke","getAccFlowId")
	response := stub.InvokeChaincode(chaincodeToCall,invokeArgs,channelId)
	if response.Status!=shim.OK{
		errStr := fmt.Sprintf("getContractHold Failed to invoke chaincode. Got error: %s", string(response.Payload))
		fmt.Printf(errStr)
		return shim.Error(errStr)
	}
	accFlowIdStr:=string(response.Payload)
	accFlowId,err:=strconv.Atoi(accFlowIdStr)
	if err!=nil{
		return shim.Error("accFlowId Atoi failed")
	}
	timeStr:=getCurrTime(stub,chaincodeToCall,channelId)//yyyyMMdd
	if contractHold.ContractStatus=="0"&&actionType=="1"{
		currTime,err:=strconv.Atoi(timeStr)
		if err!=nil{
			return shim.Error("timeStr Atoi failed")
		}
		if currTime<20170101&&currTime>20170115{
			return shim.Error("time is wrong")
		}
		if accIdStr!=contractHold.AccIdA{
			return shim.Error("accId is wrong")
		}
		if acctAssetA.AvaMoney<4000{
			return shim.Error("acctA do not have enough money")
		}
		//A给B转钱
		acctAssetA.AvaMoney-=1000
		acctAssetB.AvaMoney+=1000
		//A给C转钱
		acctAssetA.AvaMoney-=1000
		acctAssetC.AvaMoney+=1000
		//A给D转钱
		acctAssetA.AvaMoney-=1000
		acctAssetD.AvaMoney+=1000
		//A给E转钱
		acctAssetA.AvaMoney-=1000
		acctAssetE.AvaMoney+=1000
	    
	    //账户A的流水
        accFlowAAsset=AccFlow{
            "accFlow",
            accFlowIdStr,
            contractHold.AccIdA,
            "money",
            "4000",
            "1",
            contractNStr,
            timeStr,
        }
        //账户B的流水
        accFlowId=accFlowId+1
        accFlowIdStr=strconv.Itoa(accFlowId)
        accFlowBAsset=AccFlow{
            "accFlow",
            accFlowIdStr,
            contractHold.AccIdB,
            "money",
            "1000",
            "0",
            contractNStr,
            timeStr,
        }
        accFlowId=accFlowId+1
        accFlowIdStr=strconv.Itoa(accFlowId)
        accFlowCAsset=AccFlow{
            "accFlow",
            accFlowIdStr,
            contractHold.AccIdC,
            "money",
            "1000",
            "0",
            contractNStr,
            timeStr,
        }
        accFlowId=accFlowId+1
        accFlowIdStr=strconv.Itoa(accFlowId)
        accFlowDAsset=AccFlow{
            "accFlow",
            accFlowIdStr,
            contractHold.AccIdD,
            "money",
            "1000",
            "0",
            contractNStr,
            timeStr,
        }
        accFlowId=accFlowId+1
        accFlowIdStr=strconv.Itoa(accFlowId)
        accFlowEAsset=AccFlow{
            "accFlow",
            accFlowIdStr,
            contractHold.AccIdE,
            "money",
            "1000",
            "0",
            contractNStr,
            timeStr,
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
	    result=saveAccFlow(stub,accFlowCAsset,accFlowCAsset.AccFlowId,chaincodeToCall,channelId)
	    if result!=success{
	    	return shim.Error(result)
	    }
	    result=saveAccFlow(stub,accFlowDAsset,accFlowDAsset.AccFlowId,chaincodeToCall,channelId)
	    if result!=success{
	    	return shim.Error(result)
	    }
	    result=saveAccFlow(stub,accFlowEAsset,accFlowEAsset.AccFlowId,chaincodeToCall,channelId)
	    if result!=success{
	    	return shim.Error(result)
	    }
	    contractHold.ContractStatus="2"
	}else if contractHold.ContractStatus=="2"&&actionType=="3"{
		if accIdStr!=contractHold.AccIdA{
			return shim.Error("accId is wrong;"+contractHold.AccIdA+";"+accIdStr)
		}
		invokeArgs:=util.ToChaincodeArgs("invoke","getCurrTime")
		response := stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
		if response.Status!=shim.OK{
	    	errStr := fmt.Sprintf("getCurrTime  Failed to invoke chaincode. Got error: %s", string(response.Payload))
			fmt.Printf(errStr)
			return shim.Error(errStr)
	    }
	    timeBytes:=response.Payload//yyyyMMdd
	    time,err:=strconv.Atoi(string(timeBytes))
        if err!=nil{
        	return shim.Error("time Atoi failed")
        }
        if time<20170802&&time>20170815{
        	return shim.Error("time is wrong")
        }
        //获得均价
        invokeArgs=util.ToChaincodeArgs("invoke","getAveragePrice","20170710","20170725","sh0002")
		response = stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
		if response.Status!=shim.OK{
	    	errStr := fmt.Sprintf("getAveragePrice Failed to invoke chaincode. Got error: %s", string(response.Payload))
			fmt.Printf(errStr)
			return shim.Error(errStr+";"+chaincodeToCall+";"+channelId)
	    }
	    //均值 float64
	    averPrice,err:=strconv.ParseFloat(string(response.Payload),64) 
	    if err!=nil{
	    	return shim.Error("averPrice ParseFloat failed")
	    }
	    invokeArgs=util.ToChaincodeArgs("invoke","getStandDeviationPrice","20170410","20170725","sh0002")
		response = stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
		if response.Status!=shim.OK{
	    	errStr := fmt.Sprintf("getStandDeviationPrice Failed to invoke chaincode. Got error: %s", string(response.Payload))
			fmt.Printf(errStr)
			return shim.Error(errStr)
	    }
	    //标准差 float64
	    stardardDeviation,err:=strconv.ParseFloat(string(response.Payload),64) 
	     if err!=nil{
	    	return shim.Error("stardardDeviation ParseFloat failed")
	    }
	    amtTotal:=(math.Max(100-averPrice,0)+stardardDeviation*0.001)*100
	    if amtTotal>0{
	    	var acctStr string
	    	amtTotal1:=int(math.Min(amtTotal,2000))
	    	amtTotal1Str:=strconv.Itoa(amtTotal1)
	    	if acctAssetB.AvaMoney<amtTotal1{
	    		if acctAssetC.AvaMoney<amtTotal1{
	    			if acctAssetD.AvaMoney<amtTotal1{
	    				if acctAssetE.AvaMoney<amtTotal1{
	    					return shim.Error("there is not enough money")
	    				}else{
	    					acctAssetE.AvaMoney-=amtTotal1
	    			        acctAssetA.AvaMoney+=amtTotal1
	    			        acctStr=acctAssetE.AcctId
	    				}
	    			}else{
	    				acctAssetD.AvaMoney-=amtTotal1
	    			    acctAssetA.AvaMoney+=amtTotal1
	    			    acctStr=acctAssetD.AcctId
	    			}
	    		}else{
	    			acctAssetC.AvaMoney-=amtTotal1
	    			acctAssetA.AvaMoney+=amtTotal1
	    			acctStr=acctAssetC.AcctId
	    		}
	    	}else{
	    		acctAssetB.AvaMoney-=amtTotal1
	    		acctAssetA.AvaMoney+=amtTotal1
	    		acctStr=acctAssetB.AcctId
	    	}
	    	accFlowAAsset=AccFlow{
	            "accFlow",
	            accFlowIdStr,
	            contractHold.AccIdA,
	            "money",
	            amtTotal1Str,
	            "0",
	            contractNStr,
	            timeStr,
	        }
	        //账户B的流水
	        accFlowId=accFlowId+1
	        accFlowIdStr=strconv.Itoa(accFlowId)
	        accFlowBAsset=AccFlow{
	            "accFlow",
	            accFlowIdStr,
	            acctStr,
	            "money",
	            amtTotal1Str,
	            "0",
	            contractNStr,
	            timeStr,
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
		    //更新流水号ss
		    accFlowId=accFlowId+1
		    accFlowIdStr=strconv.Itoa(accFlowId)
	    }
	    if amtTotal>2000{
	    	var acctStr string
	    	amtTotal2:=int(math.Min(amtTotal-2000,2000))
	    	amtTotal2Str:=strconv.Itoa(amtTotal2)
	    	if acctAssetC.AvaMoney<amtTotal2{
	    		if acctAssetD.AvaMoney<amtTotal2{
	    			if acctAssetE.AvaMoney<amtTotal2{
	    				return shim.Error("there is not enough money")
	    			}else{
	    				acctAssetE.AvaMoney-=amtTotal2
	    		        acctAssetA.AvaMoney+=amtTotal2
	    		        acctStr=acctAssetE.AcctId
	    			}
	    		}else{
	    			acctAssetD.AvaMoney-=amtTotal2
	    		    acctAssetA.AvaMoney+=amtTotal2
	    		    acctStr=acctAssetD.AcctId
	    		}
	    	}else{
	    		acctAssetC.AvaMoney-=amtTotal2
	    		acctAssetA.AvaMoney+=amtTotal2
	    		acctStr=acctAssetC.AcctId
	    	}
	    	accFlowAAsset=AccFlow{
	            "accFlow",
	            accFlowIdStr,
	            contractHold.AccIdA,
	            "money",
	            amtTotal2Str,
	            "0",
	            contractNStr,
	            timeStr,
	        }
	        //账户B的流水
	        accFlowId=accFlowId+1
	        accFlowIdStr=strconv.Itoa(accFlowId)
	        accFlowBAsset=AccFlow{
	            "accFlow",
	            accFlowIdStr,
	            acctStr,
	            "money",
	            amtTotal2Str,
	            "0",
	            contractNStr,
	            timeStr,
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
		    //更新流水号
		    accFlowId=accFlowId+1
		    accFlowIdStr=strconv.Itoa(accFlowId)
	    }
	    if amtTotal>4000{
	    	var acctStr string
	    	amtTotal3:=int(math.Min(amtTotal-4000,2000))
	    	amtTotal3Str:=strconv.Itoa(amtTotal3)
	    	if acctAssetD.AvaMoney<amtTotal3{
	    		if acctAssetE.AvaMoney<amtTotal3{
	    			return shim.Error("there is not enough money")
	    		}else{
	    			acctAssetE.AvaMoney-=amtTotal3
	    		    acctAssetA.AvaMoney+=amtTotal3
	    		    acctStr=acctAssetE.AcctId
	    		}
	    	}else{
	    		acctAssetD.AvaMoney-=amtTotal3
	    		acctAssetA.AvaMoney+=amtTotal3
	    		acctStr=acctAssetD.AcctId
	    	}
	    	accFlowAAsset=AccFlow{
	            "accFlow",
	            accFlowIdStr,
	            contractHold.AccIdA,
	            "money",
	            amtTotal3Str,
	            "0",
	            contractNStr,
	            timeStr,
	        }
	        //账户B的流水
	        accFlowId=accFlowId+1
	        accFlowIdStr=strconv.Itoa(accFlowId)
	        accFlowBAsset=AccFlow{
	            "accFlow",
	            accFlowIdStr,
	            acctStr,
	            "money",
	            amtTotal3Str,
	            "1",
	            contractNStr,
	            timeStr,
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
	    	//更新流水号
		    accFlowId=accFlowId+1
		    accFlowIdStr=strconv.Itoa(accFlowId)
	    }
	    if amtTotal>6000{
	    	var acctStr string
	    	amtTotal4:=int(amtTotal-6000)
	    	amtTotal4Str:=strconv.Itoa(amtTotal4)
	    	if acctAssetE.AvaMoney<amtTotal4{
	    		return shim.Error("there is not enough money")
	    	}else{
	    		acctAssetE.AvaMoney-=amtTotal4
                acctAssetA.AvaMoney+=amtTotal4
                acctStr=acctAssetE.AcctId
	    	}
	    	accFlowAAsset=AccFlow{
	            "accFlow",
	            accFlowIdStr,
	            contractHold.AccIdA,
	            "money",
	            amtTotal4Str,
	            "0",
	            contractNStr,
	            timeStr,
	        }
	        //账户B的流水
	        accFlowId=accFlowId+1
	        accFlowIdStr=strconv.Itoa(accFlowId)
	        accFlowBAsset=AccFlow{
	            "accFlow",
	            accFlowIdStr,
	            acctStr,
	            "money",
	            amtTotal4Str,
	            "1",
	            contractNStr,
	            timeStr,
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
	    }
	    contractHold.ContractStatus="4"
	}else{
		return shim.Error("unvalid Status")
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
    result=saveAcctAsset(stub,acctAssetC,contractHold.AccIdC,chaincodeToCall,channelId)
    if result!=success{
    	return shim.Error(result)
    }
    result=saveAcctAsset(stub,acctAssetD,contractHold.AccIdD,chaincodeToCall,channelId)
    if result!=success{
    	return shim.Error(result)
    }
    result=saveAcctAsset(stub,acctAssetE,contractHold.AccIdE,chaincodeToCall,channelId)
    if result!=success{
    	return shim.Error(result)
    }
	//保存最新的状态
	contractHoldBytes,err:=json.Marshal(contractHold)
	if err!=nil{
		return shim.Error("contractHold marshal failed")
	}
	invokeArgs=util.ToChaincodeArgs("invoke","saveContractHold",contractNStr,string(contractHoldBytes))
    response = stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
    if response.Status!=shim.OK{
    	errStr := fmt.Sprintf("7 Failed to invoke chaincode. Got error: %s", string(response.Payload))
		fmt.Printf(errStr)
		return shim.Error(errStr)
    }
    //保存最新的流水号
    invokeArgs=util.ToChaincodeArgs("invoke","saveAccFlowId",accFlowIdStr)
	response=stub.InvokeChaincode(chaincodeToCall, invokeArgs, channelId)
	if response.Status!=shim.OK{
    	errStr := fmt.Sprintf("B saveAcctAsset Failed to invoke chaincode. Got error: %s", string(response.Payload))
		fmt.Printf(errStr)
		return shim.Error(errStr)
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
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
