[**日本語セクション**](#pjsekai-overlay-append--フォークプロセカ風動画作成補助ツール-日本語)
# pjsekai-overlay-APPEND / Forked PJSekai-style video creation tool (English)
Fork of [pjsekai-overlay](https://github.com/sevenc-nanashi/pjsekai-overlay) by [TootieJin](https://tootiejin.com), an open-sourced tool to make Project Sekai Fanmade (custom chart) videos.

- **16:9**

https://github.com/user-attachments/assets/5e0eb5c9-cc6d-4d0b-a8a6-409d141e2e8e

- **4:3**

https://github.com/user-attachments/assets/bf82fc81-ca4e-4f07-a7c7-cfe34080432e

> [!CAUTION]
> **For English users:** This tool is primary only for Japanese users, people with technical know-how and basic knowledge of AviUtl, as this repo **ONLY WORKS IN _AviUtl JP_.**\
> Only use this tool if you can figure it out yourself. **DO NOT open issues, DM me, or request help in Sonolus / Chart Cyanvas Discord servers about this**.

This is a forked version of pjsekai-overlay with additional features originally not in the main repo, including:
  - [Extra assets](./assets/extra%20assets) (thank you [ReiyuN](https://discordid.netlify.app/?id=383636820409188374) for the contribution!)
  - Added/adjusted elements to look identical to the official photography
  - Quickly make 1080p videos
  - iPad (4:3) video support

## Requirements

- [AviUtl](http://spring-fragrance.mints.ne.jp/aviutl/) + [Advanced Editing plug-in](http://spring-fragrance.mints.ne.jp/aviutl/) + [L-SMASH Works](https://github.com/Mr-Ojii/L-SMASH-Works-Auto-Builds/releases/latest)
  - (Recommended: [patch.aul](https://scrapbox.io/ePi5131/patch.aul))
- [Unmult](https://github.com/mes51/AVIUtl_Unmult)
- Basic knowledge of AviUtl

## Video Guide

1. [Make your chart first.](https://cc.sevenc7c.com)
2. Go to [Sonolus](https://sonolus.com/) to find your chart.
3. Screen record the video with **BLACK background** and「Hide UI」turned on
4. Transfer the video file to your computer.
   - Download the [ffmpeg](https://www.ffmpeg.org/) encoder if you haven't.
5. Once done, refer to the usage guide below.

## Usage Guide (pjsekai-overlay-APPEND)

0. Create an AviUtl project with specifications below:
   - **For phone users:** 1920x1080, 60fps
   - **For iPad users:** 1440x1080, 60fps
1. Download the latest version of pjsekai-overlay-APPEND [here](https://github.com/TootieJin/pjsekai-overlay-APPEND/releases/latest/).
2. Unzip the file
3. Import ALL files & folders in the [`depenencies/aviutl animation`](./depenencies/aviutl%20animation) folder into here:
```
   aviutl
      ⌞Plugins
         ⌞script
```
   - *If a folder is missing, make a new folder with said name.*
4. Open AviUtl
   - **Note: You must open AviUtl before opening pjsekai-overlay-APPEND to install objects.**
5. Open `pjsekai-overlay-APPEND.exe`
6. Input the chart ID.
   - Potato Leaves prefix: `ptlv-`, Chart Cyanvas prefix: `chcy-`
7. Import specified exo file by navigating to your `pjsekai-overlay/dist/[chart ID]` directory:
   - **For phone users:** main_16-9_1920x1080.exo
   - **For iPad users:** main_4-3_1440x1080.exo

## Usage Guide (AviUtl JP)

Refer to this [English guide](https://github.com/Khronophobia/pjsekai-overlay-english/wiki/Usage-Guide) on how to use AviUtl.

## Terms of Use

In the description of your video, please include the following:
- The name `Nanashi.`
- A link to [***the original*** repo](https://github.com/sevenc-nanashi/pjsekai-overlay) (not this fork)
- A link to `https://sevenc7c.com`

#### Example
**EN**
```
PJSekai-style video creation tool:
- Forked ver. by TootieJin (https://tootiejin.com)
   (https://github.com/TootieJin/pjsekai-overlay-APPEND)
- Original by 名無し｡ (https://sevenc7c.com) 
   (https://github.com/sevenc-nanashi/pjsekai-overlay)
```
**JP**
```
プロセカ風動画作成補助ツール：
- TootieJin (https://tootiejin.com) フォーク版
   (https://github.com/TootieJin/pjsekai-overlay-APPEND)
- 名無し｡ (https://sevenc7c.com) オリジナル版
   (https://github.com/sevenc-nanashi/pjsekai-overlay)
```

## Video Example

[![](https://res.cloudinary.com/marcomontalbano/image/upload/v1732454510/video_to_markdown/images/youtube--JVKMuscJf8c-c05b58ac6eb4c4700831b2b3070cd403.jpg)](https://youtu.be/JVKMuscJf8c "")

---------------------------------------------------------------------------------------

# pjsekai-overlay-APPEND / フォークプロセカ風動画作成補助ツール (日本語)

[TootieJin](https://tootiejin.com)氏による[pjsekai-overlay](https://github.com/sevenc-nanashi/pjsekai-overlay)用フォーク。
pjsekai-overlay(-APPEND) は、プロセカの創作譜面をプロセカ風の動画にするためのオープンソースのツールです。

これはpjsekai-overlayのフォーク版で、元々メインレポにはない以下のような追加機能があります：
  - [追加アセット](./assets/extra%20assets/) ([ReiyuN](https://discordid.netlify.app/?id=383636820409188374)、ご寄稿ありがとうございました。)
  - 本家撮影と同じように見えるように要素を追加/調整
  - 1080p動画を素早く作成
  - iPad（4:3）動画対応

## 必須事項

- [AviUtl](http://spring-fragrance.mints.ne.jp/aviutl/) + [拡張編集プラグイン](http://spring-fragrance.mints.ne.jp/aviutl/) （[導入方法](https://aviutl.info/dl-innsuto-ru/)）
  - (強く推奨：[patch.aul](https://scrapbox.io/ePi5131/patch.aul))
- [Unmult](https://github.com/mes51/AVIUtl_Unmult)
- AviUtlの基本的な知識

## 動画の作り方

1. [譜面を作る](https://wiki.purplepalette.net/create-charts)
2. [Sonolus](https://sonolus.com/)で譜面を撮影する
3. **背景を黒**にし、「Hide UI」をONにして、動画をスクリーン録画します。
4. 撮影したプレイ動画のファイルをパソコンに転送する
   - Google Drive など
5. [ffmpeg](https://www.ffmpeg.org/)で再エンコードする
   - AviUtl で読み込むため
6. 下の利用方法に従って UI を後付けする

## 利用方法

0. 以下の仕様でAviUtlプロジェクトを作成してください：
   - **スマホユーザー向け:** 1920x1080, 60fps
   - **iPadユーザー向け:** 1440x1080、60fps
1. 右の Releases から最新のバージョンの zip をダウンロードする
2. zip を解凍する
3. [`depenencies/aviutl animation`](./depenencies/aviutl%20animation)フォルダ内のすべてのファイルとフォルダをここにインポートします：
```
   aviutl
      ⌞Plugins
         ⌞script
```
   - *フォルダがない場合は、その名前で新しいフォルダを作ってください。*
4. AviUtl を起動する
   - **pjsekai-overlay が起動する前に AviUtl を起動するとオブジェクトのインストールが行われます。**
5. `pjsekai-overlay-APPEND.exe` を起動する
6. 譜面 ID を入力する
   - Potato Leaves の場合は `ptlv-` を、Chart Cyanvas の場合は `chcy-` を先頭につけたまま入力してください。
7. `pjsekai-overlay/dist/[譜面ID]`ディレクトリに移動して、指定したexoファイルをインポートします：
   - **スマホユーザー向け:** main_16-9_1920x1080.exo
   - **iPadユーザー向け:** main_4-3_1440x1080.exo

## 利用規約

動画の概要欄などに、

- `名無し｡`という名前
- このリポジトリへのリンク
- `https://sevenc7c.com`へのリンク

が含まれている文章を載せて下さい。

#### 例
**EN**
```
PJSekai-style video creation tool:
- Forked ver. by TootieJin (https://tootiejin.com)
   (https://github.com/TootieJin/pjsekai-overlay-APPEND)
- Original by 名無し｡ (https://sevenc7c.com) 
   (https://github.com/sevenc-nanashi/pjsekai-overlay)
```
**JP**
```
プロセカ風動画作成補助ツール：
- TootieJin (https://tootiejin.com) フォーク版
   (https://github.com/TootieJin/pjsekai-overlay-APPEND)
- 名無し｡ (https://sevenc7c.com) オリジナル版
   (https://github.com/sevenc-nanashi/pjsekai-overlay)
```

## 動画の例

[![](https://res.cloudinary.com/marcomontalbano/image/upload/v1732454510/video_to_markdown/images/youtube--JVKMuscJf8c-c05b58ac6eb4c4700831b2b3070cd403.jpg)](https://youtu.be/JVKMuscJf8c "")
