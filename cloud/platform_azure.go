package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/billing/armbilling"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"goutil/basic/gerrors"
	"goutil/container/gvolume"
	"goutil/encoding/gjson"
	"io"
	"net/url"
	"os"
	"sync"
)

type (
	AzureClient struct {
		accessKey  string
		secretKey  string
		lan        bool
		sessions   sync.Map // map[region]session.Session
		topics     sync.Map // map[region.topic]topicARN
		queues     sync.Map // map[region.queue]queueURL
		queueAttrs sync.Map // map[region.queue]QueueAttr
		aqsList    map[string]*azqueue.QueueClient
		aqsListMu  sync.RWMutex
		azqsc      *azqueue.ServiceClient
	}
)

const (
	azQueueEndpointExp = "https://%s.queue.core.windows.net"
)

func newAzure(accessKey string, secretKey string, LAN bool) (*AzureClient, error) {
	return &AzureClient{
		accessKey: accessKey,
		secretKey: secretKey,
		lan:       LAN,
		aqsList:   map[string]*azqueue.QueueClient{},
	}, nil
}

func (ac *AzureClient) GetBalance() (*Balance, error) {
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	clientFactory, err := armsubscriptions.NewClientFactory(cred, nil)
	if err != nil {
		return nil, err
	}
	cli, err := armbilling.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}
	cli.NewAccountsClient().Get()
}

func (ac *AzureClient) ListLocations() (map[string][]string, error) {
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	clientFactory, err := armsubscriptions.NewClientFactory(cred, nil)
	if err != nil {
		return nil, err
	}
	res := clientFactory.NewClient().NewListLocationsPager(subscriptionId, nil)
	result := map[string][]string{}
	for res.More() {
		item, err := res.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		for _, v := range item.Value {
			var zones []string
			for _, z := range v.AvailabilityZoneMappings {
				zones = append(zones, *z.LogicalZone)
			}
			result[*v.ID] = zones
		}
	}
	return result, nil
}

func (ac *AzureClient) VmListOnDemandSpecs(region string) (*InstanceSpecList, error) {

}

func (ac *AzureClient) VmListSpotSpecs(region string) (*InstanceSpecList, error) {
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	cli, err := armcompute.NewVirtualMachineSizesClient(subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}
	cli.NewListPager(region, nil)
}

func (ac *AzureClient) VmListImages(region string) ([]SysImage, error) {
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	imageClientFactory, err := armcompute.NewImagesClient(subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}
	imageClientFactory.NewListPager(&armcompute.ImagesClientListOptions{})
}

func (ac *AzureClient) VmListInstances(region string) ([]InstanceInfo, error) {
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	computeClientFactory, err := armcompute.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}

	req := &armcompute.VirtualMachinesClientListOptions{
		Expand: nil,
		Filter: nil,
	}
	computeClientFactory.NewVirtualMachinesClient().NewListPager(defaultProject, req)
}

func (ac *AzureClient) VmListSecurityGroups(region string) ([]SecurityGroup, error) {
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	networkClientFactory, err := armnetwork.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}

	networkClientFactory.NewSecurityGroupsClient().NewListAllPager(&armnetwork.SecurityGroupsClientListAllOptions{})
}

func (ac *AzureClient) VmCreateSecurityGroup(region string, sg SecurityGroup) (string, error) {
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	networkClientFactory, err := armnetwork.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}
	req := &armnetwork.SecurityGroupsClientBeginCreateOrUpdateOptions{}
	azSg := armnetwork.SecurityGroup{
		ID:         nil,
		Location:   nil,
		Properties: nil,
		Tags:       nil,
		Etag:       nil,
		Name:       nil,
		Type:       nil,
	}
	networkClientFactory.NewSecurityGroupsClient().BeginCreateOrUpdate(context.Background(), defaultProject, sg.Name, azSg, req)
}

func (ac *AzureClient) VmDeleteSecurityGroup(region, securityGroupId string) error {
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}

	networkClientFactory, err := armnetwork.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return err
	}
	networkClientFactory.NewSecurityGroupsClient().BeginDelete(context.Background(), defaultProject, securityGroupId, nil)
}

func (ac *AzureClient) VmListSwitches(region, zoneId string) ([]string, error) {

}

