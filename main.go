package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/TootieJin/pjsekai-overlay-APPEND/pkg/pjsekaioverlay"
	"github.com/TootieJin/pjsekai-overlay-APPEND/pkg/sonolus"
	"github.com/fatih/color"
	"github.com/google/go-github/v57/github"
	"github.com/srinathh/gokilo/rawmode"
	"golang.org/x/sys/windows"
)

func shouldCheckUpdate() bool {
	executablePath, err := os.Executable()
	if err != nil {
		return false
	}
	updateCheckFile, err := os.OpenFile(filepath.Join(filepath.Dir(executablePath), ".update-check"), os.O_RDONLY, 0666)
	if err != nil {
		return os.IsNotExist(err)
	}
	defer updateCheckFile.Close()

	scanner := bufio.NewScanner(updateCheckFile)
	scanner.Scan()
	lastCheckTime, err := strconv.ParseInt(scanner.Text(), 10, 64)
	if err != nil {
		return false
	}

	return time.Now().Unix()-lastCheckTime > 60*60*24
}

func checkUpdate() {
	githubClient := github.NewClient(nil)
	release, _, err := githubClient.Repositories.GetLatestRelease(context.Background(), "TootieJin", "pjsekai-overlay-APPEND")
	if err != nil {
		return
	}

	executablePath, err := os.Executable()
	if err != nil {
		return
	}
	updateCheckFile, err := os.OpenFile(filepath.Join(filepath.Dir(executablePath), ".update-check"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	defer updateCheckFile.Close()
	updateCheckFile.WriteString(strconv.FormatInt(time.Now().Unix(), 10))

	latestVersion := strings.TrimPrefix(release.GetTagName(), "v")
	if latestVersion == pjsekaioverlay.Version {
		return
	}
	fmt.Printf(color.HiCyanString("新しいバージョンがリリースされています\nNew version released: v%s -> v%s\n"), pjsekaioverlay.Version, latestVersion)
	fmt.Printf(color.HiCyanString("ダウンロード (Download Here) -> %s\n"), release.GetHTMLURL())
}

func origMain(isOptionSpecified bool) {
	Title()

	var skipAviutlInstall bool
	flag.BoolVar(&skipAviutlInstall, "no-aviutl-install", false, "AviUtlオブジェクトのインストールをスキップします。(AviUtl object installation is skipped.)")

	var outDir string
	flag.StringVar(&outDir, "out-dir", "./dist/_chartId_", "出力先ディレクトリを指定します。_chartId_ は譜面IDに置き換えられます。\nEnter the output path. _chartId_ will be replaced with the chart ID.")

	var teamPower int
	flag.IntVar(&teamPower, "team-power", 250000, "総合力を指定します。(Enter the team's power.)")

	var apCombo bool
	flag.BoolVar(&apCombo, "ap-combo", true, "コンボのAP表示を有効にします。(Enable AP display for combo.)")

	flag.Usage = func() {
		fmt.Println("Usage: pjsekai-overlay [譜面ID] [オプション]")
		flag.PrintDefaults()
	}

	flag.Parse()

	if shouldCheckUpdate() {
		checkUpdate()
	}

	if !skipAviutlInstall {
		success := pjsekaioverlay.TryInstallObject()
		if success {
			fmt.Println(color.GreenString("AviUtlオブジェクトのインストールに成功しました。(AviUtl object successfully installed.)"))
		}
	}

	var chartId string
	if flag.Arg(0) != "" {
		chartId = flag.Arg(0)
		fmt.Printf("譜面ID (Chart ID): %s\n", color.GreenString(chartId))
	} else {
		fmt.Print("譜面IDをプレフィックス込みで入力して下さい。\nEnter the chart ID including the prefix.\n\n'chcy-': Chart Cyanvas (cc.sevenc7c.com)\n'ptlv-': Potato Leaves (ptlv.sevenc7c.com)\n> ")
		fmt.Scanln(&chartId)
		fmt.Printf("\033[A\033[2K\r> %s\n", color.GreenString(chartId))
	}

	chartSource, err := pjsekaioverlay.DetectChartSource(chartId)
	if err != nil {
		fmt.Println(color.RedString("譜面のサーバーを判別できませんでした。プレフィックスも込め、正しい譜面IDを入力して下さい。\nThe specified chart doesn't exist. Please enter the correct chart ID including the prefix."))
		return
	}
	fmt.Printf("- 譜面を取得中 (Getting chart): %s%s%s ", RgbColorEscape(chartSource.Color), chartSource.Name, ResetEscape())
	chart, err := pjsekaioverlay.FetchChart(chartSource, chartId)

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}
	if chart.Engine.Version != 12 {
		fmt.Println(color.RedString(fmt.Sprintf("失敗：エンジンのバージョンが古い。\nFAIL: Unsupported engine version. - [ver.%d]", chart.Engine.Version)))
		return
	}

	fmt.Println(color.GreenString("OK"))
	fmt.Printf("  %s / %s - %s (Lv. %s)\n",
		color.CyanString(chart.Title),
		color.CyanString(chart.Artists),
		color.CyanString(chart.Author),
		color.MagentaString(strconv.Itoa(chart.Rating)),
	)

	fmt.Printf("- exeのパスを取得中 (Getting executable path)... ")
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL:%s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL:%s", err.Error())))
		return
	}

	formattedOutDir := filepath.Join(cwd, strings.Replace(outDir, "_chartId_", chartId, -1))
	fmt.Printf("- 出力先ディレクトリ (Output path): %s\n", color.CyanString(filepath.Dir(formattedOutDir)))

	fmt.Print("- ジャケットをダウンロード中 (Downloading jacket)... ")
	err = pjsekaioverlay.DownloadCover(chartSource, chart, formattedOutDir)
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL:%s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	fmt.Print("- 背景をダウンロード中 (Downloading background)... ")
	err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir)
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL:%s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	fmt.Print("- 譜面を解析中 (Analyzing chart)... ")
	levelData, err := pjsekaioverlay.FetchLevelData(chartSource, chart)

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL:%s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	if !isOptionSpecified {
		fmt.Print("総合力を指定してください。\nInput your team's power.\n> ")
		var tmpTeamPower string
		fmt.Scanln(&tmpTeamPower)
		teamPower, err = strconv.Atoi(tmpTeamPower)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL:%s", err.Error())))
			return
		}
		fmt.Printf("\033[A\033[2K\r> %s\n", color.GreenString(tmpTeamPower))

	}

	fmt.Print("- スコアを計算中 (Calculating score)... ")
	scoreData := pjsekaioverlay.CalculateScore(chart, levelData, teamPower)

	fmt.Println(color.GreenString("OK"))

	if !isOptionSpecified {
		fmt.Print("コンボのAP表示を有効にしますか？ (Enable AP indicator for combo?) [y/n]\n> ")
		before, _ := rawmode.Enable()
		tmpEnableComboApByte, _ := bufio.NewReader(os.Stdin).ReadByte()
		tmpEnableComboAp := string(tmpEnableComboApByte)
		rawmode.Restore(before)
		fmt.Printf("\n\033[A\033[2K\r> %s\n", color.GreenString(tmpEnableComboAp))
		if tmpEnableComboAp == "Y" || tmpEnableComboAp == "y" || tmpEnableComboAp == "" {
			apCombo = true
		} else {
			apCombo = false
		}
	}
	executableDir := filepath.Dir(executablePath)
	assets := filepath.Join(executableDir, "assets")

	fmt.Print("- pedファイルを生成中 (Generating ped file)... ")

	err = pjsekaioverlay.WritePedFile(scoreData, assets, apCombo, filepath.Join(formattedOutDir, "data.ped"), sonolus.LevelInfo{Rating: chart.Rating})

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL:%s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	fmt.Print("- exoファイルを生成中 (Generating exo file)... ")

	composerAndVocals := []string{chart.Artists, "？"}
	if separateAttempt := strings.Split(chart.Artists, " / "); chartSource.Id == "chart_cyanvas" && len(separateAttempt) <= 2 {
		composerAndVocals = separateAttempt
	}

	artists := fmt.Sprintf("作詞：？    作曲：%s    編曲：？\r\nVo：%s   譜面作成：%s", composerAndVocals[0], composerAndVocals[1], chart.Author)

	err = pjsekaioverlay.WriteExoFiles(assets, formattedOutDir, chart.Title, artists)

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL:%s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	fmt.Println(color.GreenString("\n全ての処理が完了しました。READMEの規約を確認した上で、exoファイルをAviUtlにインポートして下さい。\nExecution complete! Please import the exo file into AviUtl after reviewing the README terms and conditions."))
}

func main() {
	isOptionSpecified := len(os.Args) > 1
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32

	windows.GetConsoleMode(stdout, &originalMode)
	windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	origMain(isOptionSpecified)

	if !isOptionSpecified {
		fmt.Print(color.CyanString("\n- 何かキーを押すと終了します...\n- Press any key to exit..."))

		before, _ := rawmode.Enable()
		bufio.NewReader(os.Stdin).ReadByte()
		rawmode.Restore(before)
	}
}
