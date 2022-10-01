package sensors

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.withmatt.com/purpleair"
)

func MustInt(index int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if _, err := strconv.Atoi(args[index]); err != nil {
			return fmt.Errorf("argument %q must be in integer", args[index])
		}
		return nil
	}
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func split(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

var rootCmd = &cobra.Command{
	Use: "sensors",
}

var getCmd = &cobra.Command{
	Use: "get [sensor index]",
	Args: cobra.MatchAll(
		cobra.ExactArgs(1),
		MustInt(0),
	),
	Run: func(cmd *cobra.Command, args []string) {
		index := atoi(args[0])
		client := purpleair.New(viper.GetString("read-key"))
		resp, err := client.SensorData(context.Background(), index, &purpleair.SensorDataRequestOptions{
			Fields:   split(fieldsFlag),
			ReadKeys: split(readkeysFlag),
		})
		cobra.CheckErr(err)
		cobra.CheckErr(json.NewEncoder(os.Stdout).Encode(resp))
	},
}

var historyCmd = &cobra.Command{
	Use: "history [sensor index]",
	Args: cobra.MatchAll(
		cobra.ExactArgs(1),
		MustInt(0),
	),
	Run: func(cmd *cobra.Command, args []string) {
		index := atoi(args[0])
		client := purpleair.New(viper.GetString("read-key"))
		resp, err := client.SensorHistory(context.Background(), index, &purpleair.SensorHistoryRequestOptions{
			Fields:   split(fieldsFlag),
			ReadKeys: split(readkeysFlag),
		})
		cobra.CheckErr(err)
		cobra.CheckErr(json.NewEncoder(os.Stdout).Encode(resp))
	},
}

var (
	fieldsFlag   string
	readkeysFlag string
)

func Init(root *cobra.Command) {
	getCmd.Flags().StringVar(&fieldsFlag, "fields", "", "fields to select")
	historyCmd.Flags().StringVar(&fieldsFlag, "fields", "", "fields to select")
	getCmd.Flags().StringVar(&readkeysFlag, "read-keys", "", "read keys")
	historyCmd.Flags().StringVar(&readkeysFlag, "read-keys", "", "read keys")

	rootCmd.AddCommand(
		getCmd,
		historyCmd,
	)
	root.AddCommand(rootCmd)
}
