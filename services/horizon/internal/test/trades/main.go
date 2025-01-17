//Common infrastructure for testing Trades
package trades

import (
	"context"

	"github.com/stellar/go/keypair"
	"github.com/stellar/go/services/horizon/internal/db2/history"
	"github.com/stellar/go/support/time"
	"github.com/stellar/go/xdr"
)

//GetTestAsset generates an issuer on the fly and creates a CreditAlphanum4 Asset with given code
func GetTestAsset(code string) xdr.Asset {
	var codeBytes [4]byte
	copy(codeBytes[:], []byte(code))
	ca4 := xdr.AlphaNum4{Issuer: GetTestAccount(), AssetCode: codeBytes}
	return xdr.Asset{Type: xdr.AssetTypeAssetTypeCreditAlphanum4, AlphaNum4: &ca4, AlphaNum12: nil}
}

//Get generates and returns an account on the fly
func GetTestAccount() xdr.AccountId {
	var key xdr.Uint256
	kp, _ := keypair.Random()
	copy(key[:], kp.Address())
	acc, _ := xdr.NewAccountId(xdr.PublicKeyTypePublicKeyTypeEd25519, key)
	return acc
}

//IngestTestTrade mock ingests a trade
func IngestTestTrade(
	q *history.Q,
	assetSold xdr.Asset,
	assetBought xdr.Asset,
	seller xdr.AccountId,
	buyer xdr.AccountId,
	amountSold int64,
	amountBought int64,
	timestamp time.Millis,
	opCounter int64) error {

	trade := xdr.ClaimAtom{
		Type: xdr.ClaimAtomTypeClaimAtomTypeOrderBook,
		OrderBook: &xdr.ClaimOfferAtom{
			AmountBought: xdr.Int64(amountBought),
			SellerId:     seller,
			AmountSold:   xdr.Int64(amountSold),
			AssetBought:  assetBought,
			AssetSold:    assetSold,
		},
	}

	price := xdr.Price{
		N: xdr.Int32(amountBought),
		D: xdr.Int32(amountSold),
	}

	ctx := context.Background()
	accounts, err := q.CreateAccounts(ctx, []string{seller.Address(), buyer.Address()}, 2)
	if err != nil {
		return err
	}
	assets, err := q.CreateAssets(ctx, []xdr.Asset{assetBought, assetSold}, 2)
	if err != nil {
		return err
	}

	batch := q.NewTradeBatchInsertBuilder(0)
	batch.Add(ctx, history.InsertTrade{
		HistoryOperationID: opCounter,
		Order:              0,
		BuyOfferExists:     false,
		BuyOfferID:         0,
		BoughtAssetID:      assets[assetBought.String()].ID,
		SoldAssetID:        assets[assetSold.String()].ID,
		LedgerCloseTime:    timestamp.ToTime(),
		SellPrice:          price,
		Trade:              trade,
		BuyerAccountID:     accounts[buyer.Address()],
		SellerAccountID:    accounts[seller.Address()],
	})
	err = batch.Exec(ctx)
	if err != nil {
		return err
	}

	err = q.RebuildTradeAggregationTimes(context.Background(), timestamp, timestamp)
	if err != nil {
		return err
	}

	return nil
}

//PopulateTestTrades generates and ingests trades between two assets according to given parameters
func PopulateTestTrades(
	q *history.Q,
	startTs int64,
	numOfTrades int,
	delta int64,
	opStart int64) (ass1 xdr.Asset, ass2 xdr.Asset, err error) {

	acc1 := GetTestAccount()
	acc2 := GetTestAccount()
	ass1 = GetTestAsset("usd")
	ass2 = GetTestAsset("euro")
	for i := 1; i <= numOfTrades; i++ {
		timestamp := time.MillisFromInt64(startTs + (delta * int64(i-1)))
		err = IngestTestTrade(
			q, ass1, ass2, acc1, acc2, int64(i*100), int64(i*100)*int64(i), timestamp, opStart+int64(i))
		//tt.Assert.NoError(err)
		if err != nil {
			return
		}
	}
	return
}
