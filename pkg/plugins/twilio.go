package plugins

import (
	"github.com/labstack/echo/v4"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

type TwilioPlugin string
type ConfigTwilio struct {
	EnableTwilio    bool
	AccountSid      string
	AuthToken       string
	VerifyServiceID string
}

var (
	fxsTwilio    map[string]interface{} = make(map[string]interface{})
	configTwilio ConfigTwilio
	client       *twilio.RestClient
)

func (d TwilioPlugin) Run(c echo.Context,
	vars map[string]string, payload_in interface{}, dromadery_data string,
	callback chan string,
) (payload_out interface{}, next string, err error) {
	return nil, "output_1", nil
}
func (d TwilioPlugin) AddFeatureJS() map[string]interface{} {
	return fxsTwilio
}
func (d TwilioPlugin) Name() string {
	return "twilio"
}
func (d TwilioPlugin) Initialize(enable bool, accound_sid, auth_token, service_sid string) {
	configTwilio = ConfigTwilio{
		AccountSid:      accound_sid,
		AuthToken:       auth_token,
		VerifyServiceID: service_sid,
		EnableTwilio:    enable,
	}
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: configTwilio.AccountSid,
		Password: configTwilio.AuthToken,
	})
	fxsTwilio["send_otp"] = send_otp
	fxsTwilio["check_otp"] = check_otp
}
func send_otp(to string) bool {
	if !configTwilio.EnableTwilio {
		return true
	}
	params := &openapi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")
	_, err := client.VerifyV2.CreateVerification(configTwilio.VerifyServiceID, params)
	status := err == nil
	return status
}
func check_otp(to string, code string) bool {
	if !configTwilio.EnableTwilio {
		return true
	}
	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)
	resp, err := client.VerifyV2.CreateVerificationCheck(configTwilio.VerifyServiceID, params)
	status := err == nil && (*resp.Status == "approved")
	return status
}
