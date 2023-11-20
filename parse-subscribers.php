<?php

$filename = '';


$subscriber = [];


foreach (explode("\n", file_get_contents($filename)) as $line) {
    $item = explode(',', $line);

    if (count($item) < 2) {
        continue;
    }

    $userId = intval(trim($item[1], "\""));
    if ($userId <= 0) {
        continue;
    }

    $subscriber[] = [
        'user_id' => (string) $userId,
    ];
}



$body = [
    'title' => 'product',
    'id' => 'batch',
    'subscribers' => $subscriber,
];

$body = json_encode($body);


$curl = curl_init();

curl_setopt_array($curl, array(
    CURLOPT_URL => 'localhost:80/v1/topic',
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_ENCODING => '',
    CURLOPT_MAXREDIRS => 10,
    CURLOPT_TIMEOUT => 0,
    CURLOPT_FOLLOWLOCATION => true,
    CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
    CURLOPT_CUSTOMREQUEST => 'POST',
    CURLOPT_POSTFIELDS => $body,
    CURLOPT_HTTPHEADER => array(
        'Content-Type: application/json'
    ),
));

$response = curl_exec($curl);

curl_close($curl);
echo $response;
