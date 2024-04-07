package cloud

import (
	billing "cloud.google.com/go/billing/apiv1"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	storage "cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	"google.golang.org/protobuf/proto"
	"goutil/basic/gerrors"
	"goutil/container/gvolume"
	"goutil/sys/gtime"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

type (
	GcpClient struct {
		accessKey string
		secretKey string
		lan       bool
		sessions  sync.Map // map[region]session.Session
		topics    sync.Map // map[region.topic]topicARN
		queues    sync.Map // map[region.queue]queueURL
	}
)

const (
	defaultProject = "default-project-qml8394"
)

func newGcp(accessKey string, secretKey string, LAN bool) (*GcpClient, error) {
	return &GcpClient{
		accessKey: accessKey,
		secretKey: secretKey,
		lan:       LAN,
	}, nil
}

func (gc *GcpClient) GetBalance() (*Balance, error) {
	cli, err := gc.getBillingClient()
	if err != nil {
		return nil, err
	}
	info, err := cli.ListBillingAccounts(context.Background(), )
	if err != nil {
		return nil, err
	}
}

func (gc *GcpClient) ListLocations() (map[string][]string, error) {
	rgCli, err := gc.getRegionClient()
	if err != nil {
		return nil, err
	}

	iter := rgCli.List(context.Background(), nil)
	nextPageToken := (*string)(nil)

	result := map[string][]string{}
	for {
		info, errNext := iter.Next()
		if errors.Is(errNext, iterator.Done) { // FIXME
			break
		}
		if errNext != nil {
			return nil, errNext
		}
		result[*info.Name] = info.Zones

		if iter.PageInfo() != nil && iter.PageInfo().Token != "" {
			*nextPageToken = iter.PageInfo().Token
		}
		if nextPageToken == nil {
			break
		}
	}
	return result, nil
}

func (gc *GcpClient) VmListOnDemandSpecs(region string) (*InstanceSpecList, error) {

}

func (gc *GcpClient) VmListSpotSpecs(region string) (*InstanceSpecList, error) {
	mtCli, err := gc.getMtClient(region)
	if err != nil {
		return nil, err
	}

	nextPageToken := (*string)(nil)
	valResult := &InstanceSpecList{}
	for {
		req := &computepb.ListMachineTypesRequest{
			Filter:               nil,
			MaxResults:           nil,
			OrderBy:              nil,
			PageToken:            nil,
			Project:              "",
			ReturnPartialSuccess: nil,
			Zone:                 "",
		}
		iter := mtCli.List(context.Background(), req)
		for {
			info, errNext := iter.Next()
			if errors.Is(errNext, iterator.Done) { // FIXME
				break
			}
			if errNext != nil {
				return nil, errNext
			}
			item := InstanceSpec{
				RegionId:         region,
				Id:        strconv.FormatUint(*info.Id, 10),
				IsCredit:         false,
				Currency:         "",
				LogicalCpuNum:    ,
				MemoryVolume:     0,
				AvailableZoneIds: nil,
				OnDemandPrices:   nil,
				SpotPricePerHour: nil,
			}
			valResult.Specs = append(valResult.Specs, item)
		}

		if iter.PageInfo() != nil && iter.PageInfo().Token != "" {
			*nextPageToken = iter.PageInfo().Token
		}
		if nextPageToken == nil {
			break
		}
	}

	return valResult, nil
}

func (gc *GcpClient) VmListImages(region string) ([]SysImage, error) {
	imgCli, err := gc.getImageClient(region)
	if err != nil {
		return nil, err
	}

	nextPageToken := (*string)(nil)
	var valResult []SysImage = nil
	for {
		req := &computepb.ListImagesRequest{
			Filter:               nil,
			MaxResults:           100,
			OrderBy:              nil,
			PageToken:            nextPageToken,
			Project:              defaultProject,
			ReturnPartialSuccess: nil,
		}
		iter := imgCli.List(context.Background(), req)
		for {
			info, errNext := iter.Next()
			if errors.Is(errNext, iterator.Done) { // FIXME
				break
			}
			if errNext != nil {
				return nil, errNext
			}
			item := SysImage{
				Id:        strconv.FormatUint(*info.Id, 10),
				Name:      *info.Name,
				OS:        *info.Family,
				Distro:    *info.Kind,
				Arch:      *info.Architecture,
				Available: true,
			}
			valResult = append(valResult, item)
		}

		if iter.PageInfo() != nil && iter.PageInfo().Token != "" {
			*nextPageToken = iter.PageInfo().Token
		}
		if nextPageToken == nil {
			break
		}
	}

	return valResult, nil
}

func (gc *GcpClient) VmListInstances(region string) ([]InstanceInfo, error) {
	instCli, err := gc.getVmClient(region)
	if err != nil {
		return nil, err
	}

	nextPageToken := (*string)(nil)
	var valResult []InstanceInfo = nil
	for {
		req := &computepb.ListInstancesRequest{
			Filter:               nil,
			MaxResults:           100,
			OrderBy:              nil,
			PageToken:            nextPageToken,
			Project:              defaultProject,
			ReturnPartialSuccess: nil,
			Zone:                 "",
		}
		iter := instCli.List(context.Background(), req)
		for {
			info, errNext := iter.Next()
			if errors.Is(errNext, iterator.Done) { // FIXME
				break
			}
			if errNext != nil {
				return nil, errNext
			}
			var privateIPs []string
			var publicIPs []string
			for _, itf := range info.NetworkInterfaces {
				privateIPs = append(privateIPs, *itf.NetworkIP)
			}
			createionTime, err := gtime.ParseTimeStringStrict(*info.CreationTimestamp)
			if err != nil {
				return nil, err
			}
			item := InstanceInfo{
				Id:                 strconv.FormatUint(*info.Id, 10),
				Name:               *info.Name,
				Specs:              *info.MachineType,
				InstanceChargeType: "",
				InternetChargeType: "",
				SpotPriceLimit:     0,
				SpotStrategy:       "",
				SpotStartTime:      time.Time{},
				NetworkType:        "",
				PrivateIPs:         privateIPs,
				PublicIPs:          publicIPs,
				RegionId:           region,
				ZoneId:             "",
				SecurityGroupIds:   ,
				ImageId:            *info.SourceMachineImage,
				Status:             "",
				KeyPairName:        "",
				PhysicalCpuNum:     ,
				LogicalCpuNum:      0,
				GpuNum:             0,
				CreationTime:       createionTime,
				AutoReleaseTime:    time.Time{},
				MemorySize:         0,
				SysImageId:         "",
				SysImageName:       "",
				SysImageOS:         "",
			}
			valResult = append(valResult, item)
		}

		if iter.PageInfo() != nil && iter.PageInfo().Token != "" {
			*nextPageToken = iter.PageInfo().Token
		}
		if nextPageToken == nil {
			break
		}
	}

	return valResult, nil
}

func (gc *GcpClient) VmListSecurityGroups(region string) ([]SecurityGroup, error) {
	sgCli, err := gc.getSgClient(region)
	if err != nil {
		return nil, err
	}

	nextPageToken := (*string)(nil)
	var results []SecurityGroup
	for  {
		req := &computepb.ListSecurityPoliciesRequest{
			Filter:               nil,
			MaxResults:           100,
			OrderBy:              nil,
			PageToken:            nextPageToken,
			Project:              defaultProject,
			ReturnPartialSuccess: nil,
		}
		iter := sgCli.List(context.Background(), req)
		for {
			info, errNext := iter.Next()
			if errors.Is(errNext, iterator.Done) { // FIXME
				break
			}
			if errNext != nil {
				return nil, errNext
			}

			item := SecurityGroup{
				Id:          strconv.FormatUint(*info.Id, 10),
				Name:        *info.Name,
				Permissions: nil,
			}
			for _, v := range info.Rules {
				item.Permissions = append(item.Permissions, SecurityPermission{
					Description:  *v.Description,
					Direction:    ,
					Protocol:     ,
					SrcPortRange: [2]int{},
					SrcCidrIP:    "",
					DstPortRange: [2]int{},
					DstCidrIP:    "",
				})
			}
			results = append(results, item)
		}

		if iter.PageInfo() != nil && iter.PageInfo().Token != "" {
			*nextPageToken = iter.PageInfo().Token
		}
		if nextPageToken == nil {
			break
		}
	}

	return results, nil
}

func (gc *GcpClient) VmCreateSecurityGroup(region string, sg SecurityGroup) (string, error) {

}

func (gc *GcpClient) VmDeleteSecurityGroup(region, securityGroupId string) error {

}

func (gc *GcpClient) VmListSwitches(region, zoneId string) ([]string, error) {

}

// VmCreateInstances sends an instance creation request to the Compute Engine API and waits for it to complete.
// reference and sample code: https://cloud.google.com/compute/docs/instances/create-start-instance
func (gc *GcpClient) VmCreateInstances(region string, tmpl InstanceCreationTmpl, num int) ([]string, error) {
	//func createInstance(projectID, zone, instanceName, machineType, sourceImage, networkName string) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// machineType := "n1-standard-1"
	// sourceImage := "projects/debian-cloud/global/images/family/debian-10"
	// networkName := "global/networks/default"

	ctx := context.Background()
	instCli, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	/*req := &computepb.InsertInstanceRequest{
		Project: projectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name: proto.String(instanceName),
			Disks: []*computepb.AttachedDisk{
				{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskSizeGb:  proto.Int64(10),
						SourceImage: proto.String(sourceImage),
					},
					AutoDelete: proto.Bool(true),
					Boot:       proto.Bool(true),
					Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
				},
			},
			MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType)),
			NetworkInterfaces: []*computepb.NetworkInterface{
				{
					Name: proto.String(networkName),
				},
			},
		},
	}*/
	numInt64 := int64(num)
	namePattern := "inst-####"
	req := &computepb.BulkInsertInstanceRequest{
		Project: defaultProject,
		Zone:    zone,
		BulkInsertInstanceResourceResource: &computepb.BulkInsertInstanceResource{
			Count: &numInt64,
			InstanceProperties: &computepb.InstanceProperties{
				AdvancedMachineFeatures:    nil,
				CanIpForward:               nil,
				ConfidentialInstanceConfig: nil,
				Description:                nil,
				Disks: []*computepb.AttachedDisk{
					{
						InitializeParams: &computepb.AttachedDiskInitializeParams{
							DiskSizeGb:  proto.Int64(10),
							SourceImage: proto.String(sourceImage),
						},
						AutoDelete: proto.Bool(true),
						Boot:       proto.Bool(true),
						Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
					},
				},
				GuestAccelerators:       nil,
				KeyRevocationActionType: nil,
				Labels:                  nil,
				MachineType:             proto.String(fmt.Sprintf("zones/%s/machineTypes/%s", zone, tmpl.Specs)),
				Metadata:                nil,
				MinCpuPlatform:          nil,
				NetworkInterfaces: []*computepb.NetworkInterface{
					{
						Name: proto.String(networkName),
					},
				},
				NetworkPerformanceConfig: nil,
				PrivateIpv6GoogleAccess:  nil,
				ReservationAffinity:      nil,
				ResourceManagerTags:      nil,
				ResourcePolicies:         nil,
				Scheduling:               nil,
				ServiceAccounts:          nil,
				ShieldedInstanceConfig:   nil,
				Tags:                     nil,
			},
			LocationPolicy:         nil,
			MinCount:               &numInt64,
			NamePattern:            &namePattern,
			PerInstanceProperties:  nil,
			SourceInstanceTemplate: nil,
		},
	}

	op, err := instCli.BulkInsert(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("unable to create instance: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return nil, fmt.Errorf("unable to wait for the operation: %w", err)
	}
	return nil
}

func (gc *GcpClient) VmStartInstances(region string, instanceIds []string) error {
	ctx := context.Background()
	instCli, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}

	errResult := error(nil)
	errMu := sync.RWMutex{}

	wg := &sync.WaitGroup{}
	for _, v := range instanceIds {
		req := &computepb.StartInstanceRequest{
			Instance:  v,
			Project:   defaultProject,
			Zone:      ,
		}
		op, errDel := instCli.Start(ctx, req)
		if err != nil {
			errMu.Lock()
			errResult = gerrors.Join(errResult, errDel)
			errMu.Unlock()
			break
		} else {
			wg.Add(1)
			go func() {
				defer wg.Add(-1)
				errWait := op.Wait(ctx)
				if err != nil {
					errMu.Lock()
					errResult = gerrors.Join(errResult, errWait)
					errMu.Unlock()
				}
			}()
			continue
		}
	}
	wg.Wait()

	return errResult
}

func (gc *GcpClient) VmDeleteInstances(region string, instanceIds []string, force bool) error {
	ctx := context.Background()
	instCli, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}

	errResult := error(nil)
	errMu := sync.RWMutex{}

	wg := &sync.WaitGroup{}
	for _, v := range instanceIds {
		req := &computepb.DeleteInstanceRequest{
			Instance:  v,
			Project:   defaultProject,
			Zone:      ,
		}
		op, errDel := instCli.Delete(ctx, req)
		if err != nil {
			errMu.Lock()
			errResult = gerrors.Join(errResult, errDel)
			errMu.Unlock()
			break
		} else {
			wg.Add(1)
			go func() {
				defer wg.Add(-1)
				errWait := op.Wait(ctx)
				if err != nil {
					errMu.Lock()
					errResult = gerrors.Join(errResult, errWait)
					errMu.Unlock()
				}
			}()
			continue
		}
	}
	wg.Wait()

	return errResult
}

