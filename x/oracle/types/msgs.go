package types

import (
	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/ojo-network/ojo/util/checkers"
	"gopkg.in/yaml.v3"
)

var (
	_ sdk.Msg = &MsgDelegateFeedConsent{}
	_ sdk.Msg = &MsgAggregateExchangeRatePrevote{}
	_ sdk.Msg = &MsgAggregateExchangeRateVote{}
	_ sdk.Msg = &MsgGovUpdateParams{}
	_ sdk.Msg = &MsgGovCancelUpdateParamPlan{}
)

func NewMsgAggregateExchangeRatePrevote(
	hash AggregateVoteHash,
	feeder sdk.AccAddress,
	validator sdk.ValAddress,
) *MsgAggregateExchangeRatePrevote {
	return &MsgAggregateExchangeRatePrevote{
		Hash:      hash.String(),
		Feeder:    feeder.String(),
		Validator: validator.String(),
	}
}

// Type implements LegacyMsg interface
func (msg MsgAggregateExchangeRatePrevote) Type() string { return sdk.MsgTypeURL(&msg) }

// GetSignBytes implements sdk.Msg
func (msg MsgAggregateExchangeRatePrevote) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgAggregateExchangeRatePrevote) GetSigners() []sdk.AccAddress {
	return checkers.Signers(msg.Feeder)
}

// ValidateBasic Implements sdk.Msg
func (msg MsgAggregateExchangeRatePrevote) ValidateBasic() error {
	_, err := AggregateVoteHashFromHexString(msg.Hash)
	if err != nil {
		return ErrInvalidHash.Wrapf("invalid vote hash (%s)", err)
	}

	// HEX encoding doubles the hash length
	if len(msg.Hash) != tmhash.TruncatedSize*2 {
		return ErrInvalidHashLength
	}

	_, err = sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid feeder address (%s)", err)
	}

	_, err = sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address (%s)", err)
	}

	return nil
}

func NewMsgAggregateExchangeRateVote(
	salt string,
	exchangeRates string,
	feeder sdk.AccAddress,
	validator sdk.ValAddress,
) *MsgAggregateExchangeRateVote {
	return &MsgAggregateExchangeRateVote{
		Salt:          salt,
		ExchangeRates: exchangeRates,
		Feeder:        feeder.String(),
		Validator:     validator.String(),
	}
}

// Type implements LegacyMsg interface
func (msg MsgAggregateExchangeRateVote) Type() string { return sdk.MsgTypeURL(&msg) }

// GetSignBytes implements sdk.Msg
func (msg MsgAggregateExchangeRateVote) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgAggregateExchangeRateVote) GetSigners() []sdk.AccAddress {
	return checkers.Signers(msg.Feeder)
}

// ValidateBasic implements sdk.Msg
func (msg MsgAggregateExchangeRateVote) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid feeder address (%s)", err)
	}

	_, err = sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address (%s)", err)
	}

	if l := len(msg.ExchangeRates); l == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("must provide at least one oracle exchange rate")
	} else if l > 4096 {
		return sdkerrors.ErrInvalidRequest.Wrap("exchange rates string can not exceed 4096 characters")
	}

	exchangeRates, err := ParseExchangeRateDecCoins(msg.ExchangeRates)
	if err != nil {
		return sdkerrors.ErrInvalidCoins.Wrap("failed to parse exchange rates string cause: " + err.Error())
	}

	for _, exchangeRate := range exchangeRates {
		// check overflow bit length
		if exchangeRate.Amount.BigInt().BitLen() > 255+sdk.DecimalPrecisionBits {
			return ErrInvalidExchangeRate.Wrap("overflow")
		}
	}

	if len(msg.Salt) != 64 {
		return ErrInvalidSaltLength
	}
	_, err = AggregateVoteHashFromHexString(msg.Salt)
	if err != nil {
		return ErrInvalidSaltFormat.Wrap("salt must be a valid hex string")
	}

	return nil
}

// NewMsgDelegateFeedConsent creates a MsgDelegateFeedConsent instance
func NewMsgDelegateFeedConsent(operatorAddress sdk.ValAddress, feederAddress sdk.AccAddress) *MsgDelegateFeedConsent {
	return &MsgDelegateFeedConsent{
		Operator: operatorAddress.String(),
		Delegate: feederAddress.String(),
	}
}

// Type implements LegacyMsg interface
func (msg MsgDelegateFeedConsent) Type() string { return sdk.MsgTypeURL(&msg) }

// GetSignBytes implements sdk.Msg
func (msg MsgDelegateFeedConsent) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgDelegateFeedConsent) GetSigners() []sdk.AccAddress {
	operator, _ := sdk.ValAddressFromBech32(msg.Operator)
	return []sdk.AccAddress{sdk.AccAddress(operator)}
}

// ValidateBasic implements sdk.Msg
func (msg MsgDelegateFeedConsent) ValidateBasic() error {
	_, err := sdk.ValAddressFromBech32(msg.Operator)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Delegate)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegate address (%s)", err)
	}

	return nil
}

// NewMsgUpdateParams will creates a new MsgUpdateParams instance
func NewMsgUpdateParams(authority, title, description string, plan ParamUpdatePlan) *MsgGovUpdateParams {
	return &MsgGovUpdateParams{
		Authority:   authority,
		Title:       title,
		Description: description,
		Plan:        plan,
	}
}

// Type implements Msg interface
func (msg MsgGovUpdateParams) Type() string { return sdk.MsgTypeURL(&msg) }

// String implements the Stringer interface.
func (msg MsgGovUpdateParams) String() string {
	out, _ := yaml.Marshal(msg)
	return string(out)
}

