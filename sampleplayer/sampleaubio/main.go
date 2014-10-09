package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Sample struct {
	Offset       float64 // offset in milliseconds until max amplitude must be positiv
	MaxAmp       float64 // max amplitude value, must be between 0 and 1
	Channels     uint    // number of channels
	NumFrames    int     // number of frames
	SampleRate   int     // e.g. 44100
	SampleFormat string  // e.g. int16
	Duration     float64 // duration in seconds
	HeaderFormat string  // e.g. WAV
	Frequency    float64
}

/*
{
  "offset": 0.56578231292517,
  "maxAmp": 0.07525634765625,
  "numFrames": 42095,
  "sampleRate": 44100,
  "channels": 1,
  "sampleFormat": "int16",
  "duration": 0.95453514739229,
  "headerFormat": "WAV"
}
*/

// aubiopitch -i flute_A4_1_forte_normal.wav -p yinfft, delete all freqs with 0
// aubioonset -i flute_A4_1_forte_normal.wav -O mkl

// returns map position to frequency
func pitches(dir, file string) map[float64]float64 {
	res := map[float64]float64{}
	cmd := exec.Command("aubiopitch", "p", "yinfft", "-i", filepath.Join(dir, file))
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	for _, pair := range strings.Split(string(out), "\n") {
		arr := strings.SplitN(pair, " ", 2)
		pos_s, freq_s := arr[0], arr[1]
		pos, e1 := strconv.ParseFloat(pos_s, 64)
		if e1 != nil {
			fmt.Printf("Error parsing position as float: %s\n", e1)
			continue
		}
		freq, e2 := strconv.ParseFloat(freq_s, 64)
		if e2 != nil {
			fmt.Printf("Error parsing freq as float: %s\n", e2)
			continue
		}
		if freq > 0 {
			res[pos] = freq
		}
	}
}

func onsets(dir, file string) []float64 {
	res := []float64{}
	cmd := exec.Command("aubioonset", "O", "mkl", "-i", filepath.Join(dir, file))
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	for _, pos_s := range strings.Split(string(out), "\n") {
		pos, e1 := strconv.ParseFloat(pos_s, 64)
		if e1 != nil {
			fmt.Printf("Error parsing position as float: %s\n", e1)
			continue
		}
		res = append(res, pos)
	}
	return res
}

func findOnsetAndFrequency() {

}

func main() {

}

/*
	0.01 (10ms)
	0.2  (150ms)
*/

/*
0.133515 443.445251
0.139320 441.971710
0.145125 441.676453
0.150930 443.156525
0.156735 444.198395
0.162540 444.330475
0.168345 444.038025
0.174150 443.693298
0.179955 443.419983
0.185760 443.292084
0.191565 443.285736
0.197370 443.329681
0.203175 443.401215
0.208980 443.485840
*/
