package ws

// Fees2RA copies fees into permanent table Assessments.
import (
	"context"
	"fmt"
	"rentroll/bizlogic"
	"rentroll/rlib"
	"time"
)

// Fees2RA handles all the updates necessary to move the fees defined in a flow
// into the permanent tables.
//
// INPUTS
//     ctx    - db context for transactions
//     x - all the contextual info we need for performing this operation
//
// RETURNS
//     Any errors encountered
//-----------------------------------------------------------------------------
func Fees2RA(ctx context.Context, x *WriteHandlerContext) error {
	var err error

	//--------------------------------------------------
	// Handle Rentables first...
	//--------------------------------------------------
	rlib.Console("Rentables fees\n")
	for i := 0; i < len(x.raf.Rentables); i++ {
		for j := 0; j < len(x.raf.Rentables[i].Fees); j++ {
			if err = F2RASaveFee(ctx, x, &x.raf.Rentables[i].Fees[j], rlib.ELEMRENTABLE, x.raf.Rentables[i].RID, 0); err != nil {
				return err
			}
		}
	}
	//--------------------------------------------------
	// Handle pet fees...
	//--------------------------------------------------
	rlib.Console("Pet fees\n")
	for i := 0; i < len(x.raf.Pets); i++ {
		for j := 0; j < len(x.raf.Pets[i].Fees); j++ {
			if err = F2RASaveFee(ctx, x, &x.raf.Pets[i].Fees[j], rlib.ELEMPET, x.raf.Pets[i].PETID, x.raf.Pets[i].TMPTCID); err != nil {
				return err
			}
		}
	}
	//--------------------------------------------------
	// Handle vehicle fees...
	//--------------------------------------------------
	rlib.Console("Vehicle fees\n")
	for i := 0; i < len(x.raf.Vehicles); i++ {
		for j := 0; j < len(x.raf.Vehicles[i].Fees); j++ {
			if err = F2RASaveFee(ctx, x, &x.raf.Vehicles[i].Fees[j], rlib.ELEMVEHICLE, x.raf.Vehicles[i].VID, x.raf.Vehicles[i].TMPTCID); err != nil {
				return err
			}
		}
	}

	//--------------------------------------------------------------------------
	// Now clean up any assessments that are associated with the old RAID but
	// that have not been updated as part of any fee in the new RAID.
	//--------------------------------------------------------------------------
	if err = CleanUpRemainingAssessments(ctx, x); err != nil {
		return err
	}

	return nil
}

// CleanUpRemainingAssessments handles all the assessments associated with the
// original RAID that were not found while handling the amended RAID.
//
//
// INPUTS
//     ctx  - db context for transactions
//     x    - all the contextual info we need for performing this operation
//
// RETURNS
//     Any errors encountered
//-----------------------------------------------------------------------------
func CleanUpRemainingAssessments(ctx context.Context, x *WriteHandlerContext) error {
	rlib.Console("Entered CleanUpRemainingAssessments")
	if x.raf.Meta.RAID == 0 {
		rlib.Console("No cleanup necessary. x.raf.Meta.RAID is 0\n")
		return nil // nothing to do, no old RAID
	}
	//--------------------------------------------------------------------------
	// Get the list of any recurring assessments associated with the old rental
	// agreement that overlap the time range of the new rental agreement.
	//--------------------------------------------------------------------------
	m, err := rlib.GetRecurringAssessmentDefsByRAID(ctx, x.raOrig.RAID, &x.ra.RentStart, &x.ra.RentStop)
	if err != nil {
		return err
	}
	rlib.Console("Found %d recurring assessment definitions to review\n", len(m))
	for _, v := range m {
		rlib.Console("ASMID = %d\n", v.ASMID)
		if v.RentCycle == rlib.RECURNONE {
			// If it is a non-recurring assessment, reverse it.
			be := bizlogic.ReverseAssessment(ctx, &v, 0, &x.ra.RentStart)
			if len(be) > 0 {
				return bizlogic.BizErrorListToError(be)
			}
		} else {
			// If it is a recurring assessment, stop it.
			if err = bizlogic.UpdateAssessmentEndDate(ctx, &v, &x.ra.RentStart); err != nil {
				return err
			}
		}
	}
	rlib.Console("*** Completed processing recurring assessment definitions\n")
	//--------------------------------------------------------------------------
	// Anything non-recurring that happens as of the start date of the amended
	// agreement must be deleted (reversed).
	//--------------------------------------------------------------------------
	m, err = rlib.GetNorecurAssessmentsByRAIDRange(ctx, x.raOrig.RAID, &x.ra.RentStart, &x.ra.RentStop)
	if err != nil {
		return err
	}
	rlib.Console("Found %d non-recurring assessments to reverse\n", len(m))
	for _, v := range m {
		v.AppendComment(fmt.Sprintf("Reversing due to amended RAID %d", x.ra.RAID))
		bizlogic.ReverseAssessment(ctx, &v, 0, &x.ra.RentStart)
	}
	rlib.Console("*** Completed processing non-recurring\n")
	return nil
}

