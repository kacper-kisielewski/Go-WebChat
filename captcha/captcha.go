package captcha

import (
	"bytes"
	"net/http"
	"time"

	"github.com/dchest/captcha"
)

//NewCaptcha returns new captcha id
func NewCaptcha() string {
	return captcha.New()
}

//VerifyCaptcha verifies captcha
func VerifyCaptcha(captchaID, answer string) bool {
	return captcha.VerifyString(captchaID, answer)
}

//ShowCaptchaImage serves captcha over HTTP
func ShowCaptchaImage(w http.ResponseWriter, r *http.Request, id string) {
	var content bytes.Buffer

	w.Header().Set("Content-Type", "image/png")

	captcha.WriteImage(&content, id, captcha.StdWidth, captcha.StdHeight)
	http.ServeContent(w, r, id+".png", time.Now(), bytes.NewReader(content.Bytes()))
}
