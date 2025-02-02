package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	DefaultOriginAccountsCreated = false
	DefaultExpiryBlock           = uint64(5000)
	DefaultDelegationRequirement = sdk.NewDecWithPrec(1, 1)
	DefaultAirdropFactor         = sdk.NewDecWithPrec(1, 1)
)

func DefaultParams() Params {
	return Params{
		OriginAccountsCreated: DefaultOriginAccountsCreated,
		ExpiryBlock:           DefaultExpiryBlock,
		DelegationRequirement: &DefaultDelegationRequirement,
		AirdropFactor:         &DefaultAirdropFactor,
	}
}

func (p Params) Validate() error {
	if p.ExpiryBlock == 0 {
		return fmt.Errorf("expiry block cannot be 0")
	}
	if p.DelegationRequirement == nil {
		return fmt.Errorf("delegation requirement cannot be nil")
	}
	if p.AirdropFactor == nil {
		return fmt.Errorf("airdrop factor cannot be nil")
	}
	return nil
}
