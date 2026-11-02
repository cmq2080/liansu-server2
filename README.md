# 连速Server2

连速Server2（liansu-server2）是一个轻量级PHP集成环境可视化管理工具。现仅有Windows版本，可通过GUI管理。

## 起源&历程

本项目源于XServer计划（详见前版本项目[liansu-server](https://github.com/cmq2080/liansu-server)），于2026年10月30日正式立项。

GUI部分使用govcl开发

## 用法

-- * --gui: 启动图形界面窗口（仅适用于Windows版本）
-- * start: 启动环境（包括启动Nginx和php-cgi）
-- * stop: 终止环境（包括终止Nginx和php-cgi的运行）
-- * restart: 重启环境
-- * save: 保存当前配置
-- * export: 导出当前Nginx应用配置，重新生成Nginx扩展配置文件（先保存再导出）
-- * status: 查看当前环境运行状态
-- * -h/--help: 获取帮助信息