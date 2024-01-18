package task

import "github.com/sirupsen/logrus"

func (t *Task) handleIcqUpdate() error {
	if t.runForEntrustedPool {
		stackInfo, err := t.getStackInfoRes()
		if err != nil {
			return err
		}
		for _, pool := range stackInfo.Pools {
			if err := t.processPoolIcqUpdate(pool); err != nil {
				return err
			}
		}
		return nil
	}

	return t.processPoolIcqUpdate(t.poolAddr)
}

func (t *Task) processPoolIcqUpdate(poolAddr string) error {
	poolInfo, err := t.getQueryPoolInfoRes(poolAddr)
	if err != nil {
		return err
	}
	if poolInfo.EraProcessStatus != WaitQueryUpdate {
		return nil
	}

	msg := getPoolUpdateQueryExecuteMsg(poolAddr)
	txHash, err := t.neutronClient.SendContractExecuteMsg(t.stakeManager, msg, nil)
	if err != nil {
		return err
	}
	logrus.Infof("pool %s delegations icq register tx %s send success \n", poolAddr, txHash)
	return nil
}
