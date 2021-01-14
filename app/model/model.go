package model

import "time"

type File struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Dir       string    `json:"dir"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	ModTime   string    `json:"mod_time"`
	Md        string    `json:"md"`
	Level     uint8     `json:"level"`
	Ext       string    `json:"ext"`
	Status    uint8     `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type FileManager struct {
	Thumb string `json:"thumb"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Path  string `json:"path"`
	Index string `json:"index"`
}
