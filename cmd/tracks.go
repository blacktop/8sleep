/*
Copyright © 2025 blacktop

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/blacktop/clim8/pkg/eightsleep"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// tracksCmd represents the tracks command
var tracksCmd = &cobra.Command{
	Use:   "tracks",
	Short: "List audio tracks",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetBool("verbose") {
			log.SetLevel(log.DebugLevel)
		}

		cli, err := eightsleep.NewClient(
			viper.GetString("email"),
			viper.GetString("password"),
			"America/New_York",
		)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		defer cli.Stop()

		if err := cli.Start(cmd.Context()); err != nil {
			return fmt.Errorf("failed to start client: %w", err)
		}

		tracks, err := cli.GetAudioTracks(cmd.Context())
		if err != nil {
			return err
		}

		logger.Info("AUDIO TRACKS")
		jsonData, err := json.MarshalIndent(tracks, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal json: %v", err)
		}
		if err := quick.Highlight(os.Stdout, string(jsonData)+"\n", "json", "terminal256", "nord"); err != nil {
			return fmt.Errorf("failed to highlight json: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(tracksCmd)
}
