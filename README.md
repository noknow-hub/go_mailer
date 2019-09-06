# go_mailer
SMTP / SMTPS Mailer Library  

## Usage  
### 1. Import this package.  
### 2. Prepare following variables.  
- SMTP Server Host: string  
- SMTP Server Post: int  
- Mail From: string  
- Mail To: string  
- Mail Subject: string  
- Mail Body: string  
- (Optional) Mail Auth User Name: string  
- (Optional) Mail Auth Password: string  
### 3. Generate following structs.  
- Header Struct: GenHeader() function.
- (Optional) AuthConfig Struct: GenCRAMMD5Auth() or GenPlainAuth() function.  
- (Optional) tls.Config Struct: GenTlsConfig() function.  
### 4. Generate Params struct.  
Here is an example code.  
```  
import myMailer "mailer"

smtpServerHost := "noknow.info"
smtpServerPort := 465
from := "info@noknow.info"
to := "xxx@example.com"
subject := "This is a subject"
body := "This is a message."
authUserName := "noknow_auth"
authPassword := "noknow_auth_pass"
authHost := "noknow.info"
header myMailer.GenHeader(from, to, subject, myMailer.MIME_VERSION_1_0, false, myMailer.CHARSET_UTF8)
authConfig := myMailer.GenPlainAuth(authUserName, authPassword, authHost)
tlsConfig := myMailer.GenTlsConfig(smtpServerHost)
params := GenParams(smtpServerHost, smtpServerPort, header, body, authConfig, tlsConfig)
```  
### 5. Send an email.  
Here is an example code.  
```  
if err := myMailer.Send(params); err != nil {
    // Error handling
}
```  
