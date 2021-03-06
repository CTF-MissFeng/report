### 项目开发初衷

> 作为乙方渗透测试工程师，随着时间的推移，所提交的渗透测试报告大量增加，原来都是以文件目录形式保存渗透测试报告，但是存在诸多不便之处，故写了此脚本方便报告的管理。主要希望改进以下不便之处：

- 漏洞整改管理：管理还有哪些漏洞未整改

- 漏洞可视化：本周、本月等时间段漏洞情况、安全管理员属下漏洞情况等

- 渗透测试报告统一管理：快速查询

- 等等等

### 项目使用

```
1、创建mysql数据库名为assets（可自定义），格式为UTF-8,导入document/assets.sql文件

2、修改配置文件：config/config.toml，填入相应的mysql数据库信息

2、运行assets程序即可（默认用户名密码：admin/admin888@A）或者自行go build编译

linux：
前台运行：./assets
后台运行：nohup ./assets > web.log 2>&1 &

浏览器访问配置文件里的web地址即可
```



### 示例

> 不做详情解释了

#### 使用步骤

> 主要分为资产管理和业务系统模块，请先填写安全管理员和主机资产模块信息，终端资产可忽略

![index](https://github.com/CTF-MissFeng/report/blob/main/doc/1.png)

> 添加安全管理员信息

![index](https://github.com/CTF-MissFeng/report/blob/main/doc/2.png)

> 添加主机资产信息

![index](https://github.com/CTF-MissFeng/report/blob/main/doc/3.png)

> 当添加完主机资产后，可在业务系统添加主机资产对应的业务系统信息以及添加对应业务系统漏洞信息等

![index](https://github.com/CTF-MissFeng/report/blob/main/doc/4.png)

![index](https://github.com/CTF-MissFeng/report/blob/main/doc/5.png)

![index](https://github.com/CTF-MissFeng/report/blob/main/doc/7.png)

>漏洞管理

![index](https://github.com/CTF-MissFeng/report/blob/main/doc/6.png)
