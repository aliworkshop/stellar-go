package horizonclient

import (
	"testing"

	"github.com/aliworkshop/stellar-go/support/http/httptest"
	"github.com/stretchr/testify/assert"
)

func TestFund(t *testing.T) {
	friendbotFundResponse := `{
  "_links": {
    "transaction": {
      "href": "https://horizon-testnet.stellar.org/transactions/94e42f65d3ff5f30669b6109c2ce3e82c0e592c52004e3b41bb30e24df33954e"
    }
  },
  "hash": "94e42f65d3ff5f30669b6109c2ce3e82c0e592c52004e3b41bb30e24df33954e",
  "ledger": 8269,
  "envelope_xdr": "AAAAAgAAAAD2Leuk4afNVCYqxbN03yPH6kgKe/o2yiOd3CQNkpkpQwABhqAAAAFSAAAACQAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAABAAAAABB90WssODNIgi6BHveqzxTRmIpvAFRyVNM+Hm2GVuCcAAAAAAAAAABW9+rbvt6YXwwXyFszptQFlfzzFMrWObLiJmBhOzNblAAAABdIdugAAAAAAAAAAAKSmSlDAAAAQHWNbXOoVQqH0YJRr8LAtpalV+NoXb8Tv/ETkPNv2NignhN8seUSde8m2HLNLHOo+5W34BXfxfBmDXgZn8yHkwSGVuCcAAAAQDQLh1UAxYZ27sIxyYgyYFo8IUbTiANWadUJUR7K0q1eY6Q5J/BFfNlf6UqLqJ5zd8uI3TXCaBNJDkiQc1ZLEg4=",
  "result_xdr": "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=",
  "result_meta_xdr": "AAAAAgAAAAIAAAADAAAgTQAAAAAAAAAA9i3rpOGnzVQmKsWzdN8jx+pICnv6NsojndwkDZKZKUMAAAAAPDNbbAAAAVIAAAAIAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAgTQAAAAAAAAAA9i3rpOGnzVQmKsWzdN8jx+pICnv6NsojndwkDZKZKUMAAAAAPDNbbAAAAVIAAAAJAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAACBMAAAAAAAAAAAQfdFrLDgzSIIugR73qs8U0ZiKbwBUclTTPh5thlbgnAFg09HQY/uMAAAA2wAAAAoAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAACBNAAAAAAAAAAAQfdFrLDgzSIIugR73qs8U0ZiKbwBUclTTPh5thlbgnAFg07qH7ROMAAAA2wAAAAoAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAACBNAAAAAAAAAABW9+rbvt6YXwwXyFszptQFlfzzFMrWObLiJmBhOzNblAAAABdIdugAAAAgTQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAA="
}`

	hmock := httptest.NewClient()
	client := &Client{
		HorizonURL: "https://localhost/",
		HTTP:       hmock,
		isTestNet:  true,
	}

	hmock.On(
		"GET",
		"https://localhost/friendbot?addr=GBLPP2W3X3PJQXYMC7EFWM5G2QCZL7HTCTFNMONS4ITGAYJ3GNNZIQ4V",
	).ReturnString(200, friendbotFundResponse)

	tx, err := client.Fund("GBLPP2W3X3PJQXYMC7EFWM5G2QCZL7HTCTFNMONS4ITGAYJ3GNNZIQ4V")
	assert.NoError(t, err)
	assert.Equal(t, int32(8269), tx.Ledger)
}
