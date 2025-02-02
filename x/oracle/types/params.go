package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v3"
)

var (
	oneDec           = sdk.OneDec()
	minVoteThreshold = sdk.NewDecWithPrec(33, 2) // 0.33
)

// maxium number of decimals allowed for VoteThreshold
const (
	MaxVoteThresholdPrecision  = 2
	MaxVoteThresholdMultiplier = 100 // must be 10^MaxVoteThresholdPrecision
)

// Parameter keys
var (
	KeyVotePeriod                  = []byte("VotePeriod")
	KeyVoteThreshold               = []byte("VoteThreshold")
	KeyRewardBands                 = []byte("RewardBands")
	KeyRewardDistributionWindow    = []byte("RewardDistributionWindow")
	KeyAcceptList                  = []byte("AcceptList")
	KeyMandatoryList               = []byte("MandatoryList")
	KeySlashFraction               = []byte("SlashFraction")
	KeySlashWindow                 = []byte("SlashWindow")
	KeyMinValidPerWindow           = []byte("MinValidPerWindow")
	KeyHistoricStampPeriod         = []byte("HistoricStampPeriod")
	KeyMedianStampPeriod           = []byte("MedianStampPeriod")
	KeyMaximumPriceStamps          = []byte("MaximumPriceStamps")
	KeyMaximumMedianStamps         = []byte("MaximumMedianStamps")
	KeyCurrencyPairProviders       = []byte("CurrencyPairProviders")
	KeyCurrencyDeviationThresholds = []byte("CurrencyDeviationThresholds")
)

// Default parameter values
const (
	DefaultVotePeriod               = BlocksPerMinute * 3 / 10 // 18 seconds
	DefaultSlashWindow              = BlocksPerWeek            // window for a week
	DefaultRewardDistributionWindow = BlocksPerYear            // window for a year
	DefaultHistoricStampPeriod      = BlocksPerMinute * 3      // window for 3 minutes
	DefaultMaximumPriceStamps       = 60                       // retain for 3 hours
	DefaultMedianStampPeriod        = BlocksPerHour * 3        // window for 3 hours
	DefaultMaximumMedianStamps      = 24                       // retain for 3 days
)

// Default parameter values
var (
	DefaultVoteThreshold = sdk.NewDecWithPrec(50, 2) // 50%

	DefaultAcceptList = DenomList{
		{
			BaseDenom:   OjoDenom,
			SymbolDenom: OjoSymbol,
			Exponent:    OjoExponent,
		},
		{
			BaseDenom:   AtomDenom,
			SymbolDenom: AtomSymbol,
			Exponent:    AtomExponent,
		},
	}
	DefaultMandatoryList = DenomList{
		{
			BaseDenom:   AtomDenom,
			SymbolDenom: AtomSymbol,
			Exponent:    AtomExponent,
		},
	}
	DefaultSlashFraction     = sdk.NewDecWithPrec(1, 4) // 0.01%
	DefaultMinValidPerWindow = sdk.NewDecWithPrec(5, 2) // 5%
	defaultRewardBand        = sdk.NewDecWithPrec(2, 2) // 0.02

	DefaultCurrencyPairProviders = CurrencyPairProvidersList{
		CurrencyPairProviders{
			BaseDenom:  OjoSymbol,
			QuoteDenom: USDDenom,
			Providers: []string{
				"binance",
				"coinbase",
			},
		},
	}

	DefaultCurrencyDeviationThresholds = CurrencyDeviationThresholdList{
		CurrencyDeviationThreshold{
			BaseDenom: OjoSymbol,
			Threshold: "2",
		},
	}
)

var _ paramstypes.ParamSet = &Params{}

// DefaultRewardBands returns a new default RewardBandList object.
//
// This function is necessary because we cannot use a constant,
// and the reward band list is manipulated in our unit tests.
func DefaultRewardBands() RewardBandList {
	return RewardBandList{
		{
			SymbolDenom: OjoSymbol,
			RewardBand:  defaultRewardBand,
		},
		{
			SymbolDenom: AtomSymbol,
			RewardBand:  defaultRewardBand,
		},
	}
}

