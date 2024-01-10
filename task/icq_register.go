package task

import "github.com/sirupsen/logrus"

func (t *Task) handleICQRegister() error {
	if t.PoolAddr == "" {
		stackInfo, err := t.getStackInfoRes()
		if err != nil {
			return err
		}
		for _, pool := range stackInfo.Pools {
			if err := t.processICQRegister(pool); err != nil {
				return err
			}
		}
		return nil
	}

	return t.processICQRegister(t.PoolAddr)
}

func (t *Task) processICQRegister(poolAddr string) error {
	poolInfo, err := t.getQueryPoolInfoRes(poolAddr)
	if err != nil {
		return err
	}
	if poolInfo.EraProcessStatus == WaitQueryUpdate {
		return nil
	}

	msg := getDelegationICQRegisterMsg(poolAddr)
	txHash, err := t.neutronClient.SendContractExecuteMsg(t.StakeManager, msg, nil)
	if err != nil {
		return err
	}
	logrus.Infof("delegations icq register tx %s send success \n", txHash)
	return nil
}
