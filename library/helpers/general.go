package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jasonbronson/kwik-cms-engine/config"

	"github.com/google/uuid"
	"github.com/jasonbronson/kwik-cms-engine/request/response"
	"github.com/mailgun/mailgun-go/v4"
)

const ()

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func ConvertUnixTimeToTimeTime(t string) (*time.Time, error) {
	var ct time.Time
	i, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		ct = time.Now()
		return &ct, err
	}
	ct = time.Unix(i, 0)
	return &ct, nil
}

func CastRequestBody(r *http.Request, w http.ResponseWriter, obj interface{}) (err error) {
	bodyBytes := new(bytes.Buffer)
	bodyBytes.ReadFrom(r.Body)
	body := bodyBytes.String()
	if err = json.Unmarshal([]byte(body), &obj); err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err)
		// EncodeErrorTo(w, http.StatusInternalServerError, fmt.Errorf("unable to marshal JSON: %v", err), body, r.Header.Get("X-Odd-User-Agent"))
	}
	return
}

func SendMail(subject, recipient string, emailData map[string]interface{}, resetEmailImageLocation string, context context.Context) error {
	mailgunDomain := "mg.odd.io"
	mailgunAPIKey := "key-******"
	mailgun := mailgun.NewMailgun(mailgunDomain, mailgunAPIKey)
	if mailgun == nil {
		return errors.New("SendMail: Mailgun reference is nil")
	}
	sender := fmt.Sprintf("%s <support@%s>", "Email Subject", mailgunDomain)

	message := mailgun.NewMessage(sender, subject, "", recipient)
	content, err := ioutil.ReadFile("/app/static/templates/email.html")
	if err != nil {
		return errors.New("SendMail: parse content failed")
	}

	emailMessage := string(content)

	t := template.Must(template.New("email").Parse(emailMessage))
	builder := &strings.Builder{}
	if err := t.Execute(builder, emailData); err != nil {
		panic(err)
	}
	emailString := builder.String()

	message.SetHtml(emailString)
	message.AddInline(resetEmailImageLocation)
	_, _, err = mailgun.Send(context, message)
	return err
}

func SendContactEmail(context context.Context, name string, fromEmail string, subject string, content string) error {
	mg := mailgun.NewMailgun(config.Cfg.MailgunDomain, config.Cfg.MailgunSecret)
	if mg == nil {
		return errors.New("SendMail: Mailgun reference is nil")
	}
	msg := mg.NewMessage(fromEmail, subject, fmt.Sprintf("from:%s\n%s", name, content), config.Cfg.MailgunContactMeRecipient)
	_, _, err := mg.Send(context, msg)
	return err
}

type ipRange struct {
	start net.IP
	end   net.IP
}

func inRange(r ipRange, ipAddress net.IP) bool {
	if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) < 0 {
		return true
	}
	return false
}

var privateRanges = []ipRange{
	{
		start: net.ParseIP("10.0.0.0"),
		end:   net.ParseIP("10.255.255.255"),
	},
	{
		start: net.ParseIP("100.64.0.0"),
		end:   net.ParseIP("100.127.255.255"),
	},
	{
		start: net.ParseIP("172.16.0.0"),
		end:   net.ParseIP("172.31.255.255"),
	},
	{
		start: net.ParseIP("192.0.0.0"),
		end:   net.ParseIP("192.0.0.255"),
	},
	{
		start: net.ParseIP("192.168.0.0"),
		end:   net.ParseIP("192.168.255.255"),
	},
	{
		start: net.ParseIP("198.18.0.0"),
		end:   net.ParseIP("198.19.255.255"),
	},
}

func isPrivateSubnet(ipAddress net.IP) bool {
	if ipCheck := ipAddress.To4(); ipCheck != nil {
		for _, r := range privateRanges {
			if inRange(r, ipAddress) {
				return true
			}
		}
	}
	return false
}

func ReplaceString(str, old, new string) string {
	return strings.Replace(str, old, new, -1)
}

func GetUserIP(r *http.Request) string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || isPrivateSubnet(realIP) {
				continue
			}
			return ip
		}
	}
	var ip string
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		ip = forwarded
	} else {
		ip = r.RemoteAddr
	}
	return ip
}

type DefaultParameters struct {
	PageSize    int
	PageOffset  int
	SortOrder   string
	FilterBy    []string
	FilterValue []string
	FilterQuery string
	ResultTotal int
	Total       int
}
