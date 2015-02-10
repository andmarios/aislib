package ais

import (
	"errors"
	"fmt"
)

// BinaryBroadcast is a Type 8 message
type BinaryBroadcast struct {
	Repeat uint8
	MMSI   uint32
	DAC    uint16
	FID    uint8
	Data   string
}

// DecodeBinaryBroadcast decodes [the payload of] an AIS Binary Broadcast message (Type 8) but not its binary payload
func DecodeBinaryBroadcast(payload string) (BinaryBroadcast, error) {
	data := []byte(payload)
	var m BinaryBroadcast

	mType := decodeAisChar(data[0])
	if mType != 8 {
		return m, errors.New("Message isn't Binary Broadcast (type 8).")
	}

	m.Repeat = uint8(bitsToInt(6, 7, data))

	m.MMSI = uint32(bitsToInt(8, 37, data))

	m.DAC = uint16(bitsToInt(40, 49, data))
	m.FID = uint8(bitsToInt(50, 55, data))

	m.Data = payload // Data start at bit 56, but this way we simplify our code

	return m, nil
}

// PrintBinaryBroadcast returns a string with some data for a Binary Broadcast message
func PrintBinaryBroadcast(m BinaryBroadcast) string {

	message :=
		fmt.Sprintf("=== Binary Broadcast ===\n") +
			fmt.Sprintf(" Repeat       : %d\n", m.Repeat) +
			fmt.Sprintf(" MMSI         : %09d [%s]\n", m.MMSI, DecodeMMSI(m.MMSI)) +
			fmt.Sprintf(" DAC-FID      : %d-%d (%s)\n", m.DAC, m.FID, BinaryBroadcastType[int(m.DAC)][int(m.FID)])

	return message
}

// Some Binary Broadcast types. The list isn't complete but I haven't searched for a better source
var BinaryBroadcastType = map[int]map[int]string{
	1: {
		11: "Meteorological/Hydrogological Data",
		13: "Fairway closed",
		15: "Extended ship and voyage",
		17: "VTS-Generated/Synthetic targets",
		19: "Marine traffic signals",
		21: "Weather observation from ship",
		22: "Area notice (broadcast)",
		24: "Extended ship and voyage",
		26: "Environmental",
		27: "Route info broadcast",
		29: "Text description broadcast",
		31: "Meteorological and Hydrological",
	},
	200: {
		10: "Ship static and voyage related data",
		23: "EMMA warning report",
		24: "Water levels",
		40: "Signal status",
	},
	316: {
		1:  "Weather Station or Wind or Water Level",
		2:  "Lockage Order or Estimated Lock Times",
		32: "Seaway Version Message",
	},
	366: {
		1:  "Weather Station or Wind or Water Level or PAWS Hydro / Current or PAWS Hydro / Salinity Temp or PAWS Vessel Procession Order",
		2:  "Lockage Order or Estimated Lock Times",
		32: "Seaway Version Message",
	},
}
