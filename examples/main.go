package main

import (
	"encoding/json"
	"fmt"
	"nullabye"
)

type OptionalApiThing struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}
type ApiRequest struct {
	OptionalThing nullabye.OptionalStruct `json:"optional_thing"`
}

type optionalInnerThing struct {
	Id     int
	Name   string
	Amount float64
}
type innerRequest struct {
	useThing bool
	thing    optionalInnerThing
}

func main() {
	datas := [][]byte{
		[]byte(`{"optional_thing": {"id": 1, "name": "test", "amount": 9.99}}`),
		[]byte(`{"optional_thing": null}`),
		[]byte(`{}`),
		[]byte(`{"optional_thing": {"id": 1.99, "name": 7, "amount": true}}`), // Will error
	}

	for _, data := range datas {
		var request ApiRequest
		err := json.Unmarshal(data, &request)
		panicOnError(err)

		var optionalRequestThing OptionalApiThing
		_, err = request.OptionalThing.Get(&optionalRequestThing)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("%+v\n", innerRequest{
			useThing: request.OptionalThing.IsSet(),
			thing: optionalInnerThing{
				Id:     optionalRequestThing.Id,
				Name:   optionalRequestThing.Name,
				Amount: optionalRequestThing.Amount,
			},
		})
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
