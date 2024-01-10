package task

import (
	"github.com/sirupsen/logrus"
	"github.com/stafihub/lsd-relay/pkg/utils"
)

func (t *Task) handleRedeemShares() error {
	if t.runForEntrustedPool {
		stackInfo, err := t.getStackInfoRes()
		if err != nil {
			return err
		}
		for _, pool := range stackInfo.Pools {
			if err := t.processPoolRedeemShares(pool); err != nil {
				return err
			}
		}
		return nil
	}

	return t.processPoolRedeemShares(t.poolAddr)
}

func (t *Task) processPoolRedeemShares(poolAddr string) error {
	poolInfo, err := t.getQueryPoolInfoRes(poolAddr)
	if err != nil {
		return err
	}
	if !poolInfo.LsmSupport {
		return nil
	}
	if len(poolInfo.ShareTokens) > 0 {
		var coins []Coin
		for _, k := range poolInfo.ShareTokens {
			shareToken := k
			if utils.ContainsString(poolInfo.RedeemmingShareTokenDenom, shareToken.Denom) {
				continue
			}
			coins = append(coins, shareToken)
		}
		msg := getRedeemTokenForShareMsg(t.poolAddr, coins)
		txHash, err := t.neutronClient.SendContractExecuteMsg(t.stakeManager, msg, nil)
		if err != nil {
			return err
		}
		logrus.Infof("pool %s redeem tx %s send success \n", poolAddr, txHash)
	}
	return nil
}
