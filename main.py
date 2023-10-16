import requests
import time
import json

# Cloudflare 账户信息
your_email = ''
your_api_key = ''
your_zone_id = ''
your_domain = ''


# API URL
url_get_ip = "http://httpbin.org/ip"
url_get_dns_record = f"https://api.cloudflare.com/client/v4/zones/{your_zone_id}/dns_records?type=A&name={your_domain}"
url_create_dns_record = f"https://api.cloudflare.com/client/v4/zones/{your_zone_id}/dns_records"

# API 请求头
headers = {
    'X-Auth-Email': your_email,
    'X-Auth-Key': your_api_key,
    'Content-Type': 'application/json'
}

def update_dns_record():
    try:
        # 获取公网 IP 地址
        response = requests.get(url_get_ip)
        your_ip_address = response.json()['origin']
        print(f"Domain: {your_domain}, IP: {your_ip_address}")

        # 获取现有的 DNS 记录
        response = requests.get(url_get_dns_record, headers=headers)
        response_json = response.json()
        if response_json['success'] and len(response_json['result']) > 0:
            record = response_json['result'][0]

            # 更新 DNS 记录
            url_update_dns_record = f"https://api.cloudflare.com/client/v4/zones/{your_zone_id}/dns_records/{record['id']}"
            payload = {
                'type': 'A',
                'name': your_domain,
                'content': your_ip_address
            }
            response = requests.put(url_update_dns_record, headers=headers, json=payload)
            if response.json()['success']:
                print("DNS record updated successfully!")
            else:
                print("Failed to update DNS record!")
                print("Errors:", response.json()['errors'])
        else:
            # 创建 DNS 记录
            payload = {
                'type': 'A',
                'name': your_domain,
                'content': your_ip_address
            }
            response = requests.post(url_create_dns_record, headers=headers, json=payload)
            if response.json()['success']:
                print("DNS record created successfully!")
            else:
                print("Failed to create DNS record!")
                print("Errors:", response.json()['errors'])
    except requests.RequestException as e:
        print(f"An HTTP request error occurred: {str(e)}")
    except Exception as e:
        print(f"An unexpected error occurred: {str(e)}")

# 无限循环，每一分钟执行一次
while True:
    update_dns_record()
    time.sleep(60)
