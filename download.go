package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"golang.org/x/exp/slices"
)

var available = map[string]Download{
	"asn-country":				Download{ "asn-country", "csv", "COUNTRY", "https://cdn.jsdelivr.net/npm/@ip-location-db/", []string{} },
	"dbip-country": 			Download{ "dbip-country", "csv", "COUNTRY", "https://cdn.jsdelivr.net/npm/@ip-location-db/", []string{ "DBIP-LICENSE" } },
	"geo-asn-country":			Download{ "geo-asn-country", "csv", "COUNTRY", "https://cdn.jsdelivr.net/npm/@ip-location-db/", []string{} },
	"geo-whois-asn-country":	Download{ "geo-whois-asn-country", "csv", "COUNTRY", "https://cdn.jsdelivr.net/npm/@ip-location-db/", []string{} },
	"geolite2-country":			Download{ "geolite2-country", "csv", "COUNTRY", "https://cdn.jsdelivr.net/npm/@ip-location-db/", []string{ "GEOLITE2_LICENSE", "GEOLITE2_EULA" } },
	"iptoasn-country":			Download{ "iptoasn-country", "csv", "COUNTRY", "https://cdn.jsdelivr.net/npm/@ip-location-db/", []string{} },
	"webnet77-country":			Download{ "webnet77-country", "csv", "COUNTRY", "https://cdn.jsdelivr.net/npm/@ip-location-db/", []string{ "WEBNET77-LICENSE" } },

	"dbip-city":				Download{ "dbip-city", "gz", "CITY", "https://unpkg.com/@ip-location-db/", []string{ "DBIP-LICENSE" } },
	"geolite2-city":			Download{ "geolite2-city", "gz", "CITY", "https://cdn.jsdelivr.net/npm/@ip-location-db/", []string{ "GEOLITE2_LICENSE", "GEOLITE2_EULA" } },

	"asn":						Download{ "asn", "csv", "ASN", "https://cdn.jsdelivr.net/npm/@ip-location-db/", []string{ "ROUTEVIEWS-LICENSE", "DBIP-LICENSE" } },
	"dbip-asn":					Download{ "dbip-asn", "csv", "ASN", "https://cdn.jsdelivr.net/npm/@ip-location-db/", []string{ "DBIP-LICENSE" } },
	"geolite2-asn":				Download{ "geolite2-asn", "csv", "ASN", "https://cdn.jsdelivr.net/npm/@ip-location-db/", []string{ "GEOLITE2_LICENSE", "GEOLITE2_EULA" } },
	"iptoasn-asn":				Download{ "iptoasn-asn", "csv", "ASN", "https://cdn.jsdelivr.net/npm/@ip-location-db/", []string{} },
}

func downloadDataToLoad(missing []string) []DataToLoad {
	downloadPath := "./downloads"
	if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
		err := os.MkdirAll(downloadPath, 0755)
		if err != nil {
			panic(err)
		}
	}

	var dataToLoad []DataToLoad
	var downloads []Download

	downloads = downloadSelect("COUNTRY", downloads, missing)
	downloads = downloadSelect("CITY", downloads, missing)
	downloads = downloadSelect("ASN", downloads, missing)

	fmt.Println("checking for new data...")

	for _, download := range downloads {
		compression := "";
		if download.Format == "gz" {
			compression = ".gz"
		}

		urls := []string{
			fmt.Sprintf(download.CDN + "%s/%s-ipv4.csv%s", download.Folder, download.Folder, compression),
			fmt.Sprintf(download.CDN + "%s/%s-ipv6.csv%s", download.Folder, download.Folder, compression),
		}

		for _, url := range urls {
			ipVersion := 4
			if strings.Contains(url, "ipv6") {
				ipVersion = 6
			}

			fileName := path.Base(url)
			filePath := downloadPath + "/" + fileName

			changed, err := downloadFile(filePath, url)
			if err != nil {
				panic(err)
			}

			loadPath := ""
			if changed && compression != "" {
				// New file that needs decompressing first
				err := decompressFile(filePath, compression)
				if err != nil {
					panic(err)
				}
				loadPath = strings.Replace(filePath, compression, "", -1)
			} else if changed {
				// New file
				loadPath = filePath
			} else {
				if len(missing) > 0 && slices.Contains(missing, download.Type) {
					// Existing file, but our data hasn't been loaded, so re-process the old one
					loadPath = filePath
					if compression != "" {
						loadPath = strings.Replace(filePath, compression, "", -1)
					}
				}
			}

			if loadPath != "" {
				dataToLoad = append(dataToLoad, DataToLoad{ download, loadPath, ipVersion })
			}
		}
	}

	return dataToLoad
}

func downloadFile(filePath string, url string) (bool, error) {
	etagFilePath := filePath + ".etag"

	currentEtag := ""
	if fileExists(etagFilePath) {
		currentEtag = fileReadSmall(etagFilePath)
	}

	fmt.Println("downloading Etag: " + url)
	newEtag := getEtag(url)

	if newEtag != "" && currentEtag != newEtag {
		fileWriteSmall(etagFilePath, newEtag)

		fmt.Println("downloading data file: " + url)
		resp, err := http.Get(url)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()

		out, err := os.Create(filePath)
		if err != nil {
			return false, err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)

		return true, err
	} else {
		fmt.Println("Etag unchanged, skipping")
	}

	return false, nil
}

func downloadSelect(name string, downloads []Download, missing []string) []Download {
	value := os.Getenv(name)
	if len(value) > 0 && (len(missing) == 0 || slices.Contains(missing, name)) {
		download, ok := available[value]
		if ok {
			downloads = append(downloads, download)
		} else {
			panic(value + " is not a valid " + name + " option")
		}
	}

	return downloads
}