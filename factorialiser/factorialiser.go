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

package factorialiser

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"net/http"
)

type FactorialFloatOperand struct {
	V float32
}

func RPCFactorialiseFloat(factorialiserServiceAddress *string, a *float32) (float32, error) {
	conn, err := grpc.Dial(*factorialiserServiceAddress, grpc.WithInsecure())
	if err != nil {
		return 0.0, err
	}

	client := NewFactorialiserClient(conn)

	res, err := client.FactorialFloat(context.Background(), &FactorialiseFloatMessage{
		A: a,
	})
	if err != nil {
		return 0.0, err
	}
	return res.GetR(), nil
}

func GetHandler(factorialiserServiceAddress *string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		operand := FactorialFloatOperand{}
		err := decoder.Decode(&operand)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "could not unpack request body: %V", err)

			return
		}

		res, err := RPCFactorialiseFloat(factorialiserServiceAddress, &operand.V)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "could not calculate factorial: %V", err)

			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, res)
	}
}
