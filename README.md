# FALKO for foltia ANIME LOCKER

[foltia ANIME LOCKER](https://foltia.com/ANILOC/)用非公式コマンドラインツールです。

## このツールでできること

- 録画管理
    - 録画一覧の取得
    - 録画予約の実施
    - 録画予約の削除
- ファイルコピー
    - 録画したMPEG2TS or MP4ファイルを指定したフォーマット通りにリネームしてコピー
    - 同一タイトルの同一エピソードは一度のみコピー
- Slackによる通知および制御
    - 指定した時刻に当日の予約および新番組情報を通知
    - メッセージにより番組の録画予約

## インストール方法

```bash
% go get github.com/MagicalLiebe/falko
```

## 初期設定

`falko config`コマンドを実行すると、`~/.config/falko`に`config.toml`ができる。  
この設定ファイルを直接編集するか、以下のように`falko config`コマンドで一つずつパラメータを設定していく。  
なお、このコマンドを実行するPCは**foltia ANIME LOCKER**と同一LAN上にある必要がある。  

```bash
# foltia ANIME LOCKERのIPアドレスを設定
% falko config -i 192.168.xxx.xxx

# foltia ANIME LOCKERのpublicフォルダをマウントしているディレクトリを指定
% falko config -s /mnt/xxx

# 録画したファイルのコピー先のディレクトリを指定
% falko config -d /home/user/xxx

# コピーする際のファイル名のフォーマットを指定 (使用できるパラメータは後述)
% falko config -n %title%_%epnum%_%eptitle%

# コピーしたいファイルの形式を指定 ("TS" or "MP4")
% falko config -t TS

# TSパケットのドロップ数の閾値を設定
% falko config -r 10

# Slack botトークンを設定
% falko config -b xxxxxxxxxx

# Slackの定時通知時刻を設定
% falko config -c 08:00
```

## ファイル名フォーマット

ファイル名のフォーマットには以下のパラメータが使用できる。

- **%title%** : アニメのタイトル (ex: 新世紀エヴァンゲリオン)
- **%epnum%** : 話数 (ex: 01)
- **%eptitle%** : サブタイトル (ex: 使徒、襲来)

例えば、`%title%_%epnum%_%eptitle%`と指定した場合、ファイル名は`新世紀エヴァンゲリオン_01_使徒、襲来.m2t(mp4)`のようになる。

## 使い方


### ファイルのコピー

```bash
# 最初にローカルDBの更新を行う
% falko update

# コピー準備のできたファイルの確認
% falko copy -l

# ファイルコピーの実行
% falko copy
```

### Slack botの起動

```bash
# Slack botを起動
% falko slack
```

以下のような文章が表示されるので、通知して欲しいチャンネルに作成したSlack botを参加させた上で、表示されている4桁の数字を入力する。  

```bash
> 2020/05/27 00:59:55 Slackよりこのコードを入力して下さい:3653
```

`認証完了`と表示され、`Slackクライアントスタンバイ完了`となったらOK。  
指定した時刻に通知が実行される。

### Slack botからの録画予約

TIDを指定することで録画予約ができる。  
TIDは新番組の通知に示されている。  
以下の例の場合、**1730**がTID。

```
【新アニメ情報】
とある科学の超電磁砲 (1730)
BS11イレブン : 2020/5/30 (土)
http://cal.syoboi.jp/tid/1730/
```

TIDを指定して、以下のようにSlack botにメッセージを送ると録画予約ができる。  

```
rec 1730
```

以下のような返答があれば、録画予約が完了している。

```
【録画予約成功】
とある科学の超電磁砲 (1730)
```

## License

FALKO for foltia ANIME LOCKER by MagicalLiebe is licensed under the Apache License, Version2.0.  
See [LICENSE](https://github.com/MagicalLiebe/falko/blob/master/LICENSE)
