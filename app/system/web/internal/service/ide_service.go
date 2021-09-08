package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"scnu-coding/app/dao"
	"scnu-coding/app/service"
	"scnu-coding/app/system/web/internal/define"
	"scnu-coding/app/utils"
	"scnu-coding/library/role"
)

var Ide = newTheiaService()

type ideService struct {
	ide           iTheia
	ideAliveCache utils.MyCache
	lock          utils.MyMutex
}

type aliveStruct struct {
	Count     int         // front的数量
	CreatedAt *gtime.Time // 创建时间
}

func newTheiaService() *ideService {
	i := &ideService{}
	i.ide = newTheia()
	return i
}

// OpenIDE 新启ide容器
// @receiver t *ideService
// @param ctx context.Context
// @param req *define.OpenIDEReq
// @return url string
// @return err error
// @date 2021-07-17 22:23:46
func (t *ideService) OpenIDE(ctx context.Context, req *define.OpenIDEReq) (url string, err error) {
	// 当languageEnum==0，说明在完成某实验，可通过实验所属课程查的语言
	// 当languageEnum!=0，说明在打开自由工作区IDE
	if req.LanguageEnum == 0 {
		if req.LanguageEnum, err = t.getLanguageEnumByLabId(ctx, req.LabId); err != nil {
			return "", err
		}
	}
	ctxUser := service.Context.Get(ctx).User
	req.UserId = ctxUser.UserId
	// 默认挂载到自己的工作区下
	if req.MountedUserId == 0 {
		req.MountedUserId = req.UserId
	}
	if ctxUser.RoleId == role.STUDENT && req.LabId != 0 {
		deadline, err := dao.Lab.Ctx(ctx).WherePri(req.LabId).Value(dao.Lab.Columns.DeadLine)
		if err != nil {
			return "", err
		}
		if !deadline.IsNil() && !deadline.IsEmpty() && deadline.GTime().After(gtime.Now()) {
			req.IsEditAble = false
		}
	}
	url, err = t.ide.OpenTheia(ctx, req)
	if err != nil {
		return "", err
	}
	return url, nil
}

//func (t *ideService) ListContainer(courseId int) (resp []*define.ListContainerResp, err error) {
//	// 返回结果集
//	resp = make([]*define.ListContainerResp, 0)
//	keys, err := idePortCache.KeyStrings()
//	if err != nil {
//		return nil, err
//	}
//	// courseId为0是列表所有用户的ide容器,不为0表示只查出某门课的ide容器列表
//	containerKeys := make([]string, 0)
//	if courseId != 0 {
//		// 查出所有列表
//		labIds, err := dao.Lab.Where(dao.Lab.Columns.CourseId, courseId).Array(dao.Lab.Columns.LabId)
//		if err != nil {
//			return nil, err
//		}
//		for _, key := range keys {
//			split := gstr.Split(key, "-")
//			labId := split[2]
//			isExist := false
//			// 该实验在不在该课程内
//			for _, v := range labIds {
//				if v.String() == labId {
//					isExist = true
//					break
//				}
//			}
//			// 不属于该课程容器，移除
//			if isExist {
//				containerKeys = append(containerKeys, key)
//			}
//		}
//	}
//	// 返回的channel
//	retChan := make(chan *define.ListContainerResp, len(containerKeys))
//	for _, key := range containerKeys {
//		// 开协程收集信息
//		go func(retChan chan *define.ListContainerResp, key string) {
//			ret := &define.ListContainerResp{}
//			defer func() {
//				retChan <- ret
//			}()
//			// 获取语言，用户id，用户id
//			split := gstr.Split(key, "-")
//			ret.ContainerInfo.LanguageEnum = gconv.Int(split[0])
//			ret.ContainerInfo.UserId = gconv.Int(split[1])
//			ret.ContainerInfo.LabId = gconv.Int(split[2])
//			v, err := idePortCache.GetVar(key)
//			if err != nil {
//				return
//			}
//			ret.ContainerInfo.Port = v.Int()
//		}(retChan, key)
//	}
//	for i := 0; i < len(containerKeys); i++ {
//		v := <-retChan
//		resp = append(resp, v)
//	}
//	// 装填关联的字段
//	if err = dao.SysUser.WherePri(gdb.ListItemValuesUnique(resp, "ContainerInfo", "UserId")).
//		Fields(define.ListContainerResp{}.UserDetail).
//		ScanList(&resp, "UserDetail", "ContainerInfo", "user_id:UserId"); err != nil {
//		return nil, err
//	}
//	if err = dao.SysUser.WherePri(gdb.ListItemValuesUnique(resp, "ContainerInfo", "LabId")).
//		Fields(define.ListContainerResp{}.LabDetail).
//		ScanList(&resp, "LabDetail", "ContainerInfo", "lab_id:LabId"); err != nil {
//		return nil, err
//	}
//	return resp, nil
//}

