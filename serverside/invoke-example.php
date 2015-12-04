<?php
/**
 * Author: Zhangjd
 * Date: December 3rd, 2015
 * Reference: http://developer.qiniu.com/docs/v6/api/reference/fop/pfop/pfop.html
 * Description: 模拟调用七牛触发持久化处理（pfop）接口
 */
require 'vendor/autoload.php';
use Qiniu\Auth;

$accessKey = 'Access_Key';
$secretKey = 'Secret_Key';
$bucket    = 'Bucket_Name';

$pfopHost  = 'http://api.qiniu.com';
$pfopPath  = '/pfop';
$pfopURL   = $pfopHost . $pfopPath;

$notifyURL = 'Notify_URL';
$filename  = 'Target_File_Key';
$force     = '1';

// persistentOps Reference: 
//   http://developer.qiniu.com/docs/v6/api/reference/security/put-policy.html#put-policy-persistent-ops-explanation
// e.g:
//   'command1/arg1/arg2|saveas;command2/arg1;command3/arg1';
$fops = 'command1/arg1/arg2;command2/arg1;command3/arg1';

$auth = new Auth($accessKey, $secretKey);

// 构造pfop请求参数
$post = http_build_query(array(
    'bucket'    => $bucket,
    'key'       => $filename,
    'fops'      => $fops,
    'notifyURL' => $notifyURL,
    'force'     => $force
));

// 构造七牛管理凭证（Access Token）
$signingStr  = $pfopPath . "\n" . $post;
$accessToken = $auth->sign($signingStr);

// 构造pfop请求头
$header = array();
$header[] = 'Host: api.qiniu.com';
$header[] = 'Content-Type: application/x-www-form-urlencoded';
$header[] = 'Authorization: QBox ' . $accessToken;

// 发起curl请求
$curl = curl_init($pfopURL);
curl_setopt($curl, CURLOPT_HEADER, 0 );
curl_setopt($curl, CURLOPT_HTTPHEADER, $header);
curl_setopt($curl, CURLOPT_POST, 1);
curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
curl_setopt($curl, CURLOPT_POSTFIELDS, $post);
$result = curl_exec($curl);
curl_close($curl);

// 打印结果
var_dump($result);
// 如果请求成功，返回包含如下内容的JSON字符串
// {"persistentId":"<persistentId string>"}
// 如果请求失败，返回包含如下内容的JSON字符串
// {"error":  "<ErrMsg string>"}


