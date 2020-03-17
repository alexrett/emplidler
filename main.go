package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	//"github.com/alexrett/active-window"
	"github.com/alexrett/emplidler/icon"
	"github.com/alexrett/idler"
	"github.com/getlantern/systray"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

var (
	myFingerPrint        FingerPrint
	hashKey              string
	idlerInstance        *idler.Idle
	idleStartFromSecond  int
	idleCountSleepSecond int
	f                    *os.File
	logName              string
	apptime              map[string]int
	serverUrl            string
)

func init() {

}

func main() {
	myFingerPrint = GetFingerPrint()
	go encodeHashKey(myFingerPrint)
	idlerInstance = idler.NewIdle()
	idleStartFromSecond = 10
	idleCountSleepSecond = 1
	logName = "log.log"
	apptime = make(map[string]int)
	serverUrl = "" // your url here
	onExit := func() {
		//appendStatistic()
		sendMetricIdle()
	}
	go sendMetricActive()
	go loopIdler()
	// Should be called at the very beginning of main().
	systray.RunWithAppWindow("Emplidler", 1024, 524, onReady, onExit)

}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	mUrl := systray.AddMenuItem("Your statistic", "12 hours")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	for {
		select {
		case <-mUrl.ClickedCh:
			systray.ShowAppWindow(serverUrl + "/api/v1/pull?h=" + hashKey)
		case <-mQuit.ClickedCh:
			fmt.Println("Quit now...")
			systray.Quit()
			return
		}
	}

}

func activeWindowCheck(id int) {
	//activeWindowInstance := &activeWindow.ActiveWindow{}
	//program, _ := activeWindowInstance.GetActiveWindowTitle()
	//if val, ok := apptime[program]; ok {
	//	apptime[program] = val + id
	//} else {
	//	apptime[program] = val + id
	//}

	//appendStatistic()
}

func loopIdler() {
	idleTime := 0
	for {
		idleTimeCur := idlerInstance.GetIdleTime()
		if idleTimeCur > idleStartFromSecond {
			if idleTime == 0 {
				sendMetricIdle()
				debug.FreeOSMemory()
			}
			idleTime = idleTimeCur
		} else {
			if idleTime != 0 {
				activeWindowCheck(idleTime)
				sendMetricActive()
			} else {
				activeWindowCheck(idleCountSleepSecond)
			}
			idleTime = 0
		}
		time.Sleep(10 * time.Second)
	}
}

func sendMetricActive() {
	url := fmt.Sprintf(serverUrl+"/api/v1/push?h=%s&t=%d&a=1", hashKey, time.Now().Unix())
	_, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	//appendToIdleFile(1)
}

func sendMetricIdle() {
	_, err := http.Get(fmt.Sprintf(serverUrl+"/api/v1/push?h=%s&t=%d&a=0", hashKey, time.Now().Unix()))
	if err != nil {
		log.Println(err)
	}
	//appendToIdleFile(0)
}

func encodeHashKey(dataObj FingerPrint) string {
	data := fmt.Sprintf("%s:%s:%s:%s", dataObj.MachineId, dataObj.Mac, dataObj.Uid, dataObj.Username)

	sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	jsonValue, _ := json.Marshal(&dataObj)
	resp, err := http.Post(serverUrl+"/api/v1/hello?hash="+sEnc, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal(err)
	}
	if resp != nil && resp.StatusCode != 200 {
		log.Fatal(resp.StatusCode)
	}

	hashKey = sEnc

	return sEnc
}