// AddDefault adds a default reward band for the given
// denom.
func (rbl *RewardBandList) AddDefault(
	denom string,
) {
	*rbl = append(*rbl, RewardBand{
		SymbolDenom: denom,
		RewardBand:  defaultRewardBand,
	})
}

// Add adds a reward band of a given denom and
// reward band decimal.
func (rbl *RewardBandList) Add(
	denom string,
	band math.LegacyDec,
) {
	*rbl = append(*rbl, RewardBand{
		SymbolDenom: denom,
		RewardBand:  band,
	})
}

// DefaultParams creates default oracle module parameters
func DefaultParams() Params {
	return Params{
		VotePeriod:                  DefaultVotePeriod,
		VoteThreshold:               DefaultVoteThreshold,
		RewardDistributionWindow:    DefaultRewardDistributionWindow,
		AcceptList:                  DefaultAcceptList,
		MandatoryList:               DefaultMandatoryList,
		SlashFraction:               DefaultSlashFraction,
		SlashWindow:                 DefaultSlashWindow,
		MinValidPerWindow:           DefaultMinValidPerWindow,
		HistoricStampPeriod:         DefaultHistoricStampPeriod,
		MedianStampPeriod:           DefaultMedianStampPeriod,
		MaximumPriceStamps:          DefaultMaximumPriceStamps,
		MaximumMedianStamps:         DefaultMaximumMedianStamps,
		RewardBands:                 DefaultRewardBands(),
		CurrencyPairProviders:       DefaultCurrencyPairProviders,
		CurrencyDeviationThresholds: DefaultCurrencyDeviationThresholds,
	}
}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of oracle module's parameters.
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(
			KeyVotePeriod,
			&p.VotePeriod,
			validateVotePeriod,
		),
		paramstypes.NewParamSetPair(
			KeyVoteThreshold,
			&p.VoteThreshold,
			validateVoteThreshold,
		),
		paramstypes.NewParamSetPair(
			KeyRewardBands,
			&p.RewardBands,
			validateRewardBands,
		),
		paramstypes.NewParamSetPair(
			KeyRewardDistributionWindow,
			&p.RewardDistributionWindow,
			validateRewardDistributionWindow,
		),
		paramstypes.NewParamSetPair(
			KeyAcceptList,
			&p.AcceptList,
			validateDenomList,
		),
		paramstypes.NewParamSetPair(
			KeyMandatoryList,
			&p.MandatoryList,
			validateDenomList,
		),
		paramstypes.NewParamSetPair(
			KeySlashFraction,
			&p.SlashFraction,
			validateSlashFraction,
		),
		paramstypes.NewParamSetPair(
			KeySlashWindow,
			&p.SlashWindow,
			validateSlashWindow,
		),
		paramstypes.NewParamSetPair(
			KeyMinValidPerWindow,
			&p.MinValidPerWindow,
			validateMinValidPerWindow,
		),
		paramstypes.NewParamSetPair(
			KeyHistoricStampPeriod,
			&p.HistoricStampPeriod,
			validateHistoricStampPeriod,
		),
		paramstypes.NewParamSetPair(
			KeyMedianStampPeriod,
			&p.MedianStampPeriod,
			validateMedianStampPeriod,
		),
		paramstypes.NewParamSetPair(
			KeyMaximumPriceStamps,
			&p.MaximumPriceStamps,
			validateMaximumPriceStamps,
		),
		paramstypes.NewParamSetPair(
			KeyMaximumMedianStamps,
			&p.MaximumMedianStamps,
			validateMaximumMedianStamps,
		),
		paramstypes.NewParamSetPair(
			KeyCurrencyPairProviders,
			&p.CurrencyPairProviders,
			validateCurrencyPairProviders,
		),
		paramstypes.NewParamSetPair(
			KeyCurrencyDeviationThresholds,
			&p.CurrencyDeviationThresholds,
			validateCurrencyDeviationThresholds,
		),
	}
}

