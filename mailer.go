//////////////////////////////////////////////////////////////////////
// mailer.go
//
// @usage
// 
//     1. Import this package.
//
//         --------------------------------------------------
//         import myMailer "mailer"
//         --------------------------------------------------
//
//     2. Prepare following variables.
//
//         - SMTP Server Host: string
//         - SMTP Server Post: int
//         - Mail From: string
//         - Mail To: string
//         - Mail Subject: string
//         - (Optional) Mail Auth User Name: string
//         - (Optional) Mail Auth Password: string
//
//         --------------------------------------------------
//         smtpServerHost := "example.com"
//         smtpServerPort := 465
//         authUserName := "noknow"
//         authPassword := "noknow_pass"
//         authHost := "example.com"
//         from := "noknow<noreply@example.com>"
//         to := "user@example.com"
//         subject := "This is a subject."
//         --------------------------------------------------
//
//     3. Generate a mail body.
//
//         3-A. Prepare a body parameters.
//
//             --------------------------------------------------
//             bodyParams := map[string]string{
//                 "langCode": "en",
//                 "title": "noknow",
//                 "heading": "This is a heading.",
//             }
//             --------------------------------------------------
//
//         3-1. When using HTML file and text file, which means multipart/alternative.
//
//             [sample.html]
//             --------------------------------------------------
//             <!DOCTYPE html>
//             <html lang="{{ .langCode }}">
//             <head>
//               {{template "head" .}}
//               <title>{{ .title }}</title>
//             </head>
//             <body>
//               <table>
//                 <tr>
//                   <td><h1>{{ .heading }}</h1></td>
//                 </tr>
//               </table>
//             </body>
//             </html>
//             --------------------------------------------------
//
//             [head.html]
//             --------------------------------------------------
//             {{define "head"}}
//             <meta charset="UTF-8">
//             <meta name="viewport" content="width=device-width, initial-scale=1.0">
//             <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
//             {{end}}
//             --------------------------------------------------
//
//             [sample.txt]
//             --------------------------------------------------
//             {{ .title }}
//             {{ .heading }}
//             --------------------------------------------------
//
//             --------------------------------------------------
//             // HTML body.
//             htmlFiles := []string{
//                 "sample.html",
//                 "head.html",
//             }
//             htmlBody, err := myMailer.GenBodyFromFiles(
//                 myMailer.CONTENT_TYPE_TEXT_HTML,
//                 myMailer.CHARSET_UTF8,
//                 htmlFiles,
//                 bodyParams,
//             )
//             if err != nil {
//                 // Error handling.
//             }
//
//             // HTML body.
//             textFiles := []string{
//                 "sample.txt",
//             }
//             textBody, err := myMailer.GenBodyFromFiles(
//                 myMailer.CONTENT_TYPE_TEXT_PLAIN,
//                 myMailer.CHARSET_UTF8,
//                 textFiles,
//                 bodyParams,
//             )
//             if err != nil {
//                 // Error handling.
//             }
//
//             // If body are 2 or more, the last index would be used in according to RFC1341.
//             // In the following case, HTML body would be used.
//             // You can set only one body.
//             body := []*myMailer.Body{textBody, htmlBody}
//             --------------------------------------------------
//
//     4. Generate a mail header.
//
//         --------------------------------------------------
//         header := myMailer.GenHeader(from, to, subject, myMailer.MIME_VERSION_1_0)
//         --------------------------------------------------
//
//     5. Generate an authentication config. (Optional)
//
//         --------------------------------------------------
//         authConfig := myMailer.GenPlainAuth(authUserName, authPassword, authHost)
//         --------------------------------------------------
//
//     6. Generate a TLS config. (Optional)
//
//         --------------------------------------------------
//         tlsConfig := myMailer.GenTlsConfig(smtpServerHost)
//
//         // When using certificate and private key par.
//         certFile := "/etc/letsencrypt/live/noknow.info/fullchain.pem"
//         keyFile := "/etc/letsencrypt/live/noknow.info/privkey.pem"
//         tlsConfig, err = myMailer.SetCertFiles(tlsConfig, certFile, keyFile)
//         if err != nil {
//             // Error handling.
//         }
//         --------------------------------------------------
//
//     7. Generate a mail parameter.
//
//         --------------------------------------------------
//         params := myMailer.GenParams(
//             smtpServerHost,
//             smtpServerPort,
//             header,
//             body,
//             authConfig,
//             tlsConfig
//         )
//         --------------------------------------------------
//
//     8. Send an email.
//
//         --------------------------------------------------
//         if err := myMailer.Send(params); err != nil {
//             // Error handling
//         }
//         --------------------------------------------------
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
package mailer

import (
    "bytes"
    "crypto/tls"
    "errors"
    "html/template"
    "math/rand"
    "net/smtp"
    "path"
    "strconv"
    "time"
)

