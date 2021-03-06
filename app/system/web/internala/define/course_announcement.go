package define

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
)

// Fill with you ideas below.

type CourseAnnouncementResp struct {
	CourseAnnouncementId int         `orm:"course_announcement_id,primary" json:"courseAnnouncementId"` // id
	Title                string      `orm:"title" json:"title"`                                         // 标题，限30字
	Content              string      `orm:"content" json:"content"`                                     // 公告内容，无字数限制
	AttachmentSrc        string      `orm:"attachment_src"                 json:"attachmentSrc"`        // 文件url
	UpdatedAt            *gtime.Time `orm:"updated_at" json:"updatedAt"`                                // 修改时间
}

type InsertCourseAnnouncementReq struct {
	CourseId      int               // 课程id
	Title         string            // 标题，限20字
	Content       string            // 公告内容，无字数限制
	UploadFile    *ghttp.UploadFile // 附件
	AttachmentSrc string
}

type UpdateCourseAnnouncementReq struct {
	CourseAnnouncementId int               // 课程id
	Title                *string           // 标题，限20字
	Content              *string           // 公告内容，无字数限制
	IsRemoveFile         bool              // 是否把文件移除了
	UploadFile           *ghttp.UploadFile // 附件
	AttachmentSrc        string
}
