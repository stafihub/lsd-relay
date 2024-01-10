package task

import (
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/stafihub/lsd-relay/pkg/config"
	"github.com/stafihub/lsd-relay/pkg/utils"

	"github.com/stafihub/neutron-relay-sdk/client"
	"github.com/stafihub/neutron-relay-sdk/common/log"

	"github.com/sirupsen/logrus"
)

type Task struct {
	taskTicker    uint32
	stop          chan struct{}
	neutronClient *client.Client

	PoolAddr     string
	StakeManager string
}

func NewTask(cfg *config.Config) (*Task, error) {
	s := &Task{
		taskTicker:   cfg.TaskTicker,
		stop:         make(chan struct{}),
		PoolAddr:     cfg.PoolAddr,
		StakeManager: cfg.StakeManager,
	}

	kr, err := keyring.New("neutron", keyring.BackendTest, cfg.KeystorePath, os.Stdin, client.MakeEncodingConfig().Marshaler)
	if err != nil {
		return nil, err
	}

	c, err := client.NewClient(kr, cfg.KeyName, cfg.GasPrice, "neutron", cfg.EndpointList, log.NewLog("client", "neutron-relay"))
	if err != nil {
		return nil, err
	}
	s.neutronClient = c

	return s, nil
}

func (t *Task) Start() error {
	utils.SafeGoWithRestart(t.newEraHandler)
	utils.SafeGoWithRestart(t.redeemTokenHandler)
	utils.SafeGoWithRestart(t.icqRegisterHandle)
	return nil
}

func (t *Task) Stop() {
	close(t.stop)
}

func (t *Task) newEraHandler() {
	logrus.Debug("newEraHandler start -----------")
	logrus.Info("start new era Handler")

	err := t.handleNewEra()
	if err != nil {
		logrus.Warnf("newEraHandler failed, err: %s", err.Error())
	}

	ticker := time.NewTicker(time.Duration(t.taskTicker) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-t.stop:
			logrus.Info("new era task has stopped")
			return
		case <-ticker.C:
			err := t.handleNewEra()
			if err != nil {
				logrus.Warnf("newEraHandler failed, err: %s", err.Error())
				continue
			}
			logrus.Debug("newEraHandler end -----------")
		}
	}
}

func (t *Task) redeemTokenHandler() {
	logrus.Debug("redeemTokenHandler start -----------")
	logrus.Info("start redeem token Handler")

	err := t.handleRedeemShares()
	if err != nil {
		logrus.Warnf("redeemTokenHandler failed, err: %s", err.Error())
	}

	ticker := time.NewTicker(time.Duration(t.taskTicker) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-t.stop:
			logrus.Info("redeem token task has stopped")
			return
		case <-ticker.C:
			err := t.handleRedeemShares()
			if err != nil {
				logrus.Warnf("redeemTokenHandler failed, err: %s", err.Error())
				continue
			}
			logrus.Debug("redeemTokenHandler end -----------")
		}
	}
}

func (t *Task) icqRegisterHandle() {
	logrus.Debug("icqRegisterHandle start -----------")
	logrus.Info("start icq register Handler")

	err := t.handleICQRegister()
	if err != nil {
		logrus.Warnf("icqRegisterHandle failed, err: %s", err.Error())
	}

	ticker := time.NewTicker(time.Duration(t.taskTicker) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-t.stop:
			logrus.Info("icq register task has stopped")
			return
		case <-ticker.C:
			err := t.handleICQRegister()
			if err != nil {
				logrus.Warnf("icqRegisterHandle failed, err: %s", err.Error())
				continue
			}
			logrus.Debug("icqRegisterHandle end -----------")
		}
	}
}
