//////////////////////////////////////////////////////////////////////
// templates.go
//
// @usage
// 
//     1. Import this package.
//
//         --------------------------------------------------
//         import myMailerTemplates "mailer/templates"
//         --------------------------------------------------
//
//     2. Generate mail body.
//
//         When using html files.
//
//             --------------------------------------------------
//             htmlFiles := []string{"/usr/local/go/workplace/src/noknow/example.html"}
//             htmlParams := map[string]string{
//                 "title": "noknow",
//                 "heading": "This is H1.",
//                 "html": "<p>You can output a HTML.</p>",
//             }
//             body, err := myMailerTemplates.HtmlFromFiles(htmlFile, htmlParams)
//             if err != nil {
//                 // Error handling.
//             }
//             --------------------------------------------------
//
//         When using text string.
//
//             --------------------------------------------------
//             text := "{{ .title }}\r\n{{ .heading }}\r\n{{ safeHTML .html }}"
//             textParams := map[string]string{
//                 "title": "noknow",
//                 "heading": "This is H1.",
//                 "html": "<p>You can output a HTML.</p>",
//             }
//             body, err := myMailerTemplates.HtmlFromString(text, textParams)
//             if err != nil {
//                 // Error handling.
//             }
//             --------------------------------------------------
//
//         You can embed HTML when you execute safeHTML function.
//
//
// MIT License
//
// Copyright (c) 2019 noknow.info
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A 
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTW//ARE. 
//////////////////////////////////////////////////////////////////////
package templates

import (
    "bytes"
    "html/template"
    "path"
)


//////////////////////////////////////////////////////////////////////
// Generate HTML string from html files
//////////////////////////////////////////////////////////////////////
func HtmlFromFiles(fileNames []string, params map[string]string) (string, error) {
    f := template.FuncMap{
        "safeHTML": func(s string) template.HTML { return template.HTML(s) },
    }
    t, err := template.New(path.Base(fileNames[0])).Funcs(f).ParseFiles(fileNames...)
    if err != nil {
        return "", err
    }
    buffer := new(bytes.Buffer)
    if err := t.Execute(buffer, params); err != nil {
        return "", err
    }
    return buffer.String(), nil
}


//////////////////////////////////////////////////////////////////////
// Generate HTML string from text
//////////////////////////////////////////////////////////////////////
func HtmlFromString(text string, params map[string]string) (string, error) {
    f := template.FuncMap{
        "safeHTML": func(s string) template.HTML { return template.HTML(s) },
    }
    t, err := template.New("t").Funcs(f).Parse(text)
    if err != nil {
        return "", err
    }
    buffer := new(bytes.Buffer)
    if err := t.Execute(buffer, params); err != nil {
        return "", err
    }
    return buffer.String(), nil
}


//////////////////////////////////////////////////////////////////////
// Generate text string from html files
//////////////////////////////////////////////////////////////////////
func TextFromFiles(fileNames []string, params map[string]string) (string, error) {
    f := template.FuncMap{
        "safeHTML": func(s string) template.HTML { return template.HTML(s) },
    }
    t, err := template.New(path.Base(fileNames[0])).Funcs(f).ParseFiles(fileNames...)
    if err != nil {
        return "", err
    }
    buffer := new(bytes.Buffer)
    if err := t.Execute(buffer, params); err != nil {
        return "", err
    }
    return buffer.String(), nil
}


//////////////////////////////////////////////////////////////////////
// Generate text string from text
//////////////////////////////////////////////////////////////////////
func TextFromString(text string, params map[string]string) (string, error) {
    f := template.FuncMap{
        "safeHTML": func(s string) template.HTML { return template.HTML(s) },
    }
    t, err := template.New("t").Funcs(f).Parse(text)
    if err != nil {
        return "", err
    }
    buffer := new(bytes.Buffer)
    if err := t.Execute(buffer, params); err != nil {
        return "", err
    }
    return buffer.String(), nil
}

