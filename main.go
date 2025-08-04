package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	_ "net/http/pprof"

	pyroscope "github.com/grafana/pyroscope-go"
	"github.com/spf13/cobra"
)

func generateMockLoad(ctx context.Context, profileDuration time.Duration) {
	end := time.Now().Add(profileDuration)
	for time.Now().Before(end) {
		_ = fibonacci(25)
	}
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
	var serviceName string
	var ingestURL string
	var pprofPort string
	var profileDuration time.Duration

	var rootCmd = &cobra.Command{
		Use:   "profilegen",
		Short: "Generate mock profiling data and send to Pyroscope over otel (push mode), also exposes a pprof endpoint that can be scraped by otel collector",
		Run: func(cmd *cobra.Command, args []string) {
			if serviceName == "" || ingestURL == "" {
				fmt.Println("Both --service_name and --ingest_url are required.")
				os.Exit(1)
			}

			// Enable mutex/block profiling for more data
			runtime.SetMutexProfileFraction(5)
			runtime.SetBlockProfileRate(5)

			// Start Pyroscope profiler (push mode)
			_, err := pyroscope.Start(pyroscope.Config{
				ApplicationName: serviceName,
				ServerAddress:   ingestURL,
				Logger:          pyroscope.StandardLogger,
				Tags:            map[string]string{"service_name": serviceName},
				ProfileTypes: []pyroscope.ProfileType{
					pyroscope.ProfileCPU,
					pyroscope.ProfileAllocObjects,
					pyroscope.ProfileAllocSpace,
					pyroscope.ProfileInuseObjects,
					pyroscope.ProfileInuseSpace,
					pyroscope.ProfileGoroutines,
					pyroscope.ProfileMutexCount,
					pyroscope.ProfileMutexDuration,
					pyroscope.ProfileBlockCount,
					pyroscope.ProfileBlockDuration,
				},
			})
			if err != nil {
				fmt.Printf("Failed to start pyroscope profiler: %v\n", err)
				os.Exit(1)
			}

			// Expose pprof endpoints on the port specified in the flag
			go func() {
				fmt.Println("pprof endpoints available at http://localhost:" + pprofPort + "/debug/pprof/")
				http.ListenAndServe(":"+pprofPort, nil)
			}()

			ctx := context.Background()
			generateMockLoad(ctx, profileDuration)
			fmt.Println("Done generating mock load. Exiting.")
		},
	}

	rootCmd.Flags().StringVar(&pprofPort, "pprof_port", "", "Service name tag for the profile")
	rootCmd.Flags().StringVar(&serviceName, "service_name", "", "Service name tag for the profile")
	rootCmd.Flags().StringVar(&ingestURL, "ingest_url", "", "Profile ingestion API endpoint URL (for push)")
	rootCmd.Flags().DurationVar(&profileDuration, "profile_duration", 30*time.Second, "Duration of the profile")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
