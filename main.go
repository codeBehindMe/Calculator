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

package main

import (
	"Calculator/adder"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

var adderServiceAddress *string
var port *int

func BaseHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "You have reached Calculator")
}

type AddFloatOperands struct {
	a, b float32
}

func RPCAddFloats(a, b *float32) (float32, error) {
	conn, err := grpc.Dial(*adderServiceAddress, grpc.WithInsecure())
	if err != nil {
		return 0.0, err
	}
	client := adder.NewAdderClient(conn)
	res, err := client.AddFloat(context.Background(), &adder.AddFloatOperands{
		A: a,
		B: b,
	})
	if err != nil {
		return 0.0, err
	}
	return res.GetR(), nil
}

func AddFloatHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	operands := AddFloatOperands{}
	err := decoder.Decode(&operands)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "could not unpack request body: %v", err)
	}

	res, err := RPCAddFloats(&operands.a, &operands.b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "error when getting result: %v", err)
	}

	_, _ = fmt.Fprint(w, res)
}

func main() {

	adderServiceAddress = flag.String("adder", "localhost:3000", "DNS/IP of Adder service including port")
	port = flag.Int("port", 80, "Port of service")

	router := mux.NewRouter()

	router.HandleFunc("/", BaseHandler)
	router.HandleFunc("/float/add", AddFloatHandler)

	log.Printf("Starting server on port %v", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), router)
	if err != nil {
		log.Fatal("Could not start http server")
	}

}
