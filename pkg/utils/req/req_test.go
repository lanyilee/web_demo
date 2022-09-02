package req

import "testing"

func TestReqJSON(t *testing.T) {
	ReqJSON("GET", "https://baidu.com", nil, nil, nil)
}
