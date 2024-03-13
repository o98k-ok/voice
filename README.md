# VOICE
<img src="./assets/voice.png" width="100"><br>


[![build](https://github.com/o98k-ok/voice/actions/workflows/go.yml/badge.svg)](https://github.com/o98k-ok/voice/actions/workflows/go.yml)
[![Release](https://img.shields.io/badge/release-0.0.1-green.svg)](https://github.com/o98k-ok/voice/releases)

## Basic introduction
一款运行在命令行的音乐播放器.

### How to run

```shell
# go run cmd/terminal/main.go -h

go run cmd/terminal/main.go --home ./data  # ./data need exist in advance
```

home目录作用如下:
1. 存储所有的音频元数据信息
2. 存储所有的音频信息
3. 作为音频格式转化的临时目录

## Running screenshot

### 当前播放

![](./assets/Pasted%20image%2020240311154007.png)

支持的功能：
1. 歌曲基本信息
2. 音乐播放进度
3. 播放切换/暂停播放
4. 歌曲快进/倒退

### 歌曲搜索

![](./assets/Pasted%20image%2020240311154233.png)

支持的功能：
1. B站音频搜索
2. 搜索列表展示、切换
3. 歌曲快速播放

### 播放列表

![](./assets/Pasted%20image%2020240311154510.png)

1. 播放列表信息展示
2. 播放列表切换
3. 播放音乐切换
4. 歌曲删除

## Supported features

1. 终端界面
2. 命令行快捷操作
3. 音乐播放器
4. B站音频实时搜索
5. 本地音频载入

## Follow-up plan

1. 收藏功能
2. 音频下载
3. 每日推荐
4. .......

## Acknowledgments

* [bubbletea](https://github.com/charmbracelet/bubbletea): The fun, functional and stateful way to build terminal apps.
* [beep](https://github.com/faiface/beep): A little package that brings sound to any Go application.
* [lancet](https://github.com/duke-git/lancet): Lancet is a comprehensive, efficient, and reusable util function library of go.


## Star && Follow
