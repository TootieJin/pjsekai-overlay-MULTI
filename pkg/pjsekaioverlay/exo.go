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

//go:embed main_16-9_1920x1080.exo
var rawBaseExo []byte

//go:embed main_4-3_1440x1080.exo
var rawBaseExo43 []byte

func WriteExoFiles(assets string, destDir string, title string, description string) error {
	baseExo := string(rawBaseExo)
	baseExo43 := string(rawBaseExo43)
	replacedExo := baseExo
	replacedExo43 := baseExo43
	mapping := []string{
		"{assets}", strings.ReplaceAll(assets, "\\", "/"),
		"{dist}", strings.ReplaceAll(destDir, "\\", "/"),
		"{text:difficulty}", encodeString("APPEND"),
		"{text:extra}", encodeString("Made using pjsekai-overlay-APPEND"),
		"{text:title}", encodeString(title),
		"{text:description}", encodeString(description),
	}
	for i := range mapping {
		if i%2 == 0 {
			continue
		}
		if !strings.Contains(replacedExo, mapping[i-1]) {
			panic(fmt.Sprintf("exoファイルの生成に失敗しました (Failed to generate exo file) [Missing: %s]", mapping[i-1]))
		}
		if !strings.Contains(replacedExo43, mapping[i-1]) {
			panic(fmt.Sprintf("exoファイルの生成に失敗しました (Failed to generate exo file) [Missing: %s]", mapping[i-1]))
		}
		replacedExo = strings.ReplaceAll(replacedExo, mapping[i-1], mapping[i])
		replacedExo43 = strings.ReplaceAll(replacedExo43, mapping[i-1], mapping[i])
	}
	replacedExo = strings.ReplaceAll(replacedExo, "\n", "\r\n")
	replacedExo43 = strings.ReplaceAll(replacedExo43, "\n", "\r\n")

	encodedExo, err := io.ReadAll(transform.NewReader(
		strings.NewReader(replacedExo), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return fmt.Errorf("エンコードに失敗しました (Encoding failed) [%w]", err)
	}
	encodedExo43, err := io.ReadAll(transform.NewReader(
		strings.NewReader(replacedExo43), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return fmt.Errorf("エンコードに失敗しました (Encoding failed) [%w]", err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "main_16-9_1920x1080.exo"),
		encodedExo,
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "main_4-3_1440x1080.exo"),
		encodedExo43,
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	return nil
}
