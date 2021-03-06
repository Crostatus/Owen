package main

import (
    "github.com/dghubble/go-twitter/twitter"
		"github.com/dghubble/oauth1"
		"bufio"
		//"fmt"
		"log"
		"os"
		"time"
		"strconv"
		"strings"
)

var consumerKey string = ""
var consumerSecret string = ""
var accessToken string = ""
var accessSecret string = ""


var max_tweet_length int = 280
var months = [12]string {
	"gen",
	"feb",
	"mar",
	"apr",
	"mag",
	"giu",
	"lug",
	"ago",
	"set",
	"ott",
	"nov",
	"dec",
}

func append_to_last_element(slice []string, line string) []string{
	n := len(slice)
	if n == 0 {
		return slice
	}
	previous_el := slice[n - 1]
	previous_el = previous_el[: len(previous_el) - 1]
	slice[n - 1] = previous_el + " " + line[6:] + "\n"
	return slice
}

func need_to_add(prev_el string, new_el string, today_date string) int{
	new_date := new_el[0:6]
	if strings.Compare(new_date, today_date) == 0 {
		return 1
	}

	month_check := new_date[0:3]
	for i:= 0; i < len(months); i++ {
		if strings.Compare(month_check, months[i]) == 0 {
			return 0
		}
	}
	if len(prev_el) < 6 {
		return 0
	}
	if strings.Compare(today_date, prev_el[:6]) == 0 {
		return 2
	}

	return 0
}

func get_today_string() string {
	t := time.Now()
	formatted_month := months[int(t.Month()) - 1]
	tmp_day := int(t.Day())
	var day string
	if tmp_day < 10 {
		day = "0" + strconv.Itoa(tmp_day)
	} else {
		day = strconv.Itoa(tmp_day)
	}
	return formatted_month + " " + day
}

func extend(slice []string, element string) []string{
    n := len(slice)
    if n == cap(slice) {
        // Slice is full; must grow.
        // We double its size and add 1, so if the size is zero we still grow.
        newSlice := make([]string, len(slice), 2*len(slice)+1)
        copy(newSlice, slice)
        slice = newSlice
    }
    slice = slice[0 : n+1]
    slice[n] = element
		return slice
}

func get_today_events(input_file_name string) []string{
	if input_file_name == "" {
		log.Fatal("Fatal error: missing input file.")
	}
	file, err := os.Open(input_file_name)
	if err != nil {
			log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	today_date := get_today_string()
	var line string
	var result []string = make([]string, 0, 5)
	var is_new int
	var prev_line string = ""
	for scanner.Scan() {
			line = scanner.Text()
			is_new = need_to_add(prev_line, line, today_date)
			prev_line = line
			if is_new == 1 {
				result = extend(result, line + "\n")
			} else if is_new == 2 {
				result = append_to_last_element(result, line)
			}
	}
	if err := scanner.Err(); err != nil {
			log.Fatal(err)
	}
	return result
}

func main(){
	today_events := get_today_events("./events.txt")
	//Twitter client
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)


	var query string = ""
	var post_length int = 0

	for i:= 0; i < len(today_events); i++ {
		post_length += len(today_events[i])
		if post_length < max_tweet_length {
			query = query + today_events[i]
		} else {
			client.Statuses.Update(query, nil)
			time.Sleep(3 * time.Second)
			post_length = len(today_events[i])
			query = today_events[i]
		}
	}
	time.Sleep(3 * time.Second)
	client.Statuses.Update(query, nil)

}
