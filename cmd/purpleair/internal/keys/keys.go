package keys

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go.withmatt.com/purpleair"
)

var cmd = &cobra.Command{
	Use: "keys",
	Run: func(cmd *cobra.Command, args []string) {
		keys := map[string]*purpleair.KeysResponse{}
		var wg sync.WaitGroup
		var mux sync.Mutex
		for _, keyType := range []string{"read-key", "write-key"} {
			if key := viper.GetString(keyType); key != "" {
				wg.Add(1)
				go func() {
					defer wg.Done()
					kr, err := getKey(key)
					cobra.CheckErr(err)
					mux.Lock()
					keys[key] = kr
					mux.Unlock()
				}()
			}
		}
		wg.Wait()
		cobra.CheckErr(json.NewEncoder(os.Stdout).Encode(keys))
	},
}

func getKey(key string) (*purpleair.KeysResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return purpleair.New(key).Keys(ctx)
}

func Init(root *cobra.Command) {
	root.AddCommand(cmd)
}
