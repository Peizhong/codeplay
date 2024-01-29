<?php

$i++;
$cfg['Servers'][$i]['verbose'] = 'local';
$cfg['Servers'][$i]['host'] = 'mysql';
# $cfg['Servers'][$i]['port'] = 3306;
# $cfg['Servers'][$i]['auth_type'] = 'cookie';
$cfg['Servers'][$i]['user'] = 'app_user';
$cfg['Servers'][$i]['password'] = 'app_user';
$cfg['Servers'][$i]['compress'] = false;

$i++;
$cfg['Servers'][$i]['verbose'] = 'hawk test';
$cfg['Servers'][$i]['host'] = '10.131.248.27';
# $cfg['Servers'][$i]['port'] = 3306;
# $cfg['Servers'][$i]['auth_type'] = 'cookie';
$cfg['Servers'][$i]['user'] = 'root';
$cfg['Servers'][$i]['password'] = '';
$cfg['Servers'][$i]['AllowNoPassword'] = true;

$i++;
$cfg['Servers'][$i]['verbose'] = 'nightinggale test';
$cfg['Servers'][$i]['host'] = '10.131.24.20';
$cfg['Servers'][$i]['port'] = 4330;
# $cfg['Servers'][$i]['auth_type'] = 'cookie';
$cfg['Servers'][$i]['user'] = 'root';
$cfg['Servers'][$i]['password'] = 'mz.com';