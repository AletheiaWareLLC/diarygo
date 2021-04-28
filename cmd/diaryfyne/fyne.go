//go:generate fyne bundle -o bundled.go icon.svg
//go:generate fyne bundle -append -o bundled.go Icon.png
package main

import (
	"aletheiaware.com/bcfynego/ui"
	bcuidata "aletheiaware.com/bcfynego/ui/data"
	"aletheiaware.com/bcgo"
	"aletheiaware.com/diarygo"
	"aletheiaware.com/spaceclientgo"
	"aletheiaware.com/spacefynego"
	"aletheiaware.com/spacefynego/ui/data"
	"context"
	"encoding/base64"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
)

func CreateIcon() fyne.CanvasObject {
	// Renders incorrectly: https://github.com/srwiley/oksvg/issues/27
	return &canvas.Image{
		Resource: resourceIconSvg,
		FillMode: canvas.ImageFillContain,
	}
}

// CreateToolbar returns a toolbar with add, refresh, storage, account, and help options
func CreateToolbar(f spacefynego.SpaceFyne, c spaceclientgo.SpaceClient, n bcgo.Node, d diarygo.Diary, l *widget.List) *widget.Toolbar {
	return widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			go func() {
				index := d.Length()
				_, err := d.Add(n, nil)
				if err != nil {
					f.ShowError(err)
					return
				}
				if err := d.Refresh(n); err != nil {
					f.ShowError(err)
					return
				}
				l.Select(index)
			}()
		}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			go func() {
				if err := d.Refresh(n); err != nil {
					f.ShowError(err)
					return
				}
				l.Refresh()
			}()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.NewThemedResource(data.StorageIcon), func() {
			go f.ShowStorage(c)
		}),
		widget.NewToolbarAction(theme.NewThemedResource(bcuidata.AccountIcon), func() {
			go f.ShowAccount(c)
		}),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			go f.ShowHelp(c)
		}),
	)
}

// ShowEditor displays a dialog containing the content of file with the given ID.
// Modifying the content and tapping "Save" will write the update to the file.
func ShowEditor(f spacefynego.SpaceFyne, c spaceclientgo.SpaceClient, n bcgo.Node, d diarygo.Diary, id string) {
	hash, err := base64.RawURLEncoding.DecodeString(id)
	if err != nil {
		f.ShowError(err)
		return
	}

	e := &widget.Entry{
		PlaceHolder: "Dear Diary...",
		MultiLine:   true,
		Wrapping:    fyne.TextWrapWord,
	}

	ctx, cancel := context.WithCancel(context.Background())

	c.WatchFile(ctx, n, hash, func() {
		reader, err := c.ReadFile(n, hash)
		if err != nil {
			f.ShowError(err)
			return
		}
		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			f.ShowError(err)
			return
		}
		e.SetText(string(bytes))
	})

	dialog := dialog.NewCustomConfirm(bcgo.TimestampToString(d.Timestamp(id)), "Save", "Cancel", e, func(result bool) {
		if result {
			writer, err := c.WriteFile(n, nil, hash)
			if err != nil {
				f.ShowError(err)
				return
			}
			data := []byte(e.Text)
			count, err := writer.Write(data)
			if err != nil {
				f.ShowError(err)
				return
			}
			if length := len(data); count != length {
				f.ShowError(fmt.Errorf("Failed to write all data: Size %d, Wrote %d", length, count))
				return
			}
			if err := writer.Close(); err != nil {
				f.ShowError(err)
				return
			}
		}
	}, f.Window())
	dialog.SetOnClosed(cancel)
	dialog.Show()
	dialog.Resize(ui.DialogSize)
}
