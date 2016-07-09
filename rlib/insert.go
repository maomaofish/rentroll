package rlib

import "fmt"

// InsertAssessment writes a new assessmenttype record to the database
func InsertAssessment(a *Assessment) error {
	_, err := RRdb.Prepstmt.InsertAssessment.Exec(a.BID, a.RID, a.ATypeLID, a.RAID, a.Amount, a.Start, a.Stop, a.RecurCycle, a.ProrationCycle, a.InvoiceNo, a.AcctRule, a.Comment, a.LastModBy)
	return err
}

// InsertBuilding writes a new Building record to the database
func InsertBuilding(a *Building) (int64, error) {
	var rid = int64(0)
	res, err := RRdb.Prepstmt.InsertBuilding.Exec(a.BID, a.Address, a.Address2, a.City, a.State, a.PostalCode, a.Country, a.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			rid = int64(id)
		}
	} else {
		Ulog("Error inserting Building:  %v\n", err)
	}
	return rid, err
}

// InsertBuildingWithID writes a new Building record to the database with the supplied bldgid
// the Building ID must be set in the supplied Building struct ptr (a.BLDGID).
func InsertBuildingWithID(a *Building) (int64, error) {
	var rid = int64(0)
	res, err := RRdb.Prepstmt.InsertBuildingWithID.Exec(a.BLDGID, a.BID, a.Address, a.Address2, a.City, a.State, a.PostalCode, a.Country, a.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			rid = int64(id)
		}
	} else {
		Ulog("InsertBuildingWithID: error inserting Building:  %v\n", err)
		Ulog("Bldg = %#v\n", *a)
	}
	return rid, err
}

// InsertBusiness writes a new Business record.
// returns the new Business ID and any associated error
func InsertBusiness(b *Business) (int64, error) {
	var bid = int64(0)
	res, err := RRdb.Prepstmt.InsertBusiness.Exec(b.Designation, b.Name, b.DefaultRentalPeriod, b.ParkingPermitInUse, b.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			bid = int64(id)
		}
	}
	return bid, err
}

// InsertJournalEntry writes a new Journal entry to the database
func InsertJournalEntry(j *Journal) (int64, error) {
	var rid = int64(0)
	res, err := RRdb.Prepstmt.InsertJournal.Exec(j.BID, j.RAID, j.Dt, j.Amount, j.Type, j.ID, j.Comment, j.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			rid = int64(id)
		}
	}
	return rid, err
}

// InsertJournalAllocationEntry writes a new JournalAllocation record to the database
func InsertJournalAllocationEntry(ja *JournalAllocation) error {
	_, err := RRdb.Prepstmt.InsertJournalAllocation.Exec(ja.JID, ja.RID, ja.Amount, ja.ASMID, ja.AcctRule)
	return err
}

// InsertJournalMarker writes a new JournalMarker record to the database
func InsertJournalMarker(jm *JournalMarker) error {
	_, err := RRdb.Prepstmt.InsertJournalMarker.Exec(jm.BID, jm.State, jm.DtStart, jm.DtStop)
	return err
}

//======================================
//  LEDGER MARKER
//======================================

// InsertLedgerMarker writes a new LedgerMarker record to the database
func InsertLedgerMarker(l *LedgerMarker) error {
	_, err := RRdb.Prepstmt.InsertLedgerMarker.Exec(l.LID, l.BID, l.DtStart, l.DtStop, l.Balance, l.State, l.LastModBy)
	if err != nil {
		fmt.Printf("InsertLedgerMarker: err = %#v\n", err)
	}
	return err
}

// InsertLedgerEntry writes a new LedgerEntry to the database
func InsertLedgerEntry(l *LedgerEntry) (int64, error) {
	var rid = int64(0)
	res, err := RRdb.Prepstmt.InsertLedgerEntry.Exec(l.BID, l.JID, l.JAID, l.LID, l.RAID, l.Dt, l.Amount, l.Comment, l.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			rid = int64(id)
		}
	} else {
		Ulog("Error inserting LedgerEntry:  %v\n", err)
	}
	return rid, err
}

