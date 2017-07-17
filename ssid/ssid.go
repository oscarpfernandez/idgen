package ssid

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/oscarpfernandez/idgen/common"
)

const (
	// IDMaxBitLen maximum number of bits used by SSID
	IDMaxBitLen = 63

	// TimeLenBit the number of bits used for the elapsed time
	TimeLenBit = 40

	// GeneratorLenBit the number of bits used for the generator's ID
	GeneratorLenBit = 8

	// SequenceLenBit the number of bit used for the sequence ID
	SequenceLenBit = IDMaxBitLen - TimeLenBit - GeneratorLenBit

	// SequenceIdMask the mask for the sequence ID
	// For instance if SequenceLenBit is 4 on a 8 bit number we would have
	// ( 1<<4 -1 ) = ( 00010000 -1 ) = ( 00001111 )
	SequenceIdMask = uint16(1<<SequenceLenBit - 1)
)

var (
	ErrInvalidStartTime   = errors.New("start time is in the future")
	ErrInvalidGeneratorID = errors.New("Invalid generator ID")
)

// Config defines some of the basic properties that must be defined to use the generator
type Config struct {
	GeneratorID uint16    // an Id that identifies this execution
	StartTime   time.Time // the baseline starttime of this execution
}

// SSID defines the essential elements to use the ssid generator
type SSID struct {
	initialTime       int64       // the instant we start generating a batch of ids
	generatorID       uint16      // an ID that identifies this instance
	currentSequenceID uint16      // the current sequence number [0, 2^SequenceLenBit - 1]
	chronoTime        int64       // the elapsed time in milis since 'initialTime'
	mutex             *sync.Mutex // to guard concurrent access to the struct
}

// NewSSID creates a new SSIS struct provided a Config struct parameters
func NewSSID(config *Config) (*SSID, error) {
	if config.StartTime.After(time.Now()) {
		return nil, ErrInvalidStartTime
	}

	if config.GeneratorID > uint16(1<<GeneratorLenBit-1) {
		return nil, ErrInvalidGeneratorID
	}

	if config.StartTime.IsZero() {
		config.StartTime = time.Date(2016, 7, 4, 0, 0, 0, 0, time.UTC)
	}

	ssid := &SSID{
		initialTime:       common.TimeInMilis(config.StartTime),
		generatorID:       config.GeneratorID,
		currentSequenceID: 0,
		chronoTime:        0,
		mutex:             new(sync.Mutex),
	}

	return ssid, nil
}

// GenerateIDs creates n SSIDs in batch, using the same timestamp and sequence counter [0, uint16(SequenceIdMask)]
// returnin an array of the generated IDs. Error is returned if the number of IDs is bigger than the sequence counter
// or if we the maximum time allowed and encodable by TimeLenBit bits.
func (ssid *SSID) GenerateIDs(n uint16) ([]uint64, error) {
	if n > SequenceIdMask {
		return nil, fmt.Errorf("Requested Ids is greater than the max allowed: %d", SequenceIdMask)
	}

	generatedIDs := make([]uint64, 0, n)

	for i := uint16(0); i < n; i++ {
		id, err := ssid.GenerateID()
		if err != nil {
			return generatedIDs, err
		}
		generatedIDs = append(generatedIDs, id)
	}

	return generatedIDs, nil
}

// GenerateID generates a single SSID
func (ssid *SSID) GenerateID() (uint64, error) {
	ssid.mutex.Lock()
	defer ssid.mutex.Unlock()
	newElapsedTime := common.ElapsedTimeInMilis(ssid.initialTime)
	//fmt.Printf("newElapsedTime: %d, chronoTime: %d \n", newElapsedTime, ssid.chronoTime)
	if newElapsedTime > ssid.chronoTime {
		ssid.chronoTime = newElapsedTime
		ssid.currentSequenceID = (ssid.currentSequenceID + 1) & SequenceIdMask
	} else {
		ssid.currentSequenceID = (ssid.currentSequenceID + 1) & SequenceIdMask
	}

	if ssid.chronoTime > (1<<TimeLenBit - 1) {
		return 0, errors.New("we have reached the maximum time allowed")
	}

	return ssid.getID(), nil
}

// getID composes the final SSID provided their components
func (ssid *SSID) getID() uint64 {
	maskedTime := uint64(ssid.chronoTime << (GeneratorLenBit + SequenceLenBit))
	maskedGeneratorID := uint64(ssid.generatorID << SequenceLenBit)
	maskedSequenceID := uint64(ssid.currentSequenceID)

	return (maskedTime | maskedGeneratorID | maskedSequenceID)
}