func (gc *GcpClient) OsIsBucketExist(region, bucketName string) (bool, error) {
	buckets, err := gc.OsListBuckets(region)
	if err != nil {
		return false, err
	}
	for _, v := range buckets {
		if v == bucketName {
			return true, nil
		}
	}
	return false, nil
}

func (gc *GcpClient) OsCreateBucket(region, bucketName string) error {
	client, err := gc.getCloudStorageClient(region)
	if err != nil {
		return err
	}
	attr := &storage.BucketAttrs{
		Name:     bucketName,
		Location: region,
	}
	return client.Bucket(bucketName).Create(context.Background(), defaultProject, attr)
}

func (gc *GcpClient) OsDeleteBucket(region, bucketName string, deleteIfNotEmpty *bool) error {
	client, err := gc.getCloudStorageClient(region)
	if err != nil {
		return err
	}
	return client.Bucket(bucketName).Delete(context.Background())
}

func (gc *GcpClient) OsListBuckets(region string) ([]string, error) {
	client, err := gc.getCloudStorageClient(region)
	if err != nil {
		return nil, err
	}
	iter := client.Buckets(context.Background(), defaultProject)
	var results []string
	for {
		attr, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		results = append(results, attr.Name)
	}
	return results, nil
}

func (gc *GcpClient) OsListObjectKeys(region, bucketName string, keyPrefix *string, pageSize int, pageToken *string) ([]string, *string, error) {
	if pageSize <= 0 {
		pageSize = 100
	}

	client, err := gc.getCloudStorageClient(region)
	if err != nil {
		return nil, nil, err
	}
	q := &storage.Query{}
	if keyPrefix != nil {
		q.Prefix = *keyPrefix
	}
	iter := client.Bucket(bucketName).Objects(context.Background(), q)
	pi := iter.PageInfo()
	pi.MaxSize = pageSize
	if pageToken != nil {
		pi.Token = *pageToken
	}
	var results []string
	for {
		attr, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		results = append(results, attr.Name)
	}
	// reference : https://github.com/GoogleCloudPlatform/gcsfuse/blob/5f8369a1da7c34293c3aee7dc5b8795b088449a5/internal/storage/bucket_handle.go#L251
	nextPageToken := (*string)(nil)
	if iter.PageInfo() != nil && iter.PageInfo().Token != "" {
		*nextPageToken = iter.PageInfo().Token
	}
	return results, nextPageToken, nil
}

