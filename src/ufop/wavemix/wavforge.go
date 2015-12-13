/**
* Author: Zhangjd
* Date: December 12th, 2015
* Description: Sound synthesis and wave file generation in Golang
* Reference: https://github.com/sk89q/WavForge
*/

package wavemix

import (
    "errors"
    "fmt"
    "math"
    "strings"
)

type WavForge struct {
    channels      int     // Store the number of channels to be generated.
    sampleRate    float64 // The sample rate at which the sample_count will be generated at.
    bitsPerSample float64 // Maximum number of bits per sample.
    sampleCount   int     // Store the number of samples that have been generated.
    output        string  // Contains the samples.
}

func (this *WavForge) InitConfig() () {
    this.channels = 2
    this.sampleRate = 44100
    this.bitsPerSample = 16
    this.sampleCount = 0
    this.output = ""
}

func (this *WavForge) SetChannels (channels int) () {
    this.channels = channels
}

func (this *WavForge) getChannels () (int) {
    return this.channels
}

func (this *WavForge) SetSampleRate (sampleRate float64) () {
    this.sampleRate = sampleRate
}

func (this *WavForge) getSampleRate () (float64) {
    return this.sampleRate
}

func (this *WavForge) SetBitsPerSample (bitsPerSample float64) () {
    this.bitsPerSample = bitsPerSample
}

func (this *WavForge) getBitsPerSample () (float64) {
    return this.bitsPerSample
}

func (this *WavForge) getSampleCount () (int) {
    return this.sampleCount
}


func getWavData () () {

}

func getWavHeader () () {

}

// Encodes a sample
func (this *WavForge) EncodeSample (number float64) (encodedStr string, err error) {
    max := math.Pow(2, this.bitsPerSample)
    if number > 0 {
        number += max
    }
    if number >= max {
        if number == max {
            number = 0
        } else {
            err = errors.New(fmt.Sprintf("Overflow (%f won't fit into an %f-bit integer)", number, this.bitsPerSample))
            return
        }
    }
    charSlice := make([]string, 0)
    for {
        charSlice = append(charSlice, (string(rune((int(math.Floor(number))) % 256))))
        number = math.Floor(number / 256)
        if number > 0 {
            break
        }
    }
    for i := 0; i < - (-(int(this.bitsPerSample)) >> 3) - len(charSlice); i++ {
        charSlice = append(charSlice, (string(rune(0))))
    }
    encodedStr = strings.Join(charSlice, "")
    return
}

// 合成指定频率的正弦波
func (this *WavForge) synthesizeSine (frequency float64, volume float64, seconds float64) () {
    total := math.Floor(this.sampleRate * seconds)

    // add wing for decrease noise, increase/decrease voice smoothly
    raiseWing := total * 0.250
    dropWing  := total * 0.250
    b := math.Pow(2, this.bitsPerSample) / 2

    for i := 0.0; i < total; i ++ {
        var wingRatio float64
        if i < raiseWing {
            wingRatio = i / raiseWing
        } else if dropWing >= (total - i) {
            wingRatio = (total - i) / dropWing
        } else {
            wingRatio = 1.0
        }
        // Add a sample for each channel
        encodedStr, err := this.EncodeSample(volume * b * wingRatio * math.Sin(2 * math.Pi * i * frequency / this.sampleRate))
        if err != nil {
            // TODO
        }
        this.output += strings.Repeat(encodedStr, this.channels)
        this.sampleCount++
    }

}


func (this *WavForge) createWave (result string, err error) {
    baseFrequency := 18000
    characters    := "0123456789abcdefghijklmnopqrstuv"
    period        := 0.0872
    var frequency [32]float64
    for i := 0; i < len(frequency); i ++ {
        frequency[i] = float64(baseFrequency + i * 64)
    }

    testCode := "uv0123456789abcdefgh"
    for i := 0; i < len(testCode); i++ {
        char := testCode[i]
        pos := strings.Index(characters, (string(char)))
        this.synthesizeSine(17800, 0.6, period / 2.0 * 1.4)
        this.synthesizeSine(frequency[pos], 0.6, period / 2.0 * 0.6)
    }
    result = this.output
    return
}








