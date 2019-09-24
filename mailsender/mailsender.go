package mailsender

import (
    "fmt"
    "log"
    "net/smtp"
    "strings"

    m "github.com/AishwaryaRK/tweats/datamodel"
)

const senderMail = "*******@gmail.com"
const senderCredential = "******"
const smtpHost = "smtp.gmail.com"
const smtpURI = "smtp.gmail.com:587"

const mailMsgFormat = `From: %s
To: %s
Subject: %s

%s`
const mailSubject = "[Tweats] Lunch Buddy"
const mailBodyFormat = `Hi Tweeps,

#Some greetings here#

regards,
Tweats Team`

// Send sends emails to tweeps
func Send(tweeps []m.Tweep) {
    receiverArr := make([]string, len(tweeps))
    for index, tweep := range tweeps {
        receiverArr[index] = tweep.LDAP
    }

    receiverStr := strings.Join(receiverArr, ",")

    msg := fmt.Sprintf(mailMsgFormat, senderMail, receiverStr, mailSubject, "TODO: replace body")
    auth := smtp.PlainAuth("", senderMail, senderCredential, smtpHost)
    err := smtp.SendMail(smtpURI, auth, senderMail, receiverArr, []byte(msg))

    if err != nil {
        log.Printf("smtp error: %s", err)
        return
    }

    log.Printf("Mail sent to: %s", receiverStr)
}
