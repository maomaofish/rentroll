package bizlogic

import (
	"fmt"
	"rentroll/rlib"
)

// SaveDeposit validates that all the information in the deposit
// meets the business criteria for a deposit and Inserts or Updates
// the deposit if the criteria are met.
//
// @params
//	a - the depost struct
//  rcpts - an array of RCPTIDs that are being targeted for the deposit.
//
// @returns
//	errlist - an array of errors
//-----------------------------------------------------------------------
func SaveDeposit(a *rlib.Deposit, newRcpts []int64) []BizError {
	var e []BizError
	var rlist []rlib.Receipt
	tot := float64(0)
	//------------------------------------------------------------
	// First, validate that all newRcpts are eligible for inclusion
	// in this receipt
	//------------------------------------------------------------
	for i := 0; i < len(newRcpts); i++ {
		r := rlib.GetReceipt(newRcpts[i])
		tot += r.Amount
		if r.DID != 0 && r.DID != a.DID {
			s := fmt.Sprintf(BizErrors[ReceiptAlreadyDeposited].Message, rlib.IDtoShortString("RCPT", r.RCPTID), rlib.IDtoShortString("D", r.DID))
			b := BizError{Errno: ReceiptAlreadyDeposited, Message: s}
			e = append(e, b)
			continue
		}
		if r.BID != a.BID {
			s := fmt.Sprintf(BizErrors[ReceiptBizMismatch].Message, rlib.IDtoShortString("RCPT", r.RCPTID))
			b := BizError{Errno: ReceiptBizMismatch, Message: s}
			e = append(e, b)
			continue
		}
		rlist = append(rlist, r)
	}
	//------------------------------------------------------------
	// next, validate that the total of all newRcpts matches Amount
	//------------------------------------------------------------
	if tot != a.Amount {
		e = AddBizErrToList(e, DepositTotalMismatch)
		return e
	}

	//------------------------------------------------------------
	// Save the deposit
	//------------------------------------------------------------
	if a.DID == 0 {
		_, err := rlib.InsertDeposit(a)
		if err != nil {
			e = AddErrToBizErrlist(err, e)
		}
		for i := 0; i < len(newRcpts); i++ {
			var dp = rlib.DepositPart{
				DID:    a.DID,
				BID:    a.BID,
				RCPTID: newRcpts[i],
			}
			err = rlib.InsertDepositPart(&dp)
			if err != nil {
				e = AddErrToBizErrlist(err, e)
			}
			if rlist[i].DID == 0 {
				rlist[i].DID = a.DID
				err = rlib.UpdateReceipt(&rlist[i])
				if err != nil {
					e = AddErrToBizErrlist(err, e)
				}
			}
		}
	} else {
		// err := rlib.UpdateDeposit(a)
		// if err != nil {
		// 	e = AddErrToBizErrlist(err, e)
		// }
		// //---------------------------------------------------------------------------
		// // If any receipts have been removed from the previous version.  To do
		// // this we will compare the list of current Deposit's RCPTIDs to the
		// // list of newly proposed RCPTIDs.  We will compare the two lists and
		// // produce 2 new lists: addlist and removelist.  Then we will add and
		// // link the addlist, and unlink the removelist.  The new Receipts are
		// // already provided in newRcpts.
		// //---------------------------------------------------------------------------
		// curRcpts, err = rlib.GetDepositParts(a.DID)
		// if err != nil {
		// 	e = AddErrToBizErrlist(err, e)
		// 	return e
		// }

		// current := map[int64]int{}
		// for i := 0; i < len(curRcpts); i++ {
		// 	current[curRcpts[i]] = 0
		// }

		// var addlist []int64
		// for i := 0; i < len(newRcpts); i++ {
		// 	_, ok := current[newRcpts[i]]
		// 	if !ok {
		// 		addlist = append(addlist, newRcpts[i])
		// 	}
		// }

		// var removelist []int64

	}
	return e
}