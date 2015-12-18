/**
 * Author: Zhangjd
 * Date: December 17th, 2015
 * Description: Sound synthesis and wave file generation in Golang
 * Reference: https://github.com/sk89q/WavForge
 */

package wavemix

import (
    "errors"
    "fmt"
    "math"
    "strings"
    "unsafe"
)

type WavForge struct {
    channels       int      // Store the number of channels to be generated.
    sampleRate     float64  // The sample rate at which the sample_count will be generated at.
    bitsPerSample  float64  // Maximum number of bits per sample.
    sampleCount    int      // Store the number of samples that have been generated.
    output         []byte   // Contains the samples.
}

type WavHeader struct {
    flag_RIFF[4]      byte     // [0,4]  ChunkID "RIFF"
    chunkSize         uint32   // [4,4]  ChunkSize
    flag_WAVE[4]      byte     // [8,4]  Format "WAVE"
    flag_fmt[4]       byte     // [12,4] Subchunk1ID "fmt "
    subchunk_1_size   uint32   // [16,4] Subchunk1Size: 16 for PCM
    wFormatTag        uint16   // [20,2] AudioFormat: 1 for PCM
    wChannels         uint16   // [22,2] NumChannels: 1 for mono, 2 for stereo
    dwSamplesPerSec   uint32   // [24,4] SampleRate（每秒样本数）
    dwAvgBytesPerSec  uint32   // [28,4] 每秒播放字节数, 其值为通道数×每秒数据位数×每样本的数据位数／8
    wBlockAlign       uint16   // [32,2] 数据块的调整数, 其值为通道数×每样本的数据位值／8
    uiBitsPerSample   uint16   // [34,2] BitsPerSample
    flag_data[4]      byte     // [36,4] Subchunk1ID＂data＂
    subchunk_2_size   uint32   // [40,4] Subchunk2Size
}

func (this *WavForge) initConfig() () {
    this.channels = 2
    this.sampleRate = 44100
    this.bitsPerSample = 16
    this.sampleCount = 0
    return
}

func (this *WavForge) setChannels (channels int) () {
    this.channels = channels
    return
}

func (this *WavForge) getChannels () (int) {
    return this.channels
}

func (this *WavForge) setSampleRate (sampleRate float64) () {
    this.sampleRate = sampleRate
    return
}

func (this *WavForge) getSampleRate () (float64) {
    return this.sampleRate
}

func (this *WavForge) setBitsPerSample (bitsPerSample float64) () {
    this.bitsPerSample = bitsPerSample
    return
}

func (this *WavForge) getBitsPerSample () (float64) {
    return this.bitsPerSample
}

func (this *WavForge) getSampleCount () (int) {
    return this.sampleCount
}

func (this *WavForge) getWavData () ([]byte) {
    return append(this.getWavHeader(), ([]byte(this.output))...)
}

// Generate the WAV header.
func (this *WavForge) getWavHeader () (header []byte) {
    subchunk_2_size := uint32(this.getSampleCount() * this.channels * (int(this.bitsPerSample)) / 8)

    var wavHeader WavHeader
    copy(wavHeader.flag_RIFF[:], "RIFF")
    wavHeader.chunkSize        = subchunk_2_size + 36
    copy(wavHeader.flag_WAVE[:], "WAVE")
    copy(wavHeader.flag_fmt[:], "fmt ")
    wavHeader.subchunk_1_size  = 16
    wavHeader.wFormatTag       = 1
    wavHeader.wChannels        = uint16(this.channels)
    wavHeader.dwSamplesPerSec  = uint32(this.sampleRate)
    wavHeader.dwAvgBytesPerSec = uint32(this.sampleRate * (float64(this.channels)) * this.bitsPerSample / 8)
    wavHeader.wBlockAlign      = uint16((float64(this.channels)) * this.bitsPerSample / 8)
    wavHeader.uiBitsPerSample  = uint16(this.bitsPerSample)
    copy(wavHeader.flag_data[:], "data")
    wavHeader.subchunk_2_size  = uint32(subchunk_2_size)

    // Reference: http://www.golangtc.com/t/54210b56320b52379100000d
    header = (*[44]byte)(unsafe.Pointer(&wavHeader))[:]
    return 
}

// Encodes a sample.
func (this *WavForge) encodeSample (number float64) (byteArr []byte, err error) {
    max := math.Pow(2, this.bitsPerSample)
    if number < 0 {
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
    if number > 0 {
        for {
            mod := uint8((int(math.Floor(number))) % 256)
            byteArr = append(byteArr, mod)
            number = math.Floor(number / 256)
            if number == 0 {
                break
            }
        }
    }
    for i := 0; i < -(-(int(this.bitsPerSample)) >> 3) - len(byteArr); i++ {
        byteArr = append(byteArr, uint8(0))
    }
    return
}

// Generate a sine waveform.
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
        byteArr, err := this.encodeSample(volume * b * wingRatio * math.Sin(2 * math.Pi * i * frequency / this.sampleRate))
        if err != nil {
            // TODO
        }
        for j := 0; j < this.channels - 1; j++ {
            byteArr = append(byteArr, byteArr...)
        }
        this.output = append(this.output, byteArr...)
        this.sampleCount++
    }
    return
}

// Encode microwave string into sine waveform.
func (this *WavForge) encodeWave (uvpmrscode string) (err error) {
    if len(uvpmrscode) != 20 || string(uvpmrscode[:2]) != "uv" {
        err = errors.New(fmt.Sprintf("Invalid code format!"))
        return
    }
    var baseFrequency float64 = 17800
    characters    := "0123456789abcdefghijklmnopqrstuv"
    period        := 0.0872
    volume        := 0.8
    var frequency [32]float64
    for i := 0; i < len(frequency); i ++ {
        frequency[i] = float64(18000 + i * 64)
    }
    for i := 0; i < len(uvpmrscode); i++ {
        char := uvpmrscode[i]
        pos := strings.Index(characters, (string(char)))
        this.synthesizeSine(baseFrequency, volume, period * 0.7)
        this.synthesizeSine(frequency[pos], volume, period * 0.3)
    }
    return
}








