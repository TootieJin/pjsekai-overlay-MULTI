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

//go:embed main.exo
var rawBaseExo []byte

//go:embed main_en.exo
var rawBaseExoEn []byte

func WriteExoFiles(assets string, destDir string, title string, description string) error {
	baseExo := string(rawBaseExo)
	baseExoEn := string(rawBaseExoEn)
	replacedExo := baseExo
	replacedExoEn := baseExoEn
	mapping := []string{
		"{assets}", strings.ReplaceAll(assets, "\\", "/"),
		"{dist}", strings.ReplaceAll(destDir, "\\", "/"),
		"{text:difficulty}", encodeString("MASTER"),
		"{text:extra}", encodeString("EXTRA INFO"),
		"{text:title}", encodeString(title),
		"{text:description}", encodeString(description),
	}
	for i := range mapping {
		if i%2 == 0 {
			continue
		}
		if !strings.Contains(replacedExo, mapping[i-1]) {
			panic(fmt.Sprintf("[JP] exoファイルの生成に失敗しました（%sが見つかりません）", mapping[i-1]))
		}
		if !strings.Contains(replacedExoEn, mapping[i-1]) {
			panic(fmt.Sprintf("[EN] Failed to generate exo file (%s not found)", mapping[i-1]))
		}
		replacedExo = strings.ReplaceAll(replacedExo, mapping[i-1], mapping[i])
		replacedExoEn = strings.ReplaceAll(replacedExoEn, mapping[i-1], mapping[i])
	}
	replacedExo = strings.ReplaceAll(replacedExo, "\n", "\r\n")
	replacedExoEn = strings.ReplaceAll(replacedExoEn, "\n", "\r\n")
	encodedExo, err := io.ReadAll(transform.NewReader(
		strings.NewReader(replacedExo), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return fmt.Errorf("エンコードに失敗しました (Encoding failed) [%w]", err)
	}
	encodedExoEn, err := io.ReadAll(transform.NewReader(
		strings.NewReader(replacedExoEn), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return fmt.Errorf("エンコードに失敗しました (Encoding failed) [%w]", err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "main.exo"),
		encodedExo,
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "main_en.exo"),
		encodedExoEn,
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}

	return nil
}
