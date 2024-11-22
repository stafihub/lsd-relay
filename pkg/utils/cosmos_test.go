package utils_test

import (
	"testing"

	"github.com/stafihub/lsd-relay/pkg/utils"
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
