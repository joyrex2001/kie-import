package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/joyrex2001/kie-import/clone"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "kie-import",
	Short: "KIE import will import a given Drools Workbench repo.",
	Long: `KIE import is tool that will allow automation of creating KIE servers,
created with the Drools Workbench, in an OpenShift environment.

The Drools Workbench contains an internal git repository that can be cloned
and, with small adjustments on the maven pom file, can be compiled as a KI
server.

One of the issues that comes with cloning and automating this repo, is that
it is based on username / password combinations, and uses an legacy key
algorithm. This will make the configuration, and automation a bit more
challenging. KIE import will make this easier.`,
	Run: clone.Main,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().String("drools-host", "", "host of the drools git server")
	rootCmd.PersistentFlags().String("drools-git-ssh-port", "8001", "port of the drools ssh git server")
	rootCmd.PersistentFlags().String("drools-git-repo", "", "git repository to be cloned")
	rootCmd.PersistentFlags().String("user", "", "username for login to git repository")
	rootCmd.PersistentFlags().String("password", "", "password for login to git repository")
	rootCmd.PersistentFlags().String("destination", "/tmp/kie-import", "destination folder for cloned repo")
	viper.BindPFlag("drools.host", rootCmd.PersistentFlags().Lookup("drools-host"))
	viper.BindPFlag("drools.git_ssh_port", rootCmd.PersistentFlags().Lookup("drools-git-ssh-port"))
	viper.BindPFlag("drools.git_repo", rootCmd.PersistentFlags().Lookup("drools-git-repo"))
	viper.BindPFlag("git.user", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("git.password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("git.destination", rootCmd.PersistentFlags().Lookup("destination"))
	viper.BindEnv("drools.host", "DROOLS_HOST")
	viper.BindEnv("drools.git_ssh_port", "DROOLS_GIT_SSH_PORT")
	viper.BindEnv("drools.git_repo", "DROOLS_GIT_REPO")
	viper.BindEnv("git.user", "GIT_USERNAME")
	viper.BindEnv("git.password", "GIT_PASSWORD")
	viper.BindEnv("git.destination", "GIT_DESTINATION")
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name "kopenvoor" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		// fmt.Printf("not using config file: %s\n", err)
	} else {
		fmt.Printf("using config: %s\n", viper.ConfigFileUsed())
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
