package task

import (
	"encoding/json"
)

// EraProcessStatus
const (
	EraUpdateStarted = "EraUpdateStarted"
	EraUpdateEnded   = "EraUpdateEnded"
	BondStarted      = "BondStarted"
	BondEnded        = "BondEnded"
	WithdrawStarted  = "WithdrawStarted"
	WithdrawEnded    = "WithdrawEnded"
	RestakeStarted   = "RestakeStarted"
	RestakeEnded     = "RestakeEnded"
	ActiveEnded      = "ActiveEnded"
)

// ValidatorUpdateStatus
const (
	WaitQueryUpdate = "WaitQueryUpdate"
)

type PoolAddr struct {
	Addr string `json:"pool_addr"`
}

type QueryPoolInfoReq struct {
	PoolInfo PoolAddr `json:"pool_info"`
}

func getQueryPoolInfoReq(poolAddr string) []byte {
	poolReq := QueryPoolInfoReq{
		PoolInfo: PoolAddr{
			Addr: poolAddr,
		},
	}
	marshal, _ := json.Marshal(poolReq)
	return marshal
}

type QueryPoolInfoRes struct {
	Bond                      string      `json:"bond"`
	Unbond                    string      `json:"unbond"`
	Active                    string      `json:"active"`
	LsdToken                  string      `json:"lsd_token"`
	IcaId                     string      `json:"ica_id"`
	IbcDenom                  string      `json:"ibc_denom"`
	ChannelIdOfIbcDenom       string      `json:"channel_id_of_ibc_denom"`
	RemoteDenom               string      `json:"remote_denom"`
	ValidatorAddrs            []string    `json:"validator_addrs"`
	Era                       uint64      `json:"era"`
	Rate                      string      `json:"rate"`
	EraSeconds                uint64      `json:"era_seconds"`
	Offset                    uint64      `json:"offset"`
	MinimalStake              string      `json:"minimal_stake"`
	UnstakeTimesLimit         uint64      `json:"unstake_times_limit"`
	NextUnstakeIndex          uint64      `json:"next_unstake_index"`
	UnbondingPeriod           uint64      `json:"unbonding_period"`
	EraProcessStatus          string      `json:"era_process_status"`
	ValidatorUpdateStatus     string      `json:"validator_update_status"`
	UnbondCommission          string      `json:"unbond_commission"`
	PlatformFeeCommission     string      `json:"platform_fee_commission"`
	TotalPlatformFee          string      `json:"total_platform_fee"`
	PlatformFeeReceiver       string      `json:"platform_fee_receiver"`
	Admin                     string      `json:"admin"`
	ShareTokens               []Coin      `json:"share_tokens"`
	RedeemmingShareTokenDenom []string    `json:"redeemming_share_token_denom"`
	EraSnapshot               eraSnapshot `json:"era_snapshot"`
	Paused                    bool        `json:"paused"`
	LsmSupport                bool        `json:"lsm_support"`
	LsmPendingLimit           uint64      `json:"lsm_pending_limit"`
	RateChangeLimit           string      `json:"rate_change_limit"`
}

type eraSnapshot struct {
	Era           uint64 `json:"era"`
	Bond          string `json:"bond"`
	Unbond        string `json:"unbond"`
	Active        string `json:"active"`
	RestakeAmount string `json:"restake_amount"`
	BondHeight    uint64 `json:"bond_height"`
}

type Coin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

func getQueryPoolInfoRes(data []byte) (*QueryPoolInfoRes, error) {
	var res QueryPoolInfoRes
	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func getEraUpdateMsg(poolAddr string) []byte {
	eraUpdateMsg := struct {
		PoolAddr `json:"era_update"`
	}{
		PoolAddr: PoolAddr{Addr: poolAddr},
	}
	marshal, _ := json.Marshal(eraUpdateMsg)
	return marshal
}

func getEraBondMsg(poolAddr string) []byte {
	eraBondMsg := struct {
		PoolAddr `json:"era_bond"`
	}{
		PoolAddr: PoolAddr{Addr: poolAddr},
	}
	marshal, _ := json.Marshal(eraBondMsg)
	return marshal
}

func getEraCollectWithdrawMsg(poolAddr string) []byte {
	eraCollectWithdrawMsg := struct {
		PoolAddr `json:"era_collect_withdraw"`
	}{
		PoolAddr: PoolAddr{Addr: poolAddr},
	}
	marshal, _ := json.Marshal(eraCollectWithdrawMsg)
	return marshal
}

func getEraRestakeMsg(poolAddr string) []byte {
	eraRestakeMsg := struct {
		PoolAddr `json:"era_restake"`
	}{
		PoolAddr: PoolAddr{Addr: poolAddr},
	}
	marshal, _ := json.Marshal(eraRestakeMsg)
	return marshal
}

func getEraActiveMsg(poolAddr string) []byte {
	eraActiveMsg := struct {
		PoolAddr `json:"era_active"`
	}{
		PoolAddr: PoolAddr{Addr: poolAddr},
	}
	marshal, _ := json.Marshal(eraActiveMsg)
	return marshal
}

func getDelegationICQRegisterMsg(poolAddr string) []byte {
	msg := struct {
		PoolAddr `json:"pool_update_delegations_query"`
	}{
		PoolAddr: PoolAddr{Addr: poolAddr},
	}
	marshal, _ := json.Marshal(msg)
	return marshal
}

type RedeemTokenForShareMsg struct {
	PoolAddr string `json:"pool_addr"`
	Tokens   []Coin `json:"tokens"`
}

func getRedeemTokenForShareMsg(poolAddr string, tokens []Coin) []byte {
	msg := struct {
		RedeemTokenForShareMsg `json:"redeem_token_for_share"`
	}{
		RedeemTokenForShareMsg: RedeemTokenForShareMsg{
			PoolAddr: poolAddr,
			Tokens:   tokens,
		},
	}
	marshal, _ := json.Marshal(msg)
	return marshal
}
