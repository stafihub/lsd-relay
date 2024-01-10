package task

import "github.com/sirupsen/logrus"

func (t *Task) handleICQRegister() error {
	if t.PoolAddr == "" {
		// todo: read the trust pool list
		return nil
	}

	return t.processICQRegister(t.PoolAddr)
}

func (t *Task) processICQRegister(poolAddr string) error {
	poolInfoRes, err := t.neutronClient.QuerySmartContractState(t.StakeManager, getQueryPoolInfoReq(poolAddr))
	if err != nil {
		return err
	}
	res, err := getQueryPoolInfoRes(poolInfoRes.Data.Bytes())
	if err != nil {
		return err
	}
	if res.EraProcessStatus == WaitQueryUpdate {
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
