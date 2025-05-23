package task

import (
	"encoding/json"
	"fmt"
)

// Status
const (
	EraUpdateStarted = "era_update_started"
	EraUpdateEnded   = "era_update_ended"
	BondStarted      = "bond_started"
	BondEnded        = "bond_ended"
	WithdrawStarted  = "withdraw_started"
	WithdrawEnded    = "withdraw_ended"
	RebondStarted    = "rebond_started"
	RebondEnded      = "rebond_ended"
	ActiveEnded      = "active_ended"
)

// QueryKind
const (
	BalancesQueryKind    = "balances"
	DelegationsQueryKind = "delegations"
	ValidatorsQueryKind  = "validators"
)

var StatusForExecute = map[string]string{
	ActiveEnded:    "era_update",
	EraUpdateEnded: "era_bond",
	BondEnded:      "era_withdraw_collect",
	WithdrawEnded:  "era_rebond",
	RebondEnded:    "era_active",
}

// ValidatorUpdateStatus
const (
	WaitQueryUpdate = "wait_query_update"
)

type PoolAddr struct {
	Addr string `json:"pool_addr"`
}

type PoolBond struct {
	Addr       string   `json:"pool_addr"`
	SelectVals []string `json:"select_vals"`
}

type QueryPoolInfoReq struct {
	PoolInfo PoolAddr `json:"pool_info"`
}

type StackInfoReq struct{}

type QueryPoolInfoRes struct {
	IcaId                     string      `json:"ica_id"`
	Era                       uint64      `json:"era"`
	EraSeconds                uint64      `json:"era_seconds"`
	Offset                    int64       `json:"offset"`
	Bond                      string      `json:"bond"`
	Unbond                    string      `json:"unbond"`
	Active                    string      `json:"active"`
	Rate                      string      `json:"rate"`
	RateChangeLimit           string      `json:"rate_change_limit"`
	Status                    string      `json:"status"`
	ValidatorUpdateStatus     string      `json:"validator_update_status"`
	ShareTokens               []Coin      `json:"share_tokens"`
	RedeemmingShareTokenDenom []string    `json:"redeemming_share_token_denom"`
	EraSnapshot               eraSnapshot `json:"era_snapshot"`
	Paused                    bool        `json:"paused"`
	LsmSupport                bool        `json:"lsm_support"`
	ValidatorAddrs            []string    `json:"validator_addrs"`
}
type DelegationsRes struct {
	Delegations []struct {
		Delegator string `json:"delegator"`
		Validator string `json:"validator"`
		Amount    struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"amount"`
	} `json:"delegations"`
	LastSubmittedLocalHeight int `json:"last_submitted_local_height"`
}

type RegisterQueryInfoRes struct {
	RegisteredQuery struct {
		Id    int    `json:"id"`
		Owner string `json:"owner"`
		Keys  []struct {
			Path string `json:"path"`
			Key  string `json:"key"`
		} `json:"keys"`
		QueryType                       string `json:"query_type"`
		TransactionsFilter              string `json:"transactions_filter"`
		ConnectionId                    string `json:"connection_id"`
		UpdatePeriod                    uint64 `json:"update_period"`
		LastSubmittedResultLocalHeight  uint64 `json:"last_submitted_result_local_height"`
		LastSubmittedResultRemoteHeight struct {
			RevisionNumber int `json:"revision_number"`
			RevisionHeight int `json:"revision_height"`
		} `json:"last_submitted_result_remote_height"`
		Deposit []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"deposit"`
		SubmitTimeout      int `json:"submit_timeout"`
		RegisteredAtHeight int `json:"registered_at_height"`
	} `json:"registered_query"`
}

type ICAData struct {
	Admin              string `json:"admin"`
	PoolAddressIcaInfo struct {
		CtrlChannelID    string `json:"ctrl_channel_id"`
		CtrlConnectionID string `json:"ctrl_connection_id"`
		CtrlPortID       string `json:"ctrl_port_id"`
		HostChannelID    string `json:"host_channel_id"`
		HostConnectionID string `json:"host_connection_id"`
		IcaAddr          string `json:"ica_addr"`
	} `json:"pool_address_ica_info"`
	WithdrawAddressIcaInfo struct {
		CtrlChannelID    string `json:"ctrl_channel_id"`
		CtrlConnectionID string `json:"ctrl_connection_id"`
		CtrlPortID       string `json:"ctrl_port_id"`
		HostChannelID    string `json:"host_channel_id"`
		HostConnectionID string `json:"host_connection_id"`
		IcaAddr          string `json:"ica_addr"`
	} `json:"withdraw_address_ica_info"`
}

type StackInfoRes struct {
	Pools []string `json:"pools"`
}

type eraSnapshot struct {
	Era            uint64 `json:"era"`
	Bond           string `json:"bond"`
	Unbond         string `json:"unbond"`
	Active         string `json:"active"`
	RestakeAmount  string `json:"restake_amount"`
	LastStepHeight uint64 `json:"last_step_height"`
}

type Coin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type RedeemTokenForShareMsg struct {
	PoolAddr string `json:"pool_addr"`
	Tokens   []Coin `json:"tokens"`
}

type UpdateIcqUpdatePeriodMsg struct {
	Addr            string `json:"pool_addr"`
	NewUpdatePeriod uint64 `json:"new_update_period"`
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

func getEraUpdateMsg(poolAddr string) []byte {
	eraUpdateMsg := struct {
		PoolAddr `json:"era_update"`
	}{
		PoolAddr: PoolAddr{Addr: poolAddr},
	}
	marshal, _ := json.Marshal(eraUpdateMsg)
	return marshal
}

func getEraBondMsg(poolAddr string, selVals []string) []byte {
	eraBondMsg := struct {
		PoolBond `json:"era_bond"`
	}{
		PoolBond: PoolBond{Addr: poolAddr, SelectVals: selVals},
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

func getEraRebondMsg(poolAddr string, selVals []string) []byte {
	msg := struct {
		PoolBond `json:"era_rebond"`
	}{
		PoolBond: PoolBond{Addr: poolAddr, SelectVals: selVals},
	}
	marshal, _ := json.Marshal(msg)
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

func (t *Task) getRegisteredIcqQuery(icaAddr, queryKind string) (*RegisterQueryInfoRes, error) {
	msg := fmt.Sprintf("{\"get_ica_registered_query\":{\"ica_addr\":\"%s\",\"query_kind\":\"%s\"}}", icaAddr, queryKind)
	rawRes, err := t.neutronClient.QuerySmartContractState(t.stakeManager, []byte(msg))
	if err != nil {
		return nil, err
	}
	var res RegisterQueryInfoRes
	err = json.Unmarshal(rawRes.Data.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (t *Task) GetDelegations(icaAddr string) (*DelegationsRes, error) {
	msg := fmt.Sprintf(`{"delegations":{"pool_addr":"%s","sdk_greater_or_equal_v047":true}}`, icaAddr)
	rawRes, err := t.neutronClient.QuerySmartContractState(t.stakeManager, []byte(msg))
	if err != nil {
		return nil, err
	}
	var res DelegationsRes
	err = json.Unmarshal(rawRes.Data.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (t *Task) getPoolIcaInfo(icaId string) (*ICAData, error) {
	msg := fmt.Sprintf("{\"interchain_account_address_from_contract\":{\"interchain_account_id\":\"%s\"}}", icaId)
	rawRes, err := t.neutronClient.QuerySmartContractState(t.stakeManager, []byte(msg))
	if err != nil {
		return nil, err
	}
	var res ICAData
	_ = json.Unmarshal(rawRes.Data.Bytes(), &res)
	return &res, nil
}

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
