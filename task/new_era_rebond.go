package task

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/sirupsen/logrus"
	"sync"
)

var newEraRebondFuncName = "NewEraRebond"

func (t *Task) handleNewEraRebond() error {
	if t.runForEntrustedPool {
		stackInfo, err := t.getStackInfoRes()
		if err != nil {
			return err
		}
		wg := sync.WaitGroup{}
		for _, poolAddr := range stackInfo.Pools {
			wg.Add(1)
			poolAddr := poolAddr
			go func(addr string) {
				defer wg.Done()
				_ = t.processPoolNewEraRebond(addr)
			}(poolAddr)
		}
		wg.Wait()
		return nil
	}

	return t.processPoolNewEraRebond(t.poolAddr)
}

func (t *Task) processPoolNewEraRebond(poolAddr string) error {
	var err error

	poolInfo, err := t.getQueryPoolInfoRes(poolAddr)
	if err != nil {
		return err
	}

	if poolInfo.Status != WithdrawEnded {
		return nil
	}

	logger := logrus.WithFields(logrus.Fields{
		"pool":         poolAddr,
		"rebondAmount": poolInfo.EraSnapshot.RestakeAmount,
		"action":       newEraRebondFuncName,
	})

	if !t.checkIcqSubmitHeight(poolAddr, DelegationsQueryKind, poolInfo.EraSnapshot.LastStepHeight) {
		logger.Warnln("delegation interchain query not ready")
		return nil
	}
	ibcFee, err := t.neutronClient.GetTotalIbcFee()
	if err != nil {
		return err
	}
	ibcFeeCoins := types.NewCoins(types.NewCoin(t.neutronClient.GetDenom(), ibcFee))
	txHash, err := t.neutronClient.SendContractExecuteMsg(t.stakeManager, getEraRebondMsg(poolAddr), ibcFeeCoins)
	if err != nil {
		logger.Warnf("failed, err: %s \n", err.Error())
		return err
	}

	logger.WithFields(logrus.Fields{
		"txHash": txHash,
	}).Infoln("success")

	return nil
}
