package mailsender

import (
    "fmt"
    "log"
    "net/smtp"
    "strings"

    model "github.com/AishwaryaRK/tweats/datamodel"
    matcher "github.com/AishwaryaRK/tweats/matcher"
)

const senderMail = "*******@gmail.com"
const senderCredential = "******"
const smtpHost = "smtp.gmail.com"
const smtpURI = "smtp.gmail.com:587"
const twitterMailSuffix = "@twitter.com"

const mailMsgFormat = `From: %s
To: %s
Subject: %s

%s`
const mailSubject = "[Tweats] Lunch Buddy"
const mailBodyFormat = `Hello %s,

Thank you for signing up to SG TwEATS. We found your lunch buddy! Before you start: Please make sure you’ve read and understood the guidelines and best practices at go/SgTwEATS.

Your common interests are: %s.

Hope you have a fun & enjoyable time together! Have questions? Need advice? Contact the admin via jiliannew@twitter.com.

Regards,
TwEATS team
`

// Send sends emails for one Match
func Send(match matcher.Match) {
    send(match.MatchedTweeps, match.MatchedInterest)
}

// Send sends emails to tweeps
func send(tweeps []model.Tweep, interest string) {
    receiverArr := make([]string, len(tweeps))
    receiverNames := genReceiverNames(tweeps)
    for index, tweep := range tweeps {
        receiverArr[index] = tweep.LDAP + twitterMailSuffix
    }
    receiverStr := strings.Join(receiverArr, ",")

    mailBody := fmt.Sprintf(mailBodyFormat, receiverNames, interest)
    msg := fmt.Sprintf(mailMsgFormat, senderMail, receiverStr, mailSubject, mailBody)
    auth := smtp.PlainAuth("", senderMail, senderCredential, smtpHost)
    err := smtp.SendMail(smtpURI, auth, senderMail, receiverArr, []byte(msg))

    if err != nil {
        log.Printf("smtp error: %s", err)
        return
    }

    log.Printf("Mail sent to: %s", receiverStr)
}

func genReceiverNames(tweeps []model.Tweep) (nameLine string) {
    nameLine = ""
    for index, tweep := range tweeps {
        if (index != 0) {
            nameLine += ", "
        }
        nameLine += tweep.Name
    }
    return
}
