package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cezkuj/gopage/gopage"
)

var (
	dbUser string
	dbPass string
	dbHost string
	dbName string
	prod   bool
)
var rootCmd = &cobra.Command{
	Use:   "gopage",
	Short: "Start a web server with siginig up/loggin in enabled",
	Long: `Start a web server with siginig up/loggin in enabled
        Examples:
        gopage`,
	Args: cobra.NoArgs,
	Run:  startServer,
}

func startServer(cmd *cobra.Command, args []string) {
	dbCfg := gopage.NewDbCfg(dbUser, dbPass, dbHost, dbName)
	gopage.StartServer(dbCfg, prod)

}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&dbUser, "user", "u", "", "Sets user for database conneciton, required")
	rootCmd.MarkFlagRequired("user")
	rootCmd.Flags().StringVarP(&dbPass, "pass", "p", "", "Sets password for database conneciton, required")
	rootCmd.MarkFlagRequired("pass")
	rootCmd.Flags().StringVarP(&dbHost, "host", "o", "", "Sets host for database conneciton, required")
	rootCmd.MarkFlagRequired("host")
	rootCmd.Flags().StringVarP(&dbName, "name", "n", "", "Sets name for database conneciton, required")
	rootCmd.MarkFlagRequired("name")
	rootCmd.Flags().BoolVarP(&prod, "prod", "t", false, "Sets production mode with tls enabled. Default value is false.")
}
