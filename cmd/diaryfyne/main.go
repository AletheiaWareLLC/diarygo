package main

import (
	"aletheiaware.com/bcgo"
	"aletheiaware.com/diarygo"
	"aletheiaware.com/spaceclientgo"
	"aletheiaware.com/spacefynego"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Create Fyne App
	a := app.NewWithID("com.aletheiaware.diary")
	a.SetIcon(resourceIconPng)
	a.Settings().SetTheme(&diaryTheme{})

	// Create Fyne Window
	w := a.NewWindow("Diary")

	// Create Space Client
	c := spaceclientgo.NewSpaceClient()

	// Create Space Fyne
	f := spacefynego.NewSpaceFyne(a, w, c)

	// Create Diary
	d := diarygo.NewDiary(c)

	// Create Unlock Button
	u := widget.NewButton("Unlock", func() {
		// Trigger Access Flow
		go f.Account(c)
	})

	i := CreateIcon()
	m := container.NewMax(i, container.NewCenter(u))

	// Create Document List
	l := &widget.List{
		CreateItem: func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
	}

	// Populate list upon sign in
	f.AddOnSignedIn(func(bcgo.Account) {
		go func() {
			n, err := f.Node(c)
			if err != nil {
				f.ShowError(err)
				return
			}

			// Refresh Diary Entries
			if err := d.Refresh(n); err != nil {
				f.ShowError(err)
				return
			}

			// Attach list
			l.Length = d.Length
			l.UpdateItem = func(item widget.ListItemID, obj fyne.CanvasObject) {
				obj.(*widget.Label).SetText(bcgo.TimestampToString(d.Timestamp(d.ID(item))))
			}
			l.OnSelected = func(item widget.ListItemID) {
				ShowEditor(f, c, n, d, d.ID(item))
				l.Unselect(item) // TODO FIXME HACK
			}
			l.Refresh()

			// Update window content to include toolbar, and list
			w.SetContent(
				container.NewBorder(
					CreateToolbar(f, c, n, d, l),
					nil,
					nil,
					nil,
					l))
		}()
	})

	// Clear list upon sign in
	f.AddOnSignedOut(func() {
		go func() {
			d.Clear()
			l.Refresh()
			w.SetContent(m)
		}()
	})

	w.SetContent(m)
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(480, 640))
	w.ShowAndRun()
}