// GetSignBytes implements Msg
func (msg MsgGovUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgGovUpdateParams) GetSigners() []sdk.AccAddress {
	return checkers.Signers(msg.Authority)
}

// ValidateBasic implements Msg and validates params for each param key
// specified in the proposal. If one param is invalid, the whole proposal
// will fail to go through.
func (msg MsgGovUpdateParams) ValidateBasic() error {
	if err := checkers.ValidateProposal(msg.Title, msg.Description, msg.Authority); err != nil {
		return err
	}

	return msg.Plan.ValidateBasic()
}

// NewMsgGovCancelUpdateParamPlan will creates a new MsgGovCancelUpdateParamPlan instance
func NewMsgGovCancelUpdateParamPlan(authority, title, description string) *MsgGovCancelUpdateParamPlan {
	return &MsgGovCancelUpdateParamPlan{
		Authority:   authority,
		Title:       title,
		Description: description,
	}
}

// Type implements Msg interface
func (msg MsgGovCancelUpdateParamPlan) Type() string { return sdk.MsgTypeURL(&msg) }

// String implements the Stringer interface.
func (msg MsgGovCancelUpdateParamPlan) String() string {
	out, _ := yaml.Marshal(msg)
	return string(out)
}

// GetSignBytes implements Msg
func (msg MsgGovCancelUpdateParamPlan) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgGovCancelUpdateParamPlan) GetSigners() []sdk.AccAddress {
	return checkers.Signers(msg.Authority)
}

// ValidateBasic implements Msg
func (msg MsgGovCancelUpdateParamPlan) ValidateBasic() error {
	return checkers.ValidateProposal(msg.Title, msg.Description, msg.Authority)
}

// NewMsgGovAddDenoms will creates a new NewMsgGovAddDenoms instance
func NewMsgGovAddDenoms(
	authority, title, description string,
	height int64, denoms DenomList, mandatory bool,
	rewardBand sdk.Dec, currencyPairProviders CurrencyPairProvidersList,
	currencyDeviationThresholds CurrencyDeviationThresholdList,
) *MsgGovAddDenoms {
	return &MsgGovAddDenoms{
		Authority:   authority,
		Title:       title,
		Description: description,
		Height:      height,
		DenomList:   denoms,
		Mandatory:   mandatory,
		RewardBand:  &rewardBand,

		CurrencyPairProviders:       currencyPairProviders,
		CurrencyDeviationThresholds: currencyDeviationThresholds,
	}
}

// Type implements Msg interface
func (msg MsgGovAddDenoms) Type() string { return sdk.MsgTypeURL(&msg) }

// String implements the Stringer interface.
func (msg MsgGovAddDenoms) String() string {
	out, _ := yaml.Marshal(msg)
	return string(out)
}

// GetSignBytes implements Msg
func (msg MsgGovAddDenoms) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgGovAddDenoms) GetSigners() []sdk.AccAddress {
	return checkers.Signers(msg.Authority)
}

// ValidateBasic implements Msg
func (msg MsgGovAddDenoms) ValidateBasic() error {
	return checkers.ValidateProposal(msg.Title, msg.Description, msg.Authority)
}

// NewMsgGovRemoveCurrencyPairProviders will creates a new NewMsgGovRemoveCurrencyPairProviders instance
func NewMsgGovRemoveCurrencyPairProviders(
	authority, title, description string,
	height int64, currencyPairProviders CurrencyPairProvidersList,
) *MsgGovRemoveCurrencyPairProviders {
	return &MsgGovRemoveCurrencyPairProviders{
		Authority:   authority,
		Title:       title,
		Description: description,
		Height:      height,

		CurrencyPairProviders: currencyPairProviders,
	}
}

// Type implements Msg interface
func (msg MsgGovRemoveCurrencyPairProviders) Type() string { return sdk.MsgTypeURL(&msg) }

// String implements the Stringer interface.
func (msg MsgGovRemoveCurrencyPairProviders) String() string {
	out, _ := yaml.Marshal(msg)
	return string(out)
}

// GetSignBytes implements Msg
func (msg MsgGovRemoveCurrencyPairProviders) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgGovRemoveCurrencyPairProviders) GetSigners() []sdk.AccAddress {
	return checkers.Signers(msg.Authority)
}

// ValidateBasic implements Msg
func (msg MsgGovRemoveCurrencyPairProviders) ValidateBasic() error {
	return checkers.ValidateProposal(msg.Title, msg.Description, msg.Authority)
}

// NewMsgGovRemoveCurrencyDeviationThresholds will creates a new NewMsgGovRemoveCurrencyDeviationThresholds instance
func NewMsgGovRemoveCurrencyDeviationThresholds(
	authority, title, description string,
	height int64, currencies []string,
) *MsgGovRemoveCurrencyDeviationThresholds {
	return &MsgGovRemoveCurrencyDeviationThresholds{
		Authority:   authority,
		Title:       title,
		Description: description,
		Height:      height,
		Currencies:  currencies,
	}
}

// Type implements Msg interface
func (msg MsgGovRemoveCurrencyDeviationThresholds) Type() string { return sdk.MsgTypeURL(&msg) }

// String implements the Stringer interface.
func (msg MsgGovRemoveCurrencyDeviationThresholds) String() string {
	out, _ := yaml.Marshal(msg)
	return string(out)
}

// GetSignBytes implements Msg
func (msg MsgGovRemoveCurrencyDeviationThresholds) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgGovRemoveCurrencyDeviationThresholds) GetSigners() []sdk.AccAddress {
	return checkers.Signers(msg.Authority)
}

// ValidateBasic implements Msg
func (msg MsgGovRemoveCurrencyDeviationThresholds) ValidateBasic() error {
	return checkers.ValidateProposal(msg.Title, msg.Description, msg.Authority)
}
