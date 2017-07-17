package ssid

import (
	"testing"
	"time"

	"github.com/oscarpfernandez/idgen/common"
)

var (
	ssidInstance *SSID
	generatorID  int16
	startTime    int64
	err          error
)

func init() {
	config := Config{
		StartTime:   time.Now(),
		GeneratorID: 0,
	}

	ssidInstance, err = NewSSID(&config)
	if err != nil {
		panic("SSID instance could not be generated")
	}
}

func generateID(t *testing.T) uint64 {
	ssid, err := ssidInstance.GenerateID()
	if err != nil {
		t.Error("Could not generate ID")
	}
	return ssid
}

func TestGenerateOneSSID(t *testing.T) {
	ssid1 := generateID(t)
	ssid2 := generateID(t)

	if ssid2 <= ssid1 {
		t.Errorf("The id1: %d is not less than id2: %d", ssid1, ssid2)
	}
}

func TestGenerateBatch10000SSID(t *testing.T) {
	var numberOfIds uint16 = 10000

	ssids, err := ssidInstance.GenerateIDs(numberOfIds)
	if err != nil {
		t.Error("Could not generate IDs")
	}

	if len(ssids) != int(numberOfIds) {
		t.Errorf("ID array size does not match. Expected: %d, Got: %d", numberOfIds, len(ssids))
	}

	var baseVal uint64
	for _, val := range ssids {
		if val > baseVal {
			baseVal = val
		} else {
			t.Error("IDs not generated sequentially")
		}
	}
}
func TestGenerateBatchSSIDsTryOverSequenceLimit(t *testing.T) {
	numberOfIDs := SequenceIdMask + 1

	_, err := ssidInstance.GenerateIDs(numberOfIDs)
	if err == nil {
		t.Error("Max Sequence ID was allowed!")
	}
}

func TestGenerateSSIDsConcurrently2M(t *testing.T) {
	idConsumer := make(chan uint64)

	numberOfIDs := 10000
	generateID := func() {
		for i := 0; i < numberOfIDs; i++ {
			idConsumer <- generateID(t)
		}
	}

	const numberOfParallelGenerators = 200
	for i := 0; i < numberOfParallelGenerators; i++ {
		go generateID()
	}

	idMap := make(map[uint64]uint64)

	for i := 0; i < numberOfIDs*numberOfParallelGenerators; i++ {
		id := <-idConsumer
		_, exists := idMap[id]
		if exists {
			t.Error("Key already exists in the map - duplicated ID")
		} else {
			idMap[id] = id
		}
	}

	if len(idMap) != numberOfIDs*numberOfParallelGenerators {
		t.Error("Not all of the IDs were generated")
	}
}

func TestTryGenerateIDBeyondMaxTime(t *testing.T) {
	yearDuration := time.Duration(24*365) * time.Hour

	validDurationInYears := time.Duration(34) * yearDuration
	ssidInstance.initialTime -= int64(validDurationInYears) / common.TimeUnitMilis
	_, err := ssidInstance.GenerateID()
	if err != nil {
		t.Error("Valid ID within valid time frame was not created")
	}

	oneYear := time.Duration(1) * yearDuration
	ssidInstance.initialTime -= int64(oneYear) / common.TimeUnitMilis

	_, err = ssidInstance.GenerateID()
	if err == nil {
		t.Error("Out of time bounds ID was generated. ID time should be over")
	}
}
