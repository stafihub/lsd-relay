package task

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/sirupsen/logrus"
	"github.com/stafihub/lsd-relay/pkg/utils"
)

var redeemSharesFuncName = "ExecuteRedeemShares"

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
		logger := logrus.WithFields(logrus.Fields{
			"pool":   poolAddr,
			"action": redeemSharesFuncName,
		})
		msg := getRedeemTokenForShareMsg(poolAddr, coins)

		ibcFee, err := t.neutronClient.GetTotalIbcFee()
		if err != nil {
			return err
		}
		ibcFeeCoins := types.NewCoins(types.NewCoin(t.neutronClient.GetDenom(), ibcFee))

		txHash, err := t.neutronClient.SendContractExecuteMsg(t.stakeManager, msg, ibcFeeCoins)
		if err != nil {
			logger.Warnf("failed, err: %s \n", err.Error())
			return err
		}

		logger.WithFields(logrus.Fields{
			"txHash": txHash,
		}).Infoln("success")
	}
	return nil
}