// String implements fmt.Stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate performs basic validation on oracle parameters.
func (p Params) Validate() error {
	if p.VotePeriod == 0 {
		return ErrInvalidParamValue.Wrap("oracle parameter VotePeriod must be > 0")
	}
	if p.VoteThreshold.LTE(sdk.NewDecWithPrec(33, 2)) {
		return ErrInvalidParamValue.Wrap("oracle parameter VoteThreshold must be greater than 33 percent")
	}

	if p.RewardDistributionWindow < p.VotePeriod {
		return ErrInvalidParamValue.Wrap(
			"oracle parameter RewardDistributionWindow must be greater than or equal with VotePeriod",
		)
	}

	if p.SlashFraction.GT(sdk.OneDec()) || p.SlashFraction.IsNegative() {
		return ErrInvalidParamValue.Wrap("oracle parameter SlashFraction must be between [0, 1]")
	}

	if p.SlashWindow < p.VotePeriod {
		return ErrInvalidParamValue.Wrap("oracle parameter SlashWindow must be greater than or equal with VotePeriod")
	}

	if p.SlashWindow%p.VotePeriod != 0 {
		return ErrInvalidParamValue.Wrap("oracle parameter SlashWindow must be an exact multiple of VotePeriod")
	}

	if p.MinValidPerWindow.GT(sdk.OneDec()) || p.MinValidPerWindow.IsNegative() {
		return ErrInvalidParamValue.Wrap("oracle parameter MinValidPerWindow must be between [0, 1]")
	}

	if err := validateDenomList(p.AcceptList); err != nil {
		return err
	}

	if err := validateDenomList(p.MandatoryList); err != nil {
		return err
	}

	if p.HistoricStampPeriod > p.MedianStampPeriod {
		return ErrInvalidParamValue.Wrap(
			"oracle parameter MedianStampPeriod must be greater than or equal with HistoricStampPeriod",
		)
	}

	if p.HistoricStampPeriod%p.VotePeriod != 0 || p.MedianStampPeriod%p.VotePeriod != 0 {
		return ErrInvalidParamValue.Wrap(
			"oracle parameters HistoricStampPeriod and MedianStampPeriod must be exact multiples of VotePeriod",
		)
	}

	if err := validateRewardBands(p.RewardBands); err != nil {
		return err
	}

	// all denoms in mandatory list must be in accept list
	if !p.AcceptList.ContainDenoms(p.MandatoryList) {
		return ErrInvalidParamValue.Wrap("denom in MandatoryList not present in AcceptList")
	}

	return nil
}

func validateVotePeriod(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return ErrInvalidParamValue.Wrapf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return ErrInvalidParamValue.Wrap("oracle parameter VotePeriod must be > 0")
	}

	return nil
}

func validateVoteThreshold(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return ErrInvalidParamValue.Wrapf("invalid parameter type: %T", i)
	}

	if v.LTE(minVoteThreshold) || v.GT(oneDec) {
		return ErrInvalidParamValue.Wrapf("threshold must be bigger than %s and <= 1", minVoteThreshold)
	}
	val := v.MulInt64(100).TruncateInt64()
	v2 := sdk.NewDecWithPrec(val, MaxVoteThresholdPrecision)
	if !v2.Equal(v) {
		return ErrInvalidParamValue.Wrap("threshold precision must be maximum 2 decimals")
	}
	return nil
}

func validateRewardBand(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return ErrInvalidParamValue.Wrapf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return ErrInvalidParamValue.Wrap("oracle parameter RewardBand must be between [0, 1]")
	}

	if v.GT(sdk.OneDec()) {
		return ErrInvalidParamValue.Wrap("oracle parameter RewardBand must be between [0, 1]")
	}

	return nil
}

func validateRewardDistributionWindow(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return ErrInvalidParamValue.Wrapf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return ErrInvalidParamValue.Wrap("oracle parameter RewardDistributionWindow must be > 0")
	}

	return nil
}

func validateDenomList(i interface{}) error {
	v, ok := i.(DenomList)
	if !ok {
		return ErrInvalidParamValue.Wrapf("invalid parameter type: %T", i)
	}

	for _, d := range v {
		if len(d.BaseDenom) == 0 {
			return ErrInvalidParamValue.Wrap("oracle parameter AcceptList Denom must have BaseDenom")
		}
		if len(d.SymbolDenom) == 0 {
			return ErrInvalidParamValue.Wrap("oracle parameter AcceptList Denom must have SymbolDenom")
		}
	}

	return nil
}

