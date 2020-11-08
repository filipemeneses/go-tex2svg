package main

import (
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

func LatexToSvg(latexCode string) []byte {
	latexCodeByte := []byte(addLatexToTemplate(string(latexCode)))
	hashGenerator := hmac.New(sha256.New, []byte(nil))
	hashGenerator.Write(latexCodeByte)
	latexCodeShaHex := hex.EncodeToString(hashGenerator.Sum(nil))
	latexFilePath := strings.Join([]string{"tmpfs/", latexCodeShaHex}, "")
	latexFilePathPdf := strings.Join([]string{latexFilePath, ".pdf"}, "")
	latexFilePathSvg := strings.Join([]string{latexFilePath, ".svg"}, "")

	ioutil.WriteFile(latexFilePath, latexCodeByte, 0644)

	cmd := exec.Command("pdflatex", "--interaction=batchmode", "--output-directory=tmpfs", latexFilePath)
	cmd.Start()
	cmd.Wait()

	cmdSvg := exec.Command("pdf2svg", latexFilePathPdf, latexFilePathSvg, "1")
	cmdSvg.Start()
	cmdSvg.Wait()

	svgByte, _ := ioutil.ReadFile(latexFilePathSvg)

	m := minify.New()
	m.AddFunc("image/svg+xml", svg.Minify)
	svgByteMin, _ := m.Bytes("image/svg+xml", svgByte)

	go delHexFiles(latexCodeShaHex)

	return svgByteMin
}

func HandleLatex(ctx *fasthttp.RequestCtx) {
	latexSvgMinBytes := LatexToSvg(string(ctx.PostBody()[:]))
	ctx.SetContentType("image/svg+xml")
	ctx.Response.Header.Add("Content-Encoding", "gzip")
	ctx.Write(fasthttp.AppendGzipBytes(nil, latexSvgMinBytes))
}

func delHexFiles(latexCodeShaHex string) {
	latexFilePath := strings.Join([]string{"tmpfs/", latexCodeShaHex}, "")
	latexFilePathPdf := strings.Join([]string{latexFilePath, ".pdf"}, "")
	latexFilePathLog := strings.Join([]string{latexFilePath, ".log"}, "")
	latexFilePathSvg := strings.Join([]string{latexFilePath, ".svg"}, "")
	os.Remove(latexFilePath)
	os.Remove(latexFilePathPdf)
	os.Remove(latexFilePathLog)
	os.Remove(latexFilePathSvg)
}

func main() {
	r := router.New()
	r.POST("/latex-to-svg", HandleLatex)
	fasthttp.ListenAndServe(":4000", r.Handler)
}
