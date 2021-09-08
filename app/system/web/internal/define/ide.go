package define

// @Author: 陈健航
// @Date: 2021/3/5 20:09
// @Description:

type OpenIDEReq struct {
	IDEIdentifier
	//UserId        int  // 打开IDE的用户
	//LanguageEnum  int  // 语言版本
	//LabId         int  // 实验id
	IsEditAble    bool // 是否可写
	MountedUserId int  // 被挂载工作空间的用户
}

type IDEIdentifier struct {
	UserId int
	LabId int
	LanguageEnum int
	//IsEditable  bool
}

type CloseIDEReq struct {
	IDEIdentifier
	//UserId       int
	//LabId        int
	//LanguageEnum int
}

//type CheckCodeReq struct {
//	TeacherId int
//	StuId     int
//	LabId     int
//}

//type CompilerErrorLogResp struct {
//	StuId       int    `json:"stuId"`
//	StuNum      string `json:"stuNum"`
//	CompilerLog string `json:"compilerLog"`
//}

//type IdeStat struct {
//Port int // 占用的端口
//IsFrontOpen  bool        // 是否有前端示例正在打开
//CreatedAt    *gtime.Time // 创建时间
//LastActiveAt *gtime.Time // 上一次心跳的时间
//}

//type ListContainerResp struct {
//	ContainerInfo *struct {
//		UserId       int
//		LabId        int
//		Port         int
//		LanguageEnum int
//	}
//	UserDetail *struct {
//		UserId   int    `orm:"user_id"`
//		Username string `orm:"username"`
//		UserNum  string `orm:"user_num"`
//	}
//	LabDetail *struct {
//		LabId    int    `orm:"lab_id"`
//		LabTitle string `orm:"lab_title"`
//	}
//}

//type PlagiarismCheckResp struct {
//	UserId1    int    `json:"user_id_1"`
//	UserId2    int    `json:"user_id_2"`
//	RealName1  string `json:"real_name_1"`
//	RealName2  string `json:"real_name_2"`
//	Num1       string `json:"num_1"`
//	Num2       string `json:"num_2"`
//	Url        string `json:"url"`
//	Similarity int    `json:"similarity"`
//}

//type PlagiarismCheckReq struct {
//	LabId int
//}
//
//type IdeAliveReq struct {
//	UserId int // 用户Id
//	LabId  int // 实验Id
//}
//
//type FileNode struct {
//	Name string `json:"name"` // 文件/文件夹名字
//	//IsDir     bool        `json:"isDir"`     // 是否文件夹
//	Content   string      `json:"content"`   // 内容
//	ChildNode []*FileNode `json:"childNode"` // 子节点
//}
