package qr_test

import (
	"strings"
	"testing"

	"github.com/zatrano/framework/core/qr"
)

func TestQRSVG(t *testing.T) {
	svg := qr.SVG("https://zatrano.test")
	if !strings.Contains(svg, "<svg") || !strings.Contains(svg, "</svg>") {
		t.Fatal(svg)
	}
	if !strings.Contains(svg, `fill="#111"`) {
		t.Fatal("missing modules")
	}
}
