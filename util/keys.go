package util

// 全部文件、全部文件索引、重复文件、重复文件索引

const (
	PicKey                = "pic_files"
	VideoKey              = "video_files"
	DupSuccessKey         = "duplicate_success_files" // 全部文件已完成key
	DupKey                = "duplicate_files"         // 全部文件key
	DupIndexKey           = "duplicate_index_files"   // 全部文件索引key
	RepeatKey             = "repeat_files"            // 重复文件key
	RepeatIndexKey        = "repeat_index_files"      // 重复文件索引
	InitialID      uint64 = 100000001
	StatusBegin           = 1
	StatusRepeat          = 2
)

var (
	PicExt   = []string{".jpg", ".jpeg", ".png", ".bmp", ".gif", ".webp"}
	VideoExt = []string{".mp4", ".wmv", ".avi", ".rm", ".mpeg", ".flv", ".3gp", ".mov"}
)

func GetPicIndexKey(id string) string {
	return PicKey + ":" + id
}

func GetVideoIndexKey(id string) string {
	return VideoKey + ":" + id
}

func GetDupSucKey(dir string) string {
	return DupSuccessKey + ":" + dir
}

func GetDupIndexKey(id string) string {
	return DupIndexKey + ":" + id
}

func GetDupKey(md, dir string) string {
	return DupKey + ":" + md + ":" + dir
}

func GetDupPreKey(md string) string {
	return DupKey + ":" + md
}

func GetRepeatIndexKey(id string) string {
	return RepeatIndexKey + ":" + id
}

func GetRepeatKey(md, dir string) string {
	return RepeatKey + ":" + md + ":" + dir
}

func GetRepeatPreKey(md string) string {
	return RepeatKey + ":" + md
}
