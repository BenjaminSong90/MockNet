package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mocknet/utils"
)

// RData 请求数据
type RData struct {
	Uri        string                    //request real uri
	Params     utils.KVA[string, string] //request path params
	QueryArray utils.KVA[string, string] //request query params
	BodyValue  string                    //body 中读取到的值
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

func (handler *MimeJsonHandler) CollectParam(context *gin.Context, keyName string, userCardQueryArray []string) *RData {
	data := &RData{}
	data.Uri = context.Request.URL.RequestURI()
	for _, entry := range context.Params {
		data.Params = append(data.Params, utils.KV[string, string]{
			Key:   entry.Key,
			Value: entry.Value,
		})
	}

	for _, entry := range userCardQueryArray {
		q := context.Query(entry)
		if len(q) != 0 {
			data.QueryArray = append(data.QueryArray, utils.KV[string, string]{
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
