package cloud

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/account"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"goutil/basic/gerrors"
	"goutil/container/gvolume"
	"goutil/sys/gsync"
	"goutil/sys/gtime"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

/**
SNS特性：（参考https://aws.amazon.com/cn/sns/faqs/?ref=driverlayer.com/web）
1、保证接收成功。前提是要权限、网络、AWS后台等等都通畅。
2、不保证只接收一次，可能重复接收。
3、不保证按顺序接收，可能乱序接收。
*/

type (
	AwsClient struct {
		accessKey string
		secretKey string
		lan       bool
		sessions  sync.Map // map[region]session.Session
		topics    sync.Map // map[region.topic]topicARN
		queues    sync.Map // map[region.queue]queueURL
		ec2List   map[string]*ec2.EC2
		ec2ListMu sync.RWMutex
		s3List    map[string]*s3.S3
		s3ListMu  sync.RWMutex
		buckets   map[string]*s3.Bucket
		bucketsMu sync.RWMutex
		sqsList   map[string]*sqs.SQS
		sqsListMu sync.RWMutex
	}
)



func newAws(accessKey string, secretKey string, LAN bool) (*AwsClient, error) {
	return &AwsClient{
		accessKey: accessKey,
		secretKey: secretKey,
		lan:       LAN,
		ec2List:   map[string]*ec2.EC2{},
		s3List:    map[string]*s3.S3{},
		buckets:   map[string]*s3.Bucket{},
		sqsList:   map[string]*sqs.SQS{},
	}, nil
}

func (ac *AwsClient) GetBalance() (*Balance, error) {
	sess, err := ac.getSession(region)
	if err != nil {
		return nil, err
	}

	account.New(sess, &aws.Config{})
}

func (ac *AwsClient) ListLocations() (map[string][]string, error) {
	ec2cli, err := ac.getEc2Client("ap-east-1") // hong kong
	if err != nil {
		return nil, err
	}
	allZones := true
	in := &ec2.DescribeAvailabilityZonesInput{
		AllAvailabilityZones: &allZones,
	}
	out, err := ec2cli.DescribeAvailabilityZones(in)
	if err != nil {
		return nil, err
	}

	result := map[string][]string{}
	for _, v := range out.AvailabilityZones {
		val, ok := result[*v.RegionName]
		var zones []string
		if ok {
			zones = val
		}
		zones = append(zones, *v.ZoneId)
		result[*v.RegionName] = zones
	}
	return result, nil
}

func (ac *AwsClient) VmListOnDemandSpecs(region string) (*InstanceSpecList, error) {

}

func (ac *AwsClient) VmListSpotSpecs(region string) (*InstanceSpecList, error) {
	ec2cli, err := ac.getEc2Client(region)
	if err != nil {
		return nil, err
	}

	result := &InstanceSpecList{}
	nextToken := (*string)(nil)
	for {
		maxResults := int64(100)
		in := &ec2.DescribeInstanceTypesInput{
			MaxResults: &maxResults,
			NextToken:  nextToken,
		}
		out, err := ec2cli.DescribeInstanceTypes(in)
		if err != nil {
			return nil, err
		}
		for _, v := range out.InstanceTypes {
			result.Specs = append(result.Specs, InstanceSpec{
				RegionId:      region,
				Id:            *v.InstanceType,
				LogicalCpuNum: (int)(*v.VCpuInfo.DefaultCores),
				MemoryVolume:  gvolume.FromByteSizeUint64((uint64)(*v.MemoryInfo.SizeInMiB) * 1024 * 1024),
			})
		}
		nextToken = out.NextToken
		if nextToken == nil {
			break
		}
	}

	return result, nil
}

