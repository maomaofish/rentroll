package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rentroll/rlib"
)

// returns bud from cached list (rlib.RRdb.BUDlist)
func bidToBud(businessID int64) (string, error) {
	for bud, bid := range rlib.RRdb.BUDlist {
		if businessID == bid {
			return bud, nil
		}
	}
	return "", fmt.Errorf("Could not find business for bid: %d", businessID)
}

// GetAssessmentList returns all assessments for the supplied business
func GetAssessmentList(bid int64) (map[string][]IDTextMap, error) {

	// initialize list with 0-id value
	list := []IDTextMap{{ID: 0, Text: " -- Select Assessment Rule -- "}}

	// json response data
	appData := make(map[string][]IDTextMap)

	bud, err := bidToBud(bid)
	if err != nil {
		return appData, err
	}

	// get records and append in IDTextMap list
	m := rlib.GetARsByType(bid, rlib.ARASSESSMENT)
	for i := 0; i < len(m); i++ {
		list = append(list, IDTextMap{ID: m[i].ARID, Text: m[i].Name})
	}
	appData[bud] = list

	return appData, nil
}

// GetReceiptList returns all assessments for the supplied business
func GetReceiptList(bid int64) (map[string][]IDTextMap, error) {

	// initialize list with 0-id value
	list := []IDTextMap{{ID: 0, Text: " -- Select Receipt Rule -- "}}

	// json response data
	appData := make(map[string][]IDTextMap)

	bud, err := bidToBud(bid)
	if err != nil {
		return appData, err
	}

	// get records and append in IDTextMap list
	m := rlib.GetARsByType(bid, rlib.ARRECEIPT)
	for i := 0; i < len(m); i++ {
		list = append(list, IDTextMap{ID: m[i].ARID, Text: m[i].Name})
	}
	appData[bud] = list

	return appData, nil
}

// GetDepositoryList returns all assessments for the supplied business
func GetDepositoryList(bid int64) (map[string][]IDTextMap, error) {

	// initialize list with 0-id value
	list := []IDTextMap{{ID: 0, Text: " -- Select Depository -- "}}

	// json response data
	appData := make(map[string][]IDTextMap)

	bud, err := bidToBud(bid)
	if err != nil {
		return appData, err
	}

	m := rlib.GetAllDepositories(bid)
	for i := 0; i < len(m); i++ {
		list = append(list, IDTextMap{ID: m[i].DEPID, Text: m[i].Name})
	}
	appData[bud] = list

	return appData, nil
}

// SvcUIErrAndVarResponse encapsulates a lot of lines that would need to appear
// in each case of a switch.  This just makes things a lot more readable and
// it bottlenecks the handling so it is easy to extend or modify.
func SvcUIErrAndVarResponse(w http.ResponseWriter, funcname string, err error, x interface{}) {
	if err != nil {
		SvcGridErrorReturn(w, err, funcname)
		return
	}
	if err := json.NewEncoder(w).Encode(x); err != nil {
		SvcGridErrorReturn(w, err, funcname)
		return
	}
}

// SvcUIVal returns the requested variable in JSON form
//
// wsdoc {
//  @Title  Get UI Value
//	@URL /v1/uival/:BID/varname
//  @Method  GET
//	@Synopsis Return JSON representing the UI Value
//  @Desc Return data can be parsed to create the string lists used in the UI.
//	@Input
//  @Response JSONResponse
// wsdoc }
func SvcUIVal(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcUIVar"
	rlib.Console("Entered %s\n", funcname)
	switch d.DetVal {
	case "app.AssessmentRules":
		asmData, err := GetAssessmentList(d.BID)
		SvcUIErrAndVarResponse(w, funcname, err, asmData)
	case "app.ReceiptRules":
		rcptData, err := GetReceiptList(d.BID)
		SvcUIErrAndVarResponse(w, funcname, err, rcptData)
	case "app.Depositories":
		data, err := GetDepositoryList(d.BID)
		SvcUIErrAndVarResponse(w, funcname, err, data)
	case "app.depmeth":
		depmeth := GetJSDepositMethods()
		SvcUIErrAndVarResponse(w, funcname, nil, depmeth)
	default:
		e := fmt.Errorf("Unknown variable requested: %s", d.DetVal)
		SvcGridErrorReturn(w, e, funcname)
		return
	}
}
