package service

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/util/gconv"
	"scnu-coding/app/dao"
	"scnu-coding/app/service"
	define2 "scnu-coding/app/system/admin/internala/define"
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/app/utils"
	"scnu-coding/library/response"
)

var IDE = iDEService{}

type iDEService struct {
}

// ListContainer 列出所有容器信息
func (s iDEService) ListContainer(ctx context.Context) (resp *response.PageResp, err error) {
	ctxPage := service.Context.Get(ctx).PageInfo
	containers, err := utils.DockerUtil.ListContainer(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: "ide"}),
	})
	if err != nil {
		return nil, err
	}
	records := make([]*define2.ListContainerResp, 0)
	// 填入基础信息，分页这里要手动分了
	for i := (ctxPage.Current - 1) * ctxPage.PageSize; i < (ctxPage.Current)*ctxPage.PageSize && i < len(containers); i++ {
		records = append(records, &define2.ListContainerResp{
			ContainerId: containers[i].ID,
			UserId:      gconv.Int(containers[i].Labels["userId"]),
			LabId:       gconv.Int(containers[i].Labels["labId"]),
			State:       containers[i].State,
			Status:      containers[i].Status,
		})
	}
	// 绑定详细信息
	if err = dao.SysUser.Ctx(ctx).
		Where(dao.SysUser.Columns.UserId, gdb.ListItemValuesUnique(records, "UserId")).
		Fields(define2.ListContainerResp{}.UserDetail).
		ScanList(records, "UserDetail", "user_id:UserId"); err != nil {
		return nil, err
	}
	if err = dao.Lab.Ctx(ctx).
		Where(dao.Lab.Columns.LabId, gdb.ListItemValuesUnique(records, "LabId")).
		Fields(define2.ListContainerResp{}.LabDetail).
		ScanList(records, "LabDetail", "lab_id:LabId"); err != nil {
		return nil, err
	}
	// 使用协程查询容器相关使用状态
	statChan := make(chan *define.ContainerStat)
	for _, record := range records {
		go func(ID string) {
			var stat *define.ContainerStat
			defer func() {
				statChan <- stat
			}()
			stat, err = utils.DockerUtil.GetContainerStat(ctx, ID)
			if err != nil {
				return
			}
		}(record.ContainerId)
	}
	// 写入容器的运行时状态
	for i := 0; i < len(records); i++ {
		select {
		case stat := <-statChan:
			{
				if stat == nil {
					continue
				}
				for _, record := range records {
					// TODO 还有CPU的要加
					if record.ContainerId == stat.Id {
						record.Memory = stat.MemoryStats.Usage
						record.MemoryLimit = stat.MemoryStats.Limit
						break
					}
				}
			}
		}
	}
	// 组装返回集
	resp = response.GetPageResp(records, len(records), nil)
	return resp, nil
}

func (s *iDEService) RestartContainer(ctx context.Context, containerId string) (err error) {
	if err = utils.DockerUtil.RestartContainer(ctx, containerId); err != nil {
		return err
	}
	return nil
}

func (s *iDEService) RemoveContainer(ctx context.Context, containerId string) (err error) {
	if err = utils.DockerUtil.RemoveContainer(ctx, containerId); err != nil {
		return err
	}
	return nil
}

func (s *iDEService) GetServerInfo(ctx context.Context) (info *define2.ServerInfo, err error) {
	info = &define2.ServerInfo{}
	info1, err := utils.DockerUtil.ServerInfo(ctx)
	if err != nil {
		return nil, err
	}
	if err = gconv.Struct(info1, info); err != nil {
		return nil, err
	}
	return info, nil
}
