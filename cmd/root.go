/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gordonklaus/portaudio"
	"github.com/raspiantoro/vocafex/pkg/audio/pipeline"
	"github.com/raspiantoro/vocafex/pkg/audio/sink"
	"github.com/raspiantoro/vocafex/pkg/audio/source"
	"github.com/spf13/cobra"
	"github.com/zimmski/osutil"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vocafex",
	Short: "audio processing software",
	Long: `vocafex is an audio processing software.
with vocafex it will allow you to add soundfx to any sound input in realtime.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		_, err := osutil.CaptureWithCGo(initPortaudio)
		if err != nil {
			return
		}

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)

		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			<-sigint

			portaudio.Terminate()
			fmt.Println("Bye, vocafex is exit now!")

			cancel()
		}()

		sampleRate := 44100
		buffSize := 1024
		order := binary.BigEndian

		buffIn := make([]float32, buffSize)
		micConfig := source.MicConfig{
			NumChannel: 1,
			SampleRate: float64(sampleRate),
			Order:      order,
			Buffer:     buffIn,
		}

		audioSource, err := source.NewAudioSource(source.AudioInputTypeMic, source.WithMicConfig(micConfig))
		if err != nil {
			log.Fatal(err)
		}

		buffOut := make([]float32, buffSize)
		speakerConfig := sink.SpeakerConfig{
			NumChannel: 1,
			SampleRate: float64(sampleRate),
			Order:      order,
			Buffer:     buffOut,
		}

		audioSink, err := sink.NewAudioSink(sink.AudioOutputTypeSpeaker, sink.WithSpeakerConfig(speakerConfig))

		audioPipe := pipeline.Pipeline{
			Source: audioSource,
			Sink:   audioSink,
		}

		err = audioPipe.Start(ctx)
		if err != nil {
			log.Fatal(err)
		}

		<-ctx.Done()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initPortaudio() {
	portaudio.Initialize()
}
