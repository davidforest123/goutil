package cloud

import (
	"goutil/basic/gerrors"
	"goutil/container/gnum"
	"goutil/container/gvolume"
)

/**
AWS S3，新建从未访问过的Key是强一致性的，新建已访问过的Key、Delete、Modify、ListBucket等其他操作都是最终一致性eventual consistency的
https://segmentfault.com/a/1190000022079593

各大平台特性对比
https://www.cnblogs.com/fastone/p/11766161.html

除了AWS S3，其他平台的对象存储都是强一致性的。
*/

const (
	// ACLPrivate definition : private read and write
	ACLPrivate ACLType = "private"
	// ACLPublicRead definition : public read and private write
	ACLPublicRead ACLType = "public-read"
	// ACLPublicReadWrite definition : public read and public write
	ACLPublicReadWrite ACLType = "public-read-write"
)

type (
	Platform string

	Production string

	PlatformInfo struct {
		SupportPrepaid         bool
		SupportPostpaid        bool
		IsGetsStrongConsist    bool
		IsNonGetsStrongConsist bool
	}

	Balance struct {
		Currency  string
		Available gnum.Decimal
	}

	SysImage struct {
		Id        string
		Name      string
		OS        string // "linux", "windows"...
		Distro    string // "debian", "ubuntu"...
		Arch      string // "x86_64"...
		Available bool
	}

	// TODO：Region 要不要优化掉
	// AWS Lambda的demo
	// https://github.com/razeone/serverless-blog-api/blob/master/hello/main.go
	cloudBasic interface {
		GetBalance() (*Balance, error)
		ListLocations() (map[string][]string, error)
		VmListSpotSpecs(region string) (*InstanceSpecList, error)
		VmListImages(region string) ([]SysImage, error)
		VmListInstances(region string) ([]InstanceInfo, error)
		VmListSecurityGroups(region string) ([]SecurityGroup, error)
		VmCreateSecurityGroup(region string, sg SecurityGroup) (string, error)
		VmDeleteSecurityGroup(region, securityGroupId string) error
		VmListSwitches(region, zoneId string) ([]string, error)
		VmCreateInstances(region string, tmpl InstanceCreationTmpl, num int) ([]string, error)
		VmStartInstances(region string, instanceIds []string) error
		VmDeleteInstances(region string, instanceIds []string, force bool) error
		OsIsBucketExist(region, bucketName string) (bool, error)
		OsCreateBucket(region, bucketName string) error
		OsDeleteBucket(region, bucketName string, deleteIfNotEmpty *bool) error
		OsListBuckets(region string) ([]string, error)
		OsListObjectKeys(region, bucketName string, keyPrefix *string, pageSize int, pageToken *string) ([]string, *string, error)
		OsGetObjectSize(region, bucketName, objectKey string) (*gvolume.Volume, error)
		OsGetObject(region, bucketName, objectKey string) ([]byte, error)
		OsUpsertObject(region, bucketName, objectKey string, objectVal []byte) error
		OsRenameObject(region, bucketName, oldObjectKey, newObjectKey string) error
		OsDeleteObject(region, bucketName, objectKey string) error
		MqListQueues(region string, maxResults int, pageToken *string, queueNamePrefix *string) (queues []string, nextPageToken *string, err error)
		MqCreateQueue(region string, queue string, attr *QueueAttr) error
		MqSend(region string, queue string, msg string) error
		MqReceive(region string, queue string, msgNum int, waitSeconds int, deleteAfterReceived bool) ([]Msg, error)
		MqDeleteQueue(region string, queue string) error
		Close() error
	}

	Auth struct {
		AccessKey string
		SecretKey string
	}

	// ACLType bucket/object ACL
	ACLType string
)

var (
	VM  = Production("vm")  // virtual machine
	CTR = Production("ctr") // container
	FC  = Production("fc")  // function compute
	OS  = Production("os")  // object storage
	MQ  = Production("mq")  // message queue
	DNS = Production("dns") // DNS
	LB  = Production("LB")  // Load balance

	allPlatformInfos = map[Platform]PlatformInfo{}

	Aliyun = enrollPlatform("aliyun", PlatformInfo{SupportPrepaid: true, SupportPostpaid: false})
	AWS    = enrollPlatform("aws", PlatformInfo{SupportPrepaid: false, SupportPostpaid: true})
	AZURE  = enrollPlatform("azure", PlatformInfo{SupportPrepaid: false, SupportPostpaid: true})
	GCP    = enrollPlatform("gcp", PlatformInfo{SupportPrepaid: false, SupportPostpaid: true})
)

func enrollPlatform(name string, info PlatformInfo) Platform {
	allPlatformInfos[Platform(name)] = info
	return Platform(name)
}

func (p *Platform) Info() PlatformInfo {
	info, ok := allPlatformInfos[*p]
	if !ok {
		return PlatformInfo{}
	}
	return info
}

func (p Platform) String() string {
	return string(p)
}

func newBasic(platform Platform, accessKey string, secretKey string, LAN bool) (cloudBasic, error) {
	switch platform {
	case Aliyun:
		return newAliyun(accessKey, secretKey, LAN)
	default:
		return nil, gerrors.New("unsupported platform %s", platform)
	}
}

func GetAllVps(cloud cloudBasic) ([]InstanceInfo, error) {
	allLocations, err := cloud.ListLocations()
	if err != nil {
		return nil, err
	}
	var allRegions []string
	for k := range allLocations {
		allRegions = append(allRegions, k)
	}
	var result []InstanceInfo
	for _, v := range allRegions {
		items, err := cloud.VmListInstances(v)
		if err != nil {
			return nil, err
		}
		result = append(result, items...)
	}
	return result, nil
}

type (
	CheapestSpotVpsScanner struct {
		cli    cloudBasic
		tmpRes map[string]*InstanceSpecExList
	}
)

func NewCheapestSpotVpsScanner(cli cloudBasic) *CheapestSpotVpsScanner {
	return &CheapestSpotVpsScanner{
		cli:    cli,
		tmpRes: map[string]*InstanceSpecExList{},
	}
}

func (s *CheapestSpotVpsScanner) Scan() error {
	allLocations, err := s.cli.ListLocations()
	if err != nil {
		return err
	}
	var allRegions []string
	for k := range allLocations {
		allRegions = append(allRegions, k)
	}
	for _, region := range allRegions {
		if _, exist := s.tmpRes[region]; exist {
			continue
		}
		items, err := s.cli.VmListSpotSpecs(region)
		if err != nil {
			return err
		}
		s.tmpRes[region] = items.ToSpecExList()
	}
	return nil
}

func (s *CheapestSpotVpsScanner) GetCheapestNonCreditSpotVpsSpec() *InstanceSpecExList {
	res := &InstanceSpecExList{}
	for _, v := range s.tmpRes {
		res = res.Append(v)
	}
	return res.RemoveCreditInstance().Sort()
}
