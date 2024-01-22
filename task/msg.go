package task

import (
	"encoding/json"
)

// EraProcessStatus
const (
	EraUpdateStarted = "era_update_started"
	EraUpdateEnded   = "era_update_ended"
	BondStarted      = "bond_started"
	BondEnded        = "bond_ended"
	WithdrawStarted  = "withdraw_started"
	WithdrawEnded    = "withdraw_ended"
	RestakeStarted   = "restake_started"
	RestakeEnded     = "restake_ended"
	ActiveEnded      = "active_ended"
)

var StatusForExecute = map[string]string{
	ActiveEnded:    "era_update",
	EraUpdateEnded: "era_bond",
	BondEnded:      "era_withdraw_collect",
	WithdrawEnded:  "era_restake",
	RestakeEnded:   "era_active",
}

// ValidatorUpdateStatus
const (
	WaitQueryUpdate = "wait_query_update"
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
	IcaId                     string      `json:"ica_id"`
	Era                       uint64      `json:"era"`
	EraSeconds                uint64      `json:"era_seconds"`
	Offset                    uint64      `json:"offset"`
	Bond                      string      `json:"bond"`
	Unbond                    string      `json:"unbond"`
	Active                    string      `json:"active"`
	Rate                      string      `json:"rate"`
	RateChangeLimit           string      `json:"rate_change_limit"`
	EraProcessStatus          string      `json:"era_process_status"`
	ValidatorUpdateStatus     string      `json:"validator_update_status"`
	ShareTokens               []Coin      `json:"share_tokens"`
	RedeemmingShareTokenDenom []string    `json:"redeemming_share_token_denom"`
	EraSnapshot               eraSnapshot `json:"era_snapshot"`
	Paused                    bool        `json:"paused"`
	LsmSupport                bool        `json:"lsm_support"`
}

type StackInfoRes struct {
	Pools []string `json:"pools"`
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

func getPoolUpdateQueryExecuteMsg(poolAddr string) []byte {
	msg := struct {
		PoolAddr `json:"pool_update_validators_icq"`
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

func (t *Task) getQueryPoolInfoRes(poolAddr string) (*QueryPoolInfoRes, error) {
	poolInfoRes, err := t.neutronClient.QuerySmartContractState(t.stakeManager, getQueryPoolInfoReq(poolAddr))
	if err != nil {
		return nil, err
	}
	var res QueryPoolInfoRes
	err = json.Unmarshal(poolInfoRes.Data.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type StackInfoReq struct{}

func (t *Task) getStackInfoRes() (*StackInfoRes, error) {
	msg := struct {
		StackInfoReq `json:"stack_info"`
	}{}
	marshal, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	stackInfoRes, err := t.neutronClient.QuerySmartContractState(t.stakeManager, marshal)
	if err != nil {
		return nil, err
	}
	var res StackInfoRes
	err = json.Unmarshal(stackInfoRes.Data.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type UpdateIcqUpdatePeriodMsg struct {
	Addr            string `json:"pool_addr"`
	NewUpdatePeriod uint64 `json:"new_update_period"`
}

func getEraUpdatePeriodMsg(poolAddr string, period uint64) []byte {
	msg := struct {
		UpdateIcqUpdatePeriodMsg `json:"update_icq_update_period"`
	}{
		UpdateIcqUpdatePeriodMsg: UpdateIcqUpdatePeriodMsg{
			Addr:            poolAddr,
			NewUpdatePeriod: period,
		},
	}
	marshal, _ := json.Marshal(msg)
	return marshal
}
