package pjsekaioverlay

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func encodeString(str string) string {
	bytes := utf16.Encode([]rune(str))
	encoded := make([]string, 1024)
	if len(str) > 1024 {
		panic("too long string")
	}
	for i := range encoded {
		var hex string
		if i >= len(bytes) {
			hex = fmt.Sprintf("%04x", 0)
		} else {
			hex = fmt.Sprintf("%02x%02x", bytes[i]&0xff, bytes[i]>>8)
		}

		encoded[i] = hex
	}

	return strings.Join(encoded, "")
}

//go:embed main_jp_16-9_1920x1080.exo
var rawBaseExoJP []byte

//go:embed main_jp_4-3_1440x1080.exo
var rawBaseExoJP43 []byte

//go:embed main_en_16-9_1920x1080.exo
var rawBaseExoEN []byte

//go:embed main_en_4-3_1440x1080.exo
var rawBaseExoEN43 []byte

func WriteExoFiles(assets string, destDir string, title string, description string) error {
	baseExoJP := string(rawBaseExoJP)
	baseExoJP43 := string(rawBaseExoJP43)
	baseExoEN := string(rawBaseExoEN)
	baseExoEN43 := string(rawBaseExoEN43)
	replacedExoJP := baseExoJP
	replacedExoJP43 := baseExoJP43
	replacedExoEN := baseExoEN
	replacedExoEN43 := baseExoEN43
	mapping := []string{
		"{assets}", strings.ReplaceAll(assets, "\\", "/"),
		"{dist}", strings.ReplaceAll(destDir, "\\", "/"),
		"{text:difficulty}", encodeString("APPEND"),
		"{text:extra}", encodeString("動画：TootieJin"),
		"{text:title}", encodeString(title),
		"{text:description}", encodeString(description),
	}
	for i := range mapping {
		if i%2 == 0 {
			continue
		}
		if !strings.Contains(replacedExoJP, mapping[i-1]) {
			panic(fmt.Sprintf("exoファイルの生成に失敗しました (Failed to generate exo file) [Missing: %s]", mapping[i-1]))
		}
		if !strings.Contains(replacedExoJP43, mapping[i-1]) {
			panic(fmt.Sprintf("exoファイルの生成に失敗しました (Failed to generate exo file) [Missing: %s]", mapping[i-1]))
		}
		if !strings.Contains(replacedExoEN, mapping[i-1]) {
			panic(fmt.Sprintf("exoファイルの生成に失敗しました (Failed to generate exo file) [Missing: %s]", mapping[i-1]))
		}
		if !strings.Contains(replacedExoEN43, mapping[i-1]) {
			panic(fmt.Sprintf("exoファイルの生成に失敗しました (Failed to generate exo file) [Missing: %s]", mapping[i-1]))
		}
		replacedExoJP = strings.ReplaceAll(replacedExoJP, mapping[i-1], mapping[i])
		replacedExoJP43 = strings.ReplaceAll(replacedExoJP43, mapping[i-1], mapping[i])
		replacedExoEN = strings.ReplaceAll(replacedExoEN, mapping[i-1], mapping[i])
		replacedExoEN43 = strings.ReplaceAll(replacedExoEN43, mapping[i-1], mapping[i])
	}
	replacedExoJP = strings.ReplaceAll(replacedExoJP, "\n", "\r\n")
	replacedExoJP43 = strings.ReplaceAll(replacedExoJP43, "\n", "\r\n")
	replacedExoEN = strings.ReplaceAll(replacedExoEN, "\n", "\r\n")
	replacedExoEN43 = strings.ReplaceAll(replacedExoEN43, "\n", "\r\n")

	encodedExoJP, err := io.ReadAll(transform.NewReader(
		strings.NewReader(replacedExoJP), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return fmt.Errorf("エンコードに失敗しました (Encoding failed) [%w]", err)
	}
	encodedExoJP43, err := io.ReadAll(transform.NewReader(
		strings.NewReader(replacedExoJP43), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return fmt.Errorf("エンコードに失敗しました (Encoding failed) [%w]", err)
	}
	encodedExoEN, err := io.ReadAll(transform.NewReader(
		strings.NewReader(replacedExoEN), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return fmt.Errorf("エンコードに失敗しました (Encoding failed) [%w]", err)
	}
	encodedExoEN43, err := io.ReadAll(transform.NewReader(
		strings.NewReader(replacedExoEN43), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return fmt.Errorf("エンコードに失敗しました (Encoding failed) [%w]", err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "main_jp_16-9_1920x1080.exo"),
		encodedExoJP,
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "main_jp_4-3_1440x1080.exo"),
		encodedExoJP43,
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "main_en_16-9_1920x1080.exo"),
		encodedExoEN,
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "main_en_4-3_1440x1080.exo"),
		encodedExoEN43,
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	return nil
}
