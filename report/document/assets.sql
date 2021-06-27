/*用户users表*/
CREATE TABLE users
(
    id   INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    username  varchar(20)  NOT NULL UNIQUE,
    password  varchar(200) NOT NULL,
    nick_name varchar(100),
    phone    varchar(20),
    email    varchar(100),
    remark    text,
    create_at timestamp
);

/*用户登录ip锁定表*/
CREATE TABLE user_ip
(
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    ip varchar(50) NOT NULL UNIQUE,
    lock_count int,
    create_at timestamp
);

/*用户登录日志表*/
CREATE TABLE user_log
(
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    username varchar(20) NOT NULL,
    ip varchar(50) NOT NULL,
    user_agent text,
    create_at timestamp
);

/*用户操作记录表*/
CREATE TABLE user_operation
(
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    username varchar(20) NOT NULL,
    ip varchar(50) NOT NULL,
    theme varchar(200) NOT NULL,
    content text NOT NULL,
    create_at timestamp
);

/*终端信息表*/
CREATE TABLE assets_computer
(
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    department varchar(200), # 员工所属部门
    department_sub varchar(200), # 员工所属二级部门
    person_name varchar(50), # 员工姓名
    work_number varchar(50), # 工号
    computer_type varchar(100), # 计算机类型
    computer_name varchar(100), # 计算机名
    secret_level varchar(100), # 涉密级别
    address varchar(50), # 计算机ip地址或mac地址
    internet_flag varchar(20), # 互联网权限
    file_copy_flag varchar(20), # 文件拷贝权限
    email_flag varchar(20), # 外网邮件权限
    vpn_flag varchar(20), # VPN权限
    pdm_flag varchar(20), #
    remarks text, # 备注
    create_at timestamp
);

/*员工信息表*/
CREATE TABLE assets_users
(
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    work_number varchar(30),
    id_card varchar(30),
    phone varchar(30), # 手机号
    department varchar(500), # 部门
    worker_type varchar(50), # 员工类型
    sex varchar(5), # 性别
    entry_reason varchar(20), # 入职途径
    entry_date varchar(20), # 入职时间
    port_name varchar(100), # 职位
    user_id varchar(100),
    create_at timestamp
);

/*安全管理员表*/
CREATE TABLE assets_manager
(
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    manager_name varchar(20) NOT NULL UNIQUE,
    create_at timestamp
);

/*主机资产表*/
CREATE TABLE assets_type
(
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    type_name varchar(100) NOT NULL, # 厂商
    attribution varchar(500) UNIQUE, # 应用系统 xx系统
    department varchar(200), # 资产对应部门
    assets_username varchar(100), # 资产对应管理员
    subdomain text, # 子域名
    intranet_ip varchar(500), # 业务系统对应内网ip地址 可多个
    public_ip varchar(500), # 业务系统对应公网地址 可多个
    create_at timestamp
);

/*业务系统资产表*/
CREATE TABLE assets_web
(
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    attribution varchar(500), # 应用系统 xx系统
    manager_name varchar(30) NOT NULL, # 安全管理员
    assets_name varchar(500) NOT NULL UNIQUE, # 业务系统名 唯一
    urls text NOT NULL, # 业务系统URL地址
    fingerprint text, # 业务系统指纹：框架、开发语言、通用应用程序等
    webserver text, # web服务器类型：tomcat、apache、nginx等
    screenshots_path varchar(100), # 业务系统截图地址
    remarks text, # 备注：业务管理员/联系方式/测试账户
    file_name varchar(100), # 业务系统资产附件/拓扑图
    create_at timestamp
);

/*渗透测试报告表*/
CREATE TABLE assets_reports
(
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    attribution varchar(500), # 应用系统 xx系统
    manager_name varchar(20) NOT NULL, # 安全管理员
    assets_name varchar(200) NOT NULL, # 业务系统名 唯一  与web资产表关联
    level integer NOT NULL,# 漏洞等级-1高2中3低
    level_name varchar(500) NOT NULL, # 漏洞名称
    level_status integer NOT NULL, # 漏洞修复状态-1已整改2未整改3已临时关闭
    file_path varchar(100), # 渗透测试报告保存地址
    file_date DATE, # 渗透测试编写时间
    create_at timestamp
);