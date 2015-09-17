package main

import (
	"errors"
	"hash/fnv"
	"net/url"
	"strconv"
	"time"

	"github.com/codemodus/parth"
)

func (n *node) getReferer(r string) (string, error) {
	ref, err := url.Parse(r)
	if err != nil || ref == nil {
		return "", errors.New("error parsing referer: " + err.Error())
	}
	return ref.String(), nil
}

func (n *node) getIndexSegment(s string) (string, error) {
	si := 0
	if n.su.conf.FormPathPrefix != "" {
		si = 1
	}
	seg, err := parth.SegmentToString(s, si)
	if err != nil {
		return "", err
	}
	return seg, nil
}

func (n *node) getKey() string {
	t := []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
	h := fnv.New64a()
	h.Write(t)
	s := strconv.FormatUint(h.Sum64(), 10)
	return s
}

func (n *node) getConfirmHash() string {
	t := []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
	h := fnv.New64a()
	h.Write(t)
	s := strconv.FormatUint(h.Sum64(), 10) +
		"_" + strconv.FormatInt(time.Now().Unix(), 10)
	return s
}
