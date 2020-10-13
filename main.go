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
	"Calculator/factorialiser"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var adderServiceAddress *string
var factorialiserServiceAddress *string
var port *int

func BaseHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "You have reached Calculator")
}

type FactorialFloatOperand struct {
	v float32
}

func FactorialFloatHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	operand := FactorialFloatOperand{}
	err := decoder.Decode(&operand)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "could not unpack request body: %v", err)
	}

	res, err := factorialiser.RPCFactorialiseFloat(factorialiserServiceAddress, &operand.v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "could not calculate factorial: %v", err)
	}
	_, _ = fmt.Fprint(w, res)
}

type AddFloatsOperand struct {
	a, b float32
}

func AddFloatHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	operands := &AddFloatsOperand{}
	err := decoder.Decode(&operands)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "could not unpack request body: %v", err)
	}

	res, err := adder.RPCAddFloats(adderServiceAddress, &operands.a, &operands.b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "error when getting result: %v", err)
	}

	_, _ = fmt.Fprint(w, res)
}

//validateFlags checks to see approriate variables are passed in to the flag
// containers. Flag containers must be declared as program variables and
// accessible in the global scope.
func validateFlags(){

}

func main() {

	adderServiceAddress = flag.String("adder", "", "DNS/IP of Adder service including port")
	factorialiserServiceAddress = flag.String("factorialiser", "", "DNS/IP of Factorialiser service")
	port = flag.Int("port", 80, "Port of service")

	router := mux.NewRouter()

	router.HandleFunc("/", BaseHandler)
	router.HandleFunc("/float/add", AddFloatHandler)
	router.HandleFunc("/float/factorial", FactorialFloatHandler)

	log.Printf("Starting server on port %v", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), router)
	if err != nil {
		log.Fatal("Could not start http server")
	}

}
