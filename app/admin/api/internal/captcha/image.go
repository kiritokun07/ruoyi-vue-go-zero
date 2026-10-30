// Package captcha renders a challenge. Its answer is stored separately in Redis.
package captcha

import (
	"bytes"
	"encoding/base64"
	"image/color"
	"image/jpeg"
	"image/png"

	"github.com/google/uuid"
	bc "github.com/mojocn/base64Captcha"
)

func NewDriver(kind string) bc.Driver {
	background := &color.RGBA{R: 255, G: 255, B: 255, A: 255}
	fonts := []string{"wqy-microhei.ttc"}
	if kind == "char" {
		return bc.NewDriverString(60, 180, 0, 0, 4, "abcdefghjkmnpqrstuvwxyz23456789", background, nil, fonts)
	}
	return bc.NewDriverMath(60, 180, 0, 0, background, nil, fonts)
}

func Render(driver bc.Driver) (id, encoded, answer string, err error) {
	_, question, solution := driver.GenerateIdQuestionAnswer()
	item, err := driver.DrawCaptcha(question)
	if err != nil {
		return "", "", "", err
	}
	var image bytes.Buffer
	if _, err = item.WriteTo(&image); err != nil {
		return "", "", "", err
	}
	pixels, err := png.Decode(&image)
	if err != nil {
		return "", "", "", err
	}
	image.Reset()
	if err = jpeg.Encode(&image, pixels, nil); err != nil {
		return "", "", "", err
	}
	return uuid.NewString(), base64.StdEncoding.EncodeToString(image.Bytes()), solution, nil
}
