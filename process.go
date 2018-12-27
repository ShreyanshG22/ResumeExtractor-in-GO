package main

import ("fmt"
	"bufio"
	"log"
	"os"
	"encoding/json"
	"strings"
	"math"
	"sync"
	"strconv")

var mutex = &sync.Mutex{}
var keyword[] string
var po float64
var docrel map[string]int
var final_result map[string]float64
var wg sync.WaitGroup

var kjl int

func jsonconv(s string)map[string]float64 {
	var data map[string]float64
	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		panic(err)
	}
	return data
}

func search(document string) {
	var final_count, value float64 = 0.0, 1.0
	document = strings.TrimSuffix(document, "\n")
	result := strings.SplitN(document, " ", 2)
	resume_score := jsonconv(result[1])
	for _, j := range keyword {
		if strings.Contains(j, "*") {
			j := strings.Replace(j, " ", "", -1)
			i := strings.Split(j, "*")
			if v, present := resume_score[i[0]];present {
				value = v
			} else {
				value = 1.0
			}
			weight,_ := strconv.ParseFloat(i[1], 64)
			final_count += value*(po - math.Log10(float64(docrel[i[0]])))*(weight)
		} else {
			if v, present := resume_score[j];present{
				value = v
			} else {
				value = 1.0
			}
			//fmt.Println(value*(po - math.Log10(docrel[j])))
			final_count += value*(po - math.Log10(float64(docrel[j])))		
		}
	}
	//fmt.Println(final_count)
	mutex.Lock()
	final_result[result[0]] = final_count
	mutex.Unlock()
	wg.Done()
}

func reldoc(key string) {
	var tagcount int
	tagcount = 0
//	fmt.Println(key)
	file, err := os.Open("/home/shreyanshg/Desktop/WorkOnThis.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var boolean bool
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		boolean = strings.Contains(scanner.Text(), key)
		//////////check if & contains////////////////
		if boolean {
			tagcount++
		}
	}
	fmt.Println(tagcount)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	mutex.Lock()
	docrel[key] = tagcount
	mutex.Unlock()
	wg.Done()	
}

func main() {
	kjl = 0
	reader := bufio.NewReader(os.Stdin)
	var name string
	fmt.Println("What is the keyword string?")
	name, _ = reader.ReadString('\n')
	keyword = strings.Split(name, "/")

	file1, err1 := os.Open("/home/shreyanshg/Desktop/Dictionary.txt")
	if err1 != nil {
		log.Fatal(err1)
	}
	defer file1.Close()

	reader1 := bufio.NewReader(file1)
	var l1 []byte
	l1, _, err1 = reader1.ReadLine()
	var yy float64
	yy, _ = strconv.ParseFloat(string(l1), 64)
	po = math.Log10(yy)
	l1, _, err1 = reader1.ReadLine()
	//dic := string(l1)
	if err1 != nil {
		fmt.Printf(" > Failed!: %v\n", err1)
	}
	docrel = make(map[string]int)
	for _, key := range keyword {
		wg.Add(1)
		go reldoc(key)
	}
	wg.Wait()

	final_result = make(map[string]float64)
	file, err := os.Open("/home/shreyanshg/Desktop/WorkOnThis.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var no_of_doc int = 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		no_of_doc++
		wg.Add(1)
		doc_line := strconv.Itoa(no_of_doc) + " " + scanner.Text()
		go search(doc_line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	wg.Wait()
	//fmt.Println(final_result)
}
