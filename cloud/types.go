package cloud

import (
	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/container/gnum"
	"github.com/davidforest123/goutil/container/gvolume"
	"sort"
	"time"
)

type (
	// How to charge.
	ChargeType string

	// Internet charge type.
	InternetChargeType string

	// Spot strategy.
	SpotStrategy string

	// Instance specification.
	InstanceSpec struct {
		RegionId         string
		Id               string
		IsCredit         bool // Is credit instance(突发性能实例) or normal instance.
		Currency         string
		LogicalCpuNum    int
		MemoryVolume     gvolume.Volume
		AvailableZoneIds []string
		OnDemandPrices   map[string]gnum.Decimal
		SpotPricePerHour map[string]gnum.Decimal
	}

	InstanceSpecEx struct {
		RegionId         string
		ZoneId           string
		Id               string
		IsCredit         bool // Is credit instance(突发性能实例) or normal instance.
		Currency         string
		LogicalCpuNum    int
		MemoryVolume     gvolume.Volume
		Network          string
		OnDemandPrices   gnum.Decimal
		SpotPricePerHour gnum.Decimal
	}

	InstanceSpecList struct {
		Specs []InstanceSpec
	}

	InstanceSpecExList struct {
		SpecExs []InstanceSpecEx
	}

	// InstanceInfo is a nested struct in ecs response.
	InstanceInfo struct {
		Id                 string
		Name               string
		Specs              string // "ecs.g5.large"
		InstanceChargeType ChargeType
		InternetChargeType ChargeType
		SpotPriceLimit     float64
		SpotStrategy       string
		SpotStartTime      time.Time
		NetworkType        string // "classic", "vpc"
		PrivateIPs         []string
		PublicIPs          []string
		RegionId           string
		ZoneId             string
		SecurityGroupIds   []string
		ImageId            string
		Status             string
		KeyPairName        string
		PhysicalCpuNum     int
		LogicalCpuNum      int
		GpuNum             int
		CreationTime       time.Time
		AutoReleaseTime    time.Time
		MemorySize         gvolume.Volume
		SysImageId         string
		SysImageName       string
		SysImageOS         string
	}

	// InstanceCreationTmpl is template to create instance.
	InstanceCreationTmpl struct {
		// Zone Id like "cn-hangzhou-c"
		ZoneId string
		// Start a spot instance.
		IsSpot bool
		// Open spot instance protection duration, only available on aliyun for now.
		// If open "OpenSpotDuration", price will be about a little (10% on aliyun) higher than closing it.
		OpenSpotDuration bool
		// "SwitchId" required in "vpc" network type, and now only "vpc" network type supported on aliyun,
		// "classic" network type is fading away, so "SwitchId" is required.
		SwitchId string
		// Instance specs like "ecs.g5.large".
		Specs string
		// Instance name like "my-instance-temp-name".
		Name string
		// OS image Id like "ubuntu_20_04_x64_20G_alibase_20201120.vhd".
		ImageId         string
		SecurityGroupId string
		KeyPair         string
		Password        string
		SystemDiskGB    int
		VpsCharge       ChargeType
		InternetCharge  InternetChargeType
		BandWidthMbIn   int
		BandWidthMbOut  int
		SpotStrategy    SpotStrategy
		SpotMaxPrice    gnum.Decimal
		// true: unlimited performance mode, false: limited performance mode.
		// Available param option on aliyun.
		UnlimitedPerformance bool
	}

	SecurityPermission struct {
		Description  string
		Direction    string // "in", "out"
		Protocol     string // "tcp","udp"...
		SrcPortRange [2]int
		SrcCidrIP    string
		DstPortRange [2]int
		DstCidrIP    string
	}

	SecurityGroup struct {
		Id          string
		Name        string
		Permissions []SecurityPermission
	}

	Msg struct {
		Id          string
		EnqueueTime time.Time
		Data        string
	}

	QueueAttr struct {
		MsgMaxBytes              *int // The limit of how many bytes a message can contain, in bytes.
		MsgRetentionSeconds      *int // The length of time, in seconds, for which MQ retains a message, in seconds.
		VisibilityTimeoutSeconds *int // The visibility timeout for the queue, in seconds.
	}
)

var (
	PrePaid  = ChargeType("prepaid")  // Yearly package or monthly package.
	PostPaid = ChargeType("postpaid") // Pay on demand.

	PayByBandwidth = InternetChargeType("PayByBandwidth")
	PayByTraffic   = InternetChargeType("PayByTraffic")

	SpotWithMaxPrice = SpotStrategy("SpotWithMaxPrice")
	SpotAsPriceGo    = SpotStrategy("SpotAsPriceGo")
)

func (t InstanceCreationTmpl) Verify(platform Platform) error {
	switch platform {
	case Aliyun:
		if t.IsSpot {
			if t.VpsCharge != PostPaid {
				return gerrors.New("invalid spot charge: VpsCharge %s", t.VpsCharge)
			}
		}
	}
	return nil
}

func (vsl *InstanceSpecList) ToSpecExList() *InstanceSpecExList {
	res := &InstanceSpecExList{}
	for _, spec := range vsl.Specs {
		for zoneId, spotPrice := range spec.SpotPricePerHour {
			entry := InstanceSpecEx{
				RegionId:         spec.RegionId,
				ZoneId:           zoneId,
				Id:               spec.Id,
				IsCredit:         spec.IsCredit,
				Currency:         spec.Currency,
				LogicalCpuNum:    spec.LogicalCpuNum,
				MemoryVolume:     spec.MemoryVolume,
				Network:          "vpc",
				OnDemandPrices:   spec.OnDemandPrices[zoneId],
				SpotPricePerHour: spotPrice,
			}
			res.SpecExs = append(res.SpecExs, entry)
		}
	}
	return res
}

func (vl InstanceSpecExList) Len() int {
	return len(vl.SpecExs)
}

func (vl InstanceSpecExList) Less(i, j int) bool {
	iCpuNum := gnum.NewDecimalFromInt(vl.SpecExs[i].LogicalCpuNum)
	jCpuNum := gnum.NewDecimalFromInt(vl.SpecExs[j].LogicalCpuNum)
	if iCpuNum.IsZero() || jCpuNum.IsZero() {
		return false
	}
	iSpotPrice := vl.SpecExs[i].SpotPricePerHour
	jSpotPrice := vl.SpecExs[j].SpotPricePerHour
	return iSpotPrice.Div(iCpuNum).LessThan(jSpotPrice.Div(jCpuNum))
}

func (vl InstanceSpecExList) Swap(i, j int) {
	vl.SpecExs[i], vl.SpecExs[j] = vl.SpecExs[j], vl.SpecExs[i]
}

func (vl *InstanceSpecExList) RemoveCreditInstance() *InstanceSpecExList {
	res := &InstanceSpecExList{}
	for _, v := range vl.SpecExs {
		if v.IsCredit {
			continue
		}
		res.SpecExs = append(res.SpecExs, v)
	}
	return res
}

func (vl *InstanceSpecExList) Sort() *InstanceSpecExList {
	sort.Sort(vl)
	return vl
}

func (vl *InstanceSpecExList) Append(newList *InstanceSpecExList) *InstanceSpecExList {
	for _, v := range newList.SpecExs {
		vl.SpecExs = append(vl.SpecExs, v)
	}
	return vl
}