// OpenFront 新打开一个容器页面
// @params userId
// @params labId
// @return err
// @date 2021-05-03 00:06:16
// 注意该请求是由插件发出，没有token携带userId信息
func (t *ideService) OpenFront(ctx context.Context, languageEnum int, userId int, labId int) (err error) {
	//获得语言类型
	if languageEnum == 0 {
		if languageEnum, err = t.getLanguageEnumByLabId(ctx, labId); err != nil {
			return err
		}
	}
	t.lock.Lock()
	defer t.lock.UnLock()
	alive := &aliveStruct{}
	key := fmt.Sprintf("%d-%d-%d", languageEnum, userId, labId)
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
// @receiver t *ideService
// @param Id int
// @param labID int
// @return err error
// @date 2021-07-17 22:25:40
func (t *ideService) CloseFront(ctx context.Context, req *define.CloseIDEReq) (err error) {
	if req.LanguageEnum == 0 {
		if req.LanguageEnum, err = t.getLanguageEnumByLabId(ctx, req.LabId); err != nil {
			return err
		}
	}
	t.lock.Lock()
	defer t.lock.UnLock()
	alive := &aliveStruct{}
	key := fmt.Sprintf("%d-%d-%d", req.LanguageEnum, req.UserId, req.LabId)
	data, err := t.ideAliveCache.Cache.GetVar(key)
	if err != nil {
		return err
	}
	if err = data.Struct(alive); err != nil {
		return err
	}
	alive.Count--
	// 所有的前端已经关闭
	if alive.Count == 0 {
		g.Log().Debugf("删除键为%s的容器", key)
		// 移除生命周期块
		if _, err = t.ideAliveCache.Cache.Remove(key); err != nil {
			return
		}
		// 移除端口标识块
		if err = removeIdePort(req.LanguageEnum, req.UserId, req.LabId); err != nil {
			g.Log().Error(err)
		}
		// 关闭容器
		if err = t.ide.removeIDE(ctx, &define.CloseIDEReq{
			IDEIdentifier: define.IDEIdentifier{
				UserId:       req.UserId,
				LanguageEnum: req.LanguageEnum,
				LabId:        req.LabId,
			},
		}); err != nil {
			return
		}
		// 插入编码时间
		duration := gtime.Now().Sub(alive.CreatedAt).Minutes()
		if _, err = dao.CodingTime.Ctx(ctx).Data(g.Map{
			dao.CodingTime.Columns.UserId:   req.UserId,
			dao.CodingTime.Columns.Duration: duration,
			dao.CodingTime.Columns.LabId:    req.LabId,
		}).Insert(); err != nil {
			g.Log().Error(err)
		}
	}
	if err = t.ideAliveCache.Cache.Set(key, alive, 0); err != nil {
		return err
	}
	return nil
}

func (t *ideService) getLanguageEnumByLabId(ctx context.Context, labId int) (languageEnum int, err error) {
	// 查出所用语言
	courseId, err := dao.Lab.Ctx(ctx).Cache(0).WherePri(labId).Value(dao.Lab.Columns.CourseId)
	if err != nil {
		return 0, err
	}
	languageEnumV, err := dao.Course.Ctx(ctx).Cache(0).WherePri(courseId.Int()).Value(dao.Course.Columns.LanguageType)
	if err != nil {
		return 0, err
	}
	languageEnum = languageEnumV.Int()
	return languageEnum, nil
}
