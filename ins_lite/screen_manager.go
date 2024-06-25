// Please let author have a drink, usdt trc20: TEpSxaE3kexE4e5igqmCZRMJNoDiQeWx29
// tg: @fuckins996
package ins_lite

import (
	"CentralizedControl/ins_lite/proto/msg/recver"
	"fmt"
	"sync"
	"time"
)

type ScreenManager struct {
	screenEvent      *ScreenEvent
	Screen           map[string]*recver.ScreenReceived
	screenSubmitData map[string]*ScreenSubmitData
	ScreenName2Id    map[string]int32
	screenWaitEvent  *EventList[string]
	curScreenName    string
	screenLock       sync.Mutex
}

func CreateScreenManager(client *InsLiteClient) *ScreenManager {
	return &ScreenManager{
		Screen:           map[string]*recver.ScreenReceived{},
		screenSubmitData: map[string]*ScreenSubmitData{},
		ScreenName2Id:    map[string]int32{},
		screenWaitEvent:  CreateWaitEvent[string](),
		screenEvent:      CreateScreenEvent(client.DefaultScreenDealFunc),
		curScreenName:    "",
	}
}

func (this *ScreenManager) getScreenIdByName(name string) int32 {
	this.screenLock.Lock()
	defer this.screenLock.Unlock()
	return this.ScreenName2Id[name]
}

func (this *ScreenManager) getScreenNameById(id int32) string {
	this.screenLock.Lock()
	defer this.screenLock.Unlock()
	for k, v := range this.ScreenName2Id {
		if v == id {
			return k
		}
	}
	panic(fmt.Sprintf("getScreenNameById: not find screen id: %d", id))
}

func (this *ScreenManager) addScreen(screen *recver.ScreenReceived) {
	this.screenLock.Lock()
	this.Screen[screen.GetScreenName()] = screen
	this.ScreenName2Id[screen.GetScreenName()] = screen.ScreenId
	this.curScreenName = screen.GetScreenName()
	this.screenLock.Unlock()
}

func (this *ScreenManager) getScreen(screenName string) *recver.ScreenReceived {
	this.screenLock.Lock()
	defer this.screenLock.Unlock()
	return this.Screen[screenName]
}

func (this *ScreenManager) getScreenById(screenId int32) *recver.ScreenReceived {
	return this.getScreen(this.getScreenNameById(screenId))
}

func (this *ScreenManager) GetCurrentScreen() *recver.ScreenReceived {
	return this.getScreen(this.curScreenName)
}

func (this *ScreenManager) GetScreenByWindowId(id string) *recver.SubScreen {
	return this.GetCurrentScreen().GetScreenByWindowId(id)
}

func (this *ScreenManager) DelSubmitData(screenName string, windowId string) {
	if screenName == "" {
		screenName = this.curScreenName
	}
	this.screenSubmitData[screenName].DelSubmitData(windowId)
}

func (this *ScreenManager) PutSubmitString(screenName string, windowId string, data any) error {
	if screenName == "" {
		screenName = this.curScreenName
	}
	submit := this.screenSubmitData[screenName]
	if submit == nil {
		panic(fmt.Sprintf("not find screen id: %s", screenName))
	}
	submit.PutSubmitData(windowId, data)
	return nil
}

func (this *ScreenManager) MustWaitScreenRecvFinish(screenName string) {
	if !this.WaitScreenRecvFinish(screenName) {
		panic(fmt.Sprintf("wait screen id %s error!", screenName))
	}
}

func (this *ScreenManager) WaitScreenRecvFinish(screenName string) bool {
	if this.getScreen(screenName) != nil {
		return true
	}
	event := this.screenWaitEvent.GetEvent(screenName)
	defer func() {
		this.screenWaitEvent.ReleaseEvent(event)
	}()
	return event.WaitForTime(time.Second * 60)
}

func (this *ScreenManager) WaitMultiScreenRecvFinish(screenName ...string) (bool, string) {
	for _, item := range screenName {
		if this.getScreen(item) != nil {
			return true, item
		}
	}
	events := make([]*Event[string], len(screenName))
	for idx, item := range screenName {
		events[idx] = this.screenWaitEvent.GetEvent(item)
	}
	defer func() {
		for _, item := range events {
			this.screenWaitEvent.ReleaseEvent(item)
		}
	}()
	ok, idx := WaitForTime(events, time.Second*60)
	witch := ""
	if ok {
		witch = screenName[idx]
	}
	return ok, witch
}

func (this *ScreenManager) targetScreenRecvFinish(screenId string) {
	this.screenWaitEvent.TargetEvent(screenId)
}

func (this *ScreenManager) GetScreenUpdateFinishEvent(screenName string) *Event[string] {
	return this.screenWaitEvent.GetEvent(screenName)
}

func (this *ScreenManager) targetScreenUpdateFinish(screenName string) {
	this.screenWaitEvent.TargetEvent(screenName)
}

func (this *ScreenManager) updateSubmitInfo(screen *recver.ScreenReceived, all *recver.SubScreenArray) {
	this.screenLock.Lock()
	defer this.screenLock.Unlock()

	screenSubmit := CreateScreenSubmitData()
	for idx := 0; idx < all.Count(); idx++ {
		item := all.Get(idx)
		switch item.Type {
		case 3:
			screen3 := item.ToSubScreen3()
			if screen3.SubmitDataType != 0 {
				screenSubmit.NewSubmitData(SubmitItem{
					WindowId:    screen3.WindowId.Value,
					Type:        screen3.SubmitDataType,
					NeedEncrypt: screen3.IsEncrypt(),
				})
			}
		}
	}
	oldSubmit := this.screenSubmitData[screen.GetScreenName()]
	if oldSubmit != nil {
		for idx := range screenSubmit.submit {
			item := oldSubmit.GetSubmitData(screenSubmit.submit[idx].WindowId)
			if item != nil {
				screenSubmit.submit[idx].Data = item.Data
			}
		}
	}
	this.screenSubmitData[screen.GetScreenName()] = screenSubmit
}

func (this *ScreenManager) getSubmitInfo(name string) *ScreenSubmitData {
	this.screenLock.Lock()
	defer this.screenLock.Unlock()
	return this.screenSubmitData[name].Copy()
}
