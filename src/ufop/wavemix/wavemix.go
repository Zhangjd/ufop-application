/**
 * Author: Zhangjd
 * Date: December 8th, 2015
 * Reference: http://developer.qiniu.com/docs/v6/api/reference/fop/pfop/pfop.html
 * Description: 声波合成模块
 */

package wavemix

import (
    "github.com/qiniu/api.v6/auth/digest"
)

const (
    AUDIO_MERGE_MAX_FIRST_FILE_LENGTH  = 100 * 1024 * 1024
    AUDIO_MERGE_MAX_SECOND_FILE_LENGTH = 100 * 1024 * 1024
)

type AudioMerger struct {
    mac                 *digest.Mac
    maxFirstFileLength  uint64
    maxSecondFileLength uint64
}