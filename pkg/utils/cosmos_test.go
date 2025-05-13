package utils_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stafihub/lsd-relay/pkg/utils"
	"github.com/stafihub/neutron-relay-sdk/client"
	"github.com/stafihub/neutron-relay-sdk/common/log"
)

const endpoint = "https://cosmos-rest.publicnode.com"

func TestCosmos(t *testing.T) {
	p, err := utils.GetStakingParams(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p)
	tls, err := utils.GetTotalLiquidStake(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tls)

	v, err := utils.GetValidator(endpoint, "cosmosvaloper1zqgheeawp7cmqk27dgyctd80rd8ryhqs6la9wc")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)

	des, err := utils.GetDelegatorDelegations(endpoint, "cosmos1h9f42rv9tracjnc23znjw0jv4hnv4jy9hvx3549tzkvflgvxarwsyf0ehn")
	if err != nil {
		t.Fatal(err)
	}

	total := decimal.Zero
	for _, d := range des.DelegationResponses {
		amount, err := decimal.NewFromString(d.Balance.Amount)
		if err != nil {
			t.Fatal(err)
		}
		total = total.Add(amount)
	}

	t.Log(total)

	sv, err := utils.SelectVals([]string{
		"cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn",
		"cosmosvaloper1zqgheeawp7cmqk27dgyctd80rd8ryhqs6la9wc",
		"cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k",
		"cosmosvaloper1y0us8xvsvfvqkk9c6nt5cfyu5au5tww2ztve7q",
		"cosmosvaloper1n229vhepft6wnkt5tjpwmxdmcnfz55jv3vp77d",
		"cosmosvaloper1vvwtk805lxehwle9l4yudmq6mn0g32px9xtkhc",
		"cosmosvaloper1ualhu3fjgg77g485gmyswkq3w0dp7gys6qzwrv",
		"cosmosvaloper1lcwxu50rvvgf9v6jy6q5mrzyhlszwtjxhtscmp",
		"cosmosvaloper1q6d3d089hg59x6gcx92uumx70s5y5wadklue8s",
		"cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx",
	}, "1", "1", endpoint)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sv)
}

func TestGetDelegations(t *testing.T) {
	c, err := client.NewClient(nil, "", "", "neutron", []string{"https://neutron-rpc.publicnode.com:443"}, log.NewLog("", ""))
	if err != nil {
		t.Fatal(err)
	}
	icaAddr := "cosmos1h9f42rv9tracjnc23znjw0jv4hnv4jy9hvx3549tzkvflgvxarwsyf0ehn"
	msg := fmt.Sprintf(`{"delegations":{"pool_addr":"%s","sdk_greater_or_equal_v047":true}}`, icaAddr)
	rawRes, err := c.QuerySmartContractState("neutron1jzn038zknkz2cx3qkfpnxa9uztgyf36kam4akc9w489hwxnstc3s86wmvd", []byte(msg))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rawRes)

	var res DelegationsRes
	err = json.Unmarshal(rawRes.Data.Bytes(), &res)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)

	total := decimal.Zero
	for _, d := range res.Delegations {
		amount, err := decimal.NewFromString(d.Amount.Amount)
		if err != nil {
			t.Fatal(err)
		}
		total = total.Add(amount)
	}

	delegationsNative, err := utils.GetDelegatorDelegations(endpoint, icaAddr)
	if err != nil {
		t.Fatal(err)
	}
	totalNative := decimal.Zero
	for _, d := range delegationsNative.DelegationResponses {
		amount, err := decimal.NewFromString(d.Balance.Amount)
		if err != nil {
			t.Fatal(err)
		}
		totalNative = totalNative.Add(amount)
	}

	t.Log(total.Equal(totalNative))
	t.Log("total",total)
	t.Log("totalNative",totalNative)
}

type DelegationsRes struct {
	Delegations []struct {
		Delegator string `json:"delegator"`
		Validator string `json:"validator"`
		Amount    struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"amount"`
	} `json:"delegations"`
	LastSubmittedLocalHeight int `json:"last_submitted_local_height"`
}
