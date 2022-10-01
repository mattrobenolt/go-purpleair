package addsensor

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmd = &cobra.Command{
	Use: "add-sensor",
	Run: run,
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println("add sensor", viper.Get("read-key"))
}

func Init(root *cobra.Command) {
	root.AddCommand(cmd)
}