func (ac *AwsClient) VmListImages(region string) ([]SysImage, error) {
	ec2cli, err := ac.getEc2Client(region)
	if err != nil {
		return nil, err
	}

	var results []SysImage
	for {
		nextToken := (*string)(nil)
		in := &ec2.DescribeImagesInput{
			IncludeDeprecated: aws.Bool(false),
			IncludeDisabled:   aws.Bool(false),
			MaxResults:        aws.Int64(100),
			NextToken:         nextToken,
		}
		out, err := ec2cli.DescribeImages(in)
		if err != nil {
			return nil, err
		}
		for _, v := range out.Images {
			valOS := *v.Platform
			if valOS == "" {
				valOS = "linux"
			}
			results = append(results, SysImage{
				Id:        *v.ImageId,
				Name:      *v.Name,
				OS:        valOS,
				Distro:    "",
				Arch:      *v.Architecture,
				Available: true,
			})
		}
		nextToken = out.NextToken
		if nextToken == nil {
			break
		}
	}
	return results, nil
}

func (ac *AwsClient) VmListInstances(region string) ([]InstanceInfo, error) {
	ec2cli, err := ac.getEc2Client(region)
	if err != nil {
		return nil, err
	}

	var results []InstanceInfo
	for {
		nextToken := (*string)(nil)
		in := &ec2.DescribeInstancesInput{
			MaxResults: aws.Int64(100),
			NextToken:  nextToken,
		}
		out, err := ec2cli.DescribeInstances(in)
		if err != nil {
			return nil, err
		}
		for _, x := range out.Reservations {
			for _, y := range x.Instances {
				valOS := *y.Platform
				if valOS == "" {
					valOS = "linux"
				}
				results = append(results, InstanceInfo{
					Id:                 *y.InstanceId,
					Name:               "",
					Specs:              "",
					InstanceChargeType: "",
					InternetChargeType: "",
					SpotPriceLimit:     0,
					SpotStrategy:       "",
					SpotStartTime:      time.Time{},
					NetworkType:        "",
					PrivateIPs:         []string{*y.PrivateIpAddress},
					PublicIPs:          []string{*y.PublicIpAddress},
					RegionId:           region,
					ZoneId:             "",
					SecurityGroupIds:   nil,
					ImageId:            *y.ImageId,
					Status:             "",
					KeyPairName:        "",
					PhysicalCpuNum:     (int)(*y.CpuOptions.CoreCount),
					LogicalCpuNum:      (int)(*y.CpuOptions.CoreCount) * (int)(*y.CpuOptions.ThreadsPerCore),
					GpuNum:             len(y.ElasticGpuAssociations),
					CreationTime:       *y.LaunchTime,
					AutoReleaseTime:    gtime.ZeroTime,
					MemorySize:         gvolume.Zero,
					SysImageId:         "",
					SysImageName:       "",
					SysImageOS:         *y.Platform,
				})
			}

		}
		nextToken = out.NextToken
		if nextToken == nil {
			break
		}
	}
	return results, nil
}

