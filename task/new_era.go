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

func (t *Task) processPoolNewEra(poolAddr string) error {
	_, timestamp, err := t.neutronClient.GetCurrentBLockAndTimestamp()
	if err != nil {
		return err
	}
	poolInfo, err := t.getQueryPoolInfoRes(poolAddr)
	if err != nil {
		return err
	}
	targetEra := uint64(timestamp)/poolInfo.EraSeconds + poolInfo.Offset

	var msg []byte
	switch poolInfo.EraProcessStatus {
	case ActiveEnded:
		// check targetEra to skip
		if targetEra <= poolInfo.Era {
			logrus.Infof("pool %s era %d not end yet \n", poolAddr, poolInfo.Era)
			return nil
		}
		logrus.Infof("pool-%s start new era update: old era: %d new era: %d current rate: %s target era: %d \n",
			poolAddr, poolInfo.Era, poolInfo.Era+1, poolInfo.Rate, targetEra)
		msg = getEraUpdateMsg(poolAddr)
	case EraUpdateEnded:
		logrus.Infof("pool-%s start new era bond: old era: %d new era: %d snapshot bond: %s snapshot unbond: %s \n",
			poolAddr, poolInfo.Era, poolInfo.Era+1, poolInfo.EraSnapshot.Bond, poolInfo.EraSnapshot.Unbond)
		msg = getEraBondMsg(poolAddr)
	case BondEnded:
		logrus.Infof("pool-%s start new era collect withdraw to pool: old era: %d new era: %d \n",
			poolAddr, poolInfo.Era, poolInfo.Era+1)
		msg = getEraCollectWithdrawMsg(poolAddr)
	case WithdrawEnded:
		logrus.Infof("pool-%s start pool restake: old era: %d new era: %d \n",
			poolAddr, poolInfo.Era, poolInfo.Era+1)
		msg = getEraRestakeMsg(poolAddr)
	case RestakeEnded:
		logrus.Infof("pool-%s start new era active: old era: %d new era: %d target era: %d snapshot bond: %s snapshot unbond: %s snapshot active: %s real-time bond: %s real-time unbond: %s real-time active: %s \n",
			poolAddr, poolInfo.Era, poolInfo.Era+1, targetEra, poolInfo.EraSnapshot.Bond, poolInfo.EraSnapshot.Unbond, poolInfo.EraSnapshot.Active, poolInfo.Bond, poolInfo.Unbond, poolInfo.Active)
		msg = getEraActiveMsg(poolAddr)
	default:
		logrus.Infof("pool-%s era status %s \n skip", poolAddr, poolInfo.EraProcessStatus)
	}

	txHash, err := t.neutronClient.SendContractExecuteMsg(t.stakeManager, msg, nil)
	if err != nil {
		logrus.Warnf("pool-%s execute %s :failed, err: %s \n", poolAddr, StatusForExecute[poolInfo.EraProcessStatus], err.Error())
		return err
	}

	if poolInfo.EraProcessStatus == RestakeEnded {
		time.Sleep(10 * time.Second)
		poolNewInfo, err := t.getQueryPoolInfoRes(poolAddr)
		if err != nil {
			return err
		}
		if poolNewInfo.EraProcessStatus == ActiveEnded {
			logrus.Infof("pool-%s era update complete: tx: %s new era: %d target era: %d new rate: %s\n", poolAddr, txHash, poolNewInfo.Era, targetEra, poolNewInfo.Rate)
		}
	} else {
		logrus.Infof("pool-%s execute %s: tx: %s send success \n", poolAddr, StatusForExecute[poolInfo.EraProcessStatus], txHash)
	}

	return nil
}
