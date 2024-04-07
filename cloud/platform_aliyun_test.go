package cloud

import (
	"fmt"
	"github.com/davidforest123/goutil/basic/gtest"
	"github.com/davidforest123/goutil/container/gnum"
	"github.com/davidforest123/goutil/encoding/gjson"
	"github.com/davidforest123/goutil/sys/gsysinfo"
	"testing"
)

func TestAliyunClient_ListLocations(t *testing.T) {
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), true)
	gtest.Assert(t, err)
	fmt.Println(cli.ListLocations())
}

func TestAliyunClient_GetBalance(t *testing.T) {
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), true)
	gtest.Assert(t, err)
	fmt.Println(cli.GetBalance())
}

func TestAliyunClient_EcListSpotSpecs(t *testing.T) {
	gsysinfo.SetEnv("ALI_ACCESS", "")
	gsysinfo.SetEnv("ALI_SECRET", "")
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), true)
	gtest.Assert(t, err)
	/*specs, err := cli.VmListSpotSpecs("cn-shenzhen")
	gtest.Assert(t, err)
	sortedSpecs := specs.ToSpecExList().RemoveCreditInstance().Sort()
	fmt.Println(gjson.MarshalStringDefault(sortedSpecs, true))
	*/
	res := NewCheapestSpotVpsScanner(cli)
	for {
		if err := res.Scan(); err != nil {
			fmt.Println(err)
		} else {
			break
		}
	}
	fmt.Println(gjson.MarshalStringDefault(res.GetCheapestNonCreditSpotVpsSpec(), true))
}

func TestAliyunClient_EcListImages(t *testing.T) {
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), true)
	gtest.Assert(t, err)
	imgs, err := cli.VmListImages("cn-hongkong")
	gtest.Assert(t, err)
	fmt.Println(gjson.MarshalStringDefault(imgs, true))
}

func TestAliyunClient_EcListVps(t *testing.T) {
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), true)
	gtest.Assert(t, err)
	ins, err := cli.VmListInstances("cn-hongkong")
	gtest.Assert(t, err)
	fmt.Println(gjson.MarshalStringDefault(ins, true))
}

func TestAliyunClient_EcCreateVps(t *testing.T) {
	regionId := "cn-hongkong"
	zoneId := "cn-hongkong-c"
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), true)
	gtest.Assert(t, err)

	switchIds, err := cli.VmListSwitches(regionId, zoneId)
	gtest.Assert(t, err)

	sgs, err := cli.VmListSecurityGroups(regionId)
	gtest.Assert(t, err)

	tmpl := InstanceCreationTmpl{
		ZoneId:          zoneId,
		SwitchId:        switchIds[0],
		IsSpot:          true,
		Specs:           "ecs.t5-lc2m1.nano",
		Name:            "test-spot-instance-name",
		ImageId:         "ubuntu_20_04_x64_20G_alibase_20201120.vhd",
		SecurityGroupId: sgs[0].Id,
		KeyPair:         "",
		Password:        "jcnde8r74BVGF",
		SystemDiskGB:    35,
		VpsCharge:       PostPaid,
		InternetCharge:  PayByTraffic,
		BandWidthMbIn:   21,
		BandWidthMbOut:  22,
		SpotStrategy:    SpotAsPriceGo,
		SpotMaxPrice:    gnum.NewDecimalFromFloat64(0.98),
	}
	ins, err := cli.VmCreateInstances(regionId, tmpl, 1)
	gtest.Assert(t, err)
	fmt.Println(gjson.MarshalStringDefault(ins, true))
}

func TestAliyunClient_EcDeleteVps(t *testing.T) {
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), true)
	gtest.Assert(t, err)
	err = cli.VmDeleteInstances("cn-hangzhou", []string{"i-j6cfq6ofuvebfuzy8qi1"}, true)
	gtest.Assert(t, err)
}

func TestAliyunClient_OssNewBucket(t *testing.T) {
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), true)
	gtest.Assert(t, err)
	fmt.Println(cli.OsListBuckets("cn-beijing"))
}

func TestAliyunClient_ObsGetObjectSize(t *testing.T) {
	cli, err := newAliyun(gsysinfo.GetEnv("ALI_ACCESS"), gsysinfo.GetEnv("ALI_SECRET"), false)
	gtest.Assert(t, err)
	fmt.Println(cli.OsGetObjectSize("cn-hongkong", "infdbchunk", "HTZrA5opN07tUu02XVZnzZ07QPlntIuS"))
}
