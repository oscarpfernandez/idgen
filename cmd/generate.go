package cmd

import (
	"math/rand"
	"time"

	"github.com/oscarpfernandez/idgen/ssid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	flagCount           = "count"
	flagGenerator       = "generator"
	flagRandomGenerator = "generator-rand"
)

var (
	count             uint16
	generatorID       uint16
	isRandomGenerator bool
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generates ids",
	Long:  `Generate a set of ids`,
	Run: func(cmd *cobra.Command, args []string) {

		if isRandomGenerator {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			generatorID = uint16(r.Int() % (2 ^ ssid.GeneratorLenBit))
			logrus.Infof("Random Generator Id: %d", generatorID)
		}

		config := ssid.Config{GeneratorID: generatorID}
		ssidInstance, err := ssid.NewSSID(&config)
		if err != nil {
			logrus.Error(err)
			return
		}
		logrus.Infof("Configuration: %+v", config)

		startGlobal := time.Now()
		generatedIDs, err := ssidInstance.GenerateIDs(count)
		if err != nil {
			logrus.Error(err)
			return
		}
		elapsedGlobal := time.Since(startGlobal)

		for index := 0; index < len(generatedIDs); index++ {
			logrus.Infof("SSID: %d", generatedIDs[index])
		}
		logrus.Infof("Total Generation Time: %s", elapsedGlobal)
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)
	generateCmd.Flags().Uint16VarP(&count, flagCount, "c", 1, "the number of ids to be generated")
	generateCmd.Flags().Uint16VarP(&generatorID, flagGenerator, "g", 0, "the generator id")
	generateCmd.Flags().BoolVarP(&isRandomGenerator, flagRandomGenerator, "r", false, "random generator id")
}
