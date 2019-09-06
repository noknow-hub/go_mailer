//////////////////////////////////////////////////////////////////////
// mailer.go
//
// @usage
// 
//     1. Import this package.
//
//     2. Prepare following variables.
//
//         - SMTP Server Host: string
//         - SMTP Server Post: int
//         - Mail From: string
//         - Mail To: string
//         - Mail Subject: string
//         - Mail Body: string
//         - (Optional) Mail Auth User Name: string
//         - (Optional) Mail Auth Password: string
//
//     3. Generate following structs.
//
//         - Header Struct: GenHeader() function.
//         - (Optional) AuthConfig Struct: GenCRAMMD5Auth() or GenPlainAuth() function.
//         - (Optional) tls.Config Struct: GenTlsConfig() function.
//
//     4. Generate Params struct.
//
//         Here is an example code.
//
//             import myMailer "mailer"
//
//             smtpServerHost := "noknow.info"
//             smtpServerPort := 465
//             from := "info@noknow.info"
//             to := "xxx@example.com"
//             subject := "This is a subject"
//             body := "This is a message."
//             authUserName := "noknow_auth"
//             authPassword := "noknow_auth_pass"
//             authHost := "noknow.info"
//             header myMailer.GenHeader(from, to, subject, myMailer.MIME_VERSION_1_0, false, myMailer.CHARSET_UTF8)
//             authConfig := myMailer.GenPlainAuth(authUserName, authPassword, authHost)
//             tlsConfig := myMailer.GenTlsConfig(smtpServerHost)
//             params := GenParams(smtpServerHost, smtpServerPort, header, body, authConfig, tlsConfig)
//
//     5. Send an email.
//
//         Here is an example code.
//
//             if err := myMailer.Send(params); err != nil {
//                 // Error handling
//             }
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
    "crypto/tls"
    "net/smtp"
    "strconv"
)

const (
    CHARSET_UTF8 = "UTF-8"
    MIME_VERSION_1_0 = "1.0"
)

type Params struct {
    AuthConfig *AuthConfig
    Body string
    Header *Header
    SmtpServerHost string
    SmtpServerPort int
    TlsConfig *tls.Config
}

type Header struct {
    From string
    MimeVersion string
    ContentType string
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



//////////////////////////////////////////////////////////////////////
// Send Email
//////////////////////////////////////////////////////////////////////
func Send(params *Params) error {
    // Set up headers and message
    headers := make(map[string]string)
    headers["From"] = params.Header.From
    headers["To"] = params.Header.To
    headers["Subject"] = params.Header.Subject
    headers["MIME-version"] = params.Header.MimeVersion + "\r\n" + params.Header.ContentType
    body := make([]byte, 0)
    for k,v := range headers {
        body = append(body, k + ": " + v + "\r\n"...)
    }
    body = append(body, "\r\n" + params.Body...)

    // Connect to the SMTP server
    var c *smtp.Client
    var err error
    if params.TlsConfig != nil {
        conn, err := tls.Dial("tcp", params.SmtpServerHost + ":" + strconv.Itoa(params.SmtpServerPort), params.TlsConfig)
        if err != nil {
            return err
        }
        c, err = smtp.NewClient(conn, params.SmtpServerHost)
        if err != nil {
            return err
        }
    } else {
        c, err = smtp.Dial(params.SmtpServerHost + ":" + strconv.Itoa(params.SmtpServerPort))
        if err != nil {
            return err
        }
    }

    defer c.Close()

    // Authentication
    if params.AuthConfig != nil {
        if params.AuthConfig.Crammd5Auth != nil {
            auth := smtp.CRAMMD5Auth(params.AuthConfig.Crammd5Auth.UserName, params.AuthConfig.Crammd5Auth.Secret)
            if err = c.Auth(auth); err != nil {
                return err
            }
        }
        if params.AuthConfig.PlainAuth != nil {
            auth := smtp.PlainAuth("", params.AuthConfig.PlainAuth.UserName, params.AuthConfig.PlainAuth.Password, params.AuthConfig.PlainAuth.Host)
            if err = c.Auth(auth); err != nil {
                return err
            }
        }
    }

    // Mail commands
    if err = c.Mail(params.Header.From); err != nil {
        return err
    }
    if err = c.Rcpt(params.Header.To); err != nil {
        return err
    }
    wc, err := c.Data()
    if err != nil {
        return err
    }
    _, err = wc.Write(body)
    if err != nil {
        return err
    }
    if err = wc.Close(); err != nil {
        return err
    }
    if err = c.Quit(); err != nil {
        return err
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
func GenParams(smtpServerHost string, smtpServerPort int, header *Header, body string, authConfig *AuthConfig, tlsConfig *tls.Config) *Params {
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
// Generate Header Struct
//////////////////////////////////////////////////////////////////////
func GenHeader(from string, to string, subject string, mimeVersion string, isHtml bool, charset string) *Header {
    return &Header{
        From: from,
        MimeVersion: mimeVersion,
        ContentType: genContentType(isHtml, charset),
        Subject: subject,
        To: to,
    }
}


//////////////////////////////////////////////////////////////////////
// Generate Content Type
//////////////////////////////////////////////////////////////////////
func genContentType(isHtml bool, charset string) string {
    if isHtml {
        return "Content-Type: text/html; charset=\"" + charset + "\";"
    } else {
        return "Content-Type: text/plain; charset=\"" + charset + "\";"
    }
}