// F2RASaveFee handles all the updates necessary to move the
// supplied fee into the permanent tables.
//
// INPUTS
//     ctx  - db context for transactions
//     x    - all the contextual info we need for performing this operation
//     elt  - element type if is this is bound to a pet or vehicle
//     id   - RID if elt == rlib.ELEMRENTABLE, or tmpid of the element
//            (TMPPETID, TMPVID), valid if elt > 0
//     tcid - tmptcid of the transactant responsible, valid if elt > 0
//
// RETURNS
//     Any errors encountered
//-----------------------------------------------------------------------------
func F2RASaveFee(ctx context.Context, x *WriteHandlerContext, fee *rlib.RAFeesData, eltype, id, tmptcid int64) error {
	rlib.Console("F2RASaveFee processing fee = %d, ASMID = %d\n", fee.TMPASMID, fee.ASMID)
	if 0 < fee.ASMID {
		return F2RAUpdateExistingAssessment(ctx, x, fee, eltype, id, tmptcid)
	}
	return F2RASaveNewFee(ctx, x, fee, eltype, id, tmptcid)

}

// F2RASaveNewFee handles all the updates necessary to move the
// supplied fee into the permanent tables.
//
// INPUTS
//     ctx  - db context for transactions
//     x    - all the contextual info we need for performing this operation
//     elt  - element type if is this is bound to a pet or vehicle
//     id   - RID if elt == rlib.ELEMRENTABLE, or tmpid of the element
//            (TMPPETID, TMPVID), valid if elt > 0
//     tcid - tmptcid of the transactant responsible, valid if elt > 0
//
// RETURNS
//     Any errors encountered
//-----------------------------------------------------------------------------
func F2RASaveNewFee(ctx context.Context, x *WriteHandlerContext, fee *rlib.RAFeesData, eltype, id, tmptcid int64) error {
	// rlib.Console("Entered F2RASaveNewFee\n")
	//-------------------------------------------------------------------
	// Create a new assessment from this day forward...
	//-------------------------------------------------------------------
	var b rlib.Assessment
	dt := time.Time(x.raf.Dates.AgreementStart)
	if fee.ASMID > 0 {
		b.AppendComment(fmt.Sprintf("Continuation of ASMID %d from RAID %d", fee.ASMID, x.raf.Meta.RAID))
	}
	Start := time.Time(fee.Start) // the start time will be either the fee start
	if Start.Before(dt) {         // or the start of the new rental agreement
		Start = dt // whichever is later
	}
	b.Stop = time.Time(fee.Stop)
	b.BID = x.raf.Dates.BID

	//-------------------------------------------------------------------
	// Set the Element Type and ID if necessary
	//-------------------------------------------------------------------
	b.AssocElemType = eltype
	b.AssocElemID = id

	//-------------------------------------------------------------------
	// find the RID associated with this pet
	//-------------------------------------------------------------------
	switch eltype {
	case rlib.ELEMRENTABLE:
		b.RID = id
	case rlib.ELEMPET:
		if b.RID = GetRIDForTMPTCID(ctx, x, tmptcid); b.RID <= 0 {
			return fmt.Errorf("No RID associated with TMPTCID = %d", tmptcid)
		}
		// rlib.Console("GetRIDForTMPTCID( TMPTCID=%d) ===> %d\n", tmptcid, b.RID)
		// rlib.Console("    ID for this pet is %d\n", b.AssocElemID)
	case rlib.ELEMVEHICLE:
		if b.RID = GetRIDForTMPTCID(ctx, x, tmptcid); b.RID <= 0 {
			return fmt.Errorf("No RID associated with TMPTCID = %d", tmptcid)
		}
		// rlib.Console("GetRIDForTMPTCID( TMPTCID=%d) ===> %d\n", tmptcid, b.RID)
		// rlib.Console("    ID for this vehicle is %d\n", b.AssocElemID)
	}

	// rlib.Console("bid = %d, fee ARID = %d\n", b.BID, fee.ARID)
	b.Amount = fee.ContractAmount
	b.AcctRule = ""
	b.RentCycle = fee.RentCycle
	b.RAID = x.ra.RAID
	b.Start = time.Time(fee.Start)
	b.Stop = time.Time(fee.Stop)
	b.RentCycle = fee.RentCycle
	b.ProrationCycle = rlib.RRdb.BizTypes[b.BID].AR[fee.ARID].DefaultProrationCycle
	b.InvoiceNo = 0
	b.ARID = fee.ARID
	switch eltype {
	case rlib.ELEMRENTABLE:
		// nothing to do at this time
	case rlib.ELEMPET:
		b.FLAGS |= 1 << 3  // PETID required
		b.AssocElemID = id // must be the PETID
		b.AssocElemType = eltype
	case rlib.ELEMVEHICLE:
		b.FLAGS |= 1 << 4  // VID required
		b.AssocElemID = id // must be the PETID
		b.AssocElemType = eltype
	}

	_, err := rlib.InsertAssessment(ctx, &b)
	if err != nil {
		return err
	}
	return nil
}

