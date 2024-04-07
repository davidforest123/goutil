package twilio

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/container/gnum"
	"github.com/davidforest123/goutil/container/gvolume"
	"github.com/davidforest123/goutil/encoding/gjson"
	"github.com/davidforest123/goutil/i18n/gfiat"
	"github.com/sfreiberg/gotwilio"
	"net/http"
	"net/url"
)

type (
	TwilioClient struct {
		accountSid string
		authToken  string
		tw         *gotwilio.Twilio
	}
)

// New twilio client.
func newTwilio(accountSid, authToken string) (*TwilioClient, error) {
	return &TwilioClient{accountSid: accountSid,
		authToken: authToken,
		tw:        gotwilio.NewTwilioClient(accountSid, authToken),
	}, nil
}

// Reference: https://github.com/e154/smart-home/blob/08bfcca81c321b893cad87842a4ea40713f62e09/system/twilio/twilio.go#L102
// Get balance.
func (t *TwilioClient) GetBalance() (*Balance, error) {
	// Balance ...
	type twBalance struct {
		Currency   string `json:"currency"`
		Balance    string `json:"balance"`
		AccountSid string `json:"account_sid"`
	}

	// Build request.
	uri, err := url.Parse(fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Balance.json", t.accountSid))
	if err != nil {
		return nil, gerrors.New(err.Error())
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return nil, gerrors.New(err.Error())
	}
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", t.accountSid, t.authToken))))

	// Send request and decode response.
	resp, err := client.Do(req)
	if err != nil {
		return nil, gerrors.New(err.Error())
	}
	defer resp.Body.Close()
	tb := &twBalance{}
	if err = json.NewDecoder(resp.Body).Decode(tb); err != nil {
		return nil, gerrors.New(err.Error())
	}

	// Build 'Balance'.
	res := &Balance{}
	res.Currency, err = gfiat.ParseFiat(tb.Currency)
	if err != nil {
		return nil, err
	}
	res.Available, err = gnum.NewDecimalFromString(tb.Balance)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (t *TwilioClient) GlobalListRegions() ([]string, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) GlobalListZones(regionId string) ([]string, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) VmListSpotSpecs(regionId string) (*InstanceSpecList, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) VmListImages(regionId string) ([]SysImage, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) VmListInstances(regionId string) ([]InstanceInfo, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) VmListSecurityGroups(regionId string) ([]SecurityGroup, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) VmCreateSecurityGroup(regionId string, sg SecurityGroup) (string, error) {
	return "", gerrors.ErrNotSupport
}

func (t *TwilioClient) VmDeleteSecurityGroup(regionId, securityGroupId string) error {
	return gerrors.ErrNotSupport
}

func (t *TwilioClient) VmListSwitches(regionId, zoneId string) ([]string, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) VmCreateInstances(regionId string, tmpl InstanceCreationTmpl) (string, error) {
	return "", gerrors.ErrNotSupport
}

func (t *TwilioClient) VmStartInstances(regionId, vmId string) error {
	return gerrors.ErrNotSupport
}

func (t *TwilioClient) VmDeleteInstances(regionId string, instanceIds []string, force bool) error {
	return gerrors.ErrNotSupport
}

func (t *TwilioClient) OsIsBucketExist(regionId, bucketName string) (bool, error) {
	return false, gerrors.ErrNotSupport
}

func (t *TwilioClient) OsCreateBucket(regionId, bucketName string) error {
	return gerrors.ErrNotSupport
}

func (t *TwilioClient) OsDeleteBucket(regionId, bucketName string, deleteIfNotEmpty *bool) error {
	return gerrors.ErrNotSupport
}

func (t *TwilioClient) OsListBuckets(regionId string) ([]string, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) OsListObjectKeys(regionId string, bucketName string) ([]string, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) ObsScanObjectKeys(regionId, bucketName, pageToken string) ([]string, string, error) {
	return nil, "", gerrors.ErrNotSupport
}

func (t *TwilioClient) OsGetObjectSize(regionId, bucketName, objectKey string) (*gvolume.Volume, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) OsGetObject(regionId, bucketName, objectKey string) ([]byte, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) OsUpsertObject(regionId, bucketName, objectKey string, objectVal []byte) error {
	return gerrors.ErrNotSupport
}

func (t *TwilioClient) OsRenameObject(regionId, bucketName, oldObjectKey, newObjectKey string) error {
	return gerrors.ErrNotSupport
}

func (t *TwilioClient) OsDeleteObject(regionId, bucketName, objectKey string) error {
	return gerrors.ErrNotSupport
}

func (t *TwilioClient) SmsSendTmpl(regionId string, mobiles []string, sign, tmpl string, params map[string]string) error {
	return gerrors.ErrNotSupport
}

// Send SMS.
func (t *TwilioClient) SmsSendMsg(fromMobile, toMobile, message string) error {
	_, ex, err := t.tw.SendSMS(fromMobile, toMobile, message, "", "")
	if ex != nil && ex.Code != 200 {
		return gerrors.New(gjson.MarshalStringDefault(ex, false))
	}
	return err
}

func (t *TwilioClient) FcListFunctions(regionId, service string) ([]string, error) {
	return nil, gerrors.ErrNotSupport
}

func (t *TwilioClient) Close() error {
	return nil
}
