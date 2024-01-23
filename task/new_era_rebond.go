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

	if poolInfo.Status != WithdrawEnded {
		return nil
	}

	poolIca, err := t.getPoolIcaInfo(poolInfo.IcaId)
	if err != nil {
		return err
	}
	if len(poolIca) < 2 {
		return errors.New("ica data query failed")
	}

	logger := logrus.WithFields(logrus.Fields{
		"pool":   poolAddr,
		"action": newEraRebondFuncName,
	})

	if !t.checkIcqSubmitHeight(poolAddr, DelegationsQueryKind, poolInfo.EraSnapshot.LastStepHeight) {
		logger.Warnln("delegation interchain query not ready")
		return nil
	}

	txHash, err := t.neutronClient.SendContractExecuteMsg(t.stakeManager, getEraRebondMsg(poolAddr), nil)
	if err != nil {
		logger.Warnf("failed, err: %s \n", err.Error())
		return err
	}

	logger.WithFields(logrus.Fields{
		"txHash": txHash,
	}).Infoln("success")

	return nil
}
