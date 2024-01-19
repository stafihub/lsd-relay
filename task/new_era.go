package task

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
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

	logrus.Infof(
		"Pool Information:\n"+
			"  - Pool Address: %s\n"+
			"  - Pool ICA ID: %s\n"+
			"  - Current Era: %d\n"+
			"  - Target Era: %d\n"+
			"  - Era Second: %d\n"+
			"  - Rate: %s\n"+
			"  - RateChangeLimit: %s\n"+
			"  - EraProcessStatus: %s\n"+
			"  - Bond: %s\n"+
			"  - ErSnapshot.Bond: %s\n"+
			"  - Unbond: %s\n"+
			"  - ErSnapshot.Unbond: %s\n"+
			"  - Active: %s\n"+
			"  - ErSnapshot.Active: %s\n",
		poolAddr,
		poolInfo.IcaId,
		poolInfo.Era,
		targetEra,
		poolInfo.EraSeconds,
		poolInfo.Rate,
		poolInfo.RateChangeLimit,
		poolInfo.EraProcessStatus,
		poolInfo.Bond,
		poolInfo.EraSnapshot.Bond,
		poolInfo.Unbond,
		poolInfo.EraSnapshot.Unbond,
		poolInfo.Active,
		poolInfo.EraSnapshot.Active,
	)

	var msg []byte
	switch poolInfo.EraProcessStatus {
	case ActiveEnded:
		// check targetEra to skip
		if targetEra <= poolInfo.Era {
			logrus.Infof("pool %s era %d not end yet \n", poolAddr, poolInfo.Era)
			return nil
		}
		msg = getEraUpdateMsg(poolAddr)
	case EraUpdateEnded:
		msg = getEraBondMsg(poolAddr)
	case BondEnded:
		msg = getEraCollectWithdrawMsg(poolAddr)
	case WithdrawEnded:
		msg = getEraRestakeMsg(poolAddr)
	case RestakeEnded:
		msg = getEraActiveMsg(poolAddr)
	default:
		logrus.Infof("pool %s era status %s \n skip", poolAddr, poolInfo.EraProcessStatus)
	}

	retry := 0
	for {
		var err error
		if retry > 10 {
			updatePeriodMsg := getEraUpdatePeriodMsg(poolAddr, 12)
			txHash, e := t.neutronClient.SendContractExecuteMsg(t.stakeManager, updatePeriodMsg, nil)
			if e != nil {
				return errors.Wrap(err, e.Error())
			}
			logrus.Infof("pool %s execute update_icq_update_period tx %s send success \n", poolAddr, txHash)
			time.Sleep(time.Second * 20)
		}
		if retry > 15 {
			return err
		}
		txHash, err := t.neutronClient.SendContractExecuteMsg(t.stakeManager, msg, nil)
		if err != nil {
			logrus.Warnf("pool %s execute %s failed, err: %s \n", poolAddr, StatusForExecute[poolInfo.EraProcessStatus], err.Error())
			if strings.Contains(err.Error(), "SubmissionHeight") {
				retry++
				time.Sleep(time.Second * 3)
				continue
			} else {
				return err
			}
		}
		logrus.Infof("pool %s execute %s tx %s send success \n", poolAddr, StatusForExecute[poolInfo.EraProcessStatus], txHash)
		break
	}

	return nil
}
