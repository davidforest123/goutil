package cloud

import (
	"github.com/davidforest123/goutil/container/gvolume"
	"sync"
	"time"
)

type (
	retry struct {
		num      uint8
		interval time.Duration
	}

	cloud struct {
		cloud   cloudBasic
		retries sync.Map // map[funcName]retry{}
	}

	Cloud interface {
		cloudBasic
		MqIsQueueExists(region string, queue string) (bool, error)
		SetRetry(funcName string, num uint8, interval time.Duration)
	}
)

func NewCloud(platform Platform, accessKey string, secretKey string, LAN bool) (Cloud, error) {
	c, err := newBasic(platform, accessKey, secretKey, LAN)
	if err != nil {
		return nil, err
	}
	res := cloud{
		cloud: c,
	}
	return &res, nil
}

func (c *cloud) SetRetry(funcName string, num uint8, interval time.Duration) {
	if num <= 0 {
		num = 1
	}
	if interval < 0 {
		interval = time.Duration(0)
	}

	c.retries.Store(funcName, retry{num: num, interval: interval})
}

func (c *cloud) getRetryNum(funcName string) int {
	val, ok := c.retries.Load(funcName)
	if ok {
		return int(val.(retry).num)
	} else {
		return 1
	}
}

func (c *cloud) getRetryInterval(funcName string) time.Duration {
	val, ok := c.retries.Load(funcName)
	if ok {
		return val.(retry).interval
	} else {
		return time.Duration(0)
	}
}

func (c *cloud) GetBalance() (balance *Balance, err error) {
	for i := 0; i < c.getRetryNum("GetBalance"); i++ {
		balance, err = c.cloud.GetBalance()
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("GetBalance"))
	}
	return balance, err
}

func (c *cloud) ListLocations() (res map[string][]string, err error) {
	for i := 0; i < c.getRetryNum("GlobalListRegions"); i++ {
		res, err = c.cloud.ListLocations()
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("GlobalListRegions"))
	}
	return res, err
}

func (c *cloud) VmListSpotSpecs(region string) (res *InstanceSpecList, err error) {
	for i := 0; i < c.getRetryNum("VmListSpotSpecs"); i++ {
		res, err = c.cloud.VmListSpotSpecs(region)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("VmListSpotSpecs"))
	}
	return res, err
}

func (c *cloud) VmListImages(region string) (res []SysImage, err error) {
	for i := 0; i < c.getRetryNum("VmListImages"); i++ {
		res, err = c.cloud.VmListImages(region)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("VmListImages"))
	}
	return res, err
}

func (c *cloud) VmListInstances(region string) (res []InstanceInfo, err error) {
	for i := 0; i < c.getRetryNum("VmListInstances"); i++ {
		res, err = c.cloud.VmListInstances(region)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("VmListInstances"))
	}
	return res, err
}

func (c *cloud) VmListSecurityGroups(region string) (res []SecurityGroup, err error) {
	for i := 0; i < c.getRetryNum("VmListSecurityGroups"); i++ {
		res, err = c.cloud.VmListSecurityGroups(region)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("VmListSecurityGroups"))
	}
	return res, err
}

func (c *cloud) VmCreateSecurityGroup(region string, sg SecurityGroup) (securityGroupId string, err error) {
	for i := 0; i < c.getRetryNum("VmCreateSecurityGroup"); i++ {
		securityGroupId, err = c.cloud.VmCreateSecurityGroup(region, sg)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("VmCreateSecurityGroup"))
	}
	return securityGroupId, err
}

func (c *cloud) VmDeleteSecurityGroup(region, securityGroupId string) (err error) {
	for i := 0; i < c.getRetryNum("VmDeleteSecurityGroup"); i++ {
		err = c.cloud.VmDeleteSecurityGroup(region, securityGroupId)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("VmDeleteSecurityGroup"))
	}
	return err
}

func (c *cloud) VmListSwitches(region, zoneId string) (res []string, err error) {
	for i := 0; i < c.getRetryNum("VmListSwitches"); i++ {
		res, err = c.cloud.VmListSwitches(region, zoneId)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("VmListSwitches"))
	}
	return res, err
}

func (c *cloud) VmCreateInstances(region string, tmpl InstanceCreationTmpl, num int) (res []string, err error) {
	for i := 0; i < c.getRetryNum("VmCreateInstances"); i++ {
		res, err = c.cloud.VmCreateInstances(region, tmpl, num)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("VmCreateInstances"))
	}
	return res, err
}

func (c *cloud) VmStartInstances(region string, instanceIds []string) (err error) {
	for i := 0; i < c.getRetryNum("VmStartInstances"); i++ {
		err = c.cloud.VmStartInstances(region, instanceIds)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("VmStartInstances"))
	}
	return err
}

func (c *cloud) VmDeleteInstances(region string, instanceIds []string, force bool) (err error) {
	for i := 0; i < c.getRetryNum("VmDeleteInstances"); i++ {
		err = c.cloud.VmDeleteInstances(region, instanceIds, force)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("VmDeleteInstances"))
	}
	return err
}

func (c *cloud) OsIsBucketExist(region, bucketName string) (exist bool, err error) {
	for i := 0; i < c.getRetryNum("OsIsBucketExist"); i++ {
		exist, err = c.cloud.OsIsBucketExist(region, bucketName)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("OsIsBucketExist"))
	}
	return exist, err
}

func (c *cloud) OsCreateBucket(region, bucketName string) (err error) {
	for i := 0; i < c.getRetryNum("OsCreateBucket"); i++ {
		err = c.cloud.OsCreateBucket(region, bucketName)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("OsCreateBucket"))
	}
	return err
}