func validateRewardBands(i interface{}) error {
	v, ok := i.(RewardBandList)
	if !ok {
		return ErrInvalidParamValue.Wrapf("invalid parameter type: %T", i)
	}

	for _, d := range v {
		if err := validateRewardBand(d.RewardBand); err != nil {
			return err
		}
		if len(d.SymbolDenom) == 0 {
			return ErrInvalidParamValue.Wrap("oracle parameter RewardBand must have SymbolDenom")
		}
	}

	return nil
}

func validateSlashFraction(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return ErrInvalidParamValue.Wrapf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return ErrInvalidParamValue.Wrap("oracle parameter SlashFraction must be between [0, 1]")
	}

	if v.GT(sdk.OneDec()) {
		return ErrInvalidParamValue.Wrap("oracle parameter SlashFraction must be between [0, 1]")
	}

	return nil
}

func validateSlashWindow(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return ErrInvalidParamValue.Wrapf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return ErrInvalidParamValue.Wrap("oracle parameter SlashWindow must be > 0")
	}

	return nil
}

func validateMinValidPerWindow(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return ErrInvalidParamValue.Wrapf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return ErrInvalidParamValue.Wrap("oracle parameter MinValidPerWindow must be between [0, 1]")
	}

	if v.GT(sdk.OneDec()) {
		return ErrInvalidParamValue.Wrap("oracle parameter MinValidPerWindow must be between [0, 1]")
	}

	return nil
}

func validateHistoricStampPeriod(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return ErrInvalidParamValue.Wrapf("invalid parameter type: %T", i)
	}

	if v < 1 {
		return ErrInvalidParamValue.Wrap("oracle parameter HistoricStampPeriod must be > 0")
	}

	return nil
}

func validateMedianStampPeriod(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return ErrInvalidParamValue.Wrapf("invalid parameter type: %T", i)
	}

	if v < 1 {
		return ErrInvalidParamValue.Wrap("oracle parameter MedianStampPeriod must be > 0")
	}

	return nil
}

func validateMaximumPriceStamps(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return ErrInvalidParamValue.Wrapf("invalid parameter type: %T", i)
	}

	if v < 1 {
		return ErrInvalidParamValue.Wrap("oracle parameter MaximumPriceStamps must be > 0")
	}

	return nil
}

func validateMaximumMedianStamps(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return ErrInvalidParamValue.Wrapf("invalid parameter type: %T", i)
	}

	if v < 1 {
		return ErrInvalidParamValue.Wrap("oracle parameter MaximumMedianStamps must be > 0")
	}

	return nil
}

func validateCurrencyPairProviders(i interface{}) error {
	v, ok := i.(CurrencyPairProvidersList)
	if !ok {
		return ErrInvalidParamValue.Wrapf("invalid parameter type: %T", i)
	}

	for _, c := range v {
		if len(c.BaseDenom) == 0 {
			return ErrInvalidParamValue.Wrap("oracle parameter CurrencyPairProviders must have BaseDenom")
		}
		if len(c.QuoteDenom) == 0 {
			return ErrInvalidParamValue.Wrap("oracle parameter CurrencyPairProviders must have QuoteDenom")
		}
		if len(c.Providers) < 1 {
			return ErrInvalidParamValue.Wrap("oracle parameter CurrencyPairProviders must have at least 1 provider listed")
		}
	}

	return nil
}

func validateCurrencyDeviationThresholds(i interface{}) error {
	v, ok := i.(CurrencyDeviationThresholdList)
	if !ok {
		return ErrInvalidParamValue.Wrapf("invalid parameter type: %T", i)
	}

	for _, c := range v {
		if len(c.BaseDenom) == 0 {
			return ErrInvalidParamValue.Wrap("oracle parameter CurrencyDeviationThreshold must have BaseDenom")
		}
		if len(c.Threshold) == 0 {
			return ErrInvalidParamValue.Wrap("oracle parameter CurrencyDeviationThreshold must have Threshold")
		}
	}

	return nil
}
