package common

import "time"

const FileChunk = 1 * (1 << 20) // 1 MB
const ImagesDest = "./uploads/"
const LvlDBPath = "./lvldb/"
const TenDays time.Duration = 24 * time.Hour * 10 // 10 days // TESTING: 5 * time.Second
