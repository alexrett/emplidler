package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	//"github.com/alexrett/active-window"
	"github.com/alexrett/emplidler/icon"
	"github.com/alexrett/idler"
	"github.com/getlantern/systray"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

var (
	myFingerPrint        FingerPrint
	hashKey              string
	idlerInstance        *idler.Idle
	idleStartFromSecond  int
	idleCountSleepSecond int
	ServerUrl            string
	GitCommit            string
	AppKey               string
)

func init() {

}

func main() {
	myFingerPrint = GetFingerPrint()
	go encodeHashKey(myFingerPrint)

	if cliCommand() {
		return
	}

	idlerInstance = idler.NewIdle()
	idleStartFromSecond = 10
	idleCountSleepSecond = 1
	onExit := func() {
		//appendStatistic()
		sendMetricIdle()
	}

	go func() {
		for {
			if hashKey != "" {
				sendMetricActive()
				break
			}
			time.Sleep(1 * time.Second)
		}
	}()

	go loopIdler()
	// Should be called at the very beginning of main().
	systray.RunWithAppWindow("Emplidler", 1024, 524, onReady, onExit)

}

func cliCommand() bool {
	args := os.Args[1:]
	if len(args) == 1 {
		switch args[0] {
		case "version":
			fmt.Println(GitCommit)
			return true
		case "mac":
			fmt.Println(myFingerPrint.Mac)
			return true
		case "ip":
			fmt.Println(myFingerPrint.Ip)
			return true
		case "iface":
			fmt.Println(myFingerPrint.Interface)
			return true
		case "username":
			fmt.Println(myFingerPrint.Username)
			return true
		case "name":
			fmt.Println(myFingerPrint.Name)
			return true
		case "machineid":
			fmt.Println(myFingerPrint.MachineId)
			return true
		case "help":
			help := `Emplidler 
-----
version - show current application build version
mac - show your detected mac
ip - your public ip
iface - interface for public ip
username - detected username
name - detected name
machineid - your machine id
help - this help
`
			fmt.Println(help)
			return true

		}
	}

	return false
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	mUrl := systray.AddMenuItem("Your statistic", "12 hours")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	systray.AddMenuItem(GitCommit, "version")
	for {
		select {
		case <-mUrl.ClickedCh:
			url := ServerUrl + "/api/v1/pull?h=" + hashKey
			switch runtime.GOOS {
			case "linux":
				err := exec.Command("xdg-open", url).Start()
				if err != nil {
					// todo: need log error
				}
			case "windows":
				systray.ShowAppWindow(url)
			case "darwin":
				systray.ShowAppWindow(url)
			}

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
	url := fmt.Sprintf(ServerUrl+"/api/v1/push?h=%s&t=%d&a=1", hashKey, time.Now().Unix())
	_, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	//appendToIdleFile(1)
}

func sendMetricIdle() {
	_, err := http.Get(fmt.Sprintf(ServerUrl+"/api/v1/push?h=%s&t=%d&a=0", hashKey, time.Now().Unix()))
	if err != nil {
		log.Println(err)
	}
	//appendToIdleFile(0)
}

func encodeHashKey(dataObj FingerPrint) string {
	data := fmt.Sprintf("%s:%s:%s:%s%s", dataObj.MachineId, dataObj.Mac, dataObj.Uid, dataObj.Username, AppKey)

	sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	jsonValue, _ := json.Marshal(&dataObj)
	resp, err := http.Post(ServerUrl+"/api/v1/hello?hash="+sEnc, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal(err)
	}
	if resp != nil && resp.StatusCode != 200 {
		log.Fatal(resp.StatusCode)
	}

	hashKey = sEnc

	return sEnc
}
