package main

import (
	"fmt"
    "io/ioutil"

	"github.com/valyala/fasthttp"
	"github.com/fasthttp/router"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/svg"

	"crypto/sha256"
	"encoding/hex"

	"strings"
	"bytes"
	"text/template"

	"os/exec"
	"net/url"
)

func addLatexToTemplate(latexCode string) (string) {
	t, _ := template.New("foo").Parse(`{{define "T"}}
\documentclass[12pt,crop,tikz]{standalone}
\usepackage{stix2}
\usepackage{amsmath,amssymb,mathrsfs,amsfonts,amsthm,mathtools,color,ucs}
\usepackage[utf8x]{inputenc}
\usepackage{xparse}

\begin{document}
{{.}}
\end{document}{{end}}`)
    buf := &bytes.Buffer{}
    t.ExecuteTemplate(buf, "T", latexCode)
    return buf.String()
}


func LatexToSvg(ctx *fasthttp.RequestCtx) {
	latexCode, _ := url.QueryUnescape(ctx.UserValue("latex").(string))
	latexCodeByte := []byte(addLatexToTemplate(latexCode))
	latexCodeSha := sha256.Sum256(latexCodeByte)
	latexCodeShaHex := hex.EncodeToString(latexCodeSha[:])
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

	ctx.SetContentType("image/svg+xml")
	fmt.Fprintf(ctx, "%s", string(svgByteMin))
}

func main () {
	r := router.New()
	r.GET("/latex/{latex}", LatexToSvg)
	// pass plain function to fasthttp
	fasthttp.ListenAndServe(":4000", r.Handler)
}

