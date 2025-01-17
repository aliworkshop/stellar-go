package horizonclient

import (
	"encoding/json"
	"strings"

	hProtocol "github.com/aliworkshop/stellar-go/protocols/horizon"
	"github.com/aliworkshop/stellar-go/support/errors"
	"github.com/aliworkshop/stellar-go/xdr"
)

func (herr Error) Error() string {
	s := strings.Builder{}
	s.WriteString(`horizon error: "`)
	s.WriteString(herr.Problem.Title)
	s.WriteString(`" `)
	if rc, err := herr.ResultCodes(); err == nil {
		s.WriteString(`(`)
		resultCodes := append([]string{rc.TransactionCode}, rc.OperationCodes...)
		s.WriteString(strings.Join(resultCodes, `, `))
		s.WriteString(`) `)
	}
	s.WriteString(`- check horizon.Error.Problem for more information`)
	return s.String()
}

// Envelope extracts the transaction envelope that triggered this error from the
// extra fields.
func (herr *Error) Envelope() (*xdr.TransactionEnvelope, error) {
	b64, err := herr.EnvelopeXDR()
	if err != nil {
		return nil, err
	}

	var result xdr.TransactionEnvelope
	err = xdr.SafeUnmarshalBase64(b64, &result)
	return &result, errors.Wrap(err, "xdr decode failed")
}

// EnvelopeXDR returns the base 64 serialised string representation of the XDR envelope.
// This can be stored, or decoded in the Stellar Laboratory XDR viewer for example.
func (herr *Error) EnvelopeXDR() (string, error) {
	raw, ok := herr.Problem.Extras["envelope_xdr"]
	if !ok {
		return "", ErrEnvelopeNotPopulated
	}

	var b64 string
	b64, ok = raw.(string)
	if !ok {
		return "", errors.New("type assertion failed")
	}

	return b64, nil
}

// ResultString extracts the transaction result as a string.
func (herr *Error) ResultString() (string, error) {
	raw, ok := herr.Problem.Extras["result_xdr"]
	if !ok {
		return "", ErrResultNotPopulated
	}

	b64, ok := raw.(string)
	if !ok {
		return "", errors.New("type assertion failed")
	}

	return b64, nil
}

// ResultCodes extracts a result code summary from the error, if possible.
func (herr *Error) ResultCodes() (*hProtocol.TransactionResultCodes, error) {

	raw, ok := herr.Problem.Extras["result_codes"]
	if !ok {
		return nil, ErrResultCodesNotPopulated
	}

	// converts map to []byte
	dataString, err := json.Marshal(raw)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling failed")
	}

	var result hProtocol.TransactionResultCodes
	if err = json.Unmarshal(dataString, &result); err != nil {
		return nil, errors.Wrap(err, "unmarshaling failed")
	}

	return &result, nil
}
