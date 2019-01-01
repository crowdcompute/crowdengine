// Copyright 2018 The crowdcompute:crowdengine Authors
// This file is part of the crowdcompute:crowdengine library.
//
// The crowdcompute:crowdengine library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The crowdcompute:crowdengine library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the crowdcompute:crowdengine library. If not, see <http://www.gnu.org/licenses/>.

package manager

import (
	"encoding/binary"
	"errors"
	"unicode"
)

// DockerLog represents a log
type DockerLog struct {
	Type      string `json:"type"`
	Data      string `json:"data"`
	Timestamp string `json:"timestamp"`
}

// WhitespaceAt finds the index in which a whitespace is found
func WhitespaceAt(buf []byte) (int, error) {
	for j, v := range buf {
		if unicode.IsSpace(rune(v)) {
			return j, nil
		}
	}
	return -1, errors.New("No whitespace found")
}

// DockerLogDecoder parses the log produced by Docker
func DockerLogDecoder(buf []byte) ([]DockerLog, error) {
	var (
		logs       []DockerLog
		i          = 0
		streamType = ""
	)

	if len(buf) < 8 {
		return nil, errors.New("No logs available")
	}

	for {

		if i == len(buf) {
			break
		}

		header := buf[i : i+8]
		payloadLength := int(binary.BigEndian.Uint32(header[4:]))
		payload := buf[i+8 : i+8+payloadLength]

		// extract timestamp
		pos, err := WhitespaceAt(payload)
		if err != nil {
			return nil, errors.New("Unable to extract timestamp from logs")
		}
		timestamp := payload[0:pos]
		payload = payload[pos+1:]

		switch header[0] {
		case 0:
			streamType = "STDIN"
		case 1:
			streamType = "STDOUT"
		default:
			streamType = "STDERR"
		}

		logs = append(logs, DockerLog{streamType, string(payload), string(timestamp)})
		i = i + 8 + payloadLength
	}
	return logs, nil
}
