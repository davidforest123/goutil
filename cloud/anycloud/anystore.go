package anycloud

import (
	"encoding/json"
	"goutil/basic/gerrors"
	"goutil/encoding/gjson"
	"goutil/net/grpcs"
)

type (
	AnyStore struct {
		ac  *AnyCloud
		rpc *grpcs.Client
	}

	Cursor struct {
		idx      int
		dataMaps []any
	}

	CollType string
)

const (
	AnyStoreAccessKeyKey = "AnyStore.AccessKeyKey.85629437"
	CollTypeJson         = CollType("json")
	CollTypeCsv          = CollType("csv")
	CollTypeKv           = CollType("kv")
)

// NewRPCChecker used to verify function parameters between RPC client and server.
// Rules 'grpcs' package needs.
func NewRPCChecker() grpcs.ParamChecker {
	checker := grpcs.NewParamChecker()

	checker.Require("CreateDatabase", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("CreateDatabase", grpcs.In, "password", grpcs.TypeString)

	checker.Require("AuthDatabase", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("AuthDatabase", grpcs.In, "password", grpcs.TypeString)

	checker.Require("ResetDatabasePassword", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("ResetDatabasePassword", grpcs.In, "oldPassword", grpcs.TypeString)
	checker.Require("ResetDatabasePassword", grpcs.In, "newPassword", grpcs.TypeString)

	checker.Require("ListDatabases", grpcs.Out, "dbNames", grpcs.TypeStringSlice)

	checker.Require("IsDatabaseExist", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("IsDatabaseExist", grpcs.Out, "exist", grpcs.TypeBool)

	checker.Require("RenameDatabase", grpcs.In, "oldDBName", grpcs.TypeString)
	checker.Require("RenameDatabase", grpcs.In, "newDBName", grpcs.TypeString)

	checker.Require("DeleteDatabase", grpcs.In, "dbName", grpcs.TypeString)

	checker.Require("CreateCollection", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("CreateCollection", grpcs.In, "collName", grpcs.TypeString)
	checker.Require("CreateCollection", grpcs.In, "collType", grpcs.TypeString)
	checker.Require("CreateCollection", grpcs.In, "csvColumns", grpcs.TypeStringSlice)

	checker.Require("ListCollections", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("ListCollections", grpcs.Out, "collNames", grpcs.TypeStringSlice)

	checker.Require("ResetColumns", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("ResetColumns", grpcs.In, "collName", grpcs.TypeString)
	checker.Require("ResetColumns", grpcs.In, "newCsvColumns", grpcs.TypeStringSlice)

	checker.Require("IsCollectionExist", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("IsCollectionExist", grpcs.In, "collName", grpcs.TypeString)
	checker.Require("IsCollectionExist", grpcs.Out, "exist", grpcs.TypeBool)

	checker.Require("RenameCollection", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("RenameCollection", grpcs.In, "oldCollName", grpcs.TypeString)
	checker.Require("RenameCollection", grpcs.In, "newCollName", grpcs.TypeString)

	checker.Require("DeleteCollection", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("DeleteCollection", grpcs.In, "collName", grpcs.TypeString)

	checker.Require("GetCollectionMinId", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("GetCollectionMinId", grpcs.In, "collName", grpcs.TypeString)
	checker.Require("GetCollectionMinId", grpcs.Out, "minId", grpcs.TypeIF)

	checker.Require("GetCollectionMaxId", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("GetCollectionMaxId", grpcs.In, "collName", grpcs.TypeString)
	checker.Require("GetCollectionMaxId", grpcs.Out, "maxId", grpcs.TypeIF)

	checker.Require("GetDoc", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("GetDoc", grpcs.In, "collName", grpcs.TypeString)
	checker.Require("GetDoc", grpcs.In, "docId", grpcs.TypeIF)
	checker.Require("GetDoc", grpcs.Out, "doc", grpcs.TypeIF)

	checker.Require("GetDocs", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("GetDocs", grpcs.In, "collName", grpcs.TypeString)
	checker.Require("GetDocs", grpcs.In, "IdGte", grpcs.TypeIF)
	checker.Require("GetDocs", grpcs.In, "IdLte", grpcs.TypeIF)
	checker.Require("GetDocs", grpcs.Out, "docs", grpcs.TypeIFSlice)

	checker.Require("UpsertDocs", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("UpsertDocs", grpcs.In, "collName", grpcs.TypeString)
	checker.Require("UpsertDocs", grpcs.In, "docs", grpcs.TypeIFSlice)

	checker.Require("DeleteDocs", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("DeleteDocs", grpcs.In, "collName", grpcs.TypeString)
	checker.Require("DeleteDocs", grpcs.In, "IdGte", grpcs.TypeIF)
	checker.Require("DeleteDocs", grpcs.In, "IdLte", grpcs.TypeIF)

	checker.Require("ScanObjectKeys", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("ScanObjectKeys", grpcs.In, "collName", grpcs.TypeString)
	checker.Require("ScanObjectKeys", grpcs.In, "pageToken", grpcs.TypeString)
	checker.Require("ScanObjectKeys", grpcs.Out, "keys", grpcs.TypeStringSlice)
	checker.Require("ScanObjectKeys", grpcs.Out, "nextPageToken", grpcs.TypeString)

	checker.Require("UpsertObject", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("UpsertObject", grpcs.In, "collName", grpcs.TypeString)
	checker.Require("UpsertObject", grpcs.In, "key", grpcs.TypeString)
	checker.Require("UpsertObject", grpcs.In, "value", grpcs.TypeByteSlice)

	checker.Require("DeleteObject", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("DeleteObject", grpcs.In, "collName", grpcs.TypeString)
	checker.Require("DeleteObject", grpcs.In, "key", grpcs.TypeString)

	checker.Require("GetObject", grpcs.In, "dbName", grpcs.TypeString)
	checker.Require("GetObject", grpcs.In, "collName", grpcs.TypeString)
	checker.Require("GetObject", grpcs.In, "key", grpcs.TypeString)
	checker.Require("GetObject", grpcs.Out, "value", grpcs.TypeByteSlice)

	return *checker
}

func (c *Cursor) Next() bool {
	return c.idx <= len(c.dataMaps)-1
}

func (c *Cursor) Decode(valPtr any) error {
	defer func() {
		c.idx++
	}()

	if c.idx > len(c.dataMaps)-1 {
		return gerrors.New("no new data to decode")
	}
	buf, err := gjson.MarshalBytes(c.dataMaps[c.idx], false)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, valPtr)
}

// Create grpcs request, set auth token in request if username and password available.
func (c *AnyStore) newRequestWithAuth() grpcs.Request {
	request := grpcs.NewRequest()
	request.Set(AnyStoreAccessKeyKey, c.ac.accessKey)
	return request
}

// CreateDatabase creates a new database.
func (c *AnyStore) CreateDatabase(dbName string, password string) error {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	request.Set("password", password)
	return c.rpc.Call("CreateDatabase", request, nil)
}

// ResetDatabasePassword resets password for database.
// Users can perform multiple passwords for one 'dbName'.
// The last password will take effect and the previous password will be used only for decrypting the old chunks.
func (c *AnyStore) ResetDatabasePassword(dbName string, oldPassword, newPassword string) error {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	request.Set("oldPassword", oldPassword)
	request.Set("newPassword", newPassword)
	return c.rpc.Call("ResetDatabasePassword", request, nil)
}

// ListDatabases gets all database names.
func (c *AnyStore) ListDatabases() ([]string, error) {
	reply := grpcs.NewReply()
	if err := c.rpc.Call("ListDatabases", c.newRequestWithAuth(), &reply); err != nil {
		return nil, err
	}
	var res []string
	if reply.Get("dbNames") != nil {
		for _, v := range reply.Get("dbNames").([]any) {
			res = append(res, v.(string))
		}
	}
	return res, nil
}

// IsDatabaseExist checks if database exist or not.
func (c *AnyStore) IsDatabaseExist(dbName string) (bool, error) {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	reply := grpcs.NewReply()
	if err := c.rpc.Call("IsDatabaseExist", request, &reply); err != nil {
		return false, err
	}
	return reply.Get("exist").(bool), nil
}

// ListCollections gets all collections' name in specified database.
func (c *AnyStore) ListCollections(dbName string) ([]string, error) {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	reply := grpcs.NewReply()
	if err := c.rpc.Call("ListCollections", request, &reply); err != nil {
		return nil, err
	}
	var res []string
	if reply.Get("collNames") != nil {
		for _, v := range reply.Get("collNames").([]any) {
			res = append(res, v.(string))
		}
	}
	return res, nil
}

// RenameDatabase renames database.
func (c *AnyStore) RenameDatabase(oldDBName, newDBName string) error {
	request := c.newRequestWithAuth()
	request.Set("oldDBName", oldDBName)
	request.Set("newDBName", newDBName)
	if err := c.rpc.Call("RenameDatabase", request, nil); err != nil {
		return err
	}
	return nil
}

// DeleteDatabase deletes database.
func (c *AnyStore) DeleteDatabase(dbName string) error {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	return c.rpc.Call("DeleteDatabase", request, nil)
}

// CreateCollection creates new collection.
// csvColumns: user should set a valid `csvColumns` only when `collType` is CollTypeCSV,
// users can append or rename csvColumns, but CAN NOT delete any column.
func (c *AnyStore) CreateCollection(dbName, collName string, collType CollType, csvColumns []string) error {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	request.Set("collName", collName)
	request.Set("collType", collType)
	request.Set("csvColumns", csvColumns)
	reply := grpcs.NewReply()
	return c.rpc.Call("CreateCollection", request, &reply)
}

// ResetColumns resets new columns.
func (c *AnyStore) ResetColumns(dbName, collName string, newCsvColumns []string) error {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	request.Set("collName", collName)
	request.Set("newCsvColumns", newCsvColumns)
	reply := grpcs.NewReply()
	return c.rpc.Call("ResetColumns", request, &reply)
}

// IsCollectionExist checks if collection exist or not.
func (c *AnyStore) IsCollectionExist(dbName, collName string) (bool, error) {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	request.Set("collName", collName)
	reply := grpcs.NewReply()
	if err := c.rpc.Call("IsCollectionExist", request, &reply); err != nil {
		return false, err
	}
	return reply.Get("exist").(bool), nil
}

// RenameCollection renames collection.
func (c *AnyStore) RenameCollection(dbName, oldCollName, newCollName string) error {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	request.Set("oldCollName", oldCollName)
	request.Set("newCollName", newCollName)
	if err := c.rpc.Call("RenameCollection", request, nil); err != nil {
		return err
	}
	return nil
}

// DeleteCollection deletes collection.
func (c *AnyStore) DeleteCollection(dbName, collName string) error {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	request.Set("collName", collName)
	return c.rpc.Call("DeleteCollection", request, nil)
}

// GetCollectionMinId gets the smallest document Id in specified collection.
func (c *AnyStore) GetCollectionMinId(dbName, collName string) (any, error) {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	request.Set("collName", collName)
	reply := grpcs.NewReply()
	if err := c.rpc.Call("GetCollectionMinId", request, &reply); err != nil {
		return false, err
	}
	return reply.Get("minId"), nil
}

// GetCollectionMaxId gets the largest document Id in specified collection.
func (c *AnyStore) GetCollectionMaxId(dbName, collName string) (any, error) {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	request.Set("collName", collName)
	reply := grpcs.NewReply()
	if err := c.rpc.Call("GetCollectionMaxId", request, &reply); err != nil {
		return false, err
	}
	return reply.Get("maxId"), nil
}

// GetDoc gets document by document Id.
func (c *AnyStore) GetDoc(dbName, collName string, docId any) (*Cursor, error) {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	request.Set("collName", collName)
	request.Set("docId", docId)
	reply := grpcs.NewReply()
	if err := c.rpc.Call("GetDoc", request, &reply); err != nil {
		return nil, err
	}
	res := &Cursor{idx: 0}
	if reply.Get("doc") != nil {
		res.dataMaps = []any{reply.Get("doc")}
	}
	return res, nil
}

// GetDocs gets documents by document Id range.
// 'IdGte' and 'IdLte' cannot be nil at the same time.
func (c *AnyStore) GetDocs(dbName, collName string, IdGte any, IdLte any) (*Cursor, error) {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	request.Set("collName", collName)
	request.Set("IdGte", IdGte)
	request.Set("IdLte", IdLte)
	reply := grpcs.NewReply()
	if err := c.rpc.Call("GetDocs", request, &reply); err != nil {
		return nil, err
	}
	res := &Cursor{idx: 0}
	if reply.Get("docs") != nil {
		res.dataMaps = reply.Get("docs").([]any)
	}
	return res, nil
}

// UpsertDocs update/insert documents.
func (c *AnyStore) UpsertDocs(dbName, collName string, docs []any) error {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	request.Set("collName", collName)
	request.Set("docs", docs)
	return c.rpc.Call("UpsertDocs", request, nil)
}

// DeleteDocs deletes documents by document Id range.
// 'IdGte' and 'IdLte' cannot be nil at the same time.
func (c *AnyStore) DeleteDocs(dbName, collName string, IdGte any, IdLte any) error {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	request.Set("collName", collName)
	request.Set("IdGte", IdGte)
	request.Set("IdLte", IdLte)
	return c.rpc.Call("DeleteDocs", request, nil)
}

// ScanObjectKeys scans object keys in bucket.
func (c *AnyStore) ScanObjectKeys(dbName, collName, pageToken string) ([]string, string, error) {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	request.Set("collName", collName)
	request.Set("pageToken", pageToken)
	reply := grpcs.NewReply()
	if err := c.rpc.Call("ScanObjectKeys", request, &reply); err != nil {
		return nil, "", err
	}
	var res []string
	if reply.Get("keys") != nil {
		for _, v := range reply.Get("keys").([]any) {
			res = append(res, v.(string))
		}
	}
	return res, reply.Get("nextPageToken").(string), nil
}

// UpsertObject update/insert object into bucket.
func (c *AnyStore) UpsertObject(dbName, collName, key string, value []byte) error {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	request.Set("collName", collName)
	request.Set("key", key)
	request.Set("value", value)
	return c.rpc.Call("UpsertObject", request, nil)
}

// DeleteObject deletes object from bucket.
func (c *AnyStore) DeleteObject(dbName, collName, key string) error {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	request.Set("collName", collName)
	request.Set("key", key)
	return c.rpc.Call("DeleteObject", request, nil)
}

// GetObject reads object from bucket.
func (c *AnyStore) GetObject(dbName, collName, key string) ([]byte, error) {
	request := c.newRequestWithAuth()
	request.Set("dbName", dbName)
	request.Set("collName", collName)
	request.Set("key", key)
	reply := grpcs.NewReply()
	if err := c.rpc.Call("GetObject", request, &reply); err != nil {
		return nil, err
	}
	res := []byte(nil)
	if reply.Get("value") != nil {
		res = reply.Get("value").([]byte)
	}
	return res, nil
}

// Close connection to InfDB server.
func (c *AnyStore) Close() error {
	return c.rpc.Close()
}