// InsertLedger writes a new GLAccount to the database
func InsertLedger(l *GLAccount) (int64, error) {
	var rid = int64(0)
	res, err := RRdb.Prepstmt.InsertLedger.Exec(l.PLID, l.BID, l.RAID, l.GLNumber, l.Status, l.Type, l.Name, l.AcctType, l.RAAssociated, l.AllowPost, l.RARequired, l.ManageToBudget, l.Description, l.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			rid = int64(id)
		}
	} else {
		Ulog("Error inserting GLAccount:  %v\n", err)
	}
	return rid, err
}

//======================================
// NOTE
//======================================

// InsertNote writes a new Note to the database
func InsertNote(a *Note) (int64, error) {
	var rid = int64(0)
	res, err := RRdb.Prepstmt.InsertNote.Exec(a.PNID, a.NTID, a.Comment, a.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			rid = int64(id)
		}
	} else {
		Ulog("Error inserting Note:  %v\n", err)
	}
	return rid, err
}

//======================================
// NOTE TYPE
//======================================

// InsertNoteType writes a new NoteType to the database
func InsertNoteType(a *NoteType) (int64, error) {
	var rid = int64(0)
	res, err := RRdb.Prepstmt.InsertNoteType.Exec(a.BID, a.Name, a.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			rid = int64(id)
		}
	} else {
		Ulog("Error inserting NoteType:  %v\n", err)
	}
	return rid, err
}

// InsertRentableSpecialty writes a new RentableSpecialtyType record to the database
func InsertRentableSpecialty(a *RentableSpecialtyType) error {
	_, err := RRdb.Prepstmt.InsertRentableSpecialtyType.Exec(a.RSPID, a.BID, a.Name, a.Fee, a.Description)
	return err
}

// InsertRentableMarketRates writes a new marketrate record to the database
func InsertRentableMarketRates(r *RentableMarketRate) error {
	_, err := RRdb.Prepstmt.InsertRentableMarketRates.Exec(r.RTID, r.MarketRate, r.DtStart, r.DtStop)
	return err
}

// InsertRentableType writes a new RentableType record to the database
func InsertRentableType(a *RentableType) (int64, error) {
	var rid = int64(0)
	res, err := RRdb.Prepstmt.InsertRentableType.Exec(a.BID, a.Style, a.Name, a.RentCycle, a.Proration, a.GSRPC, a.ManageToBudget, a.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			rid = int64(id)
		}
	} else {
		Ulog("Error inserting RentableType:  %v\n", err)
	}
	return rid, err
}

// InsertRentable writes a new Rentable record to the database
func InsertRentable(a *Rentable) (int64, error) {
	var rid = int64(0)
	res, err := RRdb.Prepstmt.InsertRentable.Exec(a.BID, a.Name, a.AssignmentTime, a.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			rid = int64(id)
		}
	} else {
		Ulog("InsertRentable: error inserting Building:  %v\n", err)
		Ulog("Rentable = %#v\n", *a)
	}
	return rid, err
}

// InsertTransactant writes a new Transactant record to the database
func InsertTransactant(a *Transactant) (int64, error) {
	var tid = int64(0)
	res, err := RRdb.Prepstmt.InsertTransactant.Exec(a.USERID, a.PID, a.PRSPID, a.FirstName, a.MiddleName, a.LastName, a.PreferredName, a.CompanyName, a.IsCompany, a.PrimaryEmail, a.SecondaryEmail, a.WorkPhone, a.CellPhone, a.Address, a.Address2, a.City, a.State, a.PostalCode, a.Country, a.Website, a.Notes, a.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			tid = int64(id)
		}
	} else {
		Ulog("InsertTransactant: error inserting Transactant:  %v\n", err)
		Ulog("Transactant = %#v\n", *a)
	}
	return tid, err
}

// InsertUser writes a new User record to the database
func InsertUser(a *User) (int64, error) {
	var tid = int64(0)
	res, err := RRdb.Prepstmt.InsertUser.Exec(a.TCID, a.Points, a.CarMake, a.CarModel, a.CarColor, a.CarYear, a.LicensePlateState, a.LicensePlateNumber, a.ParkingPermitNumber, a.DateofBirth, a.EmergencyContactName, a.EmergencyContactAddress, a.EmergencyContactTelephone, a.EmergencyEmail, a.AlternateAddress, a.EligibleFutureUser, a.Industry, a.Source, a.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			tid = int64(id)
		}
	} else {
		Ulog("InsertUser: error inserting User:  %v\n", err)
		Ulog("User = %#v\n", *a)
	}
	return tid, err
}

