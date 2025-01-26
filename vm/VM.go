package vms

// import (
// 	"encoding/json"
// 	//"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"strings"
// 	//"sync"
// )
// //var mu sync.Mutex

// //VM에 http 요청을 보내 VM의 상태 정보를 가져옴.
// //모르겠는 내용
// //컴퓨터의 특정 API에 http 요청을 보내 VM의 상태정보를 가져온다는데.
// //내부적으로 api가 어떤식으로 연결되어 있는건지 모르겠음.
// //시각적으로 어떤식으로 api가 오가는지 알면 좋을것 같음.
// //Core와 control의 관게뿐만 아니라 벡엔드와의 관계도 궁금함.
// func (c *Core) GetVMList(VMList map[UUID]*VM) {
// 	//http나 https로 시작하지 않으면 http:// 붙이기
// 	if !strings.HasPrefix(c.IP, "http://") && !strings.HasPrefix(c.IP, "https://") {
// 		c.IP = "http://" + c.IP
// 	}

// 	//동시에 http.Get 하는것을 방지
// 	//mu.Lock()

// 	//http.Get을 통해 http://<IP>:8080/getStatus 경로로 접속
// 	//컴퓨터(Core)에서 실행 중인 가상 머신의 상태 정보(API)를 반환
// 	//JSON 형식으로 core에서 VM의 상태 정보를 받아옴.
// 	//아직 JSON파일 내부가 어떤식으로 구성 되어 있는지는 미정.
// 	//Core의 어떤 정보들이 필요할 지는 생각해봐야함.
// 	resp, err := http.Get(c.IP + ":8080" + "/getStatus")
// 	//mu.Unlock()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}

// 	//JSON 데이터 디코딩
// 	//JSON 데이터를 디코딩하여 VMInfo 구조체를 구축
// 	info := []VMInfo{}
// 	err = json.Unmarshal(body, &info)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	//fmt.Println(info)
// }

// func (i *InfraContext) UpdateList() {
// 	go func() {
// 		for _, c := range i.Computers {// c:모든VM
// 			c.GetVMList(i.VMPool)
// 		}
// 	}()
// }

// func GetVMListAll() []VMList{

// }
