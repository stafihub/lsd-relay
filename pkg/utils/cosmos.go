package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
)

func GetStakingParams(endpoint string) (*StakingParamsRes, error) {
	url := fmt.Sprintf("%s/cosmos/staking/v1beta1/params", endpoint)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("res status: %d", res.StatusCode)
	}
	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	sp := StakingParamsRes{}
	err = json.Unmarshal(bts, &sp)
	if err != nil {
		return nil, err
	}
	return &sp, nil
}

type StakingParamsRes struct {
	Params struct {
		UnbondingTime             string `json:"unbonding_time"`
		MaxValidators             int    `json:"max_validators"`
		MaxEntries                int    `json:"max_entries"`
		HistoricalEntries         int    `json:"historical_entries"`
		BondDenom                 string `json:"bond_denom"`
		MinCommissionRate         string `json:"min_commission_rate"`
		ValidatorBondFactor       string `json:"validator_bond_factor"`
		GlobalLiquidStakingCap    string `json:"global_liquid_staking_cap"`
		ValidatorLiquidStakingCap string `json:"validator_liquid_staking_cap"`
	} `json:"params"`
}

func GetTotalLiquidStake(endpoint string) (*TotalLiquidStakeRes, error) {
	url := fmt.Sprintf("%s/cosmos/staking/v1beta1/total_liquid_staked", endpoint)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("res status: %d", res.StatusCode)
	}
	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	sp := TotalLiquidStakeRes{}
	err = json.Unmarshal(bts, &sp)
	if err != nil {
		return nil, err
	}

	return &sp, nil
}

type TotalLiquidStakeRes struct {
	Tokens string `json:"tokens"`
}

func GetValidator(endpoint, val string) (*ValidatorRes, error) {
	url := fmt.Sprintf("%s/cosmos/staking/v1beta1/validators/%s", endpoint, val)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("res status: %d", res.StatusCode)
	}
	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	sp := ValidatorRes{}
	err = json.Unmarshal(bts, &sp)
	if err != nil {
		return nil, err
	}

	return &sp, nil
}

type ValidatorRes struct {
	Validator struct {
		OperatorAddress string `json:"operator_address"`
		ConsensusPubkey struct {
			Type string `json:"@type"`
			Key  string `json:"key"`
		} `json:"consensus_pubkey"`
		Jailed          bool   `json:"jailed"`
		Status          string `json:"status"`
		Tokens          string `json:"tokens"`
		DelegatorShares string `json:"delegator_shares"`
		Description     struct {
			Moniker         string `json:"moniker"`
			Identity        string `json:"identity"`
			Website         string `json:"website"`
			SecurityContact string `json:"security_contact"`
			Details         string `json:"details"`
		} `json:"description"`
		UnbondingHeight string    `json:"unbonding_height"`
		UnbondingTime   time.Time `json:"unbonding_time"`
		Commission      struct {
			CommissionRates struct {
				Rate          string `json:"rate"`
				MaxRate       string `json:"max_rate"`
				MaxChangeRate string `json:"max_change_rate"`
			} `json:"commission_rates"`
			UpdateTime time.Time `json:"update_time"`
		} `json:"commission"`
		MinSelfDelegation       string        `json:"min_self_delegation"`
		UnbondingOnHoldRefCount string        `json:"unbonding_on_hold_ref_count"`
		UnbondingIds            []interface{} `json:"unbonding_ids"`
		ValidatorBondShares     string        `json:"validator_bond_shares"`
		LiquidShares            string        `json:"liquid_shares"`
	} `json:"validator"`
}

func GetPool(endpoint string) (*PoolRes, error) {
	url := fmt.Sprintf("%s/cosmos/staking/v1beta1/pool", endpoint)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("res status: %d", res.StatusCode)
	}
	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	sp := PoolRes{}
	err = json.Unmarshal(bts, &sp)
	if err != nil {
		return nil, err
	}

	return &sp, nil
}

type PoolRes struct {
	Pool struct {
		NotBondedTokens string `json:"not_bonded_tokens"`
		BondedTokens    string `json:"bonded_tokens"`
	} `json:"pool"`
}

var ValidatorBondCapDisabled = types.NewDecFromInt(types.NewInt(-1))

