package define

import (
	"github.com/gogf/gf/os/gtime"
)

// @Author: 陈健航
// @Date: 2021/3/5 20:09
// @Description:

type OpenIDEReq struct {
	IDEIdentifier
	IsEditAble    bool // 是否可写
	MountedUserId int  // 被挂载工作空间的用户
}

type IDEIdentifier struct {
	UserId int
	LabId  int
}

type FrontAliveReq struct {
	IDEIdentifier
	IsOpen bool // 打开/关闭页面
}

// ContainerStat 有些不会怎么用到的字段我先注释了
type ContainerStat struct {
	Read      gtime.Time `json:"read"`
	Preread   gtime.Time `json:"preread"`
	PidsStats struct {
		Current int64 `json:"current"`
	} `json:"pids_stats"`
	//BlkioStats struct {
	//	IoServiceBytesRecursive []struct {
	//		Major int    `json:"major"`
	//		Minor int    `json:"minor"`
	//		Op    string `json:"op"`
	//		Value int    `json:"value"`
	//	} `json:"io_service_bytes_recursive"`
	//	IoServicedRecursive []struct {
	//		Major int    `json:"major"`
	//		Minor int    `json:"minor"`
	//		Op    string `json:"op"`
	//		Value int    `json:"value"`
	//	} `json:"io_serviced_recursive"`
	//	IoQueueRecursive       []interface{} `json:"io_queue_recursive"`
	//	IoServiceTimeRecursive []interface{} `json:"io_service_time_recursive"`
	//	IoWaitTimeRecursive    []interface{} `json:"io_wait_time_recursive"`
	//	IoMergedRecursive      []interface{} `json:"io_merged_recursive"`
	//	IoTimeRecursive        []interface{} `json:"io_time_recursive"`
	//	SectorsRecursive       []interface{} `json:"sectors_recursive"`
	//} `json:"blkio_stats"`
	NumProcs     int64 `json:"num_procs"`
	StorageStats struct {
	} `json:"storage_stats"`
	CpuStats struct {
		CpuUsage struct {
			TotalUsage        int64   `json:"total_usage"`
			PercpuUsage       []int64 `json:"percpu_usage"`
			UsageInKernelmode int64   `json:"usage_in_kernelmode"`
			UsageInUsermode   int     `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCpuUsage int64 `json:"system_cpu_usage"`
		OnlineCpus     int64 `json:"online_cpus"`
		ThrottlingData struct {
			Periods          int64 `json:"periods"`
			ThrottledPeriods int64 `json:"throttled_periods"`
			ThrottledTime    int64 `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"cpu_stats"`
	PrecpuStats struct {
		CpuUsage struct {
			TotalUsage        int64 `json:"total_usage"`
			UsageInKernelmode int64 `json:"usage_in_kernelmode"`
			UsageInUsermode   int64 `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		ThrottlingData struct {
			Periods          int64 `json:"periods"`
			ThrottledPeriods int64 `json:"throttled_periods"`
			ThrottledTime    int64 `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"precpu_stats"`
	MemoryStats struct {
		Usage    int64 `json:"usage"`
		MaxUsage int64 `json:"max_usage"`
		//Stats    struct {
		//	ActiveAnon              int   `json:"active_anon"`
		//	ActiveFile              int   `json:"active_file"`
		//	Cache                   int   `json:"cache"`
		//	Dirty                   int   `json:"dirty"`
		//	HierarchicalMemoryLimit int64 `json:"hierarchical_memory_limit"`
		//	HierarchicalMemswLimit  int64 `json:"hierarchical_memsw_limit"`
		//	InactiveAnon            int   `json:"inactive_anon"`
		//	InactiveFile            int   `json:"inactive_file"`
		//	MappedFile              int   `json:"mapped_file"`
		//	Pgfault                 int   `json:"pgfault"`
		//	Pgmajfault              int   `json:"pgmajfault"`
		//	Pgpgin                  int   `json:"pgpgin"`
		//	Pgpgout                 int   `json:"pgpgout"`
		//	Rss                     int   `json:"rss"`
		//	RssHuge                 int   `json:"rss_huge"`
		//	TotalActiveAnon         int   `json:"total_active_anon"`
		//	TotalActiveFile         int   `json:"total_active_file"`
		//	TotalCache              int   `json:"total_cache"`
		//	TotalDirty              int   `json:"total_dirty"`
		//	TotalInactiveAnon       int   `json:"total_inactive_anon"`
		//	TotalInactiveFile       int   `json:"total_inactive_file"`
		//	TotalMappedFile         int   `json:"total_mapped_file"`
		//	TotalPgfault            int   `json:"total_pgfault"`
		//	TotalPgmajfault         int   `json:"total_pgmajfault"`
		//	TotalPgpgin             int   `json:"total_pgpgin"`
		//	TotalPgpgout            int   `json:"total_pgpgout"`
		//	TotalRss                int   `json:"total_rss"`
		//	TotalRssHuge            int   `json:"total_rss_huge"`
		//	TotalUnevictable        int   `json:"total_unevictable"`
		//	TotalWriteback          int   `json:"total_writeback"`
		//	Unevictable             int   `json:"unevictable"`
		//	Writeback               int   `json:"writeback"`
		//} `json:"stats"`
		Limit int64 `json:"limit"`
	} `json:"memory_stats"`
	Name     string `json:"name"`
	Id       string `json:"id"`
	Networks struct {
		Eth0 struct {
			RxBytes   int `json:"rx_bytes"`
			RxPackets int `json:"rx_packets"`
			RxErrors  int `json:"rx_errors"`
			RxDropped int `json:"rx_dropped"`
			TxBytes   int `json:"tx_bytes"`
			TxPackets int `json:"tx_packets"`
			TxErrors  int `json:"tx_errors"`
			TxDropped int `json:"tx_dropped"`
		} `json:"eth0"`
	} `json:"networks"`
}

//type CloseIDEReq struct {
//	IDEIdentifier
//}

//type CheckCodeReq struct {
//	TeacherId int
//	UserId     int
//	LabId     int
//}

//type CompilerErrorLogResp struct {
//	UserId       int    `json:"stuId"`
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
