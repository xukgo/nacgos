package main

/*
#include <stdlib.h>
#cgo CFLAGS: -I .
#cgo LDFLAGS: -L .

typedef void (*ConfigUpdateEvent)(void *usrData, char *group, char *dataId, char *data);

typedef struct {
  void *usrData;
  char *name;
  ConfigUpdateEvent event;
} MatchVarEventHandler;

typedef struct {
  int count;
  MatchVarEventHandler* handlers;
} MatchVarEventHandlerCollection;

extern void nacdef_doUpdateEvent(ConfigUpdateEvent evt, void *usrData, char *group, char *dataId, char *data){
	evt(usrData, group,dataId,data);
}

*/
//#cgo LDFLAGS: -Wl,--allow-multiple-definition
import "C"
import (
	"fmt"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	"github.com/xukgo/gsaber/utils/fileUtil"
	"github.com/xukgo/naconfig"
)

func main() {
	fmt.Printf("hello world\n")
}

//export nacgosInitStorage
func nacgosInitStorage(filePath *C.char, matchVarHandlers *C.MatchVarEventHandlerCollection, retCode *C.int) unsafe.Pointer {
	//fmt.Printf("matchVarHandlers:%v\n", matchVarHandlers)
	var matchHandlers = convertMatchVarHandlers(matchVarHandlers)
	//fmt.Printf("matchHandlers:\n\t%v\n", len(matchHandlers))
	gpath := C.GoString(filePath)
	fp := fileUtil.GetAbsUrl(gpath)

	storage := new(naconfig.Repo)
	err := storage.InitFromXmlPath(fp, matchHandlers)
	if err != nil {
		fmt.Printf("nacos storage init from xml return error:%s", err.Error())
		*retCode = -1
		return nil
	}

	*retCode = 0
	p := gopointer.Save(storage)
	return p
	//return unsafe.Pointer(storage)
}

func convertMatchVarHandlers(cMatchVarHandlers *C.MatchVarEventHandlerCollection) []naconfig.MatchVarHandler {
	if cMatchVarHandlers == nil {
		return nil
	}
	count := int(cMatchVarHandlers.count)
	handlersPrt := cMatchVarHandlers.handlers
	var results = make([]naconfig.MatchVarHandler, 0, count)

	//用go的方式遍历C数组
	for i := 0; i < count; i++ {
		skipLen := uintptr(C.sizeof_MatchVarEventHandler * C.int(i))
		evtHandler := (*C.MatchVarEventHandler)(unsafe.Pointer(uintptr(unsafe.Pointer(handlersPrt)) + skipLen))

		name := C.GoString(evtHandler.name)
		h := initMatchHandlerFromC(evtHandler)
		fmt.Printf("event name:%s ptr:%p\n", name, h)
		results = append(results, naconfig.InitMatchVarHandler(name, h))
	}
	return results
}

func initMatchHandlerFromC(evtHandler *C.MatchVarEventHandler) func(group string, dataId string, data string) {
	usrData := unsafe.Pointer(evtHandler.usrData)
	h := func(group, dataId, data string) {
		//fmt.Printf("update get group:%s dataId:%s data:%s\n", group, dataId, data)
		var cgroup *C.char = C.CString(group)
		defer C.free(unsafe.Pointer(cgroup))
		var cdataId *C.char = C.CString(dataId)
		defer C.free(unsafe.Pointer(cdataId))
		var cdata *C.char = C.CString(data)
		defer C.free(unsafe.Pointer(cdata))
		C.nacdef_doUpdateEvent(evtHandler.event, usrData, cgroup, cdataId, cdata)
	}
	return h
}

//export nacgosStartSubscribe
func nacgosStartSubscribe(storagePointer unsafe.Pointer, block C.int) C.int {
	storage := gopointer.Restore(storagePointer).(*naconfig.Repo)
	isBlock := block == 1
	//fmt.Printf("block:%v\n", isBlock)

	err := storage.Subscribe(isBlock)
	if err != nil {
		fmt.Printf("nacos storage subscribe return error:%s", err.Error())
		return -1
	}
	return 0
}

//export nacgosPublish
func nacgosPublish(storagePointer unsafe.Pointer, group *C.char, dataID *C.char, content *C.char) C.int {
	storage := gopointer.Restore(storagePointer).(*naconfig.Repo)

	groupStr := C.GoString(group)
	dataIDStr := C.GoString(dataID)
	contentStr := C.GoString(content)
	//fmt.Printf("publish group:%s dataId:%s data:%s\n", groupStr, dataIDStr, contentStr)
	err := storage.Publish(groupStr, dataIDStr, contentStr)
	if err != nil {
		fmt.Printf("nacos storage publish return error:%s", err.Error())
		return -1
	}
	return 0
}

//export nacgosFreeStorage
func nacgosFreeStorage(storagePointer unsafe.Pointer) {
	//storage := gopointer.Restore(storagePointer).(*naconfig.Repo)
	gopointer.Unref(storagePointer)
}
