package task

import (
	"github.com/project-nano/framework"
	"github.com/project-nano/core/modules"
	"log"
)

type GetBatchDeleteGuestExecutor struct {
	Sender         framework.MessageSender
	ResourceModule modules.ResourceModule
}


func (executor *GetBatchDeleteGuestExecutor)Execute(id framework.SessionID, request framework.Message,
	incoming chan framework.Message, terminate chan bool) (err error) {
	var batchID string
	if batchID, err = request.GetString(framework.ParamKeyID);err != nil{
		return err
	}
	resp, _ := framework.CreateJsonMessage(framework.GetBatchDeleteGuestResponse)
	resp.SetFromSession(id)
	resp.SetToSession(request.GetFromSession())
	resp.SetSuccess(false)

	var respChan = make(chan modules.ResourceResult, 1)
	executor.ResourceModule.GetBatchDeleteGuestStatus(batchID, respChan)
	var result = <- respChan
	if result.Error != nil{
		err = result.Error
		log.Printf("[%08X] get batch delete status from %s.[%08X] fail: %s", id, request.GetSender(), request.GetFromSession(), err.Error())
		resp.SetError(err.Error())
		return executor.Sender.SendMessage(resp, request.GetSender())
	}

	var guestStatus []uint64
	var guestID, guestName, deleteError []string

	for _, status := range result.BatchDelete{
		guestStatus = append(guestStatus, uint64(status.Status))
		guestID = append(guestID, status.ID)
		guestName = append(guestName, status.Name)
		deleteError = append(deleteError, status.Error)
	}
	resp.SetSuccess(true)
	resp.SetStringArray(framework.ParamKeyName, guestName)
	resp.SetStringArray(framework.ParamKeyGuest, guestID)
	resp.SetStringArray(framework.ParamKeyError, deleteError)
	resp.SetUIntArray(framework.ParamKeyStatus, guestStatus)
	return executor.Sender.SendMessage(resp, request.GetSender())
}

