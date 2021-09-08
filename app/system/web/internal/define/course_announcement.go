package define

import (
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gmeta"
)

// Fill with you ideas below.

type CourseAnnouncementListResp struct {
	gmeta.Meta           `orm:"table:course_announcement"`
	CourseAnnouncementId int    `orm:"course_announcement_id,primary" json:"courseAnnouncementId"` // id
	Title                string `orm:"title" json:"title"`                                         // 标题，限30字
	Content              string `orm:"content" json:"content"`                                     // 公告内容，无字数限制
	AttachmentFileId     int    `orm:"attachment_file_id"             json:"-"`                    // 附件url
	AttachmentFileDetail *struct {
		gmeta.Meta  `orm:"table:local_file"`
		LocalFileId int    `orm:"local_file_id,primary" json:"localFileId"` //
		Filename    string `orm:"filename"              json:"filename"`    // 文件名
		Size        int    `orm:"size"                  json:"size"`        // 文件大小
		Url         string `orm:"url"                   json:"url"`         // 文件url
		ContentType string `orm:"content_type"          json:"contentType"` // 文件类型
	} `orm:"with:local_file_id=attachment_file_id" json:"attachmentFileDetail"`
	UpdatedAt *gtime.Time `orm:"updated_at" json:"updatedAt"` // 修改时间
}

type InsertCourseAnnouncementReq struct {
	CourseId         int    // 课程id
	Title            string // 标题，限20字
	Content          string // 公告内容，无字数限制
	AttachmentFileId int    // 附件url
}

type UpdateCourseAnnouncementReq struct {
	CourseAnnouncementId int    // 课程id
	Title                string // 标题，限20字
	Content              string // 公告内容，无字数限制
	AttachmentFileId     int    // 附件url
}
