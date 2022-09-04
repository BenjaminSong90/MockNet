package utils

import (
	"fmt"
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
	return join(sep, kva.K())
}

// VJoin 拼接 KVA 的 所有 Value
func (kva KVA[K, V]) VJoin(sep string) string {
	return join(sep, kva.V())
}

func join(sep string, a ...any) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a...), " ", sep, -1), "[]")
}
