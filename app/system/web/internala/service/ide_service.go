package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"scnu-coding/app/dao"
	"scnu-coding/app/service"
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/app/utils"
	"scnu-coding/library/role"
	"strings"
	"sync"
	"time"
)

var IDE = newIDEService()

type iDEService struct {
	ide               iIDE
	ideHeartBeatCache utils.MyCache
	lock              utils.MyMutex
}

func newIDEService() iDEService {
	i := iDEService{}
	i.ide = newIDE()
	i.ideHeartBeatCache = *utils.NewMyCache()
	i.lock = utils.NewMyMutex()
	// 在服务重启将所有ide关闭
	i.removeAllIDE()
	go i.shutDownExpire()
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
	// 先置入一个心跳
	if err = t.Heartbeat(ctx, &define.HeartBeatReq{
		IDEIdentifier: req.IDEIdentifier,
	}); err != nil {
		return "", err
	}
	return url, nil
}

// Heartbeat 新打开一个容器页面
// @params userId
// @params labId
// @return err
// @date 2021-05-03 00:06:16
// 注意该请求是由插件发出，没有token携带userId信息
func (t *iDEService) Heartbeat(_ context.Context, req *define.HeartBeatReq) (err error) {
	t.lock.Lock()
	defer t.lock.UnLock()
	key := fmt.Sprintf("%d-%d", req.UserId, req.LabId)
	// 取出所有心跳池的值
	v, err := t.ideHeartBeatCache.GetOrSet(key, &define.IDEHeartBeat{
		ActiveTime: gtime.Now(),
		StartTime:  gtime.Now(),
	}, 0)
	if err != nil {
		return err
	}
	heartbeat := &define.IDEHeartBeat{}
	if err = gconv.Struct(v, &heartbeat); err != nil {
		return err
	}
	heartbeat.ActiveTime = gtime.Now()
	if err = t.ideHeartBeatCache.Set(key, heartbeat, 0); err != nil {
		return err
	}
	g.Log().Debug("receive a heartbeat")
	return nil
}

func (t *iDEService) shutDownExpire() {
	// 先睡眠30秒
	// 重新取值检查
	for true {
		time.Sleep(time.Duration(3) * time.Minute)
		func() {
			t.lock.Lock()
			defer t.lock.UnLock()
			// 列出所有在运行的ide
			names := t.ide.ListIDEContainerName(context.Background())
			wg := sync.WaitGroup{}
			for _, name := range names {
				g.Log().Debug(name)
				v, err := t.ideHeartBeatCache.Get(name)
				if err != nil {
					continue
				}
				heartbeat := &define.IDEHeartBeat{}
				if v != nil {
					if err = gconv.Struct(v, &heartbeat); err != nil {
						continue
					}
				}
				// 超时，回收IDE
				if v == nil || heartbeat.StartTime.Add(time.Duration(1)*time.Minute).Before(gtime.Now()) {
					wg.Add(1)
					split := strings.Split(name, "-")
					userId := gconv.Int(split[1])
					labId := gconv.Int(split[2])
					key := fmt.Sprintf("%d-%d", userId, labId)
					// 移除
					go func() {
						defer wg.Done()
						if err = t.ide.RemoveIDE(context.Background(), &define.IDEIdentifier{
							UserId: userId,
							LabId:  labId,
						}); err != nil {
							return
						}
						// 缓存处理
						{
							// 记录下编码时间
							if _, err = dao.CodingTime.Ctx(context.Background()).Insert(g.Map{
								dao.CodingTime.Columns.UserId:   userId,
								dao.CodingTime.Columns.LabId:    labId,
								dao.CodingTime.Columns.Duration: gtime.Now().Sub(heartbeat.StartTime),
							}); err != nil {
								return
							}
							if _, err = t.ideHeartBeatCache.Remove(key); err != nil {
								return
							}
						}
					}()
				}
			}
			wg.Wait()
		}()
	}
}

// removeAllIDE 关闭所有IDE容器
// @Description
// @receiver t
// @param ctx
// @date 2022-01-13 11:19:00
func (t *iDEService) removeAllIDE() {
	// 列表还在存活
	t.lock.Lock()
	defer t.lock.UnLock()
	ideContainerNames := t.ide.ListIDEContainerName(context.Background())
	wg := &sync.WaitGroup{}
	for _, ideContainerName := range ideContainerNames {
		wg.Add(1)
		go func(container1 string) {
			defer wg.Done()
			split := strings.Split(container1, "-")
			userId := gconv.Int(split[1])
			labId := gconv.Int(split[2])
			if err := t.ide.RemoveIDE(context.Background(),
				&define.IDEIdentifier{
					UserId: userId,
					LabId:  labId,
				}); err != nil {
				panic(err)
			}
		}(ideContainerName)
	}
	wg.Wait()
}

func (t *iDEService) IsIDERunning(userId int, labId int) (isExist bool) {
	containerName := fmt.Sprintf("ide-%d-%d", userId, labId)
	isExist, err := t.ideHeartBeatCache.Contains(containerName)
	if err != nil {
		return
	}
	return isExist
}
