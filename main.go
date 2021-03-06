package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/parnurzeal/gorequest"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var isListen = false
var body *widget.Label
var tip *widget.Label
var status *widget.Label

var models = []string{
	"iphone12mini 64gb 黑色-MG7Y3CH/A",
	"iphone12mini 64gb 白色-MG803CH/A",
	"iphone12mini 64gb 蓝色-MG823CH/A",
	"iphone12mini 64gb 绿色-MG833CH/A",
	"iphone12mini 64gb 红色-MG813CH/A",
	"iphone12mini 128gb 黑色-MG843CH/A",
	"iphone12mini 128gb 白色-MG853CH/A",
	"iphone12mini 128gb 蓝色-MG873CH/A",
	"iphone12mini 128gb 绿色-MG883CH/A",
	"iphone12mini 128gb 红色-MG863CH/A",
	"iphone12mini 256gb 黑色-MG893CH/A",
	"iphone12mini 256gb 白色-MG8A3CH/A",
	"iphone12mini 256gb 蓝色-MG8D3CH/A",
	"iphone12mini 256gb 绿色-MG8E3CH/A",
	"iphone12mini 256gb 红色-MG8C3CH/A",
	"iphone12 64gb 黑色-MGGM3CH/A",
	"iphone12 64gb 白色-MGGN3CH/A",
	"iphone12 64gb 蓝色-MGGQ3CH/A",
	"iphone12 64gb 绿色-MGGT3CH/A",
	"iphone12 64gb 红色-MGGP3CH/A",
	"iphone12 128gb 黑色-MGGU3CH/A",
	"iphone12 128gb 白色-MGGV3CH/A",
	"iphone12 128gb 蓝色-MGGX3CH/A",
	"iphone12 128gb 绿色-MGGY3CH/A",
	"iphone12 128gb 红色-MGGW3CH/A",
	"iphone12 256gb 黑色-MGH13CH/A",
	"iphone12 256gb 白色-MGH23CH/A",
	"iphone12 256gb 蓝色-MGH43CH/A",
	"iphone12 256gb 绿色-MGH53CH/A",
	"iphone12 256gb 红色-MGH33CH/A",
	"iphone12pro 128gb 石墨色-MGL93CH/A",
	"iphone12pro 128gb 银色-MGLA3CH/A",
	"iphone12pro 128gb 金色-MGLC3CH/A",
	"iphone12pro 128gb 海蓝色-MGLD3CH/A",
	"iphone12pro 256gb 石墨色-MGLE3CH/A",
	"iphone12pro 256gb 银色-MGLF3CH/A",
	"iphone12pro 256gb 金色-MGLG3CH/A",
	"iphone12pro 256gb 海蓝色-MGLH3CH/A",
	"iphone12pro 512gb 石墨色-MGLJ3CH/A",
	"iphone12pro 512gb 银色-MGLK3CH/A",
	"iphone12pro 512gb 金色-MGLL3CH/A",
	"iphone12pro 512gb 海蓝色-MGLM3CH/A",
	"iphone12promax 128gb 石墨色-MGC03CH/A",
	"iphone12promax 128gb 银色-MGC13CH/A",
	"iphone12promax 128gb 金色-MGC23CH/A",
	"iphone12promax 128gb 海蓝色-MGC33CH/A",
	"iphone12promax 256gb 石墨色-MGC43CH/A",
	"iphone12promax 256gb 银色-MGC53CH/A",
	"iphone12promax 256gb 金色-MGC63CH/A",
	"iphone12promax 256gb 海蓝色-MGC73CH/A",
	"iphone12promax 512gb 石墨色-MGC93CH/A",
	"iphone12promax 512gb 银色-MGCA3CH/A",
	"iphone12promax 512gb 金色-MGCC3CH/A",
	"iphone12promax 512gb 海蓝色-MGCE3CH/A",
}

var selectStore string
var selectModel string
var listenStores map[string]string

