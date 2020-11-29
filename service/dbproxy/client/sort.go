package client

import "time"

const baseFormat = "2006-01-02 15:04:05"

//ByUploadTime 通过上传时间排序的比较器
type ByUploadTime []FileMeta

func (v ByUploadTime) Len() int {
	return len(v)
}

func (v ByUploadTime) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v ByUploadTime) Less(i, j int) bool {
	iTime, _ := time.Parse(baseFormat, v[i].UploadAt)
	jTime, _ := time.Parse(baseFormat, v[j].UploadAt)
	return iTime.UnixNano() > jTime.UnixNano()
}