// F2RAUpdateExistingAssessment handles all the updates necessary to move the
// supplied fee into the permanent tables.
//
// INPUTS
//     ctx  - db context for transactions
//     x    - all the contextual info we need for performing this operation
//     elt  - element type if is this is bound to a pet or vehicle
//     id   - RID if elt == rlib.ELEMRENTABLE, or tmpid of the element
//            (TMPPETID, TMPVID), valid if elt > 0
//     tcid - tmptcid of the transactant responsible, valid if elt > 0
//
// RETURNS
//     Any errors encountered
//-----------------------------------------------------------------------------
func F2RAUpdateExistingAssessment(ctx context.Context, x *WriteHandlerContext, fee *rlib.RAFeesData, eltype, id, tmptcid int64) error {
	rlib.Console("Entered F2RAUpdateExistingAssessment\n")
	if fee.ASMID == int64(0) {
		return fmt.Errorf("fee.ASMID must be > 0")
	}
	a, err := rlib.GetAssessment(ctx, fee.ASMID)
	if err != nil {
		return err
	}

	//-------------------------------------------------------------------
	// skip any assessments that finished prior to this Rental Agreement
	//-------------------------------------------------------------------
	dt := time.Time(x.raf.Dates.AgreementStart)
	stop := time.Time(fee.Stop)
	if stop.Before(dt) {
		return nil // don't need to process this one
	}

	//-------------------------------------------------------------------
	// skip any non-recurring assessment that has been paid...
	//-------------------------------------------------------------------
	if a.FLAGS&3 == 2 {
		return nil // don't need to process this
	}

	//-------------------------------------------------------------------
	// If it's recurring we'll just stop it on the start date of the new
	// rental agreement
	//-------------------------------------------------------------------
	if a.RentCycle > rlib.RECURNONE {
		err = bizlogic.UpdateAssessmentEndDate(ctx, &a, &dt)
		if err != nil {
			return err
		}
		// a.Stop = dt
		// if err = rlib.UpdateAssessment(ctx, &a); err != nil {
		// 	return err
		// }
	}

	err = F2RASaveNewFee(ctx, x, fee, eltype, id, tmptcid)
	if err != nil {
		return err
	}

	return nil
}

// GetRIDForTMPTCID finds the RID associated with the supplied tmptcid.
// This routine is called when we have a pet or a vehicle and we need to
// know what RID to associate it with. The RATiePeopleData datastruct is
// organized by tmptcid.  PRID is the Parent RID for that person.
//
// INPUTS
//     ctx     - db context for transactions
//     x       - all the contextual info we need for performing this operation
//     tmptcid - tmptcid for person we want the associated RID
//
// RETURNS
//     RID of associated rentable, or -1 if not found
//-----------------------------------------------------------------------------
func GetRIDForTMPTCID(ctx context.Context, x *WriteHandlerContext, tmptcid int64) int64 {
	for i := 0; i < len(x.raf.Tie.People); i++ {
		if x.raf.Tie.People[i].TMPTCID == tmptcid {
			return x.raf.Tie.People[i].PRID
		}
	}
	return -1
}

// GetTCIDForTMPTCID finds the TCID associated with the supplied tmptcid.
//
// INPUTS
//     ctx     - db context for transactions
//     x       - all the contextual info we need for performing this operation
//     tmptcid - tmptcid for person we want the associated RID
//
// RETURNS
//     TCID of associated Transactant, or -1 if not found
//-----------------------------------------------------------------------------
func GetTCIDForTMPTCID(x *WriteHandlerContext, tmptcid int64) int64 {
	for i := 0; i < len(x.raf.People); i++ {
		if x.raf.People[i].TMPTCID == tmptcid {
			return x.raf.People[i].TCID
		}
	}
	return -1
}