func (ac *AwsClient) VmListSecurityGroups(region string) ([]SecurityGroup, error) {
	ec2cli, err := ac.getEc2Client(region)
	if err != nil {
		return nil, err
	}

	var results []InstanceInfo
	for {
		nextToken := (*string)(nil)
		in := &ec2.DescribeSecurityGroupsInput{
			MaxResults: aws.Int64(100),
			NextToken:  nextToken,
		}
		out, err := ec2cli.DescribeSecurityGroups(in)
		if err != nil {
			return nil, err
		}
		for _, x := range out.SecurityGroups {
				results = append(results, )
			item := SecurityGroup{
				Id:          *x.GroupId,
				Name:        *x.GroupName,
			}
			for _, y := range x.IpPermissions {
				item.Permissions = append(item.Permissions, SecurityPermission{
					Description:  *x.Description,
					Direction:    ,
					Protocol:     *y.IpProtocol,
					SrcPortRange: [2]int{int(*y.FromPort), int(*y.ToPort)},
					SrcCidrIP:    y.IpRanges[0].CidrIp,
					DstPortRange: [2]int{},
					DstCidrIP:    ,
			}

		}
		}
		nextToken = out.NextToken
		if nextToken == nil {
			break
		}
	}
	return results, nil
}

func (ac *AwsClient) VmCreateSecurityGroup(region string, sg SecurityGroup) (string, error) {
	ec2cli, err := ac.getEc2Client(region)
	if err != nil {
		return nil, err
	}

	in := &ec2.CreateSecurityGroupInput{
		Description:       aws.String(""),
		GroupName:         &sg.Name,
		TagSpecifications: nil,
		VpcId:             sg.,
	}
	out, err := ec2cli.CreateSecurityGroup(in)
	if err !=nil {
		return "", err
	}
	return *out.GroupId, nil
}

func (ac *AwsClient) VmDeleteSecurityGroup(region, securityGroupId string) error {
	ec2cli, err := ac.getEc2Client(region)
	if err != nil {
		return err
	}

	in := &ec2.DeleteSecurityGroupInput{
		GroupId: &securityGroupId,
	}
	_, err = ec2cli.DeleteSecurityGroup(in)
	return err
}

func (ac *AwsClient) VmListSwitches(region, zoneId string) ([]string, error) {
	ec2cli, err := ac.getEc2Client(region)
	if err != nil {
		return nil, err
	}
	nextPageToken := (*string)(nil)
	var results []string
	for  {
		in := &ec2.DescribeNetworkInterfacesInput{
			MaxResults:          aws.Int64(100),
			NextToken:           nextPageToken,
		}
		out, err := ec2cli.DescribeNetworkInterfaces(in)
		if err != nil {
			return nil, err
		}
		results = append(results, out.NetworkInterfaces[0].NetworkInterfaceId)
		nextPageToken = out.NextToken
		if nextPageToken == nil {
			break
		}
	}


}

func (ac *AwsClient) VmCreateInstances(region string, tmpl InstanceCreationTmpl, num int) ([]string, error) {
	ec2cli, err := ac.getEc2Client(region)
	if err != nil {
		return nil, err
	}

	in := &ec2.RunInstancesInput{
		AdditionalInfo:                    nil,
		BlockDeviceMappings:               nil,
		CapacityReservationSpecification:  nil,
		ClientToken:                       nil,
		CpuOptions:                        nil,
		CreditSpecification:               nil,
		DisableApiStop:                    nil,
		DisableApiTermination:             nil,
		DryRun:                            nil,
		EbsOptimized:                      nil,
		ElasticGpuSpecification:           nil,
		ElasticInferenceAccelerators:      nil,
		EnablePrimaryIpv6:                 nil,
		EnclaveOptions:                    nil,
		HibernationOptions:                nil,
		IamInstanceProfile:                nil,
		ImageId:                           &tmpl.ImageId,
		InstanceInitiatedShutdownBehavior: nil,
		InstanceMarketOptions:             nil,
		InstanceType:                      &tmpl.Specs,
		Ipv6AddressCount:                  nil,
		Ipv6Addresses:                     nil,
		KernelId:                          nil,
		KeyName:                           nil,
		LaunchTemplate:                    nil,
		LicenseSpecifications:             nil,
		MaintenanceOptions:                nil,
		MaxCount:                          aws.Int64(int64(num)),
		MetadataOptions:                   nil,
		MinCount:                          aws.Int64(int64(num)),
		Monitoring:                        nil,
		NetworkInterfaces:                 nil,
		Placement:                         nil,
		PrivateDnsNameOptions:             nil,
		PrivateIpAddress:                  nil,
		RamdiskId:                         nil,
		SecurityGroupIds:                  []*string{&tmpl.SecurityGroupId},
		SecurityGroups:                    nil,
		SubnetId:                          nil,
		TagSpecifications:                 nil,
		UserData:                          nil,
	}
	out, err := ec2cli.RunInstances(in)
	var results []string = nil
	if out != nil {
		for _, v := range out.Instances {
			results = append(results, *v.InstanceId)
		}
	}
	return results, err
}

func (ac *AwsClient) VmStartInstances(region string, instanceIds []string) error {
	ec2cli, err := ac.getEc2Client(region)
	if err != nil {
		return err
	}

	in := &ec2.StartInstancesInput{
		InstanceIds: aws.StringSlice(instanceIds),
	}
	_, err = ec2cli.StartInstances(in)
	return err
}

func (ac *AwsClient) VmDeleteInstances(region string, instanceIds []string, force bool) error {
	ec2cli, err := ac.getEc2Client(region)
	if err != nil {
		return err
	}

	in := &ec2.TerminateInstancesInput{
		InstanceIds: aws.StringSlice(instanceIds),
	}
	_, err = ec2cli.TerminateInstances(in)
	return err
}

func (ac *AwsClient) OsIsBucketExist(region, bucketName string) (bool, error) {
	results, err := ac.OsListBuckets(region)
	if err != nil {
		return false, err
	}
	for _, v := range results {
		if v == bucketName {
			return true, nil
		}
	}
	return false, nil
}

func (ac *AwsClient) OsCreateBucket(region, bucketName string) error {
	s3cli, err := ac.getS3Client(region)
	if err != nil {
		return err
	}

	in := &s3.CreateBucketInput{
		Bucket: &bucketName,
	}
	_, err = s3cli.CreateBucket(in)
	return err
}

func (ac *AwsClient) OsDeleteBucket(region, bucketName string, deleteIfNotEmpty *bool) error {
	s3cli, err := ac.getS3Client(region)
	if err != nil {
		return err
	}

	in := &s3.DeleteBucketInput{
		Bucket: &bucketName,
	}
	_, err = s3cli.DeleteBucket(in)
	return err
}

func (ac *AwsClient) OsListBuckets(region string) ([]string, error) {
	s3cli, err := ac.getS3Client(region)
	if err != nil {
		return nil, err
	}

	out, err := s3cli.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}
	var results []string
	for _, v := range out.Buckets {
		results = append(results, *v.Name)
	}
	return results, nil
}

