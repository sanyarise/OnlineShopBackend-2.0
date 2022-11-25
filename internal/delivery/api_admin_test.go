package delivery

/*import (
	"bytes"
	"encoding/binary"
	"log"
	"net/http"
	"testing"
)

func TestUploadFile(t *testing.T) {
	str := []byte("12345")
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.LittleEndian, str)
	if err != nil {
		t.Error(err)
	}

	resp, err := http.Post("http://localhost:8000/items/yyy/upload", "image/jpeg", &buf)
	if resp.StatusCode != 201 {
		t.Error("Error status code, err: ", resp.Status, err.Error())
	}
	log.Print(resp.Status)
}*/
