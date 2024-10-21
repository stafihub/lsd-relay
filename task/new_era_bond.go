package task

import (
	"errors"
	"sync"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/sirupsen/logrus"
)

var newEraBondFuncName = "NewEraBond"

func (t *Task) handleNewEraBond() error {
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
				_ = t.processPoolNewEraBond(addr)
			}(poolAddr)
		}
		wg.Wait()
		return nil
	}

	return t.processPoolNewEraBond(t.poolAddr)
}

func (t *Task) processPoolNewEraBond(poolAddr string) error {
	poolInfo, err := t.getQueryPoolInfoRes(poolAddr)
	if err != nil {
		return err
	}

	if poolInfo.Status != EraUpdateEnded {
		return nil
	}

	logger := logrus.WithFields(logrus.Fields{
		"pool":           poolAddr,
		"snapshotBond":   poolInfo.EraSnapshot.Bond,
		"snapshotUnbond": poolInfo.EraSnapshot.Unbond,
		"snapshotActive": poolInfo.EraSnapshot.Active,
		"action":         newEraBondFuncName,
	})

	poolIca, err := t.getPoolIcaInfo(poolInfo.IcaId)
	if err != nil {
		return err
	}
	if len(poolIca) < 2 {
		return errors.New("ica data query failed")
	}

	if !t.checkIcqSubmitHeight(poolAddr, DelegationsQueryKind, poolInfo.EraSnapshot.LastStepHeight) {
		logger.Warnln("delegation interchain query not ready")
		return nil
	}
	ibcFee, err := t.neutronClient.GetTotalIbcFee()
	if err != nil {
		return err
	}
	ibcFeeCoins := types.NewCoins(types.NewCoin(t.neutronClient.GetDenom(), ibcFee))
	txHash, err := t.neutronClient.SendContractExecuteMsg(t.stakeManager, getEraBondMsg(poolAddr), ibcFeeCoins)
	if err != nil {
		logger.Warnf("failed, err: %s \n", err.Error())
		return err
	}

	logger.WithFields(logrus.Fields{
		"txHash": txHash,
	}).Infoln("success")

	return nil
}
