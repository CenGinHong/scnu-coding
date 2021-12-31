package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"scnu-coding/app/dao"
	"scnu-coding/app/service"
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/app/utils"
	"scnu-coding/library/role"
)

var IDE = newIDEService()

type iDEService struct {
	ide           iIDE
	ideAliveCache utils.MyCache
	lock          utils.MyMutex
}

type aliveStruct struct {
	Count     int         // front的数量
	CreatedAt *gtime.Time // 创建时间
}

func newIDEService() iDEService {
	i := iDEService{}
	i.ide = newIDE()
	i.ideAliveCache = *utils.NewMyCache()
	i.lock = utils.NewMyMutex()
	// 定时清理任务
	_, _ = gcron.Add("0 */1 * * * *", func() { i.ClearContainer(context.Background()) })
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

// OpenFront 新打开一个容器页面
// @params userId
// @params labId
// @return err
// @date 2021-05-03 00:06:16
// 注意该请求是由插件发出，没有token携带userId信息
func (t *iDEService) OpenFront(_ context.Context, req *define.IDEIdentifier) (err error) {
	t.lock.Lock()
	defer t.lock.UnLock()
	alive := &aliveStruct{}
	key := fmt.Sprintf("%d-%d", req.UserId, req.LabId)
	data, err := t.ideAliveCache.Cache.GetVar(key)
	if err != nil {
		return err
	}
	// 本来不存在
	if data.IsNil() {
		// 存入新的信息置存活信息
		alive.Count = 1
		alive.CreatedAt = gtime.Now()
	} else {
		// 更新
		if err = data.Struct(&alive); err != nil {
			return err
		}
		alive.Count++
	}
	if err = t.ideAliveCache.Cache.Set(key, alive, 0); err != nil {
		return err
	}
	return nil
}

// CloseFront 关闭一个前端页面，如果全部前端服务被关闭，启用goroutine准备关闭容器
// @Description:
// @receiver t *iDEService
// @param Id int
// @param labID int
// @return err error
// @date 2021-07-17 22:25:40
func (t *iDEService) CloseFront(_ context.Context, req *define.IDEIdentifier) (err error) {
	t.lock.Lock()
	defer t.lock.UnLock()
	alive := &aliveStruct{}
	key := fmt.Sprintf("%d-%d", req.UserId, req.LabId)
	data, err := t.ideAliveCache.Cache.GetVar(key)
	if err != nil {
		return err
	}
	if err = data.Struct(alive); err != nil {
		return err
	}
	// 所有的前端已经关闭
	alive.Count--
	if err = t.ideAliveCache.Cache.Set(key, alive, 0); err != nil {
		return err
	}
	return nil
}

func (t *iDEService) ClearContainer(ctx context.Context) {
	data, err := t.ideAliveCache.Cache.Data()
	if err != nil {
		return
	}
	for key, value := range data {
		alive := &aliveStruct{}
		if err = gconv.Struct(value, &alive); err != nil {
			g.Log().Error(err)
		}
		// 关闭该容器
		if alive.Count == 0 {
			go func() {
				ideKey := gconv.String(key)
				split := gstr.Split(ideKey, "-")
				userId := gconv.Int(split[0])
				labId := gconv.Int(split[1])
				ident := &define.CloseIDEReq{UserId: userId, LabId: labId}
				if err = t.ide.RemoveIDE(ctx, ident); err != nil {
					g.Log().Error(err)
					return
				}
				if _, err = t.ideAliveCache.Cache.Remove(key); err != nil {
					g.Log().Error(err)
					return
				}
				// 插入编码时间
				duration := gtime.Now().Sub(alive.CreatedAt).Minutes()
				if _, err = dao.CodingTime.Ctx(ctx).Data(g.Map{
					dao.CodingTime.Columns.UserId:   userId,
					dao.CodingTime.Columns.Duration: duration,
					dao.CodingTime.Columns.LabId:    labId,
				}).Insert(); err != nil {
					g.Log().Error(err)
				}
			}()
		}
	}
}
