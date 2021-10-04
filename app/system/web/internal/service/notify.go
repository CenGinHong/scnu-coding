package service

import (
	"scnu-coding/app/utils"
)

// @Author: 陈健航
// @Date: 2021/5/8 9:59
// @Description:

var Notify = &notifyService{
	isMessageNotifyNotReadCache: utils.NewMyCache(),
	isCommentNotifyNotReadCache: utils.NewMyCache(),
}

type notifyService struct {
	isMessageNotifyNotReadCache *utils.MyCache
	isCommentNotifyNotReadCache *utils.MyCache
}

//// MessageNotify 置入消息
//// @receiver receiver
//// @params ctx
//// @params messageContent
//// @params receiverIds
//// @return err
//// @date 2021-05-08 10:23:17
//func (n *notifyService) MessageNotify(ctx context.Context, messageContent string, receiverIds []int) (err error) {
//	ctxUser := service.Context.Get(ctx).User
//	data := make([]*model.MessageNotify, 0)
//	isReadData := make(map[interface{}]interface{})
//	for _, receiverId := range receiverIds {
//		data = append(data, &model.MessageNotify{
//			ReceiverId: receiverId,
//			SenderId:   ctxUser.UserId,
//			Content:    messageContent,
//		})
//		isReadData[receiverId] = true
//	}
//	if _, err = dao.MessageNotify.Batch(len(data)).Data(data).InsertLab(); err != nil {
//		return err
//	}
//	if err = n.isMessageNotifyNotReadCache.Sets(isReadData, 0); err != nil {
//		return err
//	}
//	return nil
//}
//
//// IsNotifyNotRead 是否有未读消息
//// @receiver receiver
//// @date 2021-05-08 10:10:57
//func (n *notifyService) IsNotifyNotRead(ctx context.Context) (isMessageNotifyNotRead bool, isCommentNotifyNotRead bool, err error) {
//	ctxUser := service.Context.Get(ctx).User
//	v, err := n.isMessageNotifyNotReadCache.GetVar(ctxUser.UserId)
//	if err != nil {
//		return false, false, err
//	}
//	if v.IsNil() {
//		isMessageNotifyNotRead = false
//	} else {
//		isMessageNotifyNotRead = v.Bool()
//	}
//
//	v, err = n.isCommentNotifyNotReadCache.GetVar(ctxUser.UserId)
//	if err != nil {
//		return false, false, err
//	}
//	if v.IsNil() {
//		isCommentNotifyNotRead = false
//	} else {
//		isMessageNotifyNotRead = v.Bool()
//	}
//	return isMessageNotifyNotRead, isCommentNotifyNotRead, nil
//}
//
////// ReadMessage 已经阅读消息（就是把红点消了）
////// @receiver receiver
////// @date 2021-05-08 10:21:00
////func (receiver *notifyService) ReadMessage(ctx context.Context) (err error) {
////	ctxUser := service.Context.Get(ctx).User
////	if err = receiver.isMessageNotifyNotReadCache.Set(ctxUser.UserID, false, 0); err != nil {
////		return err
////	}
////	return nil
////}
//
//func (n *notifyService) ListMessageNotify(ctx context.Context) (resp *response.PageResp, err error) {
//	ctxUser := service.Context.Get(ctx).User
//	ctxPageInfo := service.Context.Get(ctx).PageInfo
//	records := make([]*define.MessageNotifyResp, 0)
//	d := dao.MessageNotify.Where(dao.MessageNotify.Columns.ReceiverId, ctxUser.UserId)
//	if err = d.Page(ctxPageInfo.Current, ctxPageInfo.PageSize).Scan(&records); err != nil {
//		return nil, err
//	}
//	total, err := d.Count()
//	if err != nil {
//		return nil, err
//	}
//	resp = response.GetPageResp(records, total, nil)
//	return resp, nil
//}
