package messager

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func SendMessage(title string, desc string) {
	client := &http.Client{}
	message := fmt.Sprintf("title=%s&desp=%s", title, desc)
	//fmt.Println(message)
	var data = strings.NewReader(message)
	req, err := http.NewRequest("POST", "https://sctapi.ftqq.com/SCT92876TTsoGsYW3TiGLZpvNEI6CvAyi.send", data)
	log.Println("Send message", title)
	if err != nil {
		log.Fatal(err)
	}
	client.Do(req)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}
