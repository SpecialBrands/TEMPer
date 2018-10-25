// Copyright 2018, Special Brands Holding BV
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"time"

	"github.com/zserge/hid"
)

var (
	reqtemp = []byte{0x01, 0x80, 0x33, 0x01, 0x00, 0x00, 0x00, 0x00}

	timeout = 500 * time.Millisecond
)

func main() {
	hid.UsbWalk(func(dev hid.Device) {
		readbuf := make([]byte, 8)
		info := dev.Info()

		if info.Vendor != 0x0c45 && info.Product != 0x7401 {
			fmt.Println("Device not found")
			return
		}

		err := dev.Open()
		defer dev.Close()

		if err != nil {
			fmt.Println("Error[Open]: ", err)
		}

		_, err = dev.Write(reqtemp, timeout)
		if err != nil {
			fmt.Println("Error[Request 1]: ", err)
			return
		}

		_, err = dev.Write(reqtemp, timeout)
		if err != nil {
			fmt.Println("Error[Request 2]: ", err)
			return
		}

		readbuf, err = dev.Read(8, timeout*10)
		if err != nil {
			fmt.Println("Error[Read]: ", err)
			return
		}

		for itt := 0; itt < len(readbuf); itt++ {
			if readbuf[itt+0] == 0x80 && readbuf[itt+1] == 0x02 {
				Temperature := result(readbuf[itt+2], readbuf[itt+3])
				fmt.Printf("Temperature is %f\n", Temperature)
				return
			}
		}
	})
}

func result(HH byte, LL byte) (xtemp float32) {
	xtemp = float32((256*int16(HH))+int16(HH)) * (125.5 / 32000.0)
	return
}
