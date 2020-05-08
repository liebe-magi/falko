# foltia Command Line Tool

A Command Line Tool for [foltia ANIME LOCKER](https://foltia.com/ANILOC/).

## Installation

```bash
% go get github.com/reeve0930/foltia
```

## Initial settings

You should set your environment infomation using `foltia config` command.

```bash
# Set your foltia IP address
% foltia config -i 192.168.xxx.xxx

# Set the path mounted foltia
% foltia config -s /mnt/xxx

# Set the path you want to copy
% foltia config -d /home/user/xxx

# Set the filename format (Details to follow.)
% foltia config -n %title%_%epnum%_%eptitle%

# Set the file type you want to copy ("TS" or "MP4")
% foltia config -t TS

# Set the threshold of dropped TS packets (If the number of dropped TS packets exceeds this value, it will not copy.)
% foltia config -r 10
```
## Filename format

You can use the following parameters when setting the destination file name.

- **%title%** : Animation title (ex: 新世紀エヴァンゲリオン)
- **%epnum%** : Episode number (ex: 01)
- **%eptitle%** : Episode title (ex: 使徒、襲来)

If you set the file name format `%title%_%epnum%_%eptitle%`, your destination file name is **新世紀エヴァンゲリオン_01_使徒、襲来.m2t(mp4)**. 

## How to use

```bash
# First you should update local database
% foltia update

# Start copying the file by executing the following command
% foltia copy
```

## License

Foltia Command Line Tool by reeve0930 is licensed under the Apache License, Version2.0.  
See [LICENSE](https://github.com/reeve0930/foltia/blob/master/LICENSE)
