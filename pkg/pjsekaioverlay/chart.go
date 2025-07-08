package pjsekaioverlay

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"golang.org/x/image/draw"

	"github.com/TootieJin/pjsekai-overlay-MULTI/pkg/sonolus"
)

type Source struct {
	Id    string
	Name  string
	Color int
	Host  string
}

func FetchChart(source Source, chartId string) (sonolus.LevelInfo, error) {
	var url = "https://" + source.Host + "/sonolus/levels/" + chartId

	resp, err := http.Get(url)

	if err != nil {
		return sonolus.LevelInfo{}, errors.New("サーバーに接続できませんでした。(Could not connect to server.)")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return sonolus.LevelInfo{}, errors.New("譜面が見つかりませんでした。(Unable to search chart.)")
	}

	var chart sonolus.InfoResponse[sonolus.LevelInfo]
	json.NewDecoder(resp.Body).Decode(&chart)

	return chart.Item, nil
}

func FetchAPIChart(source Source, chartId string) (sonolus.LevelAPIInfo, error) {
	var url = "https://" + source.Host + "/api/charts/" + chartId

	resp, err := http.Get(url)

	if err != nil {
		return sonolus.LevelAPIInfo{}, errors.New("APIサーバーに接続できませんでした。(Could not connect to API server.)")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return sonolus.LevelAPIInfo{}, errors.New("API譜面が見つかりませんでした。(Unable to search API chart.)")
	}

	var chart sonolus.InfoAPIResponse[sonolus.LevelAPIInfo]
	json.NewDecoder(resp.Body).Decode(&chart)

	if source.Id == "untitled_sekai" {
		return chart.Data, nil
	} else {
		return chart.Chart, nil
	}
}

func DetectChartSource(chartId string) (Source, error) {
	var source Source
	if strings.HasPrefix(chartId, "ptlv-") {
		source = Source{
			Id:    "potato_leaves",
			Name:  "Potato Leaves",
			Color: 0x88cb7f,
			Host:  "ptlv.sevenc7c.com",
		}
	} else if strings.HasPrefix(chartId, "chcy-") {
		source = Source{
			Id:    "chart_cyanvas",
			Name:  "Chart Cyanvas",
			Color: 0x83ccd2,
			Host:  "cc.sevenc7c.com",
		}
	} else if strings.HasPrefix(chartId, "utsk-") {
		source = Source{
			Id:    "untitled_sekai",
			Name:  "Untitled Sekai",
			Color: 0x6a6a6a,
			Host:  "us.pim4n-net.com",
		}
	}
	if source.Id == "" {
		return Source{
			Id:    chartId,
			Name:  "",
			Color: 0,
			Host:  "",
		}, errors.New("unknown chart source")
	}
	return source, nil
}

func FetchLevelData(source Source, level sonolus.LevelInfo) (sonolus.LevelData, error) {
	url, err := sonolus.JoinUrl("https://"+source.Host, level.Data.Url)
	if err != nil {
		return sonolus.LevelData{}, fmt.Errorf("URLの解析に失敗しました。(URL parsing failed.) [%s]", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return sonolus.LevelData{}, fmt.Errorf("サーバーに接続できませんでした。(Could not connect to server.) [%s]", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return sonolus.LevelData{}, fmt.Errorf("譜面データが見つかりませんでした。(No chart data found.) [%d]", resp.StatusCode)
	}

	var data sonolus.LevelData
	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return sonolus.LevelData{}, fmt.Errorf("譜面データの読み込みに失敗しました。(Loading chart data failed.) [%s]", err)
	}

	err = json.NewDecoder(gzipReader).Decode(&data)

	if err != nil {
		return sonolus.LevelData{}, fmt.Errorf("譜面データの読み込みに失敗しました。(Loading chart data failed.) [%s]", err)
	}

	return data, nil
}

func DownloadCover(source Source, level sonolus.LevelInfo, destPath string) error {
	url, err := sonolus.JoinUrl("https://"+source.Host, level.Cover.Url)
	if err != nil {
		return fmt.Errorf("URLの解析に失敗しました。(URL parsing failed.) [%s]", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("サーバーに接続できませんでした。（%s）", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("サーバーに接続できませんでした。(Could not connect to server.) [%d]", resp.StatusCode)
	}

	os.MkdirAll(destPath, 0755)
	imageData, _, err := image.Decode(resp.Body)

	if err != nil {
		return fmt.Errorf("ジャケットの読み込みに失敗しました。(Loading jacket failed.) [%s]", err)
	}

	// 画像のリサイズ

	newImage := image.NewRGBA(image.Rect(0, 0, 512, 512))

	draw.ApproxBiLinear.Scale(newImage, newImage.Bounds(), imageData, imageData.Bounds(), draw.Over, nil)

	file, err := os.Create(path.Join(destPath, "cover.png"))

	if err != nil {
		return fmt.Errorf("ファイルの作成に失敗しました。(Failed to create file.) [%s]", err)
	}

	defer file.Close()

	err = png.Encode(file, newImage)

	if err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました。(Failed to write file.) [%s]", err)
	}

	return nil
}
func DownloadBackground(source Source, level sonolus.LevelInfo, destPath string, chartId string) error {
	backgroundUrl, err := sonolus.JoinUrl("https://"+source.Host, level.UseBackground.Item.Image.Url)
	if err != nil {
		return fmt.Errorf("URLの解析に失敗しました。(URL parsing failed.) [%s]", err)
	}

	resp, err := http.Get(backgroundUrl)
	if err != nil {
		return fmt.Errorf("サーバーに接続できませんでした。(Could not connect to server.) [%s]", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("背景が見つかりませんでした。(Background not found.) [%d]", resp.StatusCode)
	}

	var file *os.File
	var filev1 *os.File

	if strings.Contains(chartId, "?c_background=1") && source.Id == "chart_cyanvas" {
		filev1, err = os.Create(path.Join(destPath, "background-v1.png"))
		file = nil
	} else if source.Id == "potato_leaves" {
		filev1, err = os.Create(path.Join(destPath, "background-v1.png"))
		file = nil
	} else {
		filev1 = nil
		file, err = os.Create(path.Join(destPath, "background.png"))
	}

	if err != nil {
		return fmt.Errorf("ファイルの作成に失敗しました。(Failed to create file.) [%s]", err)
	}

	defer file.Close()
	defer filev1.Close()

	if file != nil {
		if _, err := io.Copy(file, resp.Body); err != nil {
			return fmt.Errorf("ファイルの書き込みに失敗しました。(Failed to write file.) [%s]", err)
		}
	}
	if filev1 != nil {
		if _, err := io.Copy(filev1, resp.Body); err != nil {
			return fmt.Errorf("v1ファイルの書き込みに失敗しました。(Failed to write v1 file.) [%s]", err)
		}
	}

	return nil
}
