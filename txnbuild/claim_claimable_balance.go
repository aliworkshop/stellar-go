//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package txnbuild

import (
	"github.com/aliworkshop/stellar-go/support/errors"
	"github.com/aliworkshop/stellar-go/xdr"
)

// ClaimClaimableBalance represents the Stellar claim claimable balance operation. See
// https://www.stellar.org/developers/guides/concepts/list-of-operations.html
type ClaimClaimableBalance struct {
	BalanceID     string
	SourceAccount string
}

// BuildXDR for ClaimClaimableBalance returns a fully configured XDR Operation.
func (cb *ClaimClaimableBalance) BuildXDR(withMuxedAccounts bool) (xdr.Operation, error) {
	var xdrBalanceID xdr.ClaimableBalanceId
	err := xdr.SafeUnmarshalHex(cb.BalanceID, &xdrBalanceID)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to set XDR 'ClaimableBalanceId' field")
	}
	xdrOp := xdr.ClaimClaimableBalanceOp{
		BalanceId: xdrBalanceID,
	}

	opType := xdr.OperationTypeClaimClaimableBalance
	body, err := xdr.NewOperationBody(opType, xdrOp)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to build XDR OperationBody")
	}
	op := xdr.Operation{Body: body}
	if withMuxedAccounts {
		SetOpSourceMuxedAccount(&op, cb.SourceAccount)
	} else {
		SetOpSourceAccount(&op, cb.SourceAccount)
	}
	return op, nil
}

// FromXDR for ClaimClaimableBalance initializes the txnbuild struct from the corresponding xdr Operation.
func (cb *ClaimClaimableBalance) FromXDR(xdrOp xdr.Operation, withMuxedAccounts bool) error {
	result, ok := xdrOp.Body.GetClaimClaimableBalanceOp()
	if !ok {
		return errors.New("error parsing claim_claimable_balance operation from xdr")
	}

	cb.SourceAccount = accountFromXDR(xdrOp.SourceAccount, withMuxedAccounts)
	balanceID, err := xdr.MarshalHex(result.BalanceId)
	if err != nil {
		return errors.New("error parsing BalanceID in claim_claimable_balance operation from xdr")
	}
	cb.BalanceID = balanceID

	return nil
}

// Validate for ClaimClaimableBalance validates the required struct fields. It returns an error if any of the fields are
// invalid. Otherwise, it returns nil.
func (cb *ClaimClaimableBalance) Validate(withMuxedAccounts bool) error {
	var xdrBalanceID xdr.ClaimableBalanceId
	err := xdr.SafeUnmarshalHex(cb.BalanceID, &xdrBalanceID)
	if err != nil {
		return NewValidationError("BalanceID", err.Error())
	}

	return nil
}

// GetSourceAccount returns the source account of the operation, or the empty string if not
// set.
func (cb *ClaimClaimableBalance) GetSourceAccount() string {
	return cb.SourceAccount
}
