package mailsender

import (
    "fmt"
    "log"
    "net/smtp"
    "strings"
    "errors"

    model "github.com/AishwaryaRK/tweats/datamodel"
    matcher "github.com/AishwaryaRK/tweats/matcher"
    "github.com/AishwaryaRK/tweats/config"
    "github.com/AishwaryaRK/tweats/question"
)

const senderMail = config.SMTP_SENDER_MAIL
const senderCredential = config.SMTP_SENDER_CRENTIAL
const smtpHost = "smtp.gmail.com"
const smtpURI = "smtp.gmail.com:587"
const twitterMailSuffix = "@twitter.com"

const mailMsgFormat = `From: %s
To: %s
Subject: %s
MIME-version: 1.0
Content-Type: text/html

%s`
const mailSubjectFormat = "[TwEATS] Lunch Buddy (%s)"
const mailBodyFormat = `<html><body>
<p>Hello %s,</p>

<p>Thank you for signing up to SG TwEATS. We found your lunch buddy!</p>
<p><b>Before you start:</b> Please make sure you’ve read and understood the guidelines and best practices at go/SgTwEATS.</p>

<div>
<hr>
<p>Your common interests are: <b>%s</b>.</p>
<p>You may answer the following question using "Reply all" to get people know you better:</p>
<p><i>%s</i></p>
</div>

<div>
<hr>
<p>Please organise a time that is suitable for all of you. Here are each of your preferred timings:</p>
%s
</div>

<hr>
<p>Hope you have a fun & enjoyable time together! Have questions? Need advice? Contact the admin via jiliannew@twitter.com.</p>

<div>
<img src="https://media2.giphy.com/media/BlVnrxJgTGsUw/giphy.gif?cid=6104955ed38e7da9420a9cfac5e9bd68879b978685b50a49&rid=giphy.gif">
</div>
<p>Regards,</p>
<p>TwEATS team</p>
</body></html>`
const timeSlotFormat = "%s %d:00 - %d:00\n"

// Send sends emails for one Match
func Send(match matcher.Match) {
    send(match.MatchedTweeps, match.MatchedInterest)
}

// Send sends emails to tweeps
func send(tweeps []model.Tweep, interest string) {
    receiverArr := make([]string, len(tweeps))
    receiverNames := genReceiverNames(tweeps)
    mailSubject := fmt.Sprintf(mailSubjectFormat, receiverNames)
    preferredTimings, err := genPreferredTimings(tweeps)
    if err != nil {
        log.Printf("genPreferredTimings error: %s", err)
        return
    }

    for index, tweep := range tweeps {
        receiverArr[index] = tweep.LDAP + twitterMailSuffix
    }
    receiverStr := strings.Join(receiverArr, ",")

    mailBody := fmt.Sprintf(mailBodyFormat, receiverNames, interest, question.GenRandomQuestion(interest), preferredTimings)
    msg := fmt.Sprintf(mailMsgFormat, senderMail, receiverStr, mailSubject, mailBody)
    auth := smtp.PlainAuth("", senderMail, senderCredential, smtpHost)
    err = smtp.SendMail(smtpURI, auth, senderMail, receiverArr, []byte(msg))

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

func genPreferredTimings(tweeps []model.Tweep) (result string, err error) {
    result = ""
    var timingOneTweep string
    for _, tweep := range tweeps {
        timingOneTweep, err = genPreferredTimingForOneTweep(tweep)
        if err != nil {
            return
        }
        result += timingOneTweep
    }
    return
}

func genPreferredTimingForOneTweep(tweep model.Tweep) (result string, err error) {
    result = "<div><p><b>" + tweep.Name + "<b></p>"
    var timeSlotStr string
    for _, ava := range tweep.Availabilities {
        timeSlotStr, err = genTimeSlotString(ava.Weekday, ava.TimeSlots)
        if err != nil {
            return
        }
        result += timeSlotStr
    }
    result += "</div>"
    return
}

func genTimeSlotString(weekDay int, timeSlots []model.TimeSlot) (result string, err error) {
    var weekStr string
    weekStr, err = convertWeekdayIntToString(weekDay)

    if err != nil {
        return
    }

    result = ""
    for _, timeSlot := range timeSlots {
        result = result + "<p>" + fmt.Sprintf(timeSlotFormat, weekStr, timeSlot.Start, timeSlot.End) + "</p>"
    }
    return
}

func convertWeekdayIntToString(num int) (weekdayStr string, err error) {
    switch num {
        case 1:
            weekdayStr = "Monday"
        case 2:
            weekdayStr = "Tuesday"
        case 3:
            weekdayStr = "Wednesday"
        case 4:
            weekdayStr = "Thursday"
        case 5:
            weekdayStr = "Friday"
        default:
            err = errors.New("Invalid weekday.")
    }
    return
}
