package main

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var allowList map[string]bool = make(map[string]bool)

type Cipher struct {
	toQian   map[rune]rune
	fromQian map[rune]rune
}

func NewCipher() *Cipher {
	const alphaNum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	const qian = "天地玄黄宇宙洪荒日月盈昃辰宿列张寒来暑往秋收冬藏闰余成岁律吕调阳云腾致雨露结为霜金生丽水玉出昆冈剑号巨阙珠称夜光果珍李柰菜重"

	toQian := make(map[rune]rune)
	fromQian := make(map[rune]rune)

	alphaRunes := []rune(alphaNum)
	qianRunes := []rune(qian)

	for i := 0; i < len(alphaRunes); i++ {
		a := alphaRunes[i]
		q := qianRunes[i]
		toQian[a] = q
		fromQian[q] = a
	}

	return &Cipher{toQian: toQian, fromQian: fromQian}
}

func (c *Cipher) Encode(input string) string {
	var builder strings.Builder
	for _, r := range input {
		if val, ok := c.toQian[r]; ok {
			builder.WriteRune(val)
		} else {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func (c *Cipher) Decode(input string) string {
	var builder strings.Builder
	for _, r := range input {
		if val, ok := c.fromQian[r]; ok {
			builder.WriteRune(val)
		} else {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func main() {
	eFile := flag.String("e", "", "csv file to encode")
	dFile := flag.String("d", "", "phiproto file to decode")
	output := flag.String("o", "", "output file (optional)")
	flag.Parse()

	if *eFile == "" && *dFile == "" {
		println("Please specify a file to encode or decode using -e or -d")
		os.Exit(1)
	}
	allowListEnv := os.Getenv("PHI_ALLOWLIST")
	for _, url := range strings.Split(allowListEnv, ",") {
		allowList[url] = true
	}
	if *eFile != "" {
		fileContent, err := os.ReadFile(*eFile)

		if err != nil {
			println("Error encoding file:", err.Error())
			os.Exit(1)
		}
		encodedContent, err := encodeFile(string(fileContent))
		if err != nil {
			println("Error encoding file:", err.Error())
			os.Exit(1)
		}
		if *output != "" {
			err = os.WriteFile(*output, encodedContent, 0644)
			if err != nil {
				println("Error writing output file:", err.Error())
				os.Exit(1)
			}
		} else {
			println(string(encodedContent))
		}
		return
	} else if *dFile != "" {
		fileContent, err := os.ReadFile(*dFile)
		if err != nil {
			println("Error decoding file:", err.Error())
			os.Exit(1)
		}
		decodedContent, err := decodeFile(fileContent)
		if err != nil {
			println("Error decoding file:", err.Error())
			os.Exit(1)
		}
		if *output != "" {
			err = os.WriteFile(*output, []byte(decodedContent), 0644)
			if err != nil {
				println("Error writing output file:", err.Error())
				os.Exit(1)
			}
		} else {
			println(decodedContent)
		}
		return
	}
}

func getCsvDims(m string) (rows int, cols int, err error) {
	r := csv.NewReader(strings.NewReader(m))
	r.FieldsPerRecord = -1

	records, err := r.ReadAll()
	if err != nil {
		return 0, 0, err
	}

	rows = len(records)
	for _, record := range records {
		if len(record) > cols {
			cols = len(record)
		}
	}

	return rows, cols, nil
}

func encodeFile(m string) ([]byte, error) {
	r, c, e := getCsvDims(m)
	if e != nil {
		return nil, e
	}
	result := make([]map[string]string, 0)
	rows := strings.Split(m, "\n")
	for i := 0; i < r; i++ {
		rowContent := make(map[string]string)
		for j := 0; j < c; j++ {
			encodedValue := base64.StdEncoding.EncodeToString([]byte(strings.Split(rows[i], ",")[j]))
			rowContent[strconv.Itoa(j)] = encodedValue
		}
		result = append(result, rowContent)
	}
	return json.Marshal(result)
}

func decodeFile(c []byte) (string, error) {
	var data []map[string]string
	cipher := NewCipher()
	runtime := cipher.Decode(rt)
	err := json.Unmarshal(c, &data)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	for _, row := range data {
		var rowValues []string
		for j := 0; j < len(row); j++ {
			field := row[strconv.Itoa(j)]
			b64decoded, err := base64.StdEncoding.DecodeString(field)
			if err != nil {
				println("Error decoding base64 field:", err.Error())
				continue
			}
			field = string(b64decoded)
			content = field
			rowValues = append(rowValues, content)
		}
		sb.WriteString(strings.Join(rowValues, ",") + "\n")
	}
	return sb.String(), nil

}
