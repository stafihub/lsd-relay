package task

import (
	"github.com/sirupsen/logrus"
)

func (t *Task) handleNewEra() error {
	_, timestamp, err := t.neutronClient.GetCurrentBLockAndTimestamp()
	if err != nil {
		return err
	}
	if t.PoolAddr == "" {
		// todo: read the trust pool list
		return nil
	}

	return t.processPoolNewEra(t.PoolAddr, uint64(timestamp))
}

func (t *Task) processPoolNewEra(poolAddr string, timestamp uint64) error {
	poolInfoRes, err := t.neutronClient.QuerySmartContractState(t.StakeManager, getQueryPoolInfoReq(t.PoolAddr))
	if err != nil {
		return err
	}
	poolInfo, err := getQueryPoolInfoRes(poolInfoRes.Data.Bytes())
	if err != nil {
		return err
	}
	logrus.Debugf("process pool: %+v \n", poolInfo)
	var msg []byte
	switch poolInfo.EraProcessStatus {
	case ActiveEnded:
		// check time to skip
		era := timestamp/poolInfo.EraSeconds + poolInfo.Offset
		if era <= poolInfo.Era {
			logrus.Warnf("pool %s era %d not end yet \n", poolAddr, era)
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
		logrus.Debugf("pool %s era status %s \n", t.PoolAddr, poolInfo.EraProcessStatus)
	}

	txHash, err := t.neutronClient.SendContractExecuteMsg(t.StakeManager, msg, nil)
	if err != nil {
		return err
	}
	logrus.Infof("tx %s send success \n", txHash)

	return nil
}