func (gc *GcpClient) OsUpsertObject(region, bucketName, blobId string, blob []byte) error {
	client, err := gc.getCloudStorageClient(region)
	if err != nil {
		return err
	}
	w := client.Bucket(bucketName).Object(blobId).NewWriter(context.Background())
	defer w.Close()
	wTotal := 0
	for wTotal < len(blob) {
		wOnce, err := w.Write(blob[wTotal:])
		if err != nil {
			return err
		}
		wTotal += wOnce
	}
	return nil
}

func (gc *GcpClient) OsGetObject(region, bucketName, blobId string) ([]byte, error) {
	client, err := gc.getCloudStorageClient(region)
	if err != nil {
		return nil, err
	}
	r, err := client.Bucket(bucketName).Object(blobId).NewRangeReader(context.Background(), 0, -1 /*toTheEnd*/)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

func (gc *GcpClient) OsGetObjectSize(region, bucketName, blobId string) (*gvolume.Volume, error) {
	client, err := gc.getCloudStorageClient(region)
	if err != nil {
		return nil, err
	}
	attrs, err := client.Bucket(bucketName).Object(blobId).Attrs(context.Background())
	if err != nil {
		return nil, err
	}
	size := gvolume.FromByteSizeUint64((uint64(attrs.Size)))
	return &size, nil
}

func (gc *GcpClient) OsDeleteObject(region, bucketName, blobId string) error {
	client, err := gc.getCloudStorageClient(region)
	if err != nil {
		return err
	}
	return client.Bucket(bucketName).Object(blobId).Delete(context.Background())
}

func (gc *GcpClient) OsRenameObject(region, bucketName, oldObjectKey, newObjectKey string) error {
	client, err := gc.getCloudStorageClient(region)
	if err != nil {
		return err
	}
	oldObj := client.Bucket(bucketName).Object(oldObjectKey)
	newObj := client.Bucket(bucketName).Object(newObjectKey)
	copier := newObj.CopierFrom(oldObj)
	_, err = copier.Run(context.Background())
	if err != nil {
		return err
	}
	return oldObj.Delete(context.Background())
}

func (gc *GcpClient) MqCreateQueue(region string, queue string, attr *QueueAttr) error {
	client, err := gc.getCloudTaskClient(region)
	if err != nil {
		return err
	}
	req := &taskspb.CreateQueueRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", defaultProject, region),
		Queue: &taskspb.Queue{
			Name: fmt.Sprintf("projects/%s/locations/%s/queues/%s", defaultProject, region, queue),
		},
	}
	_, err = client.CreateQueue(context.Background(), req)
	return err
}

