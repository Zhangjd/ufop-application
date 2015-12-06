<?php
/**
 * Author: Zhangjd
 * Date: December 3rd, 2015
 * Reference: http://developer.qiniu.com/docs/v6/api/reference/fop/pfop/pfop.html
 */
if (isset($_SERVER['REQUEST_METHOD']) && $_SERVER['REQUEST_METHOD'] != 'POST') {
    http_response_code(405);
    exit;
}
$json = file_get_contents('php://input');
$jsonobj = json_decode($json);
if (!$jsonobj) {
    http_response_code(403);
    exit;
}

$id          = $jsonobj->id;          // 持久化处理会话标识
$pipeline    = $jsonobj->pipeline;    // 云处理操作的处理队列，默认使用队列为共享队列0.default
$code        = $jsonobj->code;        // 状态码，0-成功，1-等待处理，2-正在处理，3-处理失败，4-通知提交失败。
$desc        = $jsonobj->desc;        // 与状态码相对应的详细描述
$inputKey    = $jsonobj->inputKey;    // 处理源文件的文件名
$inputBucket = $jsonobj->inputBucket; // 处理源文件所在的空间名
$items       = $jsonobj->items;       // 云处理操作列表，包含每个云处理操作的状态信息
$reqid       = $jsonobj->reqid;       // 云处理请求的请求id，主要用于七牛技术人员的问题排查

foreach ($items as $key => $value) {
    # code...
    $itemCmd       = $value->cmd;       // 所执行的云处理操作命令（fopN）
    $itemCode      = $value->code;      // 状态码
    $itemDesc      = $value->desc;      // 与状态码相对应的详细描述
    $itemReturnOld = $value->returnOld; // 默认为0。返回1，则是旧数据
    if ($itemCode === 0) {
        $itemHash  = $value->hash;      // 云处理结果保存在服务端的唯一hash标识
        $itemKey   = $value->key;       // 云处理结果的外链资源名
    } else {
        $itemError = $value->error;     // 如果处理失败，该字段会给出失败的详细原因
    }
}

switch ($code) {
    case 0:
        # code...
        break;
    case 1:
        # code...
        break;
    case 2:
        # code...
        break;
    case 3:
        # code...
        break;
    case 4:
        # code...
        break;
    default:
        # code...
        break;
}








