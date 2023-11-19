<?php

$t = '';

$t = explode("\n", $t);

$subscriber = [];

foreach ($t as $v) {
    $v = trim($v);
    var_dump($v);

    $subscriber[] = [
        'user_id' => $v
    ];
}


$body = [
    'title' => 'test',
    'id' => '4b90af8d-613e-4444-8a1b-b03466425f9b',
    'subscribers' => $subscriber,
];

echo json_encode($body);
