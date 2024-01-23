package task

import (
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
)

var eraActiveFuncName = "NewEraActive"

func (t *Task) handleNewEraActive() error {
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
				_ = t.processPoolNewEraActive(addr)
			}(poolAddr)
		}
		wg.Wait()
		return nil
	}

	return t.processPoolNewEraActive(t.poolAddr)
}

func (t *Task) processPoolNewEraActive(poolAddr string) error {
	var err error

	poolInfo, err := t.getQueryPoolInfoRes(poolAddr)
	if err != nil {
		return err
	}

	if poolInfo.EraProcessStatus != RebondEnded {
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
		"pool":    poolAddr,
		"oldRate": poolInfo.Rate,
		"action":  eraActiveFuncName,
	})

	if !t.checkIcqSubmitHeight(poolAddr, DelegationsQueryKind, poolInfo.EraSnapshot.BondHeight) {
		logger.Warnln("delegation icq query not ready")
		return nil
	}
	txHash, err := t.neutronClient.SendContractExecuteMsg(t.stakeManager, getEraActiveMsg(poolAddr), nil)
	if err != nil {
		logger.Warnf("failed, err: %s \n", err.Error())
		return err
	}

	retry := 0
	for {
		retry++
		if retry > 30 {
			logger.WithFields(logrus.Fields{
				"txHash": txHash,
			}).Warnln("tx success but result check been timeout")
			break
		}
		poolNewInfo, _ := t.getQueryPoolInfoRes(poolAddr)
		if poolNewInfo.EraProcessStatus == ActiveEnded {
			logger.WithFields(logrus.Fields{
				"active":  poolNewInfo.Active,
				"newRate": poolNewInfo.Rate,
				"txHash":  txHash,
			}).
				Infof("success(the new era task has been completed)")
			break
		}
	}

	return nil
}