func SelectVals(validatorAddrs []string, poolBond, poolUnbond, endpoint string) ([]string, error) {
	if len(endpoint) == 0 {
		return validatorAddrs, nil
	}
	bond, ok := types.NewIntFromString(poolBond)
	if !ok {
		return nil, fmt.Errorf("parse poolInfo.Bond %s failed", poolBond)
	}
	unbond, ok := types.NewIntFromString(poolUnbond)
	if !ok {
		return nil, fmt.Errorf("parse poolInfo.Unbond %s failed", poolUnbond)
	}

	totalShouldBondAmount := bond.Sub(unbond)

	stakingParams, err := GetStakingParams(endpoint)
	if err != nil {
		return nil, err
	}
	globalLiquidStakingCap, err := types.NewDecFromStr(stakingParams.Params.GlobalLiquidStakingCap)
	if err != nil {
		return nil, fmt.Errorf("parse GlobalLiquidStakingCap %s failed", stakingParams.Params.GlobalLiquidStakingCap)
	}
	validatorLiquidStakingCap, err := types.NewDecFromStr(stakingParams.Params.ValidatorLiquidStakingCap)
	if err != nil {
		return nil, fmt.Errorf("parse ValidatorLiquidStakingCap %s failed", stakingParams.Params.ValidatorLiquidStakingCap)
	}
	validatorBondFactor, err := types.NewDecFromStr(stakingParams.Params.ValidatorBondFactor)
	if err != nil {
		return nil, fmt.Errorf("parse ValidatorBondFactor %s failed", stakingParams.Params.ValidatorBondFactor)
	}

	totalLiquidStake, err := GetTotalLiquidStake(endpoint)
	if err != nil {
		return nil, err
	}
	totalLiquidStakedAmount, ok := types.NewIntFromString(totalLiquidStake.Tokens)
	if !ok {
		return nil, fmt.Errorf("parse totalLiquidStake.Tokens %s failed", totalLiquidStake.Tokens)
	}

	poolRes, err := GetPool(endpoint)
	if err != nil {
		return nil, err
	}
	poolBondedTokens, ok := types.NewIntFromString(poolRes.Pool.BondedTokens)
	if !ok {
		return nil, fmt.Errorf("parse Pool.BondedTokens %s failed", poolRes.Pool.BondedTokens)
	}

	totalStakedAmount := poolBondedTokens.Add(totalShouldBondAmount)
	totalLiquidStakedAmount = totalLiquidStakedAmount.Add(totalShouldBondAmount)

	// 0 check global liquid staking cap
	liquidStakePercent := types.NewDecFromInt(totalLiquidStakedAmount).Quo(types.NewDecFromInt(totalStakedAmount))
	if liquidStakePercent.GT(globalLiquidStakingCap) {
		return nil, fmt.Errorf("ExceedsGlobalLiquidStakingCap %s", liquidStakePercent.String())
	}

	valAddrs := make([]string, 0, 10)
	for _, val := range validatorAddrs {
		validatorInfo, err := GetValidator(endpoint, val)
		if err != nil {
			return nil, err
		}
		delegatorShares, err := types.NewDecFromStr(validatorInfo.Validator.DelegatorShares)
		if err != nil {
			return nil, fmt.Errorf("parse DelegatorShares %s failed", validatorInfo.Validator.DelegatorShares)
		}
		validatorBondShares, err := types.NewDecFromStr(validatorInfo.Validator.ValidatorBondShares)
		if err != nil {
			return nil, fmt.Errorf("parse ValidatorBondShares %s failed", validatorInfo.Validator.ValidatorBondShares)
		}
		liquidShares, err := types.NewDecFromStr(validatorInfo.Validator.LiquidShares)
		if err != nil {
			return nil, fmt.Errorf("parse LiquidShares %s failed", validatorInfo.Validator.LiquidShares)
		}

		valTokens, ok := types.NewIntFromString(validatorInfo.Validator.Tokens)
		if !ok {
			return nil, fmt.Errorf("parse Validator.Tokens %s failed", validatorInfo.Validator.Tokens)
		}

		if valTokens.IsZero() {
			return nil, fmt.Errorf("valTokens is zero")
		}
		shares := delegatorShares.MulInt(totalShouldBondAmount).QuoInt(valTokens)

		// maxValLiquidSharesLog := validatorBondShares.Mul(validatorBondFactor)
		// logrus.Debugf("val: %s maxValLiquidShares: %s", val, maxValLiquidSharesLog.String())

		// 1 check val bond cap
		if !validatorBondFactor.Equal(ValidatorBondCapDisabled) {
			maxValLiquidShares := validatorBondShares.Mul(validatorBondFactor)
			if liquidShares.Add(shares).GT(maxValLiquidShares) {
				continue
			}
		}

		// 2 check val liquid staking cap
		updatedLiquidShares := liquidShares.Add(shares)
		updatedTotalShares := delegatorShares.Add(shares)
		liquidStakePercent := updatedLiquidShares.Quo(updatedTotalShares)
		if liquidStakePercent.GT(validatorLiquidStakingCap) {
			continue
		}

		valAddrs = append(valAddrs, val)
	}

	return valAddrs, nil
}
