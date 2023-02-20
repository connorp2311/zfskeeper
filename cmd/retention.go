/*
Copyright © 2023 Connor Parsons

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
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/connorp2311/zfsTools/utils"
	"github.com/spf13/cobra"
)

type Snapshot struct {
	name     string
	creation time.Time
}

var (
	intraDailyRetention int
	dailyRetention      int
	weeklyRetention     int
	monthlyRetention    int
	dataset             string
	dryRun              bool
)

// sortAndRemoveDuplicates sorts an array of snapshots in descending order by their epoch time,
// and removes any duplicates based on their snapshot name
func sortAndRemoveDuplicates(snapshots []Snapshot) []Snapshot {
	// Sort the snapshots in descending order by epoch time
	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].creation.After(snapshots[j].creation)
	})

	// Remove duplicates based on snapshot name
	seen := make(map[string]bool)
	result := []Snapshot{}
	for _, snapshot := range snapshots {
		if _, ok := seen[snapshot.name]; !ok {
			seen[snapshot.name] = true
			result = append(result, snapshot)
		}
	}

	return result
}

// get a list of snapshots for the dataset with zfs list
// return an array of snapshots using the Snapshot struct
// if zfs list returns an error, return an empty array
func getSnapshots(dataset string) []Snapshot {
	out, err := exec.Command("zfs", "list", "-Hp", "-t", "snapshot", "-o", "name,creation", dataset).Output()
	if err != nil {
		return []Snapshot{}
	}

	snapshotLines := strings.Split(string(out), "\n")
	snapshots := []Snapshot{}
	for _, snapshotLine := range snapshotLines {
		if snapshotLine == "" {
			continue
		}
		fields := strings.Fields(snapshotLine)
		name := fields[0]
		epochStr := fields[1]
		epoch, err := strconv.ParseInt(epochStr, 10, 64)
		if err != nil {
			// Failed to parse epoch from snapshot name, skip it
			continue
		}
		snapshots = append(snapshots, Snapshot{name: name, creation: time.Unix(epoch, 0)})
	}

	return snapshots
}

//get the latest snapshot from the array of snapshots
func getLatestSnapshot(snapshots []Snapshot) Snapshot {
	latestSnapshot := Snapshot{}
	for _, snapshot := range snapshots {
		if snapshot.creation.After(latestSnapshot.creation) {
			latestSnapshot = snapshot
		}
	}
	return latestSnapshot
}

// getSnapshotsWithinWindow returns a slice of snapshots that were created within the specified time window.
func getSnapshotsWithinWindow(snapshots []Snapshot, startTime time.Time, endTime time.Time) []Snapshot {
	var keepSnapshots []Snapshot
	for _, snapshot := range snapshots {
		if snapshot.creation.Before(startTime) && snapshot.creation.After(endTime) || snapshot.creation.Equal(startTime) {
			keepSnapshots = append(keepSnapshots, snapshot)
		}
	}
	return keepSnapshots
}

// getIntraDailySnapshots gets the snapshots to keep for intra-daily retention
func getIntraDailySnapshots(snapshots []Snapshot, intraDailyRetention int) []Snapshot {
	if intraDailyRetention == 0 {
		return []Snapshot{}
	}
	// Calculate the start time for the retention window
	latestSnapshot := getLatestSnapshot(snapshots)
	//startDate is the time of the latest snapshot
	startDate := latestSnapshot.creation
	//endDate is the time of the latest snapshot minus 1 day * intraDailyRetention
	endDate := latestSnapshot.creation.AddDate(0, 0, -1*intraDailyRetention)
	// Get the snapshots within the retention window
	keepSnapshots := getSnapshotsWithinWindow(snapshots, startDate, endDate)
	return keepSnapshots
}

// findRetention returns the latest snapshot for each retention period, for a given number of retention periods.
// The function uses the latest snapshot in the snapshots slice as the starting point for calculating retention periods.
func findRetention(snapshots []Snapshot, retentionPeriodCount int, retentionPeriodDurationDays int) []Snapshot {
	if retentionPeriodCount == 0 {
		return []Snapshot{}
	}
	// Calculate the start time for the retention window
	latestSnapshot := getLatestSnapshot(snapshots)
	// Keep the latest snapshot for each day within the retention window
	var keepSnapshots []Snapshot

	// from latestSnapshot.creation cycle through each 30 day period
	// if a snapshots exists within that period, add the latest snapshot to keepSnapshots
	for i := 0; i < retentionPeriodCount; i++ {
		startTime := latestSnapshot.creation.Add(-time.Duration(i) * 24 * time.Duration(retentionPeriodDurationDays) * time.Hour)
		endTime := latestSnapshot.creation.Add(-time.Duration(i+1) * 24 * time.Duration(retentionPeriodDurationDays) * time.Hour)
		snapshotsWithinWindow := getSnapshotsWithinWindow(snapshots, startTime, endTime)
		keepSnapshots = append(keepSnapshots, getLatestSnapshot(snapshotsWithinWindow))
	}
	return keepSnapshots
}

//delete snapshots with zfs destroy, if dryRun is true use the -n flag to simulate the command
func pruneSnapshot(snapshot Snapshot, dryRun bool) (string, error) {
	if dryRun {
		out, err := exec.Command("zfs", "destroy", "-v", "-n", snapshot.name).Output()
		if err != nil {
			return string(out), err
		} else {
			return string(out), nil
		}
	} else {
		out, err := exec.Command("zfs", "destroy", "-v", snapshot.name).Output()
		if err != nil {
			return string(out), err
		} else {
			return string(out), nil
		}
	}
}

// retentionCmd represents the retention command
var retentionCmd = &cobra.Command{
	Use:   "retention <dataset>",
	Short: "This tool is designed to automate the process of deleting ZFS snapshots based on retention policies.",
	Long: `This command performs a retention policy on a ZFS datasets snapshots. 
It provides several flags that allow users to specify the number of days to keep the intra-daily, daily, weekly, and monthly snapshots.`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		dataset := args[0]

		logger, err := utils.NewLogger(logFile, "RET")
		if err != nil {
			fmt.Printf("Error creating logger: %s\n", err)
			os.Exit(1)
		}
		defer logger.Close()

		//check that the dataset exists
		_, err = exec.Command("zfs", "list", "-t", "filesystem", dataset).Output()
		if err != nil {
			logger.Log(fmt.Sprintf("Dataset %s does not exist", dataset))
			os.Exit(1)
		}

		//log if dry run
		if dryRun {
			logger.Log("Performing dry run - no snapshots will be deleted")
		} else {
			if os.Geteuid() != 0 {
				logger.Log("This command must be run as root, enabling dry run to simulate the command")
				dryRun = true
			}
		}

		//get all snapshots
		snapshots := getSnapshots(dataset)

		//log and exit if no snapshots
		if len(snapshots) == 0 {
			logger.Log(fmt.Sprintf("No snapshots found for dataset %s", dataset))
			os.Exit(0)
		}

		keepSnapshots := []Snapshot{}
		keepSnapshots = append(keepSnapshots, getIntraDailySnapshots(snapshots, intraDailyRetention)...)
		keepSnapshots = append(keepSnapshots, findRetention(snapshots, dailyRetention, 1)...)
		keepSnapshots = append(keepSnapshots, findRetention(snapshots, weeklyRetention, 7)...)
		keepSnapshots = append(keepSnapshots, findRetention(snapshots, monthlyRetention, 30)...)
		keepSnapshots = sortAndRemoveDuplicates(keepSnapshots)

		pruneSnapshots := []Snapshot{}

		// create a map to store the names of the snapshots to keep
		keepSnapshotNames := make(map[string]bool)
		for _, snapshot := range keepSnapshots {
			keepSnapshotNames[snapshot.name] = true
		}

		// add all snapshots to the pruneSnapshots array that are not in the keepSnapshots array
		for _, snapshot := range snapshots {
			if _, found := keepSnapshotNames[snapshot.name]; !found {
				pruneSnapshots = append(pruneSnapshots, snapshot)
			}
		}

		//log and exit if no snapshots to prune
		if len(pruneSnapshots) == 0 {
			logger.Log(fmt.Sprintf("No snapshots to prune for dataset %s", dataset))
			os.Exit(0)
		}

		logger.Log(fmt.Sprintf("Pruning %d snapshots for dataset %s", len(pruneSnapshots), dataset))

		for _, snapshot := range pruneSnapshots {
			time.Sleep(1 * time.Second)
			out, err := pruneSnapshot(snapshot, dryRun)
			if err != nil {
				logger.Log(fmt.Sprintf("Error deleting snapshot %s: %s", snapshot.name, err))
			} else {
				logger.Log(out)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(retentionCmd)

	retentionCmd.Flags().IntVarP(&intraDailyRetention, "intra-daily", "i", 0, "Intra-daily retention in days")
	retentionCmd.Flags().IntVarP(&dailyRetention, "daily", "d", 0, "Daily retention in days")
	retentionCmd.Flags().IntVarP(&weeklyRetention, "weekly", "w", 0, "Weekly retention in weeks")
	retentionCmd.Flags().IntVarP(&monthlyRetention, "monthly", "m", 0, "Monthly retention in months")
	retentionCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Perform a dry run and do not delete any snapshots")

	retentionCmd.MarkFlagRequired("intra-daily")
	retentionCmd.MarkFlagRequired("daily")
	retentionCmd.MarkFlagRequired("weekly")
	retentionCmd.MarkFlagRequired("monthly")
	retentionCmd.MarkFlagRequired("dataset")
}
