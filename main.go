package main

// 使用:: MailHog -smtp-bind-addr 0.0.0.0:1030 -api-bind-addr 127.0.0.1:8026 -ui-bind-addr 127.0.0.1:8026 我这样就行了

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xiaorui/reddit-async/reddit-backend/cmd"
	"github.com/xiaorui/reddit-async/reddit-backend/controller"
	"github.com/xiaorui/reddit-async/reddit-backend/dao/mysql"
	"github.com/xiaorui/reddit-async/reddit-backend/dao/redis"
	"github.com/xiaorui/reddit-async/reddit-backend/logger"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/async"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/console"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/rabbitmq"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/snowflake"
	"github.com/xiaorui/reddit-async/reddit-backend/settings"
	"go.uber.org/zap"
)

// @title			热点论坛
// @version		1.0
// @description	这是一个热点论坛项目， 能够根据当下热点来向用户展示论坛帖子
// @termsOfService	http://swagger.io/terms/
// @contact.name	xiaorui zheng
// @contact.url	http://www.swagger.io/support
// @contact.email	1298453249@qq.com
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:9000
// @BasePath		/api/v1
func main() {
	var rootCmd = &cobra.Command{
		Use:   "bluebell",
		Short: "[Start] bluebell...",
		Long:  `Default will run "serve" command, you can use "-h" flag to see all subcommands`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(os.Args) < 2 {
				os.Args = append(os.Args, "conf/config.yaml")
				// fmt.Println("Please set your config file!")
			}

			//1. 加载配置文件
			if err := settings.Init(os.Args[1]); err != nil {
				fmt.Printf("init settings failed, err:%v", err)
				return
			}
			//2 初始化日志文件
			if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
				fmt.Printf("init logger failed, err:%v", err)
				return
			}
			zap.L().Sync()
			zap.L().Debug("logger init success...")

			//4 初始化redis
			if err := redis.Init(settings.Conf.RedisConfig); err != nil {
				fmt.Printf("init redis failed, err:%v", err)
				return
			}

			//3. 初始化mysql
			if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
				fmt.Printf("init mysql failed, err:%v", err)
				return
			}

			//初始化雪花算法， 用于创建用户id
			if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
				fmt.Printf("snowflake.Init err:%v", err)
				return
			}

			// 初始化消费者
			go rabbitmq.Consumer()
			// TODO: 发起一个定时任务， 每周会生成当下的所有热点信息， 将热点信息投递给所有的已经订阅周报的邮箱， 默认订阅周报
			if err := async.SendWeekReport(); err != nil {
				fmt.Println("async.SendWeekReport error...")
				return
			}

			//注册gin中的validator校验器
			if err := controller.InitTrans("zh"); err != nil {
				fmt.Printf("controller.InitTrans err:%v", err)
				return
			}
		},
	}
	defer mysql.Close()
	defer redis.Close()
	fmt.Printf("你好，我是代币")
	// 注册子命令
	rootCmd.AddCommand(
		cmd.CmdServe,
	)

	// 配置默认运行 Web 服务， 就是说默认就会执行这个命令了
	cmd.RegisterDefaultCmd(rootCmd, cmd.CmdServe)

	// 注册全局参数，--env
	cmd.RegisterGlobalFlags(rootCmd)

	// 执行主命令
	if err := rootCmd.Execute(); err != nil {
		console.Exit(fmt.Sprintf("Failed to run app with %v: %s", os.Args, err.Error()))
	}
}
