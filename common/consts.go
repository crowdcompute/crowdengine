package common

import "time"

// FileChunk is the size of a chunk when uploading a file
const FileChunk = 1 * (1 << 20)

// ImagesDest is the destination folder for storing images
const ImagesDest = "./uploads/"

// TenDays represents 10 days in time
const TenDays time.Duration = 24 * time.Hour * 10
