package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	ipAPI = "http://httpbin.org/ip"
	cfAPI = "https://api.cloudflare.com/client/v4/zones"
)

type IPResponse struct {
	Origin string `json:"origin"`
}

type DNSRecord struct {
	ID      string `json:"id,omitempty"` // 添加 ID 字段
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl,omitempty"` // Add ttl

}

type RecordResponse struct {
	Result  []DNSRecord `json:"result"`
	Success bool        `json:"success"`
	Errors  []string    `json:"errors"`
}

func main() {
	email, apiKey, zoneID, domain, err := readCredentials("ddns.txt")
	if err != nil {
		fmt.Println("Error reading credentials:", err)
		return
	}
	for {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("%s DNS record  !\n", currentTime)
		updateDNSRecord(email, apiKey, zoneID, domain)
		time.Sleep(1 * time.Minute)
	}
}

func readCredentials(filePath string) (email, apiKey, zoneID, domain string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) < 4 {
		return "", "", "", "", fmt.Errorf("ddns.txt does not contain enough lines")
	}

	return strings.TrimSpace(lines[0]), strings.TrimSpace(lines[1]), strings.TrimSpace(lines[2]), strings.TrimSpace(lines[3]), nil
}

func updateDNSRecord(email string, apiKey string, zoneID string, domain string) {
	ipAddress := getPublicIP()
	fmt.Printf("Domain: %s, IP: %s\n", domain, ipAddress)

	client := &http.Client{}
	headers := map[string]string{
		"X-Auth-Email": email,
		"X-Auth-Key":   apiKey,
		"Content-Type": "application/json",
	}

	// Get existing DNS records
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/dns_records?type=A&name=%s", cfAPI, zoneID, domain), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	for key, val := range headers {
		req.Header.Add(key, val)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Check and update/create DNS record
	recordResponse := RecordResponse{}
	json.Unmarshal(body, &recordResponse)
	if recordResponse.Success && len(recordResponse.Result) > 0 {
		// Update record
		recordID := recordResponse.Result[0].ID
		updateData := DNSRecord{
			Type:    "A",
			Name:    domain,
			Content: ipAddress,
			TTL:     60,
		}
		data, err := json.Marshal(updateData)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s/dns_records/%s", cfAPI, zoneID, recordID), bytes.NewBuffer(data))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		for key, val := range headers {
			req.Header.Add(key, val)
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		json.Unmarshal(body, &recordResponse)
		if recordResponse.Success {
			fmt.Println("DNS record updated successfully!")
		} else {
			fmt.Println("Failed to update DNS record!")
		}
	} else {
		// Create record...
		// Similar to the update logic, but using POST method
		// Create record
		createData := DNSRecord{
			Type:    "A",
			Name:    domain,
			Content: ipAddress,
			TTL:     60,
		}
		data, err := json.Marshal(createData)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/dns_records", cfAPI, zoneID), bytes.NewBuffer(data))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		for key, val := range headers {
			req.Header.Add(key, val)
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		json.Unmarshal(body, &recordResponse)
		if recordResponse.Success {
			fmt.Println("DNS record created successfully!")
		} else {
			fmt.Println("Failed to create DNS record!")
		}
	}
}

func getPublicIP() string {
	resp, err := http.Get(ipAPI)
	if err != nil {
		fmt.Println("Error getting public IP:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return ""
	}

	ipResponse := IPResponse{}
	json.Unmarshal(body, &ipResponse)
	return ipResponse.Origin
}
