package fly

import (
	"io"
	"net/http"
	"os"
)

// return the stream and size of pac file
func GetPAC() ([]byte, int) {
	// run from the same dir with flynet.pac
	file, err := os.Open(`flynet.pac`)
	if err != nil {
		// run from cmd/client/client.go
		file, err = os.Open("../../flynet.pac")
		if err != nil {
			// run from fly/pac_test.go
			file, err = os.Open("../flynet.pac")
			if err != nil {
				logger.Println("pac file not found! start downlaod from github...")
				err := downloadPAC()
				logger.Println("pac file download completed!")
				if err != nil {
					return nil, 0
				}
				return GetPAC()
			}
		}
	}

	index := 0
	fileBuff := make([]byte, 200*1024)
	for {
		logger.Println("start reading...")
		n, err := file.Read(fileBuff[index:])
		if err != nil {
			if err == io.EOF {
				logger.Println("read pac file ok.")
				break
			}
			logger.Println("read pac file error --->", err)
		}
		index += n
	}
	return fileBuff, index
}

func downloadPAC() error {
	fileName := "flynet.pac"
	resp, err := http.Get("https://raw.githubusercontent.com/petronny/gfwlist2pac/master/gfwlist.pac")
	if err != nil {
		logger.Println("download pac file failed --->", err)
		logger.Panicln("please check your network or disable pac mode!")
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(fileName)
	if err != nil {
		logger.Println("create pac file failed --->", err)
		return err
	}
	io.Copy(file, resp.Body)
	return nil
}