func main() {
	// 本地测试时手动加载字体
	// loadFont()

	a := app.NewWithID("ip12")
	// 打包时自动加载字体
	a.Settings().SetTheme(&myTheme{})
	w := a.NewWindow("iPhone12 Pro")
	w.Resize(fyne.NewSize(610, 500))

	body = widget.NewLabel("")
	tip = widget.NewLabel("请选择门店和型号")
	status = widget.NewLabel("暂停")
	stores := stores()

	listenStores = make(map[string]string)

	w.SetContent(widget.NewVBox(
		widget.NewLabel("1.首次运行请先点击测试按钮，确保能正确打开网页\n" +
			"2.选择门店和型号，点击添加按钮\n" +
			"3.点击开始\n" +
			"4.匹配到之后会直接进入门店预购页面，输入注册码选择预约时间即可",
		),
		widget.NewHBox(
			widget.NewSelect(stores, func(b string) {
				selectStore = b
			}),
			widget.NewSelect(models, func(b string) {
				selectModel = b
			}),
			widget.NewButton("添加", func() {
				if selectModel != "" && selectStore != "" {
					md := strings.Split(selectStore, "-")[0]+"."+strings.Split(selectModel, "-")[1]
					mdText := strings.Split(selectStore, "-")[1]+" "+strings.Split(selectModel, "-")[0]

					if !inArray(listenStores, md) {
						listenStores[md] = mdText
					}

					body.SetText(strings.Join(getValues(listenStores), "\n"))
				}
			}),
			widget.NewButton("清空", func() {
				listenStores = map[string]string{}
				body.SetText("")
			}),
		),
		tip,
		body,
		layout.NewSpacer(),
		widget.NewHBox(
			widget.NewButton("开始", func() {
				if len(listenStores) < 1 {
					tip.SetText("请添加要监听的门店和型号")
					return
				}

				isListen = true
				status.SetText("监听中")
			}),
			widget.NewButton("暂停", func() {
				isListen = false
				status.SetText("暂停")
			}),
			widget.NewButton("测试", func() {
				go openBrowser()
			}),
			widget.NewButton("退出", func() {
				a.Quit()
			}),
			layout.NewSpacer(),
			widget.NewLabel("状态: "),
			status,
		),
	))
	go listen()
	w.ShowAndRun()
	_ = os.Unsetenv("FYNE_FONT")
}

func listen() {
	// 库存接口
	availability := "https://reserve-prime.apple.com/CN/zh_CN/reserve/A/availability.json"

	for  {
		time.Sleep(time.Second*1)

		if !isListen {
			continue
		}

		_, bd, errs := gorequest.New().Get(availability).End()
		if len(errs) != 0 {
			log.Println(errs)
			continue
		}
		str := ""
		t := time.Now().Format("2006-01-02 15:04:05")
		for model, title := range listenStores {
			value := gjson.Get(bd, "stores."+model+".availability")
			if value.Map()["contract"].Bool() && value.Map()["unlocked"].Bool() {
				openBrowser()

				tip.SetText("已匹配到: " + title+ ", 暂停监听")
				status.SetText("暂停")
				isListen = false
			} else {
				str += t+" "+title+"无货\n"
			}
		}
		body.SetText(str)
	}
}

func openBrowser(args ...string) {
	var err error
	var url string
	if len(args) > 0 {
		params := strings.Split(args[0], ".")
		url = "https://reserve-prime.apple.com/CN/zh_CN/reserve/F?quantity=1&anchor-store="+params[0]+
			"&store="+params[0]+"&partNumber="+params[1]+"&plan=unlocked"
	} else {
		url = "https://reserve-prime.apple.com/CN/zh_CN/reserve/A/availability?iUP=N"
	}

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatalln("自动打开网页失败，请自行手动操作")
	}
}

func loadFont() {
	//dir, err := os.Getwd()
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	//switch runtime.GOOS {
	//case "linux":
	//	path = "./zh-cn.ttf"
	//case "windows":
	//	path = "./zh-cn.ttf"
	//case "darwin":
	//	path = "./zh-cn.ttf"
	//default:
	//	path = "./zh-cn.ttf"
	//}

	_ = os.Setenv("FYNE_FONT", dir+"/zh-cn.ttf")
}

func stores() []string {
	// 门店列表
	var stores []string

	availability := "https://reserve-prime.apple.com/CN/zh_CN/reserve/A/stores.json"
	_, bd, errs := gorequest.New().Get(availability).End()
	if len(errs) != 0 {
		log.Fatalln(errs)
	}

	for _, store := range gjson.Get(bd, "stores").Array() {
		str := store.Get("storeNumber").String()+
			"-"+store.Get("city").String()+
			" "+store.Get("storeName").String()

		stores = append(stores, str)
	}

	return stores
}

func inArray(slice map[string]string, s string) bool {
	for key := range slice {
		if key == s {
			return true
		}
	}
	return false
}

func getValues(slice map[string]string) []string {
	var values []string
	for _, value := range slice {
		values = append(values, value)
	}

	return values
}