package services
import (
	"encoding/json"
	"fmt"
	"strconv"

)

type j struct {
	Cl []string `json:"cl"`
	Gr []string `json:"gr"`
	Cr []string `json:"cr"`
}

//				*************************************************************************
// 				*****************       IGNORE THIS CLASS *******************************
//				***********     ITS JUST FOR ME TO TRY OUT RANDOM THINGS *********************
// 				********************************************************************************




func DoStuff() {

	var data []j

	for i := 0; i < 1; i++ {
		v := strconv.Itoa(i)
		data = append(data, j{
			Cl: []string{"foo " + v},
			Gr: []string{"bar " + v},
			Cr: []string{"foobar " + v},
		})
	}

	b, _ := json.MarshalIndent(data, "", " ")
	fmt.Println(string(b))



	//_ = ioutil.WriteFile("file.json", b, 0644)


	//f, err := os.OpenFile(os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/config/file.json", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	//if err != nil {
	//     log.Fatal(err)
	//}

	//f.Write(b)



}