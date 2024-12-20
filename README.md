[**日本語セクション**](#pjsekai-overlay-append--フォークプロセカ風動画作成補助ツール-日本語)
https://github.com/user-attachments/assets/148e315c-b737-4ccc-90cc-b5311a00b07e
# pjsekai-overlay-APPEND / Forked PJSekai-style video creation tool (English)

Fork of [pjsekai-overlay](https://github.com/sevenc-nanashi/pjsekai-overlay) by [TootieJin](https://tootiejin.com).

> [!CAUTION]
> **For English users:** This tool is primary only for Japanese users, people with technical know-how and basic knowledge of AviUtl, as this repo **ONLY WORKS IN _AviUtl JP_.**
> Only use this tool if you can figure it out yourself. **DO NOT open issues, DM me, or request help in Sonolus / Chart Cyanvas Discord servers about this**.

This is a forked version of pjsekai-overlay with additional features originally not in the main repo, including:
  - [Extra assets](./assets/extra%20assets) (thank you [ReiyuN](https://discordid.netlify.app/?id=383636820409188374) for the contribution!)
  - Added/adjusted elements to look identical to the official photography
  - Quickly make 1080p videos

## Requirements

- - [AviUtl](http://spring-fragrance.mints.ne.jp/aviutl/) + [Advanced Editing plug-in](http://spring-fragrance.mints.ne.jp/aviutl/) + [L-SMASH Works](https://github.com/Mr-Ojii/L-SMASH-Works-Auto-Builds/releases/latest)
    (Recommended: [patch.aul](https://scrapbox.io/ePi5131/patch.aul))
- [Unmult](https://github.com/mes51/AVIUtl_Unmult)
- Basic knowledge of AviUtl

## Video Guide

1. [Make your chart first.](https://cc.sevenc7c.com)
2. Go to [Sonolus](https://sonolus.com/) to find your chart.
   - Turn on「Hide UI」
3. Transfer the video file to your computer.
   - Download the [ffmpeg](https://www.ffmpeg.org/) encoder if you haven't.
4. Once done, refer to the usage guide below.

## Usage Guide (pjsekai-overlay-APPEND)

0. Create an AviUtl project with 1920x1080, 60fps
1. Download the latest version of pjsekai-overlay-APPEND [here](https://github.com/TootieJin/pjsekai-overlay-APPEND/releases/latest/).
2. Unzip the file
3. Open AviUtl
   - **Note: You must open AviUtl before opening pjsekai-overlay-APPEND to install objects.**
4. Open `pjsekai-overlay.exe`
5. Input the chart ID.
   - Potato Leaves prefix: `ptlv-`, Chart Cyanvas prefix: `chcy-`
6. Import object file by navigating to your `pjsekai-overlay/dist/[chart ID]` directory, and select the exo file depending on which AviUtl you're running at:
   - `main.exo` is for the JP original version
   - `main_en.exo` is for the EN modded version

## Usage Guide (AviUtl EN)

Refer to this [guide](https://github.com/Khronophobia/pjsekai-overlay-english/wiki/Usage-Guide) on how to use AviUtl (English).

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

これはpjsekai-overlayのフォーク版で、元々メインレポにはない以下のような追加機能があります：
  - [追加アセット](./assets/extra%20assets/) ([ReiyuN](https://discordid.netlify.app/?id=383636820409188374)、ご寄稿ありがとうございました。)
  - 本家撮影と同じように見えるように要素を追加/調整
  - 1080p動画を素早く作成

## 必須事項

- [AviUtl](http://spring-fragrance.mints.ne.jp/aviutl/) + [拡張編集プラグイン](http://spring-fragrance.mints.ne.jp/aviutl/) （[導入方法](https://aviutl.info/dl-innsuto-ru/)）
  （強く推奨：[patch.aul](https://scrapbox.io/ePi5131/patch.aul)）
- [AVIUtl_Unmult](https://github.com/mes51/AVIUtl_Unmult)
- AviUtlの基本的な知識

## 動画の作り方

1. [譜面を作る](https://wiki.purplepalette.net/create-charts)
2. [Sonolus](https://sonolus.com/)で譜面を撮影する
   - [Potato Leaves](https://github.com/sevenc-nanashi/potato_leaves)、または [Chart Cyanvas](https://cc.sevenc7c.com)で撮影してください。
   - 「Hide UI」をオンにしてください。
3. 撮影したプレイ動画のファイルをパソコンに転送する
   - Google Drive など
4. [ffmpeg](https://www.ffmpeg.org/)で再エンコードする
   - AviUtl で読み込むため
5. 下の利用方法に従って UI を後付けする

## 利用方法

0. 1280x720, 60fps で aviutl のプロジェクトを作成する
1. 右の Releases から最新のバージョンの zip をダウンロードする
2. zip を解凍する
3. AviUtl を起動する
   - pjsekai-overlay が起動する前に AviUtl を起動するとオブジェクトのインストールが行われます。
4. `pjsekai-overlay.exe` を起動する
5. 譜面 ID を入力する
   - Potato Leaves の場合は `ptlv-` を、Chart Cyanvas の場合は `chcy-` を先頭につけたまま入力してください。
6. `pjsekai-overlay/dist/[譜面ID]`ディレクトリに移動して、オブジェクトファイルをインポートします：
   - `main.exo`はJPバージョン用です。
   - `main_en.exo`はENモディファイ版です。

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