const (
    CHARSET_ISO_2022_JP = "iso-2022-jp"
    CHARSET_US_ASCII = "us-ascii"
    CHARSET_UTF8 = "UTF-8"
    CONTENT_TYPE_TEXT_HTML = "text/html"
    CONTENT_TYPE_TEXT_PLAIN = "text/plain"
    CONTENT_TYPE_TEXT_RICHTEXT = "text/richtext"
    CONTENT_TYPE_TEXT_X_WHATEVER = "text/x-whatever"
    MIME_VERSION_1_0 = "1.0"
)

type Params struct {
    AuthConfig *AuthConfig
    Body []*Body
    Header *Header
    SmtpServerHost string
    SmtpServerPort int
    TlsConfig *tls.Config
}

type Header struct {
    From string
    MimeVersion string
    Subject string
    To string
}

type AuthConfig struct {
    Crammd5Auth *CRAMMD5Auth
    PlainAuth *PlainAuth
}

type CRAMMD5Auth struct {
    UserName string
    Secret string
}

type PlainAuth struct {
    UserName string
    Password string
    Host string
}

type Body struct {
    ContentType string
    Charset string
    Data string
}

//////////////////////////////////////////////////////////////////////
// Send Email
//////////////////////////////////////////////////////////////////////
func Send(params *Params) error {
    // Set up headers and message.
    headers := make(map[string]string)
    headers["From"] = params.Header.From
    headers["To"] = params.Header.To
    headers["Subject"] = params.Header.Subject
    headers["MIME-version"] = params.Header.MimeVersion
    body := make([]byte, 0)
    for k,v := range headers {
        body = append(body, k + ": " + v + "\r\n"...)
    }
    var boundary string
    if len(params.Body) > 1 {
        boundary = genBoundary()
        body = append(body, "Content-Type: multipart/alternative; boundary=\"" + boundary + "\"\r\n"...)
    }
    for _, b := range params.Body {
        if len(params.Body) > 1 {
            body = append(body, "--" + boundary + "\r\nContent-Type: " + b.ContentType + "; charset=\"" + b.Charset + "\"\r\n" + b.Data + "\r\n"...)
        } else {
            body = append(body, "Content-Type: " + b.ContentType + "; charset=\"" + b.Charset + "\"\r\n" + b.Data + "\r\n"...)
        }
    }
    if len(params.Body) > 1 {
        body = append(body, "--" + boundary + "--\r\n"...)
    }

    // Connect to the SMTP server
    var c *smtp.Client
    var err error
    if params.TlsConfig != nil {
        conn, err := tls.Dial("tcp", params.SmtpServerHost + ":" + strconv.Itoa(params.SmtpServerPort), params.TlsConfig)
        if err != nil {
            return errors.New("tls.Dial() error. err=" + err.Error())
        }
        c, err = smtp.NewClient(conn, params.SmtpServerHost)
        if err != nil {
            return errors.New("smtp.NewClient() error. err=" + err.Error())
        }
    } else {
        c, err = smtp.Dial(params.SmtpServerHost + ":" + strconv.Itoa(params.SmtpServerPort))
        if err != nil {
            return errors.New("smtp.Dial() error. err=" + err.Error())
        }
    }
    defer c.Close()

    // Authentication
    if params.AuthConfig != nil {
        if params.AuthConfig.Crammd5Auth != nil {
            auth := smtp.CRAMMD5Auth(params.AuthConfig.Crammd5Auth.UserName, params.AuthConfig.Crammd5Auth.Secret)
            if err = c.Auth(auth); err != nil {
                return errors.New("(*Client) Auth() error. err=" + err.Error())
            }
        }
        if params.AuthConfig.PlainAuth != nil {
            auth := smtp.PlainAuth("", params.AuthConfig.PlainAuth.UserName, params.AuthConfig.PlainAuth.Password, params.AuthConfig.PlainAuth.Host)
            if err = c.Auth(auth); err != nil {
                return errors.New("(*Client) Auth() error. err=" + err.Error())
            }
        }
    }

    // Mail commands
    if err = c.Mail(params.Header.From); err != nil {
        return errors.New("(*Client) Mail() error. err=" + err.Error())
    }
    if err = c.Rcpt(params.Header.To); err != nil {
        return errors.New("(*Client) Rcpt() error. err=" + err.Error())
    }
    wc, err := c.Data()
    if err != nil {
        return errors.New("(*Client) Data() error. err=" + err.Error())
    }
    _, err = wc.Write(body)
    if err != nil {
        return errors.New("(io.WriteCloser) Write() error. err=" + err.Error())
    }
    if err = wc.Close(); err != nil {
        return errors.New("(*Client) Quit() error. err=" + err.Error())
    }
    if err = c.Quit(); err != nil {
        return errors.New("(*Client) Quit() error. err=" + err.Error())
    }
    return nil
}


