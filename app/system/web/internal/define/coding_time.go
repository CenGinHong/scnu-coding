package define

import "github.com/gogf/gf/os/gtime"

// Fill with you ideas below.

// ListCodingTimeResp 编码时间
type ListCodingTimeResp struct {
	UserId     int                 `json:"userId"`
	CodingTime []*CodingTimeRecord `json:"codingTime"`
}

type CodingTimeRecord struct {
	Duration  int         `orm:"duration"               json:"duration"`  // 编码时间，分钟为单位
	CreatedAt *gtime.Time `orm:"created_at"             json:"createdAt"` // 创建时间
}
