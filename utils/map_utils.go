package utils

import "fmt"

//flat map info
func FlatMap(originMap map[string]interface{}, resultMap map[string]interface{}) {
	for key, val := range originMap {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			resultMap[key]=val
			//fmt.Println(key)
			FlatMap(val.(map[string]interface{}), resultMap)
		case []interface{}:
			resultMap[key]=val
			//fmt.Println(key)
			FlatArray(val.([]interface{}), resultMap)
		default:
			resultMap[key]=concreteVal
			//fmt.Println(key, ":", concreteVal)
		}
	}
}

//flat array map info
func FlatArray(anArray []interface{}, resultMap map[string]interface{}) {
	for i, val := range anArray {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			//fmt.Println("Index:", i)
			FlatMap(val.(map[string]interface{}), resultMap)
		case []interface{}:
			//fmt.Println("Index:", i)
			FlatArray(val.([]interface{}),resultMap)
		default:
			fmt.Println("Index", i, ":", concreteVal)
		}
	}
}
