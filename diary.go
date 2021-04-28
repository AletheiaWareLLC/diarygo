package diarygo

import (
	"aletheiaware.com/bcgo"
	"aletheiaware.com/spaceclientgo"
	"aletheiaware.com/spacego"
	"encoding/base64"
	"io"
	"sort"
	"strings"
)

const (
	DiaryName = "Diary"
	DiaryMime = spacego.MIME_TYPE_TEXT_PLAIN
)

type MetaFilter struct{}

func (f *MetaFilter) Filter(m *spacego.Meta) bool {
	return m.Name == DiaryName && m.Type == DiaryMime
}

type Diary interface {
	Add(bcgo.Node, io.Reader) (string, error)
	Clear()
	Length() int
	FindID(string) string
	ID(int) string
	Meta(string) *spacego.Meta
	Timestamp(string) uint64
	Refresh(bcgo.Node) error
}

type diary struct {
	client     spaceclientgo.SpaceClient
	ids        []string
	metas      map[string]*spacego.Meta
	timestamps map[string]uint64
}

func NewDiary(c spaceclientgo.SpaceClient) Diary {
	return &diary{
		client:     c,
		metas:      make(map[string]*spacego.Meta),
		timestamps: make(map[string]uint64),
	}
}

func (d *diary) Add(n bcgo.Node, r io.Reader) (string, error) {
	ref, err := d.client.Add(n, nil, DiaryName, DiaryMime, r)
	if err != nil {
		return "", err
	}
	id := base64.RawURLEncoding.EncodeToString(ref.RecordHash)
	if err := d.Refresh(n); err != nil {
		return "", err
	}
	return id, nil
}

func (d *diary) Clear() {
	for k := range d.metas {
		delete(d.metas, k)
	}
	d.ids = nil
}

func (d *diary) Length() int {
	return len(d.ids)
}

func (d *diary) FindID(prefix string) string {
	for _, i := range d.ids {
		if strings.HasPrefix(i, prefix) {
			return i
		}
	}
	return ""
}

func (d *diary) ID(index int) string {
	return d.ids[index]
}

func (d *diary) Meta(id string) *spacego.Meta {
	return d.metas[id]
}

func (d *diary) Timestamp(id string) uint64 {
	return d.timestamps[id]
}

func (d *diary) Refresh(n bcgo.Node) error {
	if err := d.client.SearchMeta(n, &MetaFilter{}, func(e *bcgo.BlockEntry, m *spacego.Meta) error {
		id := base64.RawURLEncoding.EncodeToString(e.RecordHash)
		if _, ok := d.metas[id]; !ok {
			d.metas[id] = m
			d.timestamps[id] = e.Record.Timestamp
			d.ids = append(d.ids, id)
		}
		return nil
	}); err != nil {
		return err
	}
	sort.Slice(d.ids, func(i, j int) bool {
		return d.timestamps[d.ids[i]] < d.timestamps[d.ids[j]]
	})
	return nil
}
