/*
   Calculator
   Copyright (C) 2020  aarontillekeratne

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

/*
  Author: aarontillekeratne
  Contact: github.com/codeBehindMe
*/

package multiplier

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

func RPCMultiplyFloat(multiplierServiceAddress *string, a, b float32) (float32, error) {
	conn, err := grpc.Dial(*multiplierServiceAddress, grpc.WithInsecure())
	if err != nil {
		return 0.0, err
	}

	client := NewMultiplierClient(conn)

	res, err := client.MultiplyFloat(context.Background(), &MultiplyFloatMessage{
		A: &a,
		B: &b,
	})
	if err != nil {
		return 0.0, err
	}
	return res.GetR(), nil
}

type MultiplyFloatsOperand struct {
	A float32 `json:"A"`
	B float32 `json:"B"`
}

func GetHandler(multiplierServiceAddress *string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("recieved multiply float request")
		decoder := json.NewDecoder(r.Body)

		var operands MultiplyFloatsOperand

		err := decoder.Decode(&operands)
		if err != nil {
			log.Printf("error occured when decoding message: %V", err)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "Could not unpack request body: %V", err)

			return
		}

		log.Printf("recieved operands A=%V , B=%V", operands.A, operands.B)

		res, err := RPCMultiplyFloat(multiplierServiceAddress, operands.A, operands.B)
		if err != nil {
			log.Printf("error when remote calling multiply float: %V", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "Could not multiply values: %V", err)

			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, res)
	}
}