func (ac *AwsClient) OsListObjectKeys(region, bucketName string, keyPrefix *string, pageSize int, pageToken *string) ([]string, *string, error) {
	s3cli, err := ac.getS3Client(region)
	if err != nil {
		return nil, nil, err
	}

	if pageSize <= 0 {
		pageSize = 100
	}
	pageSizeI64 := (int64)(pageSize)

	in := &s3.ListObjectsV2Input{
		Bucket:  &bucketName,
		MaxKeys: &pageSizeI64,
	}
	if keyPrefix != nil {
		in.Prefix = keyPrefix
	}
	if pageToken != nil {
		in.ContinuationToken = pageToken
	}
	out, err := s3cli.ListObjectsV2(in)
	if err != nil {
		return nil, nil, err
	}
	var results []string
	resultNextPageToken := (*string)(nil)
	for _, v := range out.Contents {
		results = append(results, *v.Key)
	}
	if out.NextContinuationToken != nil && *out.NextContinuationToken != "" {
		*resultNextPageToken = *out.NextContinuationToken
	}
	return results, resultNextPageToken, nil
}

func (ac *AwsClient) OsUpsertObject(region, bucketName, blobId string, blob []byte) error {
	s3cli, err := ac.getS3Client(region)
	if err != nil {
		return err
	}

	size := int64(len(blob))
	in := &s3.UploadPartInput{
		Body:          bytes.NewReader(blob),
		Bucket:        &bucketName,
		ContentLength: &size,
		Key:           &blobId,
	}
	_, err = s3cli.UploadPart(in)
	return err
}

func (ac *AwsClient) OsGetObject(region, bucketName, blobId string) ([]byte, error) {
	s3cli, err := ac.getS3Client(region)
	if err != nil {
		return nil, err
	}

	in := &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &blobId,
	}
	out, err := s3cli.GetObject(in)
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()
	return io.ReadAll(out.Body)
}

func (ac *AwsClient) OsGetObjectSize(region, bucketName, blobId string) (*gvolume.Volume, error) {
	s3cli, err := ac.getS3Client(region)
	if err != nil {
		return nil, err
	}

	in := &s3.GetObjectAttributesInput{
		Bucket: &bucketName,
		Key:    &blobId,
	}
	out, err := s3cli.GetObjectAttributes(in)
	if err != nil {
		return nil, err
	}
	size := out.ObjectSize
	result := gvolume.FromByteSizeUint64(uint64(*size))
	return &result, nil
}

