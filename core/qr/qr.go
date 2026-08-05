package qr

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/zatrano/framework/core/http"
)

// SVG renders a deterministic QR-like matrix SVG for payload (visual stub, not scannable ISO QR).
func SVG(payload string, size ...int) string {
	dim := 21
	if len(size) > 0 && size[0] >= 11 {
		dim = size[0]
		if dim%2 == 0 {
			dim++
		}
	}
	cells := matrix(payload, dim)
	scale := 8
	width := dim * scale
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" shape-rendering="crispEdges">`, width, width, dim, dim))
	b.WriteString(`<rect width="100%" height="100%" fill="#fff"/>`)
	for y := 0; y < dim; y++ {
		for x := 0; x < dim; x++ {
			if cells[y][x] {
				b.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="1" height="1" fill="#111"/>`, x, y))
			}
		}
	}
	b.WriteString(`</svg>`)
	return b.String()
}

// Response returns an image/svg+xml response.
func Response(payload string, size ...int) *http.Response {
	svg := SVG(payload, size...)
	resp := http.Text(svg)
	resp.SetContent([]byte(svg), "image/svg+xml; charset=utf-8")
	return resp
}

func matrix(payload string, dim int) [][]bool {
	h := fnv.New64a()
	_, _ = h.Write([]byte(payload))
	seed := h.Sum64()
	cells := make([][]bool, dim)
	for y := 0; y < dim; y++ {
		cells[y] = make([]bool, dim)
		for x := 0; x < dim; x++ {
			// finder-ish corners
			if inFinder(x, y, dim) {
				cells[y][x] = finderCell(x, y, dim)
				continue
			}
			v := seed ^ uint64(x*131) + uint64(y*17) + uint64(payload[x%len(payloadOr(payload))])
			cells[y][x] = (v%2 == 0)
		}
	}
	return cells
}

func payloadOr(s string) string {
	if s == "" {
		return " "
	}
	return s
}

func inFinder(x, y, dim int) bool {
	return inCorner(x, y, 0, 0) || inCorner(x, y, dim-7, 0) || inCorner(x, y, 0, dim-7)
}

func inCorner(x, y, ox, oy int) bool {
	return x >= ox && x < ox+7 && y >= oy && y < oy+7
}

func finderCell(x, y, dim int) bool {
	corners := [][2]int{{0, 0}, {dim - 7, 0}, {0, dim - 7}}
	for _, c := range corners {
		ox, oy := c[0], c[1]
		if x < ox || x >= ox+7 || y < oy || y >= oy+7 {
			continue
		}
		lx, ly := x-ox, y-oy
		if lx == 0 || ly == 0 || lx == 6 || ly == 6 {
			return true
		}
		if lx >= 2 && lx <= 4 && ly >= 2 && ly <= 4 {
			return true
		}
		return false
	}
	return false
}
