package ui

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"sort"
)

type channel struct {
	layout *tview.Grid

	name string
	chat *tview.TextView
	users *tview.List
	input *tview.TextView
	info *tview.TextView
}

type View struct {
	pages *tview.Pages
	app *tview.Application
	activeChannel int
	channels []channel
}

func NewView() (*View) {
	view := &View{
		pages: tview.NewPages(),
	}

	return view
}

func (v *View) Start() {
	v.app = tview.NewApplication()

	v.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		total_active := v.pages.GetPageCount()

		if event.Key() == tcell.KeyLeft {
			v.activeChannel = (v.activeChannel - 1 + total_active) % total_active
			page := v.channels[v.activeChannel]
			v.pages.SwitchToPage(page.name)
		}

		if event.Key() == tcell.KeyRight {
			v.activeChannel = (v.activeChannel + 1) % total_active
			page := v.channels[v.activeChannel]
			v.pages.SwitchToPage(page.name)
		}

		return event
	})

	if err := v.app.SetRoot(v.pages, true).SetFocus(v.pages).Run(); err != nil {
		panic(err)
	}
}

func (v *View) AddChannel(name string, user_list []string) {
	new_chan := channel{
		layout: tview.NewGrid().SetRows(1, 0, 1).SetColumns(20, 0, 20).SetBorders(true),
		name: name,
		chat: newTextView("text here"),
		users: newListView(),
		input: newTextView("input text here"),
		info: newTextView(name),
	}

	for id, user := range user_list {
		new_chan.users.AddItem(user, user, rune(id), nil)
	}

	// Layout
	new_chan.layout.AddItem(new_chan.users, 0, 2, 3, 1, 0, 0, false)
	new_chan.layout.AddItem(new_chan.chat,  1, 0, 1, 2, 0, 0, false)
	new_chan.layout.AddItem(new_chan.input, 2, 0, 1, 2, 0, 0, false)
	new_chan.layout.AddItem(new_chan.info,  0, 0, 1, 2, 0, 0, false)

	v.pages.AddPage(name, new_chan.layout, true, true)

	// FIXME: Might not be the best idea to keep this counter
	v.activeChannel = v.pages.GetPageCount()

	v.channels = append(v.channels, new_chan)
}

func (v *View) getChannel(name string) *channel {
	result := sort.Search(len(v.channels), func(i int) bool {
		return v.channels[i].name == name
	})

	if result < len(v.channels) {
		return &v.channels[result]
	}

	return nil
}

func (v *View) AddBufferMsg(channel, from, msg string) {
	c := v.getChannel(channel)

	if c != nil {
		line := fmt.Sprintf("<%s> %s\n", from, msg)
		c.chat.Write([]byte(line))
	}
}

func newListView() *tview.List {
	return tview.NewList().ShowSecondaryText(false)
}

func newTextView(text string) *tview.TextView {
	return tview.NewTextView().SetText(text)
}