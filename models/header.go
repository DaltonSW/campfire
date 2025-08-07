package models

import (
	"fmt"

	"github.com/dustin/go-humanize"
)

// Header gets the above-viewport content. Title and file stats
func (m model) Header() string {
	cContent := titleStyle.Render("Campfire")

	rContent := ""
	if m.fileExists {
		filesize := humanize.Bytes(uint64(m.prevFileInfo.Size()))

		rContent = fmt.Sprintf(
			"%v %v",
			fileNameStyle.Render(m.filename),
			fmt.Sprintf("(Size: %v)", filesize),
		)
	} else {
		rContent = "File not found..."
	}
	rContent = statsStyle.Render(rContent)

	lContent := statsStyle.Italic(true).Render("https://github.com/daltonsw/campfire")

	return align(m.width, lContent, cContent, rContent)
}