// InsertPayor writes a new User record to the database
func InsertPayor(a *Payor) (int64, error) {
	var tid = int64(0)
	res, err := RRdb.Prepstmt.InsertPayor.Exec(a.TCID, a.CreditLimit, a.TaxpayorID, a.AccountRep, a.EligibleFuturePayor, a.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			tid = int64(id)
		}
	} else {
		Ulog("InsertPayor: error inserting Payor:  %v\n", err)
		Ulog("Payor = %#v\n", *a)
	}
	return tid, err
}

// InsertProspect writes a new User record to the database
func InsertProspect(a *Prospect) (int64, error) {
	var tid = int64(0)
	res, err := RRdb.Prepstmt.InsertProspect.Exec(a.TCID, a.EmployerName, a.EmployerStreetAddress, a.EmployerCity, a.EmployerState,
		a.EmployerPostalCode, a.EmployerEmail, a.EmployerPhone, a.Occupation, a.ApplicationFee, a.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			tid = int64(id)
		}
	} else {
		Ulog("InsertProspect: error inserting Prospect:  %v\n", err)
		Ulog("Prospect = %#v\n", *a)
	}
	return tid, err
}

//=======================================================
//  R E C E I P T
//=======================================================

// InsertReceipt writes a new Receipt record to the database
func InsertReceipt(r *Receipt) (int64, error) {
	var tid = int64(0)
	res, err := RRdb.Prepstmt.InsertReceipt.Exec(r.BID, r.RAID, r.PMTID, r.Dt, r.DocNo, r.Amount, r.AcctRule, r.Comment, r.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			tid = int64(id)
		}
	} else {
		Ulog("InsertReceipt: error inserting Receipt:  %v\n", err)
		Ulog("Receipt = %#v\n", *r)
	}
	return tid, err
}

//=======================================================
//  R E C E I P T   A L L O C A T I O N
//=======================================================

// InsertReceiptAllocation writes a new ReceiptAllocation record to the database
func InsertReceiptAllocation(r *ReceiptAllocation) (int64, error) {
	var tid = int64(0)
	res, err := RRdb.Prepstmt.InsertReceiptAllocation.Exec(r.RCPTID, r.Amount, r.ASMID, r.AcctRule)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			tid = int64(id)
		}
	} else {
		Ulog("InsertReceiptAllocation: error inserting ReceiptAllocation:  %v\n", err)
		Ulog("ReceiptAllocation = %#v\n", *r)
	}
	return tid, err
}

//=======================================================
//  R E N T A L   A G R E E M E N T
//=======================================================

// InsertRentalAgreement writes a new User record to the database
func InsertRentalAgreement(a *RentalAgreement) (int64, error) {
	var tid = int64(0)
	res, err := RRdb.Prepstmt.InsertRentalAgreement.Exec(a.RATID, a.BID, a.NID, a.RentalStart, a.RentalStop, a.PossessionStart, a.PossessionStop, a.Renewal, a.SpecialProvisions, a.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			tid = int64(id)
		}
	} else {
		Ulog("InsertRentalAgreement: error inserting RentalAgreement:  %v\n", err)
		Ulog("RentalAgreement = %#v\n", *a)
	}
	return tid, err
}

// InsertRentalAgreementRentable writes a new User record to the database
func InsertRentalAgreementRentable(a *RentalAgreementRentable) (int64, error) {
	var tid = int64(0)
	res, err := RRdb.Prepstmt.InsertRentalAgreementRentable.Exec(a.RAID, a.RID, a.ContractRent, a.DtStart, a.DtStop)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			tid = int64(id)
		}
	} else {
		Ulog("InsertRentalAgreementRentable: error inserting RentalAgreementRentable:  %v\n", err)
		Ulog("RentalAgreementRentable = %#v\n", *a)
	}
	return tid, err
}

// InsertRentalAgreementPayor writes a new User record to the database
func InsertRentalAgreementPayor(a *RentalAgreementPayor) (int64, error) {
	var tid = int64(0)
	res, err := RRdb.Prepstmt.InsertRentalAgreementPayor.Exec(a.RAID, a.PID, a.DtStart, a.DtStop)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			tid = int64(id)
		}
	} else {
		Ulog("InsertRentalAgreementPayor: error inserting RentalAgreementRentable:  %v\n", err)
		Ulog("RentalAgreementPayor = %#v\n", *a)
	}
	return tid, err
}

