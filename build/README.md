
# Tomtom Crawler
Crawl data from Tomtom

## Installation
#### Install golang 1.21

- Macos:
  ```bash  
  brew install go
  ```  
- Windows:
  ```  
  https://dl.google.com/go/go1.21.5.windows-amd64.msi  
  ```  
- Docs:
  ```
  https://go.dev/doc/install
  ```

#### Install playwright

```  
go run github.com/playwright-community/playwright-go/cmd/playwright@latest install --with-deps  
```  
Docs:
```
https://playwright.dev/
```

## Usage

Crawl the address data from Tomtom

#### Windows
```shell  
./tomtom-addr-crawler-windows-amd64.exe c  
OPTIONS:  
 --input value, -i value       Input CSV file (default: "data/data.csv")  
  --output value, -o value      Ouput CSV file  --domain-url value, -d value  Tomtom Domain (default: "https://plan.tomtom.com/en/?p=10.82734,106.66315,9.55z&q=10.76397248,106.6881186")  
  --help, -h                    show help  
```  

## Contributing

Pull requests are welcome. For major changes, please open an issue first  
to discuss what you would like to change.

Please make sure to update tests as appropriate.
