package main

import (
	"bufio"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/alexrett/active-window"
	"github.com/alexrett/emplidler/icon"
	"github.com/alexrett/idler"
	"github.com/getlantern/systray"
	"io/ioutil"
	"log"
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
)

func init() {
	myFingerPrint = GetFingerPrint()
	hashKey = encodeHashKey(myFingerPrint)
	idlerInstance = idler.NewIdle()
	idleStartFromSecond = 10
	idleCountSleepSecond = 1
	logName = "log.log"
	apptime = make(map[string]int)
}

func main() {
	var err error
	f, err = os.OpenFile(logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	readStatistic()
	go loopIdler()
	onExit := func() {
		appendStatistic()
	}
	// Should be called at the very beginning of main().
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("EmplIdler")
	systray.SetTooltip("EmplIdler")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()
}

func readStatistic() {
	jsonFile, err := os.Open(fmt.Sprintf("./%s.json", time.Now().Format("20060102")))
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &apptime)
	if err != nil {
		log.Println(err)
	}

	defer jsonFile.Close()
}

func appendStatistic() error {
	file, err := os.OpenFile(fmt.Sprintf("./%s.json", time.Now().Format("20060102")), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	w := bufio.NewWriter(file)
	data, _ := json.Marshal(apptime)
	_, err = w.Write(data)
	if err != nil {
		return err
	}

	return w.Flush()
}

func appendToIdleFile(t int) {
	ut := time.Now().Unix()
	switch t {
	case 0:
		if _, err := f.WriteString(fmt.Sprintf("%d|0\n", ut)); err != nil {
			log.Println(err)
		}
	case 1:
		if _, err := f.WriteString(fmt.Sprintf("%d|1\n", ut)); err != nil {
			log.Println(err)
		}
	}
}

func activeWindowCheck(id int) {
	activeWindowInstance := &activeWindow.ActiveWindow{}
	program, _ := activeWindowInstance.GetActiveWindowTitle()
	if val, ok := apptime[program]; ok {
		apptime[program] = val + id
	} else {
		apptime[program] = val + id
	}

	appendStatistic()
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

		time.Sleep(time.Duration(idleCountSleepSecond) * time.Second)
	}
}

func sendMetricActive() {
	appendToIdleFile(1)
}

func sendMetricIdle() {
	appendToIdleFile(0)
}

func encodeHashKey(dataObj FingerPrint) string {
	data := fmt.Sprintf("%s:%s:%s:%s", dataObj.MachineId, dataObj.Mac, dataObj.Uid, dataObj.Username)

	sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	return sEnc
}
