Crontab 定時器專案建置展示.

## 前言
 Crontab定時器 : 採用robfig/cron/v3、config採用工廠方法(設計模型)、docker-compose容器啟動

## 檔案樹說明
```bash
├─deploy #容器化專案所需檔案放置處
├─internal #main.main()之私用package
│  └─system #OS相關調整包,例如:CPU核心數調整
├─docs #專案文檔放置處
├─configs #環境設定檔集
├─build #專案執行檔放置處(包括環境設定檔)
│  └─configs
└─pkg #main.main()之公用package
    ├─common #共用函式放置處
    │  └─customVar #字串輸出至任意變數包(struct繼承使用)
    ├─cronjob #定時任務放置處
    │  └─printExample #展示用打印測試
    └─global #共用全域變數放置處
        └─config #環境設定變數處理包
```

## 環境設定值
` ./configs/env.conf `
```bash
#UTF-8
# 設定CPU核心數 0=開啟所有核心,4=限制使用四顆核心,-1=所有核心數減一
cpu_core=0
# 打印測試初始值(Demo用)
print_example_index=10
```

## 效果展示
<!-- <img src="./diagram.svg">-->
![Alt text](./docs/Demo.gif)
