package task

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
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

	if poolInfo.EraProcessStatus != WithdrawEnded {
		return nil
	}

	_, timestamp, err := t.neutronClient.GetCurrentBLockAndTimestamp()
	if err != nil {
		return err
	}
	targetEra := uint64(timestamp)/poolInfo.EraSeconds + poolInfo.Offset

	poolIca, err := t.getPoolIcaInfo(poolInfo.IcaId)
	if err != nil {
		return err
	}
	if len(poolIca) < 2 {
		return errors.New("ica data query failed")
	}

	logger := logrus.WithFields(logrus.Fields{
		"pool":           poolAddr,
		"target era":     targetEra,
		"old era":        poolInfo.Era - 1,
		"new era":        poolInfo.Era,
		"current status": poolInfo.EraProcessStatus,
		"current rate":   poolInfo.Rate,
		"action":         newEraRebondFuncName,
	})

	if !t.checkIcqSubmitHeight(poolAddr, DelegationsQueryKind, poolInfo.EraSnapshot.BondHeight) {
		logger.Warnln("delegation icq query not ready")
		return nil
	}

	txHash, err := t.neutronClient.SendContractExecuteMsg(t.stakeManager, getEraRestakeMsg(poolAddr), nil)
	if err != nil {
		logger.Warnf("failed, err: %s \n", err.Error())
		return err
	}

	logger.WithFields(logrus.Fields{
		"tx hash": txHash,
	}).Infoln("success")

	return nil
}
