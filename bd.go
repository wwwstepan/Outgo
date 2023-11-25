package main

import (
	"io/ioutil"
	"mylogs"
	"time"
)

func dtToInt(d time.Time, ind int32) int32 {
	return (ind << 24) | (int32(d.Year()) << 10) | (int32(d.Month()) << 5) | int32(d.Day())
}

func intToDt(i int32) (time.Time, rune) {
	ind := i >> 24
	var r rune
	for k, v := range ogClassifierInd {
		if v == ind {
			r = k
			break
		}
	}
	y := (i & 0xffffff) >> 10
	md := i & 1023
	m := md >> 5
	d := md & 31
	dt := time.Date(int(y), time.Month(m), int(d), 0, 0, 0, 0, time.UTC)
	return dt, r
}

func outgoToBin() []byte {
	a := make([]byte, 0, len(outgo)*8)
	for _, o := range outgo {
		if o.sum != 0 {
			ind := ogClassifierInd[o.code]
			dtc := dtToInt(o.dt, ind)
			appendIntToByteArray(dtc, &a)
			appendIntToByteArray(o.sum, &a)
		}
	}
	return a
}

func appendIntToByteArray(i int32, a *[]byte) {
	*a = append(*a, byte(i&0xff))
	*a = append(*a, byte((i>>8)&0xff))
	*a = append(*a, byte((i>>16)&0xff))
	*a = append(*a, byte((i>>24)&0xff))
}

func saveOutgo() error {
	t := getTimer()
	a := outgoToBin()
	err := ioutil.WriteFile("outgo.dat", a, 0600)
	dur := t()
	mylogs.Log("Запись БД. %d записей, %v", len(outgo), dur)
	return err
}

func loadOutgo() error {
	t := getTimer()
	stat.clear()
	a, err := ioutil.ReadFile("outgo.dat")
	if err != nil {
		return err
	}
	outgoFromBin(a)
	dur := t()
	mylogs.Log("Загрузка БД. %d записей, %v", len(outgo), dur)
	return nil
}

func outgoFromBin(a []byte) {
	veryEarlyDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	veryFutureDate := time.Now().AddDate(10, 0, 0)

	n := len(a) >> 3
	n8 := n << 3
	for i := 0; i < n8; i += 8 {
		dtc := extractInt(a[i : i+4])
		sum := extractInt(a[i+4 : i+8])

		dt, icod := intToDt(dtc)
		if dt.Before(veryEarlyDate) || dt.After(veryFutureDate) || // suspicious date
			getOutgo(icod) == "" || // unknown code
			sum <= 0 || sum > 9_999_999 { // suspicious summ
			continue
		}
		o := outgoItem{sum: sum, code: rune(icod), dt: dt}
		outgo = append(outgo, o)
		if isCurDate(dt) {
			outgoToday = append(outgoToday, o.copy())
		}
		stat.addSum(dt, sum)
	}
}

func extractInt(q []byte) int32 {
	i := int32(q[0]) | (int32(q[1]) << 8) | (int32(q[2]) << 16) | (int32(q[3]) << 24)
	return i
}