func (ac *AwsClient) OsDeleteObject(region, bucketName, blobId string) error {
	s3cli, err := ac.getS3Client(region)
	if err != nil {
		return err
	}

	in := &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &blobId,
	}
	_, err = s3cli.DeleteObject(in)
	return err
}

func (ac *AwsClient) OsRenameObject(region, bucketName, oldObjectKey, newObjectKey string) error {
	s3cli, err := ac.getS3Client(region)
	if err != nil {
		return err
	}

	in := &s3.CopyObjectInput{
		Bucket:     &bucketName,
		CopySource: &oldObjectKey,
		Key:        &newObjectKey,
	}
	_, err = s3cli.CopyObject(in)
	if err != nil {
		return err
	}

	in2 := &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &oldObjectKey,
	}
	_, err = s3cli.DeleteObject(in2)
	return err
}

func (ac *AwsClient) MqCreateQueue(region string, queue string, attr *QueueAttr) error {
	sqsClient, err := ac.getSqsClient(region)
	if err != nil {
		return err
	}

	in := &sqs.CreateQueueInput{QueueName: aws.String(queue)}
	if attr != nil {
		if attr.MsgMaxBytes != nil {
			in.Attributes["MaximumMessageSize"] = aws.String(strconv.FormatInt(int64(*attr.MsgMaxBytes), 10))
		}
		if attr.MsgRetentionSeconds != nil {
			in.Attributes["MessageRetentionPeriod"] = aws.String(strconv.FormatInt(int64(*attr.MsgRetentionSeconds), 10))
		}
		if attr.VisibilityTimeoutSeconds != nil {
			in.Attributes["VisibilityTimeout"] = aws.String(strconv.FormatInt(int64(*attr.VisibilityTimeoutSeconds), 10))
		}
	}
	out, err := sqsClient.CreateQueue(in)
	if err != nil {
		return err
	}
	ac.queues.Store(region+"."+queue, out.QueueUrl)
	return nil
}

func (ac *AwsClient) MqListQueues(region string, maxResults int, pageToken *string, queueNamePrefix *string) (queues []string, nextPageToken *string, err error) {
	in := &sqs.ListQueuesInput{
		MaxResults:      aws.Int64(int64(maxResults)),
		NextToken:       pageToken,
		QueueNamePrefix: queueNamePrefix,
	}
	sqsClient, err := ac.getSqsClient(region)
	if err != nil {
		return nil, nil, err
	}
	out, err := sqsClient.ListQueues(in)
	if err != nil {
		return nil, nil, err
	}
	for _, v := range out.QueueUrls {
		if v != nil {
			result, err := arn.Parse(*v)
			if err != nil {
				return nil, nil, err
			}
			if result.Region != region {
				continue
			}
			queues = append(queues, result.Resource)
			ac.queues.Store(region+"."+result.Resource, *v)
		}
	}

	return queues, out.NextToken, nil
}

func (ac *AwsClient) MqSend(region string, queue string, msg string) error {
	sqsClient, err := ac.getSqsClient(region)
	if err != nil {
		return err
	}
	queueURL, err := ac.getSqsQueueURL(region, queue)
	if err != nil {
		return err
	}

	in := &sqs.SendMessageInput{
		MessageBody: aws.String(msg),
		QueueUrl:    aws.String(queueURL),
	}
	_, err = sqsClient.SendMessage(in)
	return err
}

