package task

import (
	"github.com/sirupsen/logrus"
	"github.com/stafihub/lsd-relay/pkg/utils"
)

func (t *Task) handleRedeemShares() error {
	if t.PoolAddr == "" {
		// todo: read the trust pool list
		return nil
	}

	return t.processPoolRedeemShares(t.PoolAddr)
}

func (t *Task) processPoolRedeemShares(poolAddr string) error {
	poolInfoRes, err := t.neutronClient.QuerySmartContractState(t.StakeManager, getQueryPoolInfoReq(poolAddr))
	if err != nil {
		return err
	}
	res, err := getQueryPoolInfoRes(poolInfoRes.Data.Bytes())
	if err != nil {
		return err
	}
	if len(res.ShareTokens) > 0 {
		var coins []Coin
		for _, k := range res.ShareTokens {
			shareToken := k
			if utils.ContainsString(res.RedeemmingShareTokenDenom, shareToken.Denom) {
				continue
			}
			coins = append(coins, shareToken)
		}
		msg := getRedeemTokenForShareMsg(t.PoolAddr, coins)
		txHash, err := t.neutronClient.SendContractExecuteMsg(t.StakeManager, msg, nil)
		if err != nil {
			return err
		}
		logrus.Infof("pool %s redeem tx %s send success \n", poolAddr, txHash)
	}
	return nil
}
