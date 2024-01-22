package task

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func (t *Task) handleNewEra() error {
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
				_ = t.processPoolNewEra(addr)
			}(poolAddr)
		}
		wg.Wait()
		return nil
	}

	return t.processPoolNewEra(t.poolAddr)
}

func (t *Task) checkIcqSubmitHeight(icaAddr, queryKind string, lastStepHeight uint64) bool {

	query, err := t.getRegisteredIcqQuery(icaAddr, queryKind)
	if err != nil {
		return false
	}
	if query.RegisteredQuery.LastSubmittedResultLocalHeight < lastStepHeight {
		return false
	}

	return true
}

func (t *Task) processPoolNewEra(poolAddr string) error {
	var err error
	_, timestamp, err := t.neutronClient.GetCurrentBLockAndTimestamp()
	if err != nil {
		return err
	}
	poolInfo, err := t.getQueryPoolInfoRes(poolAddr)
	if err != nil {
		return err
	}

	targetEra := uint64(timestamp)/poolInfo.EraSeconds + poolInfo.Offset

	logger := logrus.WithFields(logrus.Fields{
		"pool":       poolAddr,
		"old era":    poolInfo.Era,
		"new era":    poolInfo.Era + 1,
		"target era": targetEra,
	})

	poolIca, err := t.getPoolIcaInfo(poolInfo.IcaId)
	if err != nil {
		return err
	}
	if len(poolIca) < 2 {
		logger.Warnln("ica data query failed")
		return nil
	}

	var txHash string
	switch poolInfo.EraProcessStatus {
	case ActiveEnded:
		// check targetEra to skip
		if targetEra <= poolInfo.Era {
			return nil
		}

		txHash, err = t.neutronClient.SendContractExecuteMsg(t.stakeManager, getEraUpdateMsg(poolAddr), nil)
		logger.WithFields(logrus.Fields{
			"current status": poolInfo.EraProcessStatus,
			"current rate":   poolInfo.Rate,
			"tx hash":        txHash,
		}).Infoln("start era-update")

	case EraUpdateEnded:
		if !t.checkIcqSubmitHeight(poolAddr, DelegationsQueryKind, poolInfo.EraSnapshot.BondHeight) {
			logger.Warnln("delegation icq query not ready")
			return nil
		}
		txHash, err = t.neutronClient.SendContractExecuteMsg(t.stakeManager, getEraBondMsg(poolAddr), nil)
		logger.WithFields(logrus.Fields{
			"current status":  poolInfo.EraProcessStatus,
			"current rate":    poolInfo.Rate,
			"snapshot bond":   poolInfo.EraSnapshot.Bond,
			"snapshot unbond": poolInfo.EraSnapshot.Unbond,
			"tx hash":         txHash,
		}).Infoln("start bond")
	case BondEnded:
		if !t.checkIcqSubmitHeight(poolIca[1].IcaAddr, BalancesQueryKind, poolInfo.EraSnapshot.BondHeight) {
			logger.Warnln("withdraw address balance icq query not ready")
			return nil
		}
		txHash, err = t.neutronClient.SendContractExecuteMsg(t.stakeManager, getEraCollectWithdrawMsg(poolAddr), nil)
		logger.WithFields(logrus.Fields{
			"current status": poolInfo.EraProcessStatus,
			"current rate":   poolInfo.Rate,
			"tx hash":        txHash,
		}).Infoln("start withdraw-collect")
	case WithdrawEnded:
		if !t.checkIcqSubmitHeight(poolAddr, DelegationsQueryKind, poolInfo.EraSnapshot.BondHeight) {
			logger.Warnln("withdraw address balance icq query not ready")
			return nil
		}
		txHash, err = t.neutronClient.SendContractExecuteMsg(t.stakeManager, getEraRestakeMsg(poolAddr), nil)
		logger.WithFields(logrus.Fields{
			"current status": poolInfo.EraProcessStatus,
			"current rate":   poolInfo.Rate,
			"tx hash":        txHash,
		}).Infoln("start restake")
	case RestakeEnded:
		if !t.checkIcqSubmitHeight(poolAddr, DelegationsQueryKind, poolInfo.EraSnapshot.BondHeight) {
			logger.Warnln("delegation icq query not ready")
			return nil
		}
		txHash, err = t.neutronClient.SendContractExecuteMsg(t.stakeManager, getEraActiveMsg(poolAddr), nil)
		logger.WithFields(logrus.Fields{
			"current status":  poolInfo.EraProcessStatus,
			"current rate":    poolInfo.Rate,
			"snapshot bond":   poolInfo.EraSnapshot.Bond,
			"snapshot unbond": poolInfo.EraSnapshot.Unbond,
			"snapshot active": poolInfo.EraSnapshot.Active,
			"bond":            poolInfo.Active,
			"unbond":          poolInfo.Active,
			"active":          poolInfo.Active,
			"tx hash":         txHash,
		}).Infoln("start era-active")
	default:
		return nil
	}

	if err != nil {
		logger.Warnf("failed, err: %s \n", err.Error())
		return err
	}

	if poolInfo.EraProcessStatus == RestakeEnded {
		retry := 0
		for {
			retry++
			if retry > 10 {
				return nil
			}
			time.Sleep(10 * time.Second)
			poolNewInfo, _ := t.getQueryPoolInfoRes(poolAddr)
			if poolNewInfo.EraProcessStatus == ActiveEnded {
				logger.WithFields(logrus.Fields{
					"new rate": poolNewInfo.Rate,
				}).
					Infof("new era task has complete")
				break
			}
		}
	}

	return nil
}
