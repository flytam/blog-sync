package cli

import (
	"blog-sync/core"
	"blog-sync/log"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

type Cli struct {
	config *Config
}

func NewCli() *Cli {
	cli := &Cli{
		config: &Config{},
	}

	cli.initCli()

	return cli
}

func (cli *Cli) initCli() {
	//
	logger := log.GetLogger()
	rootCmd := &cobra.Command{
		Use:   "hsync",
		Short: "csdn->hexo",
		Long:  "csdn一键生成hexo源文件工具",
		Run: func(cmd *cobra.Command, args []string) {
			// 执行逻辑

			var cookie, csdn, output string

			if len(cli.config.Config) > 0 {
				// 配置文件执行
				buf, err := ioutil.ReadFile(cli.config.Config)
				if err != nil {
					logger.Error(err.Error())
				}
				var config Config
				json.Unmarshal(buf, &config)
				fmt.Println(config)
				cookie = config.Cookie
				csdn = config.Csdn
				output = config.Output
			} else {
				// 提供参数执行
				cookie = cli.config.Cookie
				csdn = cli.config.Csdn
				output = cli.config.Output

				if len(cookie) == 0 || len(csdn) == 0 || len(output) == 0 {
					logger.Error("参数不足")
					os.Exit(1)
				}
			}

			(&core.Core{
				cookie, csdn, output,
			}).Run()
		},
	}

	// 命令行直接配置参数
	rootCmd.PersistentFlags().StringVar(&cli.config.Csdn, "csdn", "", "csdn用户名")
	rootCmd.PersistentFlags().StringVar(&cli.config.Output, "output", "", "输出目录")
	rootCmd.PersistentFlags().StringVar(&cli.config.Cookie, "cookie", "", "cookie信息")

	// 指定文件目录
	rootCmd.PersistentFlags().StringVar(&cli.config.Config, "config", "", "配置文件目录")

	rootCmd.Execute()
}