// reference: https://github.com/MicrosoftDocs/azure-docs/blob/main/articles/public-multi-access-edge-compute-mec/tutorial-create-vm-using-go-sdk.md#provision-a-virtual-machine
func (ac *AzureClient) VmCreateInstances(region string, tmpl InstanceCreationTmpl, num int) ([]string, error) {

	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	// client factory
	resourcesClientFactory, err := armresources.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}

	computeClientFactory, err := armcompute.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}

	networkClientFactory, err := armnetwork.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}

	// Step 1: Provision a resource group
	_, err = resourcesClientFactory.NewResourceGroupsClient().CreateOrUpdate(
		context.Background(),
		"<resourceGroupName>",
		armresources.ResourceGroup{
			Location: to.Ptr("westus"),
		},
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Step 2: Provision a virtual network
	virtualNetworksClientCreateOrUpdateResponsePoller, err := networkClientFactory.NewVirtualNetworksClient().BeginCreateOrUpdate(
		context.Background(),
		"<resourceGroupName>",
		"<virtualNetworkName>",
		armnetwork.VirtualNetwork{
			Location: to.Ptr("westus"),
			ExtendedLocation: &armnetwork.ExtendedLocation{
				Name: to.Ptr("<edgezoneid>"),
				Type: to.Ptr(armnetwork.ExtendedLocationTypesEdgeZone),
			},
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						to.Ptr("10.0.0.0/16"),
					},
				},
				Subnets: []*armnetwork.Subnet{
					{
						Name: to.Ptr("test-1"),
						Properties: &armnetwork.SubnetPropertiesFormat{
							AddressPrefix: to.Ptr("10.0.0.0/24"),
						},
					},
				},
			},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	virtualNetworksClientCreateOrUpdateResponse, err := virtualNetworksClientCreateOrUpdateResponsePoller.PollUntilDone(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	subnetID := *virtualNetworksClientCreateOrUpdateResponse.Properties.Subnets[0].ID

	// Step 3: Provision an IP address
	publicIPAddressesClientCreateOrUpdateResponsePoller, err := networkClientFactory.NewPublicIPAddressesClient().BeginCreateOrUpdate(
		context.Background(),
		"<resourceGroupName>",
		"<publicIPName>",
		armnetwork.PublicIPAddress{
			Name:     to.Ptr("<publicIPName>"),
			Location: to.Ptr("westus"),
			ExtendedLocation: &armnetwork.ExtendedLocation{
				Name: to.Ptr("<edgezoneid>"),
				Type: to.Ptr(armnetwork.ExtendedLocationTypesEdgeZone),
			},
			SKU: &armnetwork.PublicIPAddressSKU{
				Name: to.Ptr(armnetwork.PublicIPAddressSKUNameStandard),
			},
			Properties: &armnetwork.PublicIPAddressPropertiesFormat{
				PublicIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodStatic),
			},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	publicIPAddressesClientCreateOrUpdateResponse, err := publicIPAddressesClientCreateOrUpdateResponsePoller.PollUntilDone(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	// Step 4: Provision the network interface client
	interfacesClientCreateOrUpdateResponsePoller, err := networkClientFactory.NewInterfacesClient().BeginCreateOrUpdate(
		context.Background(),
		"<resourceGroupName>",
		"<networkInterfaceName>",
		armnetwork.Interface{
			Location: to.Ptr("westus"),
			ExtendedLocation: &armnetwork.ExtendedLocation{
				Name: to.Ptr("<edgezoneid>"),
				Type: to.Ptr(armnetwork.ExtendedLocationTypesEdgeZone),
			},
			Properties: &armnetwork.InterfacePropertiesFormat{
				EnableAcceleratedNetworking: to.Ptr(true),
				IPConfigurations: []*armnetwork.InterfaceIPConfiguration{
					{
						Name: to.Ptr("<ipConfigurationName>"),
						Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
							Subnet: &armnetwork.Subnet{
								ID: to.Ptr(subnetID),
							},
							PublicIPAddress: &armnetwork.PublicIPAddress{
								ID: publicIPAddressesClientCreateOrUpdateResponse.ID,
							},
						},
					},
				},
			},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	interfacesClientCreateOrUpdateResponse, err := interfacesClientCreateOrUpdateResponsePoller.PollUntilDone(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	// Step 5: Provision the virtual machine
	virtualMachinesClientCreateOrUpdateResponsePoller, err := computeClientFactory.NewVirtualMachinesClient().BeginCreateOrUpdate(
		context.Background(),
		"<resourceGroupName>",
		"<vmName>",
		armcompute.VirtualMachine{
			Location: to.Ptr("westus"),
			ExtendedLocation: &armcompute.ExtendedLocation{
				Name: to.Ptr("<edgezoneid>"),
				Type: to.Ptr(armcompute.ExtendedLocationTypesEdgeZone),
			},
			Properties: &armcompute.VirtualMachineProperties{
				StorageProfile: &armcompute.StorageProfile{
					ImageReference: &armcompute.ImageReference{
						Publisher: to.Ptr("<publisher>"),
						Offer:     to.Ptr("<offer>"),
						SKU:       to.Ptr("<sku>"),
						Version:   to.Ptr("<version>"),
					},
				},
				HardwareProfile: &armcompute.HardwareProfile{
					VMSize: to.Ptr(armcompute.VirtualMachineSizeTypesStandardD2SV3),
				},
				OSProfile: &armcompute.OSProfile{
					ComputerName:  to.Ptr("<computerName>"),
					AdminUsername: to.Ptr("<adminUsername>"),
					AdminPassword: to.Ptr("<adminPassword>"),
				},
				NetworkProfile: &armcompute.NetworkProfile{
					NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
						{
							ID: interfacesClientCreateOrUpdateResponse.ID,
							Properties: &armcompute.NetworkInterfaceReferenceProperties{
								Primary: to.Ptr(true),
							},
						},
					},
				},
			},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	_, err = virtualMachinesClientCreateOrUpdateResponsePoller.PollUntilDone(context.Background(), nil)
	if err != nil {
		return nil, err
	}
}

func (ac *AzureClient) VmStartInstances(region string, instanceIds []string) error {
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}

	computeClientFactory, err := armcompute.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return err
	}

	req := &armcompute.VirtualMachinesClientBeginStartOptions{}
	for _, v := range instanceIds {
		computeClientFactory.NewVirtualMachinesClient().BeginStart(context.Background(), defaultProject, v, req)
	}
}

func (ac *AzureClient) VmDeleteInstances(region string, instanceIds []string, force bool) error {
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}

	computeClientFactory, err := armcompute.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return err
	}

	req := &armcompute.VirtualMachinesClientBeginDeleteOptions{ForceDeletion: &force}
	for _, v := range instanceIds {
		computeClientFactory.NewVirtualMachinesClient().BeginDelete(context.Background(), defaultProject, v, req)
	}
}

func (ac *AzureClient) OsIsBucketExist(region, bucketName string) (bool, error) {

}

func (ac *AzureClient) OsCreateBucket(region, bucketName string) error {
	bucketClient, err := ac.getBucket(bucketName)
	if err != nil {
		return err
	}
	_, err = bucketClient.Create(context.Background(), nil, azblob.PublicAccessContainer)
	return err
}

func (ac *AzureClient) OsDeleteBucket(region, bucketName string, deleteIfNotEmpty *bool) error {
	bucketClient, err := ac.getBucket(bucketName)
	if err != nil {
		return err
	}
	_, err = bucketClient.Delete(context.Background(), azblob.ContainerAccessConditions{})
	return err
}

func (ac *AzureClient) OsListBuckets(region string) ([]string, error) {
	credential, err := azblob.NewSharedKeyCredential(ac.accessKey, ac.secretKey)
	if err != nil {
		return nil, err
	}
	cURL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", ac.accessKey, bucket))
	if err != nil {
		return nil, err
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	s := azblob.NewServiceURL(*cURL, p)

	nextPageToken := (*string)(nil)
	var results []string
	for {
		resp, err := s.ListContainersSegment(context.Background(), azblob.Marker{Val: nextPageToken}, azblob.ListContainersSegmentOptions{MaxResults: 500})
		if err != nil {
			return nil, err
		}
		for _, v := range resp.ContainerItems {
			results = append(results, v.Name)
		}
		nextPageToken = resp.NextMarker.Val
		if nextPageToken == nil {
			break
		}
	}
	return results, nil
}

func (ac *AzureClient) OsListObjectKeys(region, bucketName string, keyPrefix *string, pageSize int, pageToken *string) ([]string, *string, error) {
	if pageSize <= 0 {
		pageSize = 100
	}

	bucketClient, err := ac.getBucket(bucketName)
	if err != nil {
		return nil, nil, err
	}
	opts := azblob.ListBlobsSegmentOptions{
		MaxResults: int32(pageSize),
	}
	if keyPrefix != nil {
		opts.Prefix = *keyPrefix
	}
	marker := azblob.Marker{}
	if pageToken != nil {
		marker.Val = pageToken
	}
	resp, err := bucketClient.ListBlobsHierarchySegment(context.Background(), marker, "", opts)
	if err != nil {
		return nil, nil, err
	}

	var results []string
	for _, v := range resp.Segment.BlobItems {
		results = append(results, v.Name)
	}
	return results, resp.NextMarker.Val, nil
}

func (ac *AzureClient) OsUpsertObject(region, bucketName, blobId string, blob []byte) error {
	bucketClient, err := ac.getBucket(bucketName)
	if err != nil {
		return err
	}
	r := bytes.NewReader(blob)
	_, err = bucketClient.NewBlockBlobURL(blobId).Upload(context.Background(), r, azblob.BlobHTTPHeaders{}, azblob.Metadata{}, azblob.BlobAccessConditions{}, azblob.AccessTierNone, azblob.BlobTagsMap{}, azblob.ClientProvidedKeyOptions{}, azblob.ImmutabilityPolicyOptions{})
	return err
}

func (ac *AzureClient) OsGetObject(region, bucketName, blobId string) ([]byte, error) {
	bucketClient, err := ac.getBucket(bucketName)
	if err != nil {
		return nil, err
	}
	resp, err := bucketClient.NewBlobURL(blobId).Download(context.Background(), 0, 0 /*toTheEnd*/, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return nil, err
	}
	r := resp.Body(azblob.RetryReaderOptions{})
	defer r.Close()
	return io.ReadAll(r)
}

func (ac *AzureClient) OsGetObjectSize(region, bucketName, blobId string) (*gvolume.Volume, error) {
	bucketClient, err := ac.getBucket(bucketName)
	if err != nil {
		return nil, err
	}
	resp, err := bucketClient.NewBlobURL(blobId).GetProperties(context.Background(), azblob.BlobAccessConditions{}, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return nil, err
	}
	result := gvolume.FromByteSizeUint64(uint64(resp.ContentLength()))
	return &result, nil
}

func (ac *AzureClient) OsDeleteObject(region, bucketName, blobId string) error {
	bucketClient, err := ac.getBucket(bucketName)
	if err != nil {
		return err
	}
	_, err = bucketClient.NewBlobURL(blobId).Delete(context.Background(), azblob.DeleteSnapshotsOptionOnly, azblob.BlobAccessConditions{})
	return err
}

func (ac *AzureClient) OsRenameObject(region, bucketName, oldObjectKey, newObjectKey string) error {
	bucketClient, err := ac.getBucket(bucketName)
	if err != nil {
		return err
	}
	oldObj := bucketClient.NewBlobURL(oldObjectKey)
	newObj := bucketClient.NewBlobURL(newObjectKey)
	_, err = newObj.StartCopyFromURL(context.Background(), oldObj.URL(), azblob.Metadata{}, azblob.ModifiedAccessConditions{}, azblob.BlobAccessConditions{}, azblob.AccessTierNone, azblob.BlobTagsMap{})
	if err != nil {
		return err
	}
	_, err = oldObj.Delete(context.Background(), azblob.DeleteSnapshotsOptionOnly, azblob.BlobAccessConditions{})
	return err
}

func (ac *AzureClient) MqCreateQueue(region string, queue string, attr *QueueAttr) error {
	sClient, err := ac.getAqsQueueServiceClient()
	if err != nil {
		return err
	}
	opts := &azqueue.CreateOptions{Metadata: map[string]*string{}}
	attrVal := gjson.MarshalStringDefault(attr, false)
	opts.Metadata["attr"] = &attrVal
	_, err = sClient.NewQueueClient(queue).Create(context.Background(), opts)
	return err
}

func (ac *AzureClient) MqListQueues(region string, maxResults int, pageToken *string, queueNamePrefix *string) (queues []string, nextPageToken *string, err error) {
	sClient, err := ac.getAqsQueueServiceClient()
	if err != nil {
		return nil, nil, err
	}
	maxNum := int32(maxResults)
	opts := &azqueue.ListQueuesOptions{
		MaxResults: &maxNum,
		Marker:     pageToken,
		Prefix:     queueNamePrefix,
	}
	pager := sClient.NewListQueuesPager(opts)
	var results []string
	nextMarker := (*string)(nil)
	if pager.More() {
		currPage, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, nil, err
		}
		for _, v := range currPage.Queues {
			qAttr := QueueAttr{}
			err = json.Unmarshal([]byte(*v.Metadata["attr"]), &qAttr)
			if err != nil {
				return nil, nil, err
			}
			qName := *v.Name
			ac.queueAttrs.Store(qName, &qAttr)
			results = append(results, *v.Name)
		}
		nextMarker = currPage.NextMarker
	}

	return results, nextMarker, nil
}

func (ac *AzureClient) MqSend(region string, queue string, msg string) error {
	sClient, err := ac.getAqsQueueServiceClient()
	if err != nil {
		return err
	}
	qAttr, err := ac.getAqsQueueAttr()
	if err != nil {
		return err
	}
	if len(msg) > *qAttr.MsgMaxBytes {
		return gerrors.New("msg size %d > max limit %d", len(msg), qAttr.MsgMaxBytes)
	}
	ttl := int32(*qAttr.MsgRetentionSeconds)
	vt := int32(*qAttr.VisibilityTimeoutSeconds)
	opts := &azqueue.EnqueueMessageOptions{
		TimeToLive:        &ttl,
		VisibilityTimeout: &vt,
	}
	_, err = sClient.NewQueueClient(queue).EnqueueMessage(context.Background(), msg, opts)
	return err
}

func (ac *AzureClient) MqReceive(region string, queue string, msgNum int, waitSeconds int, deleteAfterReceived bool) ([]Msg, error) {
	sClient, err := ac.getAqsQueueServiceClient()
	if err != nil {
		return nil, err
	}

	var results []Msg
	if deleteAfterReceived {
		opts := azqueue.PeekMessagesOptions{}
		num := (int32)(msgNum)
		opts.NumberOfMessages = &num
		resp, err := sClient.NewQueueClient(queue).PeekMessages(context.Background(), nil)
		if err != nil {
			return nil, err
		}
		for _, v := range resp.Messages {
			results = append(results, Msg{
				Id:          *v.MessageID,
				EnqueueTime: *v.InsertionTime,
				Data:        *v.MessageText,
			})
		}
	} else {
		opts := azqueue.DequeueMessagesOptions{}
		num := (int32)(msgNum)
		opts.NumberOfMessages = &num
		resp, err := sClient.NewQueueClient(queue).DequeueMessages(context.Background(), nil)
		if err != nil {
			return nil, err
		}
		for _, v := range resp.Messages {
			results = append(results, Msg{
				Id:          *v.MessageID,
				EnqueueTime: *v.InsertionTime,
				Data:        *v.MessageText,
			})
		}
	}

	return results, nil
}

func (ac *AzureClient) MqDeleteQueue(region string, queue string) error {
	sClient, err := ac.getAqsQueueServiceClient()
	if err != nil {
		return err
	}
	_, err = sClient.NewQueueClient(queue).Delete(context.Background(), nil)
	return err
}

func (ac *AzureClient) Close() error {

}

func (ac *AzureClient) getBucket(bucket string) (*azblob.ContainerURL, error) {
	credential, err := azblob.NewSharedKeyCredential(ac.accessKey, ac.secretKey)
	if err != nil {
		return nil, err
	}
	cURL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", ac.accessKey, bucket))
	if err != nil {
		return nil, err
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	containerURL := azblob.NewContainerURL(*cURL, p)
	return &containerURL, nil
}

func (ac *AzureClient) getAqsQueueServiceClient() (*azqueue.ServiceClient, error) {
	if ac.azqsc != nil {
		return ac.azqsc, nil
	}

	cred, err := azqueue.NewSharedKeyCredential(ac.accessKey, ac.secretKey)
	if err != nil {
		return nil, fmt.Errorf("error creating shared key credential: %w", err)
	}
	serviceURL := fmt.Sprintf(azQueueEndpointExp, ac.accessKey)
	ac.azqsc, err = azqueue.NewServiceClientWithSharedKeyCredential(serviceURL, cred, nil)
	if err != nil {
		return nil, err
	}
	return ac.azqsc, nil
}

func (ac *AzureClient) getAqsQueueAttr() (*QueueAttr, error) {

}

func (ac *AzureClient) getVmClient() (*az, error) {
	groupsClient := resources.NewGroupsClient(clientData.SubscriptionID)
	groupsClient.Authorizer = authorizer

	group, err := groupsClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		resources.Group{
			Location: to.StringPtr(resourceGroupLocation)})

	deploymentsClient := resources.NewDeploymentsClient(clientData.SubscriptionID)
	deploymentsClient.Authorizer = authorizer

	deploymentFuture, err := deploymentsClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		deploymentName,
		resources.Deployment{
			Properties: &resources.DeploymentProperties{
				Template:   template,
				Parameters: params,
				Mode:       resources.Incremental,
			},
		},
	)
	if err != nil {
		return
	}
}