//////////////////////////////////////////////////////////////////////
// Generate Params
// @param smtpServerHost string: SMTP server Host.
// @param smtpServerPort string: SMTP server Port.
// @param header *Header: Mail Header.
// @param body string: Mail message.
// @param authConfig *AuthConfig: Authentication configuration.
// @param tlsConfig *tls.Config: TLS configuration.
//////////////////////////////////////////////////////////////////////
func GenParams(smtpServerHost string, smtpServerPort int, header *Header, body []*Body, authConfig *AuthConfig, tlsConfig *tls.Config) *Params {
    return &Params{
        AuthConfig: authConfig,
        Body: body,
        Header: header,
        SmtpServerHost: smtpServerHost,
        SmtpServerPort: smtpServerPort,
        TlsConfig: tlsConfig,
    }
}


//////////////////////////////////////////////////////////////////////
// Generate CRAMMD5Auth Struct
//////////////////////////////////////////////////////////////////////
func GenCRAMMD5Auth(userName string, secret string) *AuthConfig {
    a := &CRAMMD5Auth{
        UserName: userName,
        Secret: secret,
    }
    return &AuthConfig{
        Crammd5Auth: a,
    }
}


//////////////////////////////////////////////////////////////////////
// Generate PlainAuth Struct
//////////////////////////////////////////////////////////////////////
func GenPlainAuth(userName string, password string, host string) *AuthConfig {
    a := &PlainAuth{
        UserName: userName,
        Password: password,
        Host: host,
    }
    return &AuthConfig{
        PlainAuth: a,
    }
}


//////////////////////////////////////////////////////////////////////
// Generate TLS Configuration Struct
//////////////////////////////////////////////////////////////////////
func GenTlsConfig(serverName string) *tls.Config {
    return &tls.Config{
        ServerName: serverName,
    }
}


//////////////////////////////////////////////////////////////////////
// Set certificate files into the TLS config.
//////////////////////////////////////////////////////////////////////
func SetCertFiles(tlsConfig *tls.Config, certFile string, keyFile string) (*tls.Config, error) {
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        return tlsConfig, err
    }
    tlsConfig.Certificates = []tls.Certificate{cert}
    return tlsConfig, nil
}


//////////////////////////////////////////////////////////////////////
// Set certificate bytes into the TLS config.
//////////////////////////////////////////////////////////////////////
func SetCertBytes(tlsConfig *tls.Config, certPem []byte, keyPem []byte) (*tls.Config, error) {
    cert, err := tls.X509KeyPair(certPem, keyPem)
    if err != nil {
        return tlsConfig, err
    }
    tlsConfig.Certificates = []tls.Certificate{cert}
    return tlsConfig, nil
}


//////////////////////////////////////////////////////////////////////
// Generate Header Struct
//////////////////////////////////////////////////////////////////////
func GenHeader(from string, to string, subject string, mimeVersion string) *Header {
    return &Header{
        From: from,
        MimeVersion: mimeVersion,
        Subject: subject,
        To: to,
    }
}


//////////////////////////////////////////////////////////////////////
// Generate a mail body from files.
//////////////////////////////////////////////////////////////////////
func GenBodyFromFiles(contentType string, charset string, fileNames []string, params map[string]string) (*Body, error) {
    f := template.FuncMap{
        "safeHTML": func(s string) template.HTML { return template.HTML(s) },
    }
    t, err := template.New(path.Base(fileNames[0])).Funcs(f).ParseFiles(fileNames...)
    if err != nil {
        return nil, err
    }
    buffer := new(bytes.Buffer)
    if err := t.Execute(buffer, params); err != nil {
        return nil, err
    }
    body := &Body{
        ContentType: contentType,
        Charset: charset,
        Data: buffer.String(),
    }
    return body, nil
}


//////////////////////////////////////////////////////////////////////
// Generate a mail body from strings.
//////////////////////////////////////////////////////////////////////
func GenBodyFromString(contentType string, charset string, text string, params map[string]string) (*Body, error) {
    f := template.FuncMap{
        "safeHTML": func(s string) template.HTML { return template.HTML(s) },
    }
    t, err := template.New("t").Funcs(f).Parse(text)
    if err != nil {
        return nil, err
    }
    buffer := new(bytes.Buffer)
    if err := t.Execute(buffer, params); err != nil {
        return nil, err
    }
    body := &Body{
        ContentType: contentType,
        Charset: charset,
        Data: buffer.String(),
    }
    return body, nil
}


//////////////////////////////////////////////////////////////////////
// Generate a radom value for boundary.
//////////////////////////////////////////////////////////////////////
func genBoundary() string {
    charset := "1234567890abcdefghijklmnopqrstuvwxyz"
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    b := make([]byte, 32)
    for i := range b {
        b[i] = charset[r.Intn(len(charset))]
    }
    return string(b)
}
