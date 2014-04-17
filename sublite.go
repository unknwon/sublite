// Copyright 2014 Unknown
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// SubLite is a theme converter from Sublime Text to LiteIDE.
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const VERSION = "0.0.1"

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Incorrect aargument number: %d", len(os.Args))
	}
	fileName := os.Args[1]

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Fail to load Sublime theme: %v", err)
	}

	data, err = generateTheme(data, fileName)
	if err != nil {
		log.Fatalf("Fail to generate LiteIDE theme: %v", err)
	}

	if err = ioutil.WriteFile(strings.TrimSuffix(fileName, ".tmTheme")+".xml", data, os.ModePerm); err != nil {
		log.Fatalf("Fail to save LiteIDE theme: %v", err)
	}
}

func warpKey(name string) string {
	return "<key>" + name + "</key>"
}

func warpString(name string) string {
	return "<string>" + name + "</string>"
}

func generateTheme(data []byte, name string) ([]byte, error) {
	buf := bytes.NewBufferString(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<style-scheme version="1.0" name="%s">
	<!--
	Auto-generated based on %s
	https://github.com/Unknwon/sublite, Unknown, 2014
	-->
`, strings.TrimSuffix(name, ".tmTheme"), name))

	var fore, back string

	data, back = getKeyContent(data, warpKey("background"))
	data, fore = getKeyContent(data, warpKey("foreground"))
	buf.WriteString(fmt.Sprintf("\t<style name=\"Text\" foreground=\"%s\" background=\"%s\"/>\n", fore, back))

	data, back = getKeyContent(data, warpKey("lineHighlight"))
	buf.WriteString(fmt.Sprintf("\t<style name=\"CurrentLine\" background=\"%s\"/>\n", back))

	data, back = getKeyContent(data, warpKey("selection"))
	buf.WriteString(fmt.Sprintf("\t<style name=\"Selection\" background=\"%s\"/>\n", back))

	data, fore = getKeyContent(data, warpKey("bracketsForeground"))
	buf.WriteString(fmt.Sprintf("\t<style name=\"Symbol\" foreground=\"%s\"/>\n", fore))

	data, fore = getKeyContent(data, warpString("Comment"), warpKey("foreground"))
	buf.WriteString(fmt.Sprintf("\t<style name=\"Comment\" foreground=\"%s\"/>\n", fore))

	data, fore = getKeyContent(data, warpString("String"), warpKey("foreground"))
	buf.WriteString(fmt.Sprintf("\t<style name=\"String\" foreground=\"%s\"/>\n", fore))

	data, fore = getKeyContent(data, warpString("Number"), warpKey("foreground"))
	buf.WriteString(fmt.Sprintf("\t<style name=\"Decimal\" foreground=\"%s\"/>\n", fore))

	data, fore = getKeyContent(data, warpString("Built-in constant"), warpKey("foreground"))
	buf.WriteString(fmt.Sprintf("\t<style name=\"BuiltinFunc\" foreground=\"%s\"/>\n", fore))
	buf.WriteString(fmt.Sprintf("\t<style name=\"Predeclared\" foreground=\"%s\"/>\n", fore))
	buf.WriteString(fmt.Sprintf("\t<style name=\"Char\" foreground=\"%s\"/>\n", fore))

	data, fore = getKeyContent(data, warpString("Keyword"), warpKey("foreground"))
	buf.WriteString(fmt.Sprintf("\t<style name=\"Keyword\" foreground=\"%s\"/>\n", fore))

	data, fore = getKeyContent(data, warpString("Storage type"), warpKey("foreground"))
	buf.WriteString(fmt.Sprintf("\t<style name=\"DataType\" foreground=\"%s\"/>\n", fore))

	data, fore = getKeyContent(data, warpString("Function name"), warpKey("foreground"))
	buf.WriteString(fmt.Sprintf("\t<style name=\"FuncDecl\" foreground=\"%s\"/>\n", fore))

	// Umatch items: Extra, IndentLine, VisualWhitespace, BaseN, Float, Alert, Error, RegionMarker, Placeholder, ToDo

	buf.WriteString("</style-scheme>")
	return buf.Bytes(), nil
}

func getKeyContent(data []byte, stringNames ...string) ([]byte, string) {
	i := 0
	for _, strName := range stringNames {
		i = bytes.Index(data, []byte(strName))
		if i == -1 {
			return data, ""
		}
		data = data[i:]
	}
	return getNextString(data)
}

func getNextString(data []byte) ([]byte, string) {
	i := bytes.Index(data, []byte("</key>")) + 6
	return data[i:], string(bytes.TrimPrefix(bytes.TrimSpace(data[i:bytes.Index(data, []byte("</string>"))]), []byte("<string>")))
}