func (c *cloud) OsDeleteBucket(region, bucketName string, deleteIfNotEmpty *bool) (err error) {
	for i := 0; i < c.getRetryNum("OsDeleteBucket"); i++ {
		err = c.cloud.OsDeleteBucket(region, bucketName, deleteIfNotEmpty)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("OsDeleteBucket"))
	}
	return err
}

func (c *cloud) OsListBuckets(region string) (names []string, err error) {
	for i := 0; i < c.getRetryNum("OsListBuckets"); i++ {
		names, err = c.cloud.OsListBuckets(region)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("OsListBuckets"))
	}
	return names, err
}

func (c *cloud) OsListObjectKeys(region, bucketName string, keyPrefix *string, pageSize int, pageToken *string) (keys []string, nextPageToken *string, err error) {
	for i := 0; i < c.getRetryNum("OsListObjectKeys"); i++ {
		keys, nextPageToken, err = c.cloud.OsListObjectKeys(region, bucketName, keyPrefix, pageSize, pageToken)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("OsListObjectKeys"))
	}
	return keys, nextPageToken, err
}

func (c *cloud) OsGetObjectSize(region, bucketName, objectKey string) (vol *gvolume.Volume, err error) {
	for i := 0; i < c.getRetryNum("OsGetObjectSize"); i++ {
		vol, err = c.cloud.OsGetObjectSize(region, bucketName, objectKey)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("OsGetObjectSize"))
	}
	return vol, err
}

func (c *cloud) OsGetObject(region, bucketName, blobId string) (val []byte, err error) {
	for i := 0; i < c.getRetryNum("OsGetObject"); i++ {
		val, err = c.cloud.OsGetObject(region, bucketName, blobId)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("OsGetObject"))
	}
	return val, err
}

func (c *cloud) OsUpsertObject(region, bucketName, blobId string, blob []byte) (err error) {
	for i := 0; i < c.getRetryNum("OsUpsertObject"); i++ {
		err = c.cloud.OsUpsertObject(region, bucketName, blobId, blob)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("OsUpsertObject"))
	}
	return err
}

func (c *cloud) OsRenameObject(region, bucketName, oldKey, newKey string) (err error) {
	for i := 0; i < c.getRetryNum("OsRenameObject"); i++ {
		err = c.cloud.OsRenameObject(region, bucketName, oldKey, newKey)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("OsRenameObject"))
	}
	return err
}

func (c *cloud) OsDeleteObject(region, bucketName, blobId string) (err error) {
	for i := 0; i < c.getRetryNum("OsDeleteObject"); i++ {
		err = c.cloud.OsDeleteObject(region, bucketName, blobId)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("OsDeleteObject"))
	}
	return err
}

/*
func (c *cloud) SmsSendTmpl(region string, mobiles []string, sign, tmpl string, params map[string]string) (err error) {
	for i := 0; i < c.getRetryNum(""); i++ {
		err = c.cloud.SmsSendTmpl(region, mobiles, sign, tmpl, params)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval(""))
	}
	return err
}

func (c *cloud) SmsSendMsg(fromMobile, toMobile, message string) (err error) {
	for i := 0; i < c.getRetryNum(""); i++ {
		err = c.cloud.SmsSendMsg(fromMobile, toMobile, message)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval(""))
	}
	return err
}*/

func (c *cloud) MqListQueues(region string, maxResults int, nextToken *string, queueNamePrefix *string) (queues []string, outNextToken *string, err error) {
	for i := 0; i < c.getRetryNum("MqListQueues"); i++ {
		queues, outNextToken, err = c.cloud.MqListQueues(region, maxResults, nextToken, queueNamePrefix)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("MqListQueues"))
	}
	return queues, outNextToken, err
}

func (c *cloud) _MqIsQueueExists(region string, queue string) (bool, error) {
	nextToken := (*string)(nil)
	for {
		results, outNextToken, err := c.MqListQueues(region, 100, nextToken, &queue)
		if err != nil {
			return false, err
		}
		for _, v := range results {
			if v == queue {
				return true, nil
			}
		}
		if outNextToken == nil {
			return false, nil
		} else {
			nextToken = outNextToken
		}
	}
}

func (c *cloud) MqIsQueueExists(region string, queue string) (exists bool, err error) {
	for i := 0; i < c.getRetryNum("MqIsQueueExists"); i++ {
		exists, err = c._MqIsQueueExists(region, queue)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("MqIsQueueExists"))
	}
	return exists, err
}

func (c *cloud) MqCreateQueue(region string, queue string, attr *QueueAttr) (err error) {
	for i := 0; i < c.getRetryNum("MqCreateQueue"); i++ {
		err = c.cloud.MqCreateQueue(region, queue, attr)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("MqCreateQueue"))
	}
	return err
}

func (c *cloud) MqSend(region string, queue string, msg string) (err error) {
	for i := 0; i < c.getRetryNum("MqSend"); i++ {
		err = c.cloud.MqSend(region, queue, msg)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("MqSend"))
	}
	return err
}

func (c *cloud) MqReceive(region string, queue string, maxResult int, waitSeconds int, deleteAfterReceived bool) (results []Msg, err error) {
	for i := 0; i < c.getRetryNum("MqReceive"); i++ {
		results, err = c.cloud.MqReceive(region, queue, maxResult, waitSeconds, deleteAfterReceived)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("MqReceive"))
	}
	return results, err
}

func (c *cloud) MqDeleteQueue(region string, queue string) (err error) {
	for i := 0; i < c.getRetryNum("MqDeleteQueue"); i++ {
		err = c.cloud.MqDeleteQueue(region, queue)
		if err == nil {
			break
		}
		time.Sleep(c.getRetryInterval("MqDeleteQueue"))
	}
	return err
}

func (c *cloud) Close() (err error) {
	return c.Close()
}
