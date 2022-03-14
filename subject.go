package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func parseSubject(message *Message) error {

	// TODO: much better parsing to better account for all the very different subjects formats used for file posts...
	pattern1 := regexp.MustCompile(`^(?P<reminder>.+)(?:[\[\(] *(?P<segmentNo>\d+) */ *(?P<totalSegments>\d+) *[\)\]])`)
	pattern2 := regexp.MustCompile(`^(?P<header>.*)?(?:[\[\(] *(?P<segmentNo>\d+) */ *(?P<totalSegments>\d+) *[\)\]])(?P<reminder>.*)?`)
	pattern3 := regexp.MustCompile(`^(?P<header>.*?)?(?: *"(?P<filename>(?P<basefilename>[^"]+?)(?:\.[^\."]+){0,2})")`)
	pattern4 := regexp.MustCompile(`^(?P<filename>(?P<basefilename>.+?)(?:\.[^\.]+){0,2})$`)

	if matches := findNamedMatches(pattern1, message.subject); matches != nil {
		message.segmentNo, _ = strconv.Atoi(matches["segmentNo"])
		message.totalSegments, _ = strconv.Atoi(matches["totalSegments"])
		reminder := matches["reminder"]
		if matches := findNamedMatches(pattern2, reminder); matches != nil {
			message.fileNo, _ = strconv.Atoi(matches["segmentNo"])
			message.totalFiles, _ = strconv.Atoi(matches["totalSegments"])
			message.header = strings.Trim(matches["header"], " -")
			reminder = matches["reminder"]
		}
		if matches := findNamedMatches(pattern3, reminder); matches != nil {
			message.header = message.header + " " + strings.Trim(matches["header"], " -")
			message.filename = matches["filename"]
			message.basefilename = matches["basefilename"]
		} else if matches := findNamedMatches(pattern4, reminder); matches != nil {
			message.filename = strings.TrimSpace(matches["filename"])
			message.basefilename = strings.TrimSpace(matches["basefilename"])
		}
		if message.header == "" {
			message.header = message.basefilename
		}
		if message.header != "" {
			message.headerHash = GetMD5Hash(message.header + message.from + strconv.Itoa(message.totalFiles))
		} else {
			return errors.New("no header found")
		}
		if message.filename != "" {
			message.fileHash = GetMD5Hash(message.headerHash + message.filename + strconv.Itoa(message.totalSegments))
		} else {
			return errors.New("no filename found")
		}
	} else {
		return errors.New("subject did not match")
	}

	return nil

}

func findNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)
	if match == nil {
		return nil
	}
	results := map[string]string{}
	for i, name := range match {
		results[regex.SubexpNames()[i]] = name
	}
	return results
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
