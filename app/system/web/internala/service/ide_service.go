package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"scnu-coding/app/dao"
	"scnu-coding/app/service"
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/app/utils"
	"scnu-coding/library/role"
	"sync"
	"time"
)

var IDE = newIDEService()

type iDEService struct {
	ide           iIDE
	ideAliveCache utils.MyCache
	lock          utils.MyMutex
}

type aliveStruct struct {
	Count            int         // front的数量
	CreatedAt        *gtime.Time // 创建时间
	CountResetZeroAt *gtime.Time // count最近清零的时间
}

func newIDEService() iDEService {
	i := iDEService{}
	i.ide = newIDE()
	i.ideAliveCache = *utils.NewMyCache()
	i.lock = utils.NewMyMutex()
	//i.removeAllIDE()
	return i
}

// OpenIDE 新启ide容器
// @receiver t *iDEService
// @param ctx context.Context
// @param req *define.OpenIDEReq
// @return url string
// @return err error
// @date 2021-07-17 22:23:46
func (t *iDEService) OpenIDE(ctx context.Context, req *define.OpenIDEReq) (url string, err error) {
	ctxUser := service.Context.Get(ctx).User
	req.UserId = ctxUser.UserId
	// 默认挂载到自己的工作区下
	if req.MountedUserId == 0 {
		req.MountedUserId = req.UserId
	}
	// 作业类型要看有没有过期
	if ctxUser.RoleId == role.STUDENT && req.LabId != 0 {
		deadline, err := dao.Lab.Ctx(ctx).WherePri(req.LabId).Value(dao.Lab.Columns.Deadline)
		if err != nil {
			return "", err
		}
		// 过期了不可编辑
		if !deadline.IsNil() && !deadline.IsEmpty() && deadline.GTime().After(gtime.Now()) {
			req.IsEditAble = false
		}
	}
	url, err = t.ide.OpenIDE(ctx, req)
	if err != nil {
		return "", err
	}
	return url, nil
}

// FrontAlive 新打开一个容器页面
// @params userId
// @params labId
// @return err
// @date 2021-05-03 00:06:16
// 注意该请求是由插件发出，没有token携带userId信息
func (t *iDEService) FrontAlive(_ context.Context, req *define.FrontAliveReq) (err error) {
	t.lock.Lock()
	defer t.lock.UnLock()
	alive := &aliveStruct{}
	key := fmt.Sprintf("%d-%d", req.UserId, req.LabId)
	data, err := t.ideAliveCache.Cache.GetOrSet(key, &aliveStruct{
		Count:            0,
		CreatedAt:        gtime.Now(),
		CountResetZeroAt: nil,
	}, 0)
	if err != nil {
		return err
	}
	if err = gconv.Struct(data, &alive); err != nil {
		return err
	}
	// 看看要不要
	if req.IsOpen {
		alive.Count++
		alive.CountResetZeroAt = nil
	} else {
		alive.Count--
	}
	// 触发关闭检查
	if alive.Count == 0 {
		currentTime := gtime.Now()
		alive.CountResetZeroAt = currentTime
		go t.readyToShutDown(key, currentTime)
	}
	if err = t.ideAliveCache.Cache.Set(key, alive, 0); err != nil {
		return err
	}
	return nil
}

func (t *iDEService) readyToShutDown(key string, atZeroTime *gtime.Time) {
	// 先睡眠30秒
	time.Sleep(time.Duration(30) * time.Second)
	t.lock.Lock()
	defer t.lock.UnLock()
	// 重新取值检查
	data, err := t.ideAliveCache.Cache.GetOrSet(key, &aliveStruct{
		Count:            0,
		CreatedAt:        gtime.Now(),
		CountResetZeroAt: nil,
	}, 0)
	if err != nil {
		g.Log().Error(err)
		return
	}
	alive := &aliveStruct{}
	if err = gconv.Struct(data, alive); err != nil {
		g.Log().Error(err)
		return
	}
	// 触发关闭
	if alive.CountResetZeroAt != nil && alive.CountResetZeroAt == atZeroTime {
		split := gstr.Split(key, "-")
		userId := gconv.Int(split[0])
		labId := gconv.Int(split[1])
		// 移除容器
		if err = t.ide.RemoveIDE(context.Background(), &define.IDEIdentifier{UserId: userId, LabId: labId}); err != nil {
			return
		}
		// 移除缓存
		if _, err = t.ideAliveCache.Cache.Remove(key); err != nil {
			g.Log().Error(err)
			return
		}
		// 插入编码时间
		duration := alive.CountResetZeroAt.Sub(alive.CreatedAt).Minutes()
		if _, err = dao.CodingTime.Ctx(context.Background()).Data(g.Map{
			dao.CodingTime.Columns.UserId:   userId,
			dao.CodingTime.Columns.Duration: duration,
			dao.CodingTime.Columns.LabId:    labId,
		}).Insert(); err != nil {
			g.Log().Error(err)
		}
	}
	return
}

// removeAllIDE 关闭所有IDE容器
// @Description
// @receiver t
// @param ctx
// @date 2022-01-13 11:19:00
func (t *iDEService) removeAllIDE() {
	// 列表还在存活
	containers, err := utils.DockerUtil.ListContainer(context.Background(), types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: "ide"})})
	if err != nil {
		panic(err)
	}
	wg := &sync.WaitGroup{}
	for _, container := range containers {
		wg.Add(1)
		go func(container1 types.Container) {
			defer wg.Done()
			if err = utils.DockerUtil.RemoveContainer(context.Background(), container1.ID); err != nil {
				panic(err)
			}
		}(container)
	}
	wg.Wait()
}

func (t *iDEService) IsIDERunning(userId int, labId int) (isExist bool) {
	containerName := fmt.Sprintf("ide-%d-%d", userId, labId)
	isExist, err := t.ideAliveCache.Contains(containerName)
	if err != nil {
		return
	}
	return isExist
}
