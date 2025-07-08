package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/TootieJin/pjsekai-overlay-MULTI/pkg/pjsekaioverlay"
	"github.com/TootieJin/pjsekai-overlay-MULTI/pkg/sonolus"
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
	release, _, err := githubClient.Repositories.GetLatestRelease(context.Background(), "TootieJin", "pjsekai-overlay-MULTI")
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

func checkSubstrings(str []string, subs ...string) string {
	for _, s := range str {
		for _, sub := range subs {
			if strings.EqualFold(s, sub) {
				return sub
			}
		}
	}
	return ""
}

func origMain(isOptionSpecified bool) {
	Title()

	var skipAviutlInstall bool
	flag.BoolVar(&skipAviutlInstall, "no-aviutl-install", false, "AviUtlオブジェクトのインストールをスキップします。(AviUtl object installation is skipped.)")

	var outDir string
	flag.StringVar(&outDir, "out-dir", "./dist/_chartId_", "出力先ディレクトリを指定します。_chartId_ は譜面IDに置き換えられます。\nEnter the output path. _chartId_ will be replaced with the chart ID.")

	var exScore bool
	flag.BoolVar(&exScore, "ex-score", false, "大会モードを有効にします。 (Enable Tournament Mode.)")

	var teamPower float64
	flag.Float64Var(&teamPower, "team-power", 250000, "総合力を指定します。(Enter the team's power.)")

	var powerWeight []float64

	var playerPos int

	var missRate []float64

	var enUI bool
	flag.BoolVar(&enUI, "en-ui", false, "英語版を使う(イントロ + v3 UI) - Use English version (Intro + v3 UI)")

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
		fmt.Print("\n譜面IDを接頭辞込みで入力して下さい。\nEnter the chart ID including the prefix.\n\n'chcy-': Chart Cyanvas (cc.sevenc7c.com)\n'ptlv-': Potato Leaves (ptlv.sevenc7c.com)\n'utsk-': Untitled Sekai (us.pim4n-net.com)\n> ")
		fmt.Scanln(&chartId)
		fmt.Printf("\033[A\033[2K\r> %s\n", color.GreenString(chartId))
	}

	chartSource, err := pjsekaioverlay.DetectChartSource(chartId)
	if err != nil {
		fmt.Println(color.RedString("FAIL: 譜面が見つかりません。接頭辞も込め、正しい譜面IDを入力して下さい。\nChart not found. Please enter the correct chart ID including the prefix."))
		return
	}
	fmt.Printf("- 譜面を取得中 (Getting chart): %s%s%s ", RgbColorEscape(chartSource.Color), chartSource.Name, ResetEscape())
	chart, err := pjsekaioverlay.FetchChart(chartSource, chartId)
	chartv1, errv1 := pjsekaioverlay.FetchChart(chartSource, chartId+"?c_background=1")

	var chart_api sonolus.LevelAPIInfo
	if chartSource.Id == "chart_cyanvas" {
		chart_api, err = pjsekaioverlay.FetchAPIChart(chartSource, chartId[5:])
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}
	}
	if chartSource.Id == "untitled_sekai" {
		chart_api, err = pjsekaioverlay.FetchAPIChart(chartSource, chartId)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}
	}

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}
	if errv1 != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", errv1.Error())))
		return
	}

	if chart.Engine.Version != 13 {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL (ver.%d):エンジンのバージョンが古い。pjsekai-overlay-MULTIを最新版に更新してください。\nUnsupported engine version. Please update pjsekai-overlay-MULTI to latest version.", chart.Engine.Version)))
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
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	formattedOutDir := filepath.Join(cwd, strings.Replace(outDir, "_chartId_", chartId+"-multi", -1))
	fmt.Printf("- 出力先ディレクトリ (Output path): %s\n", color.CyanString(filepath.Dir(formattedOutDir)))

	fmt.Print("- ジャケットをダウンロード中 (Downloading jacket)... ")
	err = pjsekaioverlay.DownloadCover(chartSource, chart, formattedOutDir)
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	fmt.Print("- 背景をダウンロード中 (Downloading background)... ")
	err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir, chartId)
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}
	err = pjsekaioverlay.DownloadBackground(chartSource, chartv1, formattedOutDir, chartId+"?c_background=1")
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	fmt.Print("- 譜面を解析中 (Analyzing chart)... ")
	levelData, err := pjsekaioverlay.FetchLevelData(chartSource, chart)

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	var tmpTeamPower string
	if !isOptionSpecified {
		fmt.Print("\n総合力を指定してください。 (Input your team power.)\n\n- 小数と科学的記数法が使える (Accepts decimals & scientific notation)\n- おすすめ (Recommended): 250000 - 300000\n- 制限 (Limit): ???\n> ")
		fmt.Scanln(&tmpTeamPower)
		teamPower, err = strconv.ParseFloat(tmpTeamPower, 64)
		if err != nil {
			if numErr, ok := err.(*strconv.NumError); ok && numErr.Err == strconv.ErrRange {
				fmt.Println(color.RedString("FAIL: あなたのPCがその総合力で計算できないのは残念だ。説明書を読んで再実行してください。\nToo bad your PC can't calculate with that team power. Read the instructions and rerun it."))
				return
			} else {
				fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
				return
			}
		}
		if teamPower >= math.Abs(1e+33) {
			fmt.Printf("\033[A\033[2K\r> %s\n", color.YellowString(tmpTeamPower))
			fmt.Println(color.YellowString("WARNING: スコアは大きすぎると精度が落ちる可能性がある。Score may decrease precision if it's too large.\n"))
		} else {
			fmt.Printf("\033[A\033[2K\r> %s\n", color.GreenString(tmpTeamPower))
		}
	}

	if !isOptionSpecified {
		fmt.Printf("\n総合力の重さを指定してください。(Input player's power weight.)\n式/Formula: powerWeight * %s\n例/Example: 1,3,2.5,1.5,0.5\nランダムに生成するには空のままにする。(Leave empty to generate randomly.)\n> ", tmpTeamPower)
		var tmpPowerWeightString string
		fmt.Scanln(&tmpPowerWeightString)
		tmpPowerWeight := strings.Split(tmpPowerWeightString, ",")

		powerWeight = make([]float64, len(tmpPowerWeight))
		if tmpPowerWeightString == "" {
			powerWeight = []float64{float64(rand.Intn(3)) + rand.Float64(), float64(rand.Intn(3)) + rand.Float64(), float64(rand.Intn(3)) + rand.Float64(), float64(rand.Intn(3)) + rand.Float64(), float64(rand.Intn(3)) + rand.Float64()}
		} else {
			for i, v := range tmpPowerWeight {
				val, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
				if err != nil {
					fmt.Println(color.RedString(fmt.Sprintf("FAIL: 総合力の重さが無効です。(Invalid power weight): %s", err.Error())))
					return
				}
				powerWeight[i] = val
			}
		}

		if len(powerWeight) != 5 {
			fmt.Println(color.RedString("FAIL: 入力不足 (Not enough input)"))
			return
		}
		fmt.Printf("\033[A\033[2K\r> %s\n", color.GreenString(fmt.Sprint(powerWeight)))
	}

	if !isOptionSpecified {
		fmt.Print("\n選手のポジションを入力してください。(Enter your player position.)\nランダムに生成するには空のままにする。(Leave empty to generate randomly.)\n> ")
		before, _ := rawmode.Enable()
		tmpPlayerPosByte, _ := bufio.NewReader(os.Stdin).ReadByte()
		tmpPlayerPos := string(tmpPlayerPosByte)
		rawmode.Restore(before)
		switch tmpPlayerPos {
		case "\r", "\n", "":
			playerPos = 1 + rand.Intn(5)
		case "1":
			playerPos = 1
		case "2":
			playerPos = 2
		case "3":
			playerPos = 3
		case "4":
			playerPos = 4
		case "5":
			playerPos = 5
		}
		if playerPos < 1 || playerPos > 5 {
			fmt.Println(color.RedString("FAIL: ポジションは1から5の間でなければなりません。\nPosition must be between 1 and 5."))
			return
		}
		fmt.Printf("\n\033[A\033[2K\r> %s\n", color.GreenString(fmt.Sprint(playerPos)))
	}

	if !isOptionSpecified {
		fmt.Print("\n選手のミスチャンスを指定してください。(Input player's combo break chance.)\n例/Example: 0.01,0.001,0.08,0.002,0.05\nランダムに生成するには空のままにする。(Leave empty to generate randomly.)\n> ")
		var tmpMissRateString string
		fmt.Scanln(&tmpMissRateString)
		tmpMissRate := strings.Split(tmpMissRateString, ",")

		missRate = make([]float64, len(tmpMissRate))
		if tmpMissRateString == "" {
			missRate = []float64{1 / float64(rand.Intn(1000)), 1 / float64(rand.Intn(1000)), 1 / float64(rand.Intn(1000)), 1 / float64(rand.Intn(1000)), 1 / float64(rand.Intn(1000))}
		} else {
			for i, v := range tmpMissRate {
				val, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
				if err != nil {
					fmt.Println(color.RedString(fmt.Sprintf("FAIL: ミスチャンスが無効です。(Invalid combo break chance): %s", err.Error())))
					return
				}
				if val > 1 {
					fmt.Println(color.RedString(fmt.Sprintf("FAIL: ミスチャンスが無効です。(Invalid combo break chance): %f > 1", val)))
					return
				}
				missRate[i] = val
			}
		}

		if len(missRate) != 5 {
			fmt.Println(color.RedString("FAIL: 入力不足 (Not enough input)"))
			return
		}
		fmt.Printf("\033[A\033[2K\r> %s\n", color.GreenString(fmt.Sprint(missRate)))
	}

	fmt.Print("- スコアを計算中 (Calculating score)... ")
	scoreData := pjsekaioverlay.CalculateScore(chart, levelData, teamPower, powerWeight, missRate, playerPos)

	fmt.Println(color.GreenString("OK"))
	if !isOptionSpecified {
		fmt.Print("\n英語版を使う？(イントロ + v3 UI) - Use English version? (Intro + v3 UI) [y/n]\n> ")
		before, _ := rawmode.Enable()
		tmpEnableENByte, _ := bufio.NewReader(os.Stdin).ReadByte()
		tmpEnableEN := string(tmpEnableENByte)
		rawmode.Restore(before)
		if tmpEnableEN == "Y" || tmpEnableEN == "y" {
			enUI = true
			fmt.Printf("\n\033[A\033[2K\r> %s\n", color.GreenString(tmpEnableEN))
			fmt.Println(color.GreenString("TOGGLE: ON"))
		} else {
			enUI = false
			fmt.Printf("\n\033[A\033[2K\r> %s\n", color.RedString(tmpEnableEN))
			fmt.Println(color.RedString("TOGGLE: OFF"))
		}
	}

	if !isOptionSpecified {
		fmt.Print("\nコンボのAP表示を有効にしますか？(これは後でAviUtlで変更できます)\nEnable AP indicator for combo? (You can change this later in AviUtl) [y/n]\n> ")
		before, _ := rawmode.Enable()
		tmpEnableComboApByte, _ := bufio.NewReader(os.Stdin).ReadByte()
		tmpEnableComboAp := string(tmpEnableComboApByte)
		rawmode.Restore(before)

		if tmpEnableComboAp == "Y" || tmpEnableComboAp == "y" || tmpEnableComboAp == "" {
			apCombo = true
			fmt.Printf("\n\033[A\033[2K\r> %s\n", color.GreenString(tmpEnableComboAp))
			fmt.Println(color.GreenString("TOGGLE: ON"))
		} else {
			apCombo = false
			fmt.Printf("\n\033[A\033[2K\r> %s\n", color.RedString(tmpEnableComboAp))
			fmt.Println(color.RedString("TOGGLE: OFF"))
		}
	}
	executableDir := filepath.Dir(executablePath)
	assets := filepath.Join(executableDir, "assets")

	fmt.Print("\n- pedファイルを生成中 (Generating ped file)... ")

	err = pjsekaioverlay.WritePedFile(scoreData, assets, filepath.Join(formattedOutDir, "data.ped"), sonolus.LevelInfo{Rating: chart.Rating}, levelData, powerWeight, enUI)
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}
	err = pjsekaioverlay.WritePedMultiFile(scoreData, assets, filepath.Join(formattedOutDir, "data-multi.ped"), levelData, playerPos)
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	fmt.Print("- exoファイルを生成中 (Generating exo file)... ")

	var difficulty string
	difficultyStrings := []string{"EASY", "NORMAL", "HARD", "EXPERT", "MASTER", "APPEND", "ETERNAL"}
	if tags := checkSubstrings(chart_api.Tags, difficultyStrings...); tags != "" {
		difficulty = tags
	} else if title := checkSubstrings(strings.Fields(chart.Title), difficultyStrings...); title != "" {
		difficulty = title
	} else {
		difficulty = "APPEND"
	}

	composerAndVocals := []string{chart.Artists, "-"}
	if separateAttempt := strings.Split(chart.Artists, " / "); chartSource.Id == "chart_cyanvas" && len(separateAttempt) <= 2 {
		composerAndVocals = separateAttempt
	}

	charter := []string{chart.Author, "-"}
	if charterTag := strings.Split(chart.Author, "#"); (chartSource.Id == "chart_cyanvas" || chartSource.Id == "untitled_sekai") && len(charterTag) <= 2 {
		charter = charterTag
	}

	description := fmt.Sprintf("作詞：-    作曲：%s    編曲：-\r\nVo：%s    譜面作成：%s", composerAndVocals[0], composerAndVocals[1], charter[0])
	descriptionv1 := fmt.Sprintf("作詞：-   作曲：%s   編曲：-\r\n歌：%s   譜面作成：%s", composerAndVocals[0], composerAndVocals[1], charter[0])
	extra := "[追加情報]"
	ap := "0.00"

	if enUI {
		description = fmt.Sprintf("Lyrics: -    Music: %s    Arrangement: -\r\nVo: %s    Chart Design: %s", composerAndVocals[0], composerAndVocals[1], charter[0])
		descriptionv1 = fmt.Sprintf("Lyrics: -   Music: %s   Arrangement: -\r\n歌：%s   Chart Design: %s", composerAndVocals[0], composerAndVocals[1], charter[0])
		extra = "[Additional Info]"
	}
	if apCombo {
		ap = "1.00"
	}

	err = pjsekaioverlay.WriteExoFiles(assets, formattedOutDir, chart.Title, description, descriptionv1, difficulty, extra, ap)

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	fmt.Println(color.GreenString("\n全ての処理が完了しました。READMEの規約を確認した上で、exoファイルをAviUtlにインポートして下さい。\nExecution complete! Please import the exo file into AviUtl after reviewing the README Terms of Use."))
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
