package cmd

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/mauricioabreu/thumbinator/internal/app/thumbinator"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var thumbsPath string
var streamsFile string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate thumbs from live streamings and videos on demand",
	Run: func(cmd *cobra.Command, args []string) {
		Main()
	},
}

// Run thumbinator, run
func Run() {
	if err := generateCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// Main run thumbinator, run
func Main() {
	streams, err := thumbinator.GetStreams(thumbinator.JSONSource{File: streamsFile})
	if err != nil {
		log.Fatalf("Could not retrieve streams to process: %s", err)
	}
	log.Infof("Retrieved the following streams: %+v", streams)
	for _, s := range streams {
		log.Debugf("Generating thumbs for %s with URL %s", s.Name, s.URL)
		thumbinator.GenerateThumb(s.URL, s.Name, thumbsPath)
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error(err)
	}
	collector := thumbinator.Collector{Watcher: watcher, Store: thumbinator.NewRedisStore(), Path: thumbsPath}
	collector.CollectThumbs(streams)
}

func init() {
	generateCmd.Flags().StringVar(&thumbsPath, "thumbsPath", "thumbnails", "Path where all thumbs will be written to")
	generateCmd.Flags().StringVar(&streamsFile, "streamsFile", "streams.json", "File with streams to extract thumbs")
	rootCmd.AddCommand(generateCmd)
}