func (ac *AwsClient) MqReceive(region string, queue string, msgNum int, waitSeconds int, deleteAfterReceived bool) ([]Msg, error) {
	if msgNum <= 0 {
		msgNum = 1
	}
	if waitSeconds <= 0 {
		waitSeconds = 30
	}

	sqsClient, err := ac.getSqsClient(region)
	if err != nil {
		return nil, err
	}
	queueURL, err := ac.getSqsQueueURL(region, queue)
	if err != nil {
		return nil, err
	}
	//visibilityTimeout := int64(0)
	//if deleteAfterReceived {
	//	visibilityTimeout = 1000
	//}

	in := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: aws.Int64(int64(msgNum)),
		WaitTimeSeconds:     aws.Int64(int64(waitSeconds)),
	}

	out, err := sqsClient.ReceiveMessage(in)
	if err != nil {
		return nil, err
	}
	var results []Msg
	inBatch := &sqs.DeleteMessageBatchInput{
		QueueUrl: aws.String(queueURL),
	}
	for _, v := range out.Messages {
		timestamp := v.Attributes["SentTimestamp"]
		epochMillis, err := strconv.ParseInt(*timestamp, 10, 64)
		if err != nil {
			return nil, err
		}

		results = append(results, Msg{
			Id:          *v.MessageId,
			EnqueueTime: gtime.EpochMillisToTime(epochMillis),
			Data:        *v.Body,
		})
		inBatch.Entries = append(inBatch.Entries, &sqs.DeleteMessageBatchRequestEntry{
			Id:            v.MessageId,
			ReceiptHandle: v.ReceiptHandle,
		})
	}

	if deleteAfterReceived {
		return results, err
	}
	_, err = sqsClient.DeleteMessageBatch(inBatch)
	return results, err
}

func (ac *AwsClient) MqDeleteQueue(region string, queue string) error {
	sqsClient, err := ac.getSqsClient(region)
	if err != nil {
		return err
	}
	queueURL, err := ac.getSqsQueueURL(region, queue)
	if err != nil {
		return err
	}
	in := &sqs.DeleteQueueInput{
		QueueUrl: &queueURL,
	}
	_, err = sqsClient.DeleteQueue(in)
	return err
}

func (ac *AwsClient) Close() error {
	return nil
}

func (ac *AwsClient) getSession(region string) (*session.Session, error) {
	val, ok := ac.sessions.Load(region)
	if ok {
		return val.(*session.Session), nil
	}

	provider := &credentials.StaticProvider{Value: credentials.Value{
		AccessKeyID:     ac.accessKey,
		SecretAccessKey: ac.secretKey,
		//SessionToken:    "SESSION", // FIXME 这个参数啥意思
	}}
	sess, err := session.NewSession(&aws.Config{
		HTTPClient:  &http.Client{},
		Region:      aws.String(region),
		Credentials: credentials.NewCredentials(provider),
		MaxRetries:  aws.Int(-1),
	})
	if err != nil {
		return nil, err
	}
	ac.sessions.Store(region, sess)
	return sess, nil
}

func (ac *AwsClient) getEc2Client(region string) (*ec2.EC2, error) {
	sess, err := ac.getSession(region)
	if err != nil {
		return nil, err
	}

	return ec2.New(sess, &aws.Config{}), nil
}

func (ac *AwsClient) getS3Client(region string) (*s3.S3, error) {
	sess, err := ac.getSession(region)
	if err != nil {
		return nil, err
	}

	return s3.New(sess, &aws.Config{}), nil
}

func (ac *AwsClient) getSnsClient(region string) (*sns.SNS, error) {
	sess, err := ac.getSession(region)
	if err != nil {
		return nil, err
	}

	return sns.New(sess, &aws.Config{}), nil
}

func (ac *AwsClient) getSqsClient(region string) (*sqs.SQS, error) {
	sess, err := ac.getSession(region)
	if err != nil {
		return nil, err
	}

	return sqs.New(sess, &aws.Config{}), nil
}

func (ac *AwsClient) getSqsQueueURL(region string, queue string) (string, error) {
	queueURL, ok := ac.queues.Load(region + "." + queue)
	if ok {
		return queueURL.(string), nil
	}

	nextToken := (*string)(nil)
	for {
		_, outNextToken, err := ac.MqListQueues(region, 100, nextToken, &queue)
		if err != nil {
			return "", err
		}
		if outNextToken == nil {
			break
		} else {
			nextToken = outNextToken
			continue
		}
	}

	queueURL, ok = ac.queues.Load(region + "." + queue)
	if ok {
		return queueURL.(string), nil
	} else {
		return "", gerrors.ErrNotExist
	}

}
