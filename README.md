# jhttp
jhttp, which is a http client tool may have good experience like postman

## Example code
**http get**

~~~go
package main

import (
	"fmt"
	"io/ioutil"

	http "github.com/jiuhuche120/jhttp"
)

func main() {
	client := http.NewClient(
		http.AddHeader("Accept", "application/vnd.github.v3+json"),
	)
	resp, err := client.Get("https://api.github.com/repos/jiuhuche120/jhttp", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(result))
}
~~~

It also supports carrying **url parameters** and **body** in get requests

~~~golang
package main

import (
	"fmt"
	"io/ioutil"

	http "github.com/jiuhuche120/jhttp"
)

func main() {
	client := http.NewClient(
		http.AddHeader("Accept", "application/vnd.github.v3+json"),
	)
	opts := []http.ParamsOption{
		http.AddParams("page", "1"),
		http.AddParams("per_page", "30"),
	}
	resp, err := client.Get("https://api.github.com/repos/jiuhuche120/jhttp/tags", nil, opts...)
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(result))
}
~~~

**http post**

~~~go
package main

import (
	"fmt"
	"io/ioutil"

	"github.com/tidwall/gjson"

	http "github.com/jiuhuche120/jhttp"
)

func main() {
	client := http.NewClient(
		http.AddHeader("Accept", "application/vnd.github.v3+json"),
		http.AddHeader("Authorization", "token XXX"),
	)
	value := gjson.Parse("{\"new_name\":\"main\"}").Value()
	resp, err := client.Post("https://api.github.com/repos/jiuhuche120/jhttp/branches/master/rename", value)
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(result))
}
~~~

It also supports post request sending **form**

~~~go
package main

import (
	"fmt"
	"io/ioutil"

	http "github.com/jiuhuche120/jhttp"
)

func main() {
	client := http.NewClient()
	form := http.NewForm(
		http.AddFormParams("key", "value", http.Text),
		http.AddFormParams("file", "file path", http.File),
	)
	resp, err := client.PostForm("https://api.github.com/repos/jiuhuche120/jhttp/branches/master/rename",form)
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(result))
}
~~~