// InsertRentalAgreementPet writes a new User record to the database
func InsertRentalAgreementPet(a *RentalAgreementPet) (int64, error) {
	var tid = int64(0)
	res, err := RRdb.Prepstmt.InsertRentalAgreementPet.Exec(a.RAID, a.Type, a.Breed, a.Color, a.Weight, a.Name, a.DtStart, a.DtStop, a.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			tid = int64(id)
		}
	} else {
		Ulog("InsertRentalAgreementPet: error inserting RentalAgreementRentable:  %v\n", err)
		Ulog("RentalAgreementPet = %#v\n", *a)
	}
	return tid, err
}

// InsertRentableUser writes a new User record to the database
func InsertRentableUser(a *RentableUser) (int64, error) {
	var tid = int64(0)
	res, err := RRdb.Prepstmt.InsertRentableUser.Exec(a.RID, a.USERID, a.DtStart, a.DtStop)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			tid = int64(id)
		}
	} else {
		Ulog("InsertRentableUser: error inserting RentableUser:  %v\n", err)
		Ulog("RentableUser = %#v\n", *a)
	}
	return tid, err
}

//=======================================================
//  RENTAL AGREEMENT TEMPLATE
//=======================================================

// InsertRentalAgreementTemplate writes a new User record to the database
func InsertRentalAgreementTemplate(a *RentalAgreementTemplate) (int64, error) {
	var tid = int64(0)
	res, err := RRdb.Prepstmt.InsertRentalAgreementTemplate.Exec(a.BID, a.RentalTemplateNumber, a.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			tid = int64(id)
		}
	} else {
		Ulog("InsertRentalAgreementTemplate: error inserting RentalAgreementTemplate:  %v\n", err)
		Ulog("RentalAgreementTemplate = %#v\n", *a)
	}
	return tid, err
}

// InsertPaymentType writes a new assessmenttype record to the database
func InsertPaymentType(a *PaymentType) error {
	_, err := RRdb.Prepstmt.InsertPaymentType.Exec(a.BID, a.Name, a.Description, a.LastModBy)
	return err
}

// InsertCustomAttribute writes a new User record to the database
func InsertCustomAttribute(a *CustomAttribute) (int64, error) {
	var tid = int64(0)
	res, err := RRdb.Prepstmt.InsertCustomAttribute.Exec(a.Type, a.Name, a.Value, a.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			tid = int64(id)
		}
	} else {
		Ulog("InsertCustomAttribute: error inserting CustomAttribute:  %v\n", err)
		Ulog("CustomAttribute = %#v\n", *a)
	}
	return tid, err
}

// InsertCustomAttributeRef writes a new assessmenttype record to the database
func InsertCustomAttributeRef(a *CustomAttributeRef) error {
	_, err := RRdb.Prepstmt.InsertCustomAttributeRef.Exec(a.ElementType, a.ID, a.CID)
	return err
}

//=======================================================
//  R E N T A B L E   R T I D
//=======================================================

// InsertRentableTypeRef writes a new RentableTypeRef record to the database
func InsertRentableTypeRef(a *RentableTypeRef) error {
	_, err := RRdb.Prepstmt.InsertRentableTypeRef.Exec(a.RID, a.RTID, a.RentCycle, a.ProrationCycle, a.DtStart, a.DtStop, a.LastModBy)
	return err
}

//=======================================================
//  RENTABLE SPECIALTY TYPE
//=======================================================

// InsertRentableSpecialtyRef writes a new RentableSpecialtyType record to the database
func InsertRentableSpecialtyRef(a *RentableSpecialtyRef) error {
	_, err := RRdb.Prepstmt.InsertRentableSpecialtyRef.Exec(a.BID, a.RID, a.RSPID, a.DtStart, a.DtStop, a.LastModBy)
	return err
}

//=======================================================
//  R E N T A B L E   S T A T U S
//=======================================================

// InsertRentableStatus writes a new RentableStatus record to the database
func InsertRentableStatus(a *RentableStatus) error {
	_, err := RRdb.Prepstmt.InsertRentableStatus.Exec(a.RID, a.DtStart, a.DtStop, a.Status, a.LastModBy)
	return err
}