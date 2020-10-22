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

   You should have received A copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

/*
  Author: aarontillekeratne
  Contact: github.com/codeBehindMe
*/

package adder

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"net/http"
)

type AddFloatsOperand struct {
	A, B float32
}

func RPCAddFloats(adderServiceAddress *string, a, b *float32) (float32, error) {
	conn, err := grpc.Dial(*adderServiceAddress, grpc.WithInsecure())
	if err != nil {
		return 0.0, err
	}
	client := NewAdderClient(conn)
	res, err := client.AddFloat(context.Background(), &AddFloatOperands{
		A: a,
		B: b,
	})
	if err != nil {
		return 0.0, err
	}
	return res.GetR(), nil
}

func GetHandler(adderServiceAddress *string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		operands := &AddFloatsOperand{}
		err := decoder.Decode(&operands)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "could not unpack request body: %V", err)

			return
		}

		res, err := RPCAddFloats(adderServiceAddress, &operands.A, &operands.B)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "error when getting result: %V", err)

			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, res)
	}
}
