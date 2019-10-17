package fly

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func StartPAC() {
	http.HandleFunc("/flynet.pac", func (w http.ResponseWriter, r *http.Request){
		// the pac file should be placed at the same dir with current running file.
		// but it is placed at the parent dir when development
		file, err := os.Open(`flynet.pac`)
		if err != nil {
			fmt.Println(err)
			// run from cmd/client/
			file, _ = os.Open("../../flynet.pac")
		}
		w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
		io.Copy(w, file)
		w.WriteHeader(200)

	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("start failed --->", err)
	}
}
