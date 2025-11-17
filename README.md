<!--
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-17 20:18:14
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2024-10-29 15:38:59
 * @FilePath: \go_809_converter\README.md
 * @Description: 
 * 
-->

# go_809_converter

把第三方推送的json数据，转换成809协议

安装809官方规定，如果主链路断开，可以使用从链路发送数据。本程序没有实现此功能。链路端口只有等待重连，重连成功后继续推送数据。

## 使用说明

### 日志输出配置

| 环境变量名                      | 介绍                                                                                                                                                  |
|----------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------|
| Name                       | 服务名称，给prometheus提供命名空间                                                                                                                              |
| GIN_MODE                   | gin 模块日志输出。release 用于生产环境; debug(默认值)用于开发环境                                                                                                         |
| LOGLEVEL                   | 最低日志显示级别 debug info warn error 。对日志级别影响。优先级高                                                                                                        |
| IS_DEV                     | 是否是开发环境.是开发环境,会把日志默认输出到 <br>`app.log`(普通日志) +<br>`app.error.log`(错误日志) +<br>`控制台`<br>默认是 false<br>开发环境 默认日志级别是debug<br>非开发环境 默认日志级别是info，且日志不会输出到文件 |
| LOGFILE                    | 日志文件名，仅在 IS_DEV=true时有用                                                                                                                             |
| LOG_IN_JSON                | 是否以json格式输出，默认false。json格式和非json格式输出的区别，请运行 `zloger_test.go`中的测试用例                                                                                  |
| EXTRA_TRACER               | 日志额外跟踪字段                                                                                                                                            |
| CLOSE_SERVER_INFO          | 日志中是否关闭服务器和span信息，以节省日志空间。默认是false。 建议，部署多实例的情况下，不要设置成 true                                                                                         |
| CONCISE_BLANK_TRACE        | 日志中是否精简空trace输出。如果是true，则当traceId全0时替换显示a,spanId全0时显示b。                                                                                             |
| SLS_OTEL_SERVICE_VERSION   | 阿里云日志服务版本号(可以不配)                                                                                                                                    |
| SLS_OTEL_PROJECT           | 阿里云链路跟踪项目名称(可以不配)                                                                                                                                   |
| SLS_OTEL_INSTANCE_ID       | 阿里云链路跟踪实例ID(可以不配)                                                                                                                                   |
| SLS_OTEL_TRACE_ENDPOINT    | 阿里云链路跟踪grpc入口(可以不配)                                                                                                                                 |
| SLS_OTEL_METRIC_ENDPOINT   | 阿里云链路跟踪监控指标入口(可以不配)                                                                                                                                 |
| SLS_OTEL_ACCESS_KEY_ID     | 阿里云链路跟踪服务的访问key(可以不配)                                                                                                                               |
| SLS_OTEL_ACCESS_KEY_SECRET | 阿里云链路跟踪服务的访问密码(可以不配)                                                                                                                                |

### mysql 配置(环境变量优先级高)

| 环境变量名                 | toml配置键                  | 介绍                 | 默认值       |
|-----------------------|--------------------------|--------------------|-----------|
| MYSQL_HOST            | mysql_db.host            | MySQL 主机地址         | localhost |
| MYSQL_PORT            | mysql_db.port            | MySQL访问端口          | 3306      |
| MYSQL_USER            | mysql_db.user            | MySQL用户名           | root      |
| MYSQL_PWD             | mysql_db.password        | MySQL密码            | root      |
| MYSQL_DATABASE        | mysql_db.database        | MySQL 访问的数据库       | ""        |
| MYSQL_POOL_SIZE       | mysql_db.pool_size       | MySQL 连接池中连接个数     | 1         |
| MYSQL_POOL_IDLE_CONNS | mysql_db.pool_idle_conns | MySQL 空闲连接个数       | 10        |
| MYSQL_SHOW_SQL        | mysql_db.showSQL         | MySQL 是否打印执行的SQL语句 | false     |

### 服务链接配置

| toml配置键                             | 介绍                                  | 默认值     |
|-------------------------------------|-------------------------------------|---------|
| {env}.converter.consolePort         | web访问和prometheus监控端口                | 13031   |
| {env}.converter.cryptoPacket        | 需要加密的报文数组，可多选 S10,S13,S99,S991,S106 | []      |
| {env}.converter.encryptKey          | 报文加密使用的key,不使用加密时可以不配置              | ""      |
| {env}.converter.extendVersion       | 是否使用809扩展版的协议                       | false   |
| {env}.converter.govServerIP         | 上级服务IP地址或域名                         | ""      |
| {env}.converter.govServerPort       | 上级服务端口                              | 0       |
| {env}.converter.localServerIP       | 本地对外暴露的ip,上级服务能够通过此ip连接到本服务         | ""      |
| {env}.converter.localServerPort     | 本地对外暴露的端口，上级服务能够通过此端口连接到本服务         | 0       |
| {env}.converter.openCrypto          | 是否对所有报文进行加密                         | false   |
| {env}.converter.platformId          | 本平台连接上级服务时传递的ID                     | 0       |
| {env}.converter.platformPassword    | 本平台连接上级服务时传递的连接密码                   | ""      |
| {env}.converter.platformUserId      | 本平台连接上级服务时传递的用户ID                   | 0       |
| {env}.converter.protocolVersion     | 本平台连接上级服务时使用的协议版本                   | "1.0.0" |
| {env}.converter.thirdpartPort       | 本服务接受第三方推送连接的端口                     | 11223   |
| {env}.converter.useLocationInterval | 是否有分钟只推送一个位置报文给上级平台                 | false   |

### 提示

如果需要使用docker部署，最好是外挂 config/configuration.toml 和 config/config.history 