func (gc *GcpClient) MqListQueues(region string, maxResults int, pageToken *string, queueNamePrefix *string) (queues []string, nextPageToken *string, err error) {
	client, err := gc.getCloudTaskClient(region)
	if err != nil {
		return nil, nil, err
	}
	req := &taskspb.ListQueuesRequest{
		Parent:   fmt.Sprintf("projects/%s/locations/%s", defaultProject, region),
		PageSize: int32(maxResults),
	}
	if pageToken != nil {
		req.PageToken = *pageToken
	}
	if queueNamePrefix != nil {
		req.Filter = fmt.Sprintf("Name:%s*", *queueNamePrefix) // FIXME
	}
	resp, err := client.ListQueues(context.Background(), req)
	if err != nil {
		return nil, nil, err
	}
	for _, v := range resp.Queues {
		splitName := strings.Split(v.Name, "/")
		queues = append(queues, splitName[len(splitName)-1])
	}
	nextPageToken = &resp.NextPageToken
	return queues, nextPageToken, nil
}

func (gc *GcpClient) MqSend(region string, queue string, msg string) error {
	client, err := gc.getCloudTaskClient(region)
	if err != nil {
		return err
	}

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", defaultProject, region, queue)
	// Build the Task payload.
	// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#CreateTaskRequest
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#HttpRequest
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        "",
				},
			},
		},
	}
	// Add a payload message if one is present.
	req.Task.GetHttpRequest().Body = []byte(msg)

	_, err = client.CreateTask(context.Background(), req)
	if err != nil {
		return fmt.Errorf("cloudtasks.CreateTask: %w", err)
	}
	return nil
}

