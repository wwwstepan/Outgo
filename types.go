package main

import "time"

type statusInf struct {
	dt             time.Time
	outgoDay       int32
	outgoMonth     int32
	outgoPrevMonth int32
	outgoYear      int32
}

type outgoItem struct {
	sum  int32
	code rune
	dt   time.Time
}

type outgoStr struct {
	code     rune
	sumToday int32
	sumMonth int32
	name     string
	dt       time.Time
}

type outgoType struct {
	code rune
	name string
	ind  int32
}

type outgoTypeJSON struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Code        int32  `json:"code"`
}

type outgoClassifierJSON struct {
	Items []outgoTypeJSON `json:"outgoclassifier"`
}

//////////////////////////////////////////////////////////
// for sort []outgoStr
//////////////////////////////////////////////////////////

type byCodeOutgoStr []outgoStr

func (o byCodeOutgoStr) Len() int {
	return len(o)
}
func (o byCodeOutgoStr) Less(i, j int) bool {
	return o[i].code < o[j].code
}
func (o byCodeOutgoStr) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

//////////////////////////////////////////////////////////
