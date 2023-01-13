/*
	Copyright 2022 ankangan

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	response, err := http.DefaultClient.Get("http://127.0.0.1:8000/a?echo=hello")
	if err != nil {
		log.Fatalln(err)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("client go data:%s", data)
}
