package main

import (
	"fmt"
	"io/ioutil"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/svg"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"bytes"
	"strings"
	"text/template"

	"os"
	"os/exec"

	"encoding/base64"
)

func addLatexToTemplate(latexCode string) string {
	t, _ := template.New("foo").Parse(`{{define "T"}}
\documentclass{standalone}
\nofiles
\usepackage{stix2}
\usepackage{amsmath,amssymb,mathrsfs,amsfonts,amsthm,mathtools,color,ucs}
\usepackage[utf8x]{inputenc}
\usepackage{xparse}
\usepackage{chemfig}
\pdfcompresslevel    = 0
\pdfobjcompresslevel = 0
\begin{document}
{{.}}
\end{document}{{end}}`)
	buf := &bytes.Buffer{}
	t.ExecuteTemplate(buf, "T", latexCode)
	return buf.String()
}

func LatexToSvg(ctx *fasthttp.RequestCtx) {
	latexBase64 := ctx.UserValue("latex").(string)
	latexCode, _ := base64.StdEncoding.DecodeString(latexBase64)
	latexCodeByte := []byte(addLatexToTemplate(string(latexCode)))
	hashGenerator := hmac.New(sha256.New, []byte("1"))
	hashGenerator.Write(latexCodeByte)
	latexCodeShaHex := hex.EncodeToString(hashGenerator.Sum(nil))
	latexFilePath := strings.Join([]string{"tmpfs/", latexCodeShaHex}, "")
	latexFilePathPdf := strings.Join([]string{latexFilePath, ".pdf"}, "")
	latexFilePathSvg := strings.Join([]string{latexFilePath, ".svg"}, "")

	ioutil.WriteFile(latexFilePath, latexCodeByte, 0644)

	cmd := exec.Command("pdflatex", "--interaction=batchmode", "--output-directory=tmpfs", latexFilePath)
	cmd.Start()
	cmd.Wait()

	cmdSvg := exec.Command("pdf2svg", latexFilePathPdf, latexFilePathSvg)
	cmdSvg.Start()
	cmdSvg.Wait()

	svgByte, _ := ioutil.ReadFile(latexFilePathSvg)

	m := minify.New()
	m.AddFunc("image/svg+xml", svg.Minify)
	svgByteMin, _ := m.Bytes("image/svg+xml", svgByte)

	ioutil.WriteFile(latexFilePathSvg, svgByteMin, 0644)

	go delHexFiles(latexCodeShaHex)

	ctx.SetContentType("image/svg+xml")
	fmt.Fprintf(ctx, "%s", string(svgByteMin))
}

func delHexFiles(latexCodeShaHex string) {
	latexFilePath := strings.Join([]string{"tmpfs/", latexCodeShaHex}, "")
	latexFilePathPdf := strings.Join([]string{latexFilePath, ".pdf"}, "")
	latexFilePathLog := strings.Join([]string{latexFilePath, ".log"}, "")
	os.Remove(latexFilePath)
	os.Remove(latexFilePathPdf)
	os.Remove(latexFilePathLog)
}

func main() {
	r := router.New()
	r.GET("/latex/{latex}", LatexToSvg)
	fasthttp.ListenAndServe(":4000", r.Handler)
}
