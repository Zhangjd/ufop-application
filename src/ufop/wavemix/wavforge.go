/**
 * Author: Zhangjd
 * Date: December 12th, 2015
 * Description: Sound synthesis and wave file generation in Golang
 * Reference: https://github.com/sk89q/WavForge
 */

package wavemix

type WavForge struct {
    config WavForgeConfig
}

type WavForgeConfig struct {
    channels      int     // Store the number of channels to be generated.
    sampleRate    int     // The sample rate at which the sample_count will be generated at.
    bitsPerSample int     // Maximum number of bits per sample.
    sampleCount   int     // Store the number of samples that have been generated.
    output        string  // Contains the samples.
}

func (this *WavForge) InitConfig() () {
    this.config.channels = 2
    this.config.sampleRate = 44100
    this.config.bitsPerSample = 16
    this.config.sampleCount = 0
    this.config.output = ""
}

func (this *WavForge) SetChannels (channels int) () {
    this.config.channels = channels
}

func (this *WavForge) getChannels () (int) {
    return this.config.channels
}

func (this *WavForge) SetSampleRate (sampleRate int) () {
    this.config.sampleRate = sampleRate
}

func (this *WavForge) getSampleRate () (int) {
    return this.config.sampleRate
}

func (this *WavForge) SetBitsPerSample (bitsPerSample int) () {
    this.config.bitsPerSample = bitsPerSample
}

func (this *WavForge) getBitsPerSample () (int) {
    return this.config.bitsPerSample
}

func (this *WavForge) getSampleCount () (int) {
    return this.config.sampleCount
}


func getWavData () () {

}

func getWavHeader () () {

}

func synthesizeSine () () {

}