func (gc *GcpClient) MqReceive(region string, queue string, msgNum int, waitSeconds int, deleteAfterReceived bool) ([]Msg, error) {
	client, err := gc.getCloudTaskClient(region)
	if err != nil {
		return nil, err
	}

	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", defaultProject, region, queue)
	req := &taskspb.GetTaskRequest{
		Name: queuePath,
	}

	var results []Msg
	var idToDel []string
	for i := 0; i < msgNum; i++ {
		resp, err := client.GetTask(context.Background(), req)
		if err != nil {
			return nil, fmt.Errorf("cloudtasks.GetTask: %w", err)
		}
		splitName := strings.Split(resp.Name, "/")
		results = append(results, Msg{
			Id:          splitName[len(splitName)-1],
			EnqueueTime: resp.CreateTime.AsTime(),
			Data:        resp.String(),
		})
		idToDel = append(idToDel, splitName[len(splitName)-1])
	}

	if deleteAfterReceived {
		for _, v := range idToDel {
			req := &taskspb.DeleteTaskRequest{Name: fmt.Sprintf(`projects/%s/locations/%s/queues/%s/tasks/%s`, defaultProject, region, queue, v)}
			_, err = client.DeleteTask(context.Background(), req)
			if err != nil {
				return results, err
			}
		}
	}

	return results, nil
}

