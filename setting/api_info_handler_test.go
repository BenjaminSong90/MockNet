package setting

import (
	"testing"
)

func TestMerge(t *testing.T) {
	apiData := &ApiData{
		Path:         "old_path",
		QueryKey:     []string{"ok", "name"},
		Method:       "GET",
		BodyKey:      "get_info",
		NeedRedirect: true,
	}
	newApiData := &ApiData{
		Path:         "old_path",
		QueryKey:     []string{"ok", "name", "hh"},
		Method:       "POST",
		BodyKey:      "get_info_new",
		NeedRedirect: true,
	}

	if apiData.Merge(*newApiData) {
		t.Error()
	}

	newApiData.Method = "GET"

	if !apiData.Merge(*newApiData) {
		t.Error()
	}

}
