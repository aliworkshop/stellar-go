//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package txnbuild

import (
	"github.com/aliworkshop/stellar-go/support/errors"
	"github.com/aliworkshop/stellar-go/xdr"
)

// EndSponsoringFutureReserves represents the Stellar begin sponsoring future reserves operation. See
// https://www.stellar.org/developers/guides/concepts/list-of-operations.html
type EndSponsoringFutureReserves struct {
	SourceAccount string
}

// BuildXDR for EndSponsoringFutureReserves returns a fully configured XDR Operation.
func (es *EndSponsoringFutureReserves) BuildXDR(withMuxedAccounts bool) (xdr.Operation, error) {
	opType := xdr.OperationTypeEndSponsoringFutureReserves
	body, err := xdr.NewOperationBody(opType, nil)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to build XDR OperationBody")
	}
	op := xdr.Operation{Body: body}
	if withMuxedAccounts {
		SetOpSourceMuxedAccount(&op, es.SourceAccount)
	} else {
		SetOpSourceAccount(&op, es.SourceAccount)
	}
	return op, nil
}

// FromXDR for EndSponsoringFutureReserves initializes the txnbuild struct from the corresponding xdr Operation.
func (es *EndSponsoringFutureReserves) FromXDR(xdrOp xdr.Operation, withMuxedAccounts bool) error {
	if xdrOp.Body.Type != xdr.OperationTypeEndSponsoringFutureReserves {
		return errors.New("error parsing end_sponsoring_future_reserves operation from xdr")
	}

	es.SourceAccount = accountFromXDR(xdrOp.SourceAccount, withMuxedAccounts)
	return nil
}

// Validate for EndSponsoringFutureReserves validates the required struct fields. It returns an error if any of the fields are
// invalid. Otherwise, it returns nil.
func (es *EndSponsoringFutureReserves) Validate(withMuxedAccounts bool) error {
	return nil
}

// GetSourceAccount returns the source account of the operation, or the empty string if not
// set.
func (es *EndSponsoringFutureReserves) GetSourceAccount() string {
	return es.SourceAccount
}