func (gc *GcpClient) MqDeleteQueue(region string, queue string) error {
	client, err := gc.getCloudTaskClient(region)
	if err != nil {
		return err
	}

	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", defaultProject, region, queue)
	req := &taskspb.DeleteQueueRequest{
		Name: queuePath,
	}

	_, err = client.DeleteQueue(context.Background(), req)
	return err
}

func (gc *GcpClient) Close() error {

}

// DefaultAuthScopes reports the default set of authentication scopes to use with this package.
func DefaultAuthScopes() []string {
	return []string{
		"https://www.googleapis.com/auth/cloud-platform",
	}
}

func defaultClientOptions() []option.ClientOption {
	return []option.ClientOption{
		option.WithEndpoint("cloudtasks.googleapis.com:443"),
		option.WithScopes(DefaultAuthScopes()...),
	}
}

func (gc *GcpClient) getCloudTaskClient(region string) (taskspb.CloudTasksClient, error) {
	conn, err := transport.DialGRPC(context.Background(), append(defaultClientOptions())...)
	if err != nil {
		return nil, err
	}
	return taskspb.NewCloudTasksClient(conn), nil
}

func (gc *GcpClient) getCloudStorageClient(region string) (*storage.Client, error) {
	storage.NewClient(context.Background(), option.WithoutAuthentication())
}

func (gc *GcpClient) getVmClient(region string) (*compute.InstancesClient, error) {
	ctx := context.Background()
	return compute.NewInstancesRESTClient(ctx)
}

func (gc *GcpClient) getImageClient(region string) (*compute.ImagesClient, error) {
	ctx := context.Background()
	return compute.NewImagesRESTClient(ctx)
}

func (gc *GcpClient) getSgClient(region string) (*compute.SecurityPoliciesClient, error) {
	ctx := context.Background()
	return compute.NewSecurityPoliciesRESTClient(ctx)
}

func (gc *GcpClient) getMtClient(region string) (*compute.MachineTypesClient, error) {
	ctx := context.Background()
	return compute.NewMachineTypesRESTClient(ctx)
}

func (gc *GcpClient) getRegionClient() (*compute.RegionsClient, error) {
	ctx := context.Background()
	return compute.NewRegionsRESTClient(ctx)
}

func (gc *GcpClient) getBillingClient() (*billing.CloudBillingClient, error) {
	ctx := context.Background()
	return billing.NewCloudBillingClient(ctx)
}