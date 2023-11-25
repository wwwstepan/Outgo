package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

func dtToInt(d time.Time, ind int32) int32 {
	return (ind << 24) | (int32(d.Year()) << 10) | (int32(d.Month()) << 5) | int32(d.Day())
}

func outgoToBin(num_days int) []byte {
	a := make([]byte, 0, num_days*3*8)
	dt := time.Now().AddDate(0, 0, -num_days)
	for nd := 0; nd < num_days; nd++ {

		n_og := 3 + rand.Intn(6)

		for nog := 0; nog < n_og; nog++ {
			ind := 1 + rand.Intn(10)
			sum := rand.Intn(100) + rand.Intn(ind*5)
			if ind == 2 || ind == 6 || ind == 9 {
				sum *= 10
			}
			if rand.Intn(100) < 20 {
				sum += rand.Intn(2000)
			}
			dtc := dtToInt(dt, int32(ind))
			appendIntToByteArray(dtc, &a)
			appendIntToByteArray(int32(sum), &a)
		}
		dt = dt.AddDate(0, 0, 1)
	}
	return a
}

func appendIntToByteArray(i int32, a *[]byte) {
	*a = append(*a, byte(i&0xff))
	*a = append(*a, byte((i>>8)&0xff))
	*a = append(*a, byte((i>>16)&0xff))
	*a = append(*a, byte((i>>24)&0xff))
}

func main() {

	num_days := 100 * 365

	a := outgoToBin(num_days)
	err := ioutil.WriteFile("outgo.dat", a, 0600)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("OK")
	}
}
