package main

import (
	system "crontab/internal/system"
	"crontab/pkg/cronjob/printExample"
	"crontab/pkg/global/config"

	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

const (
	SERVICE_NAME        = "crontab"
	SERVICE_TIME_BEFORE = 1
	CONFIG_PATH         = "./configs/env.conf"
)

func main() {

	config.InitLoad(CONFIG_PATH) //載入環境變數至config.Cfgs

	//設定使用CPU核心數
	system.SetCpuCore(config.Cfgs.CpuCore)

	//crontab 操作緩衝時間(秒)
	for timer := SERVICE_TIME_BEFORE; timer > 0; timer-- {
		log.Printf("剩餘 %v 秒後開始執行任務！\n", timer)
		<-time.After(time.Second * 1)
	}

	//cronJob主體，設定定時器忽略已存在任務(若定時任務還在執行中則會忽略，待下輪在執行)
	cronJob := cron.New(cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger),
	))

	//重新讀取環境變數設置
	cronAddJob(cronJob, fmt.Sprintf("@every %ds", 10), config.ConfigTitle, config.Reload, false)

	cronAddJob(cronJob, fmt.Sprintf("@every %ds", 5), "打印測試", func() {
		printExample.StartJob()
	}, true)

	cronJob.Start()
	select {} //永久阻塞，避免gorotinue提前結束

}

/*
添加定時器任務排程
@param cronJob: *cron.Cron
@param spec: 定時器執行時間/周期(參考robfig說明)
@param title:  給予任務名稱
@param cmd:  任務具體執行函式
@param runBeforeTask:  是否在任務排程前先執行一次
*/
func cronAddJob(cronJob *cron.Cron, spec, title string, cmd func(), runBeforeTask bool) {
	if runBeforeTask {
		go cmd()
	}

	cronJob.AddFunc(spec, func() {
		log.Printf("%v : 開始執行時間戳(ms) => %v", title, time.Now().UnixNano()/1e6)
		cmd()
		log.Printf("%v : 結束執行時間戳(ms) => %v", title, time.Now().UnixNano()/1e6)
	})
}

// @yearly (or @annually)→1月1日午夜运行一次→"0 0 0 1 1 *""
// @monthly→每个月的午夜，每个月的第一个月运行一次→"0 0 0 1  ""
// @weekly→每周一次，周日午夜运行一次→"0 0 0   0"
// @daily (or @midnight)→每天午夜运行一次→"0 0 0   *""
// @hourly→每小时运行一次→"0 0  ""
// @every 1h30m10s→每小时30分又10秒後執行一次
//c.AddFunc("CRON_TZ=Asia/Tokyo 30 04 * * *", func() { fmt.Println("Runs at 04:30 Tokyo time every day") })

// seconds = bounds{0, 59, nil}
// minutes = bounds{0, 59, nil}
// hours   = bounds{0, 23, nil}
// dom     = bounds{1, 31, nil}
// months  = bounds{1, 12, map[string]uint{
// 	"jan": 1,
// 	"feb": 2,
// 	"mar": 3,
// 	"apr": 4,
// 	"may": 5,
// 	"jun": 6,
// 	"jul": 7,
// 	"aug": 8,
// 	"sep": 9,
// 	"oct": 10,
// 	"nov": 11,
// 	"dec": 12,
// }}
// dow = bounds{0, 6, map[string]uint{
// 	"sun": 0,
// 	"mon": 1,
// 	"tue": 2,
// 	"wed": 3,
// 	"thu": 4,
// 	"fri": 5,
// 	"sat": 6,
// }}

// 與Linux 中crontab命令相似，cron庫支持用 5 個空格分隔的域來表示時間。這 5 個域含義依次為：
// Minutes：分鐘，取值範圍[0-59]，支持特殊字符* / , -；
// Hours：小時，取值範圍[0-23]，支持特殊字符* / , -；
// Day of month：每月的第幾天，取值範圍[1-31]，支持特殊字符* / , - ?；
// Month：月，取值範圍[1-12]或者使用月份名字縮寫[JAN-DEC]，支持特殊字符* / , -；
// Day of week：週曆，取值範圍[0-6]或名字縮寫[JUN-SAT]，支持特殊字符* / , - ?。
// 注意，月份和周歷名稱都是不區分大小寫的，也就是說SUN/Sun/sun表示同樣的含義（都是周日）。

// 特殊字符含义如下：
// *：使用*的域可以匹配任何值，例如將月份域（第 4 個）設置為*，表示每個月；
// /：用來指定範圍的步長，例如將小時域（第2 個）設置為3-59/15表示第3 分鐘觸發，以後每隔15 分鐘觸發一次，因此第2 次觸發為第18分鐘，第3 次為33 分鐘。 。 。直到分鐘大於 59；
// ,：用來列舉一些離散的值和多個範圍，例如將周歷的域（第 5 個）設置為MON,WED,FRI表示週一、三和五；
// -：用來表示範圍，例如將小時的域（第 1 個）設置為9-17表示上午 9 點到下午 17 點（包括 9 和 17）；
// ?：只能用在月曆和周歷的域中，用來代替*，表示每月/週的任意一天。
// 了解規則之後，我們可以定義任意時間：

// 30 * * * *：分鐘域為 30，其他域都是*表示任意。每小時的 30 分觸發；
// 30 3-6,20-23 * * *：分鐘域為 30，小時域的3-6,20-23表示 3 點到 6 點和 20 點到 23 點。 3,4,5,6,20,21,22,23 時的 30 分觸發；
// 0 0 1 1 *：1（第 4 個） 月 1（第 3 個） 號的 0（第 2 個） 時 0（第 1 個） 分觸發。

// 注意DelayIfStillRunning與SkipIfStillRunning是有本質上的區別的，前者DelayIfStillRunning只要時間足夠長，
// 所有的任務都會按部就班地完成，只是可能前一個任務耗時過長，導致後一個任務的執行時間推遲了一點。 SkipIfStillRunning會跳過一些執行。

// SkipIfStillRunning 會控制你所有的任務池中，在同一時刻只能有一個在運行，這樣會導致有些任務永遠無法被運行，這顯然不能滿足需求，我想要的是同一個任務只有一個在運行，而不同的任務各自不受影響。
// 而如果使用 DelayIfStillRunning 的話，不同任務可以並行，相同任務串行，衝突的任務會延遲執行，這會導致任務可能不是設置的時刻執行的，且延遲超過1分鐘會有日誌。
