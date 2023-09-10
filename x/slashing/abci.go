package slashing

import (
	"context"
	"time"

	"cosmossdk.io/core/comet"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	"github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx context.Context, k keeper.Keeper) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// Iterate over all the validators which *should* have signed this block
	// store whether or not they have actually signed it and slash/unbond any
	// which have missed too many blocks in a row (downtime slashing)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	for i := 0; i < sdkCtx.CometInfo().GetLastCommit().Votes().Len(); i++ {
		vote := sdkCtx.CometInfo().GetLastCommit().Votes().Get(i)
		err := k.HandleValidatorSignature(ctx, vote.Validator().Address(), vote.Validator().Power(), comet.BlockIDFlag(vote.GetBlockIDFlag()))
		if err != nil {
			return err
		}
	}
	return nil
}
