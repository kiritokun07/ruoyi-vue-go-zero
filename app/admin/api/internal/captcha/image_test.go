package captcha

import (
	"encoding/base64"
	"image/jpeg"
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	for _, kind := range []string{"math", "char"} {
		t.Run(kind, func(t *testing.T) {
			id, img, answer, err := Render(NewDriver(kind))
			if err != nil || id == "" || answer == "" {
				t.Fatalf("challenge generation failed: %v", err)
			}
			if _, err = jpeg.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(img))); err != nil {
				t.Fatal(err)
			}
		})
	}
}
