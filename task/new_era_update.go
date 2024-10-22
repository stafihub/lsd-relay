package task

import (
	"sync"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/sirupsen/logrus"
)

var newEraUpdateFuncName = "NewEraUpdate"

func (t *Task) handleNewEraUpdate() error {
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
				_ = t.processPoolNewEraUpdate(addr)
			}(poolAddr)
		}
		wg.Wait()
		return nil
	}

	return t.processPoolNewEraUpdate(t.poolAddr)
}

func (t *Task) processPoolNewEraUpdate(poolAddr string) error {
	var err error
	poolInfo, err := t.getQueryPoolInfoRes(poolAddr)
	if err != nil {
		return err
	}

	if poolInfo.Status != ActiveEnded {
		return nil
	}

	_, timestamp, err := t.neutronClient.GetCurrentBLockAndTimestamp()
	if err != nil {
		return err
	}
	targetEra := timestamp/int64(poolInfo.EraSeconds) + poolInfo.Offset

	// check targetEra to skip
	if targetEra <= int64(poolInfo.Era) {
		return nil
	}

	logger := logrus.WithFields(logrus.Fields{
		"pool":      poolAddr,
		"targetEra": targetEra,
		"newEra":    poolInfo.Era + 1,
		"action":    newEraUpdateFuncName,
	})

	ibcFee, err := t.neutronClient.GetTotalIbcFee()
	if err != nil {
		return err
	}
	ibcFeeCoins := types.NewCoins(types.NewCoin(t.neutronClient.GetDenom(), ibcFee))

	txHash, err := t.neutronClient.SendContractExecuteMsg(t.stakeManager, getEraUpdateMsg(poolAddr), ibcFeeCoins)
	if err != nil {
		logger.Warnf("failed, err: %s \n", err.Error())
		return nil
	}

	logger.WithFields(logrus.Fields{
		"txHash": txHash,
	}).Infoln("success")

	return nil
}
