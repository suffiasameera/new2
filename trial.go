package main

import "fmt"
import "errors"
import "encoding/json"
import "github.com/hyperledger/fabric/core/chaincode/shim"
//import "strconv"
//import "encoding/base64"

type SampleChaincode struct {
}

//Incident struct
type Incident struct {
IncidentID string `json:"iid"`
IName string `json:"iname"`
Desc string `json:"desc"`
Orig string `json:"orig"`
Status string `json:"status"`
}

type customEvent struct {
	Type       string `json:"type"`
	Decription string `json:"description"`
}

func Create(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("Entering Create function")

    if len(args) != 5 {
        fmt.Println("Invalid number of args")
        return nil, errors.New("Expected at least two arguments for loan application creation")
    }

    var Id = args[0]
    var In = args[1]
    var De = args[2]
    var Or = args[3]
    var St = args[4]

    var incidentInfo Incident
    incidentInfo = Incident{Id, In, De, Or, St}
     piBytes, err := json.Marshal(&incidentInfo)
     if err != nil {
       fmt.Println("Error in marshaling ",err)
       return nil, err
     }


    err = stub.PutState(Id, piBytes)
    if err != nil {
        fmt.Println("Could not save changes", err)
        return nil, err
    }

    var customEvent = "{eventType: 'Creation', description:" + Id + "' Successfully created'}"
    err = stub.SetEvent("evtSender", []byte(customEvent))
    if err != nil {
      return nil, err
    }

    fmt.Println("Successfully saved changes")
    return nil, nil
}

func Get(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("Entering Get function")

    /*if len(args) < 1 {
        fmt.Println("Invalid number of arguments")
        return nil, errors.New("Missing ID")
    }*/

    var ccMapDone map[string]string
		ccMapDone = make(map[string]string)
		var ccMapO1 map[string]string
		ccMapO1 = make(map[string]string)
		var ccMapL1 map[string]string
		ccMapL1 = make(map[string]string)

    for  i := 0; i < len(args); i++ {
    var Id = args[i]
    piBytes, err := stub.GetState(Id)
    if err != nil {
        fmt.Println("Could not fetch data with id "+Id+" from ledger", err)
        return nil, errors.New("Could not fetch data with id " + Id + err.Error())
    }

    //piBytes2 := piBytes
    var incidentInfo Incident
    err = json.Unmarshal(piBytes, &incidentInfo)
    if err != nil {
      fmt.Println("Error in unmarshaling",err)
      return nil, err
    }
    fmt.Println(incidentInfo.IName)
    fmt.Println(incidentInfo.Desc)
    fmt.Println(incidentInfo.Orig)
    fmt.Println(incidentInfo.Status)

    if incidentInfo.Status == "done" {
			ccMapDone[Id] =  incidentInfo.IName
		}

		if incidentInfo.Orig == "O1" && incidentInfo.Status == "done" {
			ccMapO1[Id] =  incidentInfo.IName
		}

		if incidentInfo.Orig == "L1" && incidentInfo.IName == "abc" && incidentInfo.Status != "done" {
			ccMapO1[Id] =  incidentInfo.IName
		}
		}

	 for c := range ccMapDone {
	 fmt.Println("Name of",c,"is",ccMapDone[c])
   }
	 for d := range ccMapO1 {
	 fmt.Println("Name of",d,"is",ccMapO1[d])
   }
	 for e := range ccMapL1 {
	 fmt.Println("Name of",e,"is",ccMapL1[e])
   }

    var incidentInfo Incident
    piBytes2, err := json.Marshal(&incidentInfo)
    if err != nil {
      fmt.Println("Error in marshaling",err)
      return nil, err
    }

    s := incidentInfo.IncidentID + " " + incidentInfo.IName + " " + incidentInfo.Desc + " " + incidentInfo.Orig + " " + incidentInfo.Status
    piBytes2 = []byte(s)
    return piBytes2, nil
}

func Delete(stub shim.ChaincodeStubInterface, args []string) (error) {
		fmt.Println("Entering Delete")
	  var Id = args[0]

			err := stub.DelState(Id)
			if err != nil {
				fmt.Println("Could not delete from ledger", err)
				return  err
			}

		fmt.Println("Successfully deleted ")
		return nil


	}

func Update(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Update")

	if len(args) < 3 {
		fmt.Println("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for update")
	}

	var Id = args[0]
	var changes = args[1]
  var op = args[2]

	laBytes, err := stub.GetState(Id)
	if err != nil {
		fmt.Println("Could not fetch data from ledger", err)
		return nil, err
	}

	var incidentInfo Incident
	err = json.Unmarshal(laBytes, &incidentInfo)
  if( op == "1" ) {
	//incidentInfo.IncidentId = changes
  } else if( op == "2" ) {
    incidentInfo.IName = changes
  } else if ( op == "3" ) {
    incidentInfo.Desc = changes
  } else if ( op == "4" ) {
    incidentInfo.Orig = changes
  } else {
    incidentInfo.Status = changes
  }

	laBytes, err = json.Marshal(&incidentInfo)
	if err != nil {
		fmt.Println("Could not marshal ", err)
		return nil, err
	}

	err = stub.PutState(Id, laBytes)
	if err != nil {
		fmt.Println("Could not save update", err)
		return nil, err
	}
/*
  var customEvent = "{eventType: 'Update', description:" + Id + "' Successfully updated status'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return nil, err
	}
*/
	fmt.Println("Successfully updated changes")
	return nil, nil

}


func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

   /* if len(args) != 4 {
     return nil, errors.New("Incorrect number of arguments. Expecting 4")
   }
   */

     var incidentInfo Incident
      incidentInfo = Incident{"1", "One", "d1", "o1", "s1"}
      bytes, err := json.Marshal (&incidentInfo)
      if err != nil {
             fmt.Println("Could not marshal incident info object", err)
             return nil, err
      }

   /*bytes, err := json.Marshal (&args[0])
   if err != nil {
          fmt.Println("Could not marshal incident info object", err)
          return nil, err
   }
  */

   err = stub.PutState("1", bytes)
   if err != nil {
     fmt.Println("Could not save ", err)
     return nil, err
   }

    return nil, nil
}

func (t *SampleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
fmt.Println("Inside Query function")
  if function == "get" {
    return Get(stub, args)
  }
  fmt.Println("Query could not find: " + function)
    return nil, nil
}

func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
/*  if function == "init" {
		return Init(stub, "init", args)
  }
  */
  if function == "create" {
    return Create(stub, args)
  }
  if function=="update" {
    return Update(stub,args)
  }
	if function=="delete" {
		Delete(stub,args)
	}
  fmt.Println("Invoke did not find func: " + function)
    return nil, nil
}

//Main function

func main() {
    err := shim.Start(new(SampleChaincode))
    if err != nil {
        fmt.Println("Could not start SampleChaincode")
    } else {
        fmt.Println("SampleChaincode successfully started")
    }

}
