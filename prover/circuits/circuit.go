package tradingFeesRebate

import (
	"github.com/brevis-network/brevis-sdk/sdk"
	"github.com/ethereum/go-ethereum/common"
)

type AppCircuit struct{}

var _ sdk.AppCircuit = &AppCircuit{}

func (c *AppCircuit) Allocate() (maxReceipts, maxStorage, maxTransactions int) {
	// This demo app is only going to use two storage data at a time so
	// we can simply limit the max number of data for storage to 1 and
	// 0 for all others
	return 0, 5, 0
}

var TRADING_FEE_TOKEN = sdk.ConstUint248(
	common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7"))

var UserAddress = sdk.ConstUint248(
	common.HexToAddress("0xB14a13ddaEa5df325732DB991F1A766ae0DbD75a"))

func (c *AppCircuit) Define(api *sdk.CircuitAPI, in sdk.DataInput) error {
	slots := sdk.NewDataStream(api, in.StorageSlots)

	var u248 = api.Uint248
	var b32 = api.Bytes32
	sdk.AssertEach(slots, func(current sdk.StorageSlot) sdk.Uint248 {
		contractIsEq := u248.IsEqual(current.Contract, TRADING_FEE_TOKEN)

		// mapping(address => uint) public tradingFees;
		// balance slot location 0x0000000000000000000000000000000000000000000000000000000000000002
		// slot key = keccak(u256(userAddress).u256(location)
		tradingFeeSlot := api.SlotOfStructFieldInMapping(2, 0, api.ToBytes32(UserAddress))
		tradingFeeSlotKeyIsEq := b32.IsEqual(current.Slot, tradingFeeSlot)

		return u248.And(
			contractIsEq,
			tradingFeeSlotKeyIsEq,
		)
	})

	tradingFees := sdk.Map(slots, func(current sdk.StorageSlot) sdk.Uint248 {
		return api.ToUint248(current.Value)
	})
	totalTradingFee := sdk.Sum(tradingFees)
	avgTradingFee, _ := u248.Div(totalTradingFee, sdk.ConstUint248(5))

	api.OutputAddress(UserAddress)
	api.OutputUint(248, totalTradingFee)
	api.OutputUint(248, avgTradingFee)
	return nil
}
