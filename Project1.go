package main

import (
  "crypto/tls"   
  "flag"         
  "fmt"
  "github.com/jackdanger/collectlinks"
  "net/http"
  "strconv"
  "os"
  "log"
  "github.com/PuerkitoBio/goquery"
)

func main() {
  flag.Parse()
  price := ""
  address := ""
  beds := ""
  piece := "http://www.daft.ie"
  var i int = 20 // counter for offset of webpage
  var x string = "http://www.daft.ie/ireland/property-for-sale/?offset=" // Website offset to scrape
  tmp := ""	// Temp to concatenate i to the end of x
  count := 0	// Count for how many links the webpage has
  intcount := 0  // Counter for how many ints there are in a row
  skip := 3	// Counter to get every third link 
  file, fileErr := os.Create("output.txt")
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}
  
  tlsConfig := &tls.Config{                
                 InsecureSkipVerify: true, 
               }                           
  transport := &http.Transport{   
    TLSClientConfig: tlsConfig,    
  }                                

  client := http.Client{Transport: transport}  
                                               

											   
  resp, err := client.Get("http://www.daft.ie/ireland/property-for-sale/")  
  if err != nil {                   
    return                          
  }
  defer resp.Body.Close()
  
  links := collectlinks.All(resp.Body)

  for _, link := range(links) {
	for _,r := range link {
		c := string(r)
		if _, err := strconv.Atoi(c); err == nil {
			intcount += 1
		} else {
			intcount = 0
		}
		if intcount == 7{
			if skip % 3 == 0 {
				link = piece + link
				doc, err := goquery.NewDocument(link)
					if err != nil {
						log.Fatal(err)
					}
					address, _ = (doc.Find("meta[property='twitter:title']").Attr("content"))
					beds, _ = (doc.Find("meta[property='twitter:data2']").Attr("content"))
					price, _ = (doc.Find("meta[property='twitter:data1']").Attr("content"))
					fmt.Fprintf(file,"%s|%s|%s|%s\n",address,beds,price,link)
			}
			skip++
		}
	}
	count += 1
  }
  
  for count > 64{
	 tmp = strconv.Itoa(i)
	 tmp = x + tmp
	 resp, err := client.Get(tmp)  
	  if err != nil {                   
		return                          
	  }
	  defer resp.Body.Close()
	  count = 0
	  links := collectlinks.All(resp.Body)
	  for _, link := range(links) {
		for _,r := range link {
			c := string(r)
			if _, err := strconv.Atoi(c); err == nil {
				intcount += 1
			} else {
				intcount = 0
			}
			if intcount == 7{
				if skip % 3 == 0 {
					link = piece + link
					doc, err := goquery.NewDocument(link)
					if err != nil {
						log.Fatal(err)
					}
					address, _ = (doc.Find("meta[property='twitter:title']").Attr("content"))
					beds, _ = (doc.Find("meta[property='twitter:data2']").Attr("content"))
					price, _ = (doc.Find("meta[property='twitter:data1']").Attr("content"))
					fmt.Fprintf(file,"%s|%s|%s|%s\n",address,beds,price,link)
				}
				skip++
			}
		}
		count += 1
	  }
	  i += 20
  }
}