package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mocknet/utils"
	"strings"
)

// KV key value 数据
type KV[K string | int, V string | int] struct {
	Key   K
	Value V
}

// KVA KV array
type KVA[K string | int, V string | int] []KV[K, V]

// K 返回 KVA 的 所有的Key
func (kva KVA[K, V]) K() (ks []K) {
	for _, entry := range kva {
		ks = append(ks, entry.Key)
	}
	return
}

// V 返回 KVA 的 所有的Value
func (kva KVA[K, V]) V() (vs []V) {
	for _, entry := range kva {
		vs = append(vs, entry.Value)
	}
	return
}

// KJoin 拼接 KVA 的 所有 Key
func (kva KVA[K, V]) KJoin(sep string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(kva.K()), " ", sep, -1), "[]")
}

// VJoin 拼接 KVA 的 所有 Value
func (kva KVA[K, V]) VJoin(sep string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(kva.V()), " ", sep, -1), "[]")
}

// RData 请求数据
type RData struct {
	Params     KVA[string, string]
	QueryArray KVA[string, string]
	BodyValue  string
}

func (data RData) GenerateKey() string {
	key := ""
	params := data.Params.VJoin(",")
	query := data.QueryArray.VJoin(",")

	if len(params) != 0 {
		key = params
	}

	if len(query) != 0 {
		if len(key) == 0 {
			key = fmt.Sprintf("%s", query)
		} else {
			key = fmt.Sprintf("%s,%s", key, query)
		}

	}

	if len(data.BodyValue) != 0 {
		if len(key) == 0 {
			key = fmt.Sprintf("%s", data.BodyValue)
		} else {
			key = fmt.Sprintf("%s,%s", key, data.BodyValue)
		}
	}

	return key
}

type MimeParamHandler interface {
	GetMimeType() string
	CollectParam(context *gin.Context, keyName string, queryArray []string) *RData
}

type MimeJsonHandler struct{}

var _ MimeParamHandler = &MimeJsonHandler{}

func (handler *MimeJsonHandler) GetMimeType() string {
	return "json"
}

func (handler *MimeJsonHandler) CollectParam(context *gin.Context, keyName string, queryArray []string) *RData {
	data := &RData{}
	for _, entry := range context.Params {
		data.Params = append(data.Params, KV[string, string]{
			Key:   entry.Key,
			Value: entry.Value,
		})
	}

	for _, entry := range queryArray {
		q := context.Query(entry)
		if len(q) != 0 {
			data.QueryArray = append(data.Params, KV[string, string]{
				Key:   entry,
				Value: q,
			})
		}
	}

	body := context.Request.Body
	if body != nil {
		bodyByte, err := io.ReadAll(body)
		defer func() {
			context.Request.Body = io.NopCloser(bytes.NewBuffer(bodyByte))
		}()
		if err == nil {
			jsonBody := map[string]interface{}{}
			decoder := json.NewDecoder(bytes.NewReader(bodyByte))
			decoder.UseNumber()
			err = decoder.Decode(&jsonBody)
			if err == nil {
				bodKV := make(map[string]interface{})
				utils.FlatMap(jsonBody, bodKV)
				keyValue, ok := bodKV[keyName]
				kvn := ""
				if kvn = utils.I2S(keyValue); (len(kvn) != 0) && ok {
					data.BodyValue = kvn
				}
			}
		}

	}

	return data

}
