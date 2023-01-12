package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

var actions = []string{
	"logged in",
	"logged out",
	"create record",
	"update record",
	"delete recoed",
}

const countResult = 100
const  countWorker = 20

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getInfo() string {
	out := fmt.Sprintf("ID:%d | Email: %s\nActivity Log:\n", u.id, u.email)
	for i, item := range u.logs {
		out += fmt.Sprintf("%d. [%s] at %s\n", i+1, item.action, item.timestamp)
	}
	return out
}

func main() {
	rand.Seed(time.Now().Unix())

	t := time.Now()
	jobs := make(chan int, countResult)
	users := make(chan User, countResult)

	wg := &sync.WaitGroup{}
	// users := generateUser(1000)
	// for _, user := range users {
	// 	wg.Add(1)
	// 	go saveUserInfo(user,wg,jobs,result)
	// }
	generateUser(countWorker,jobs,users)
	generateJobs(countResult, jobs,wg)
	saveUsersInfo(countWorker,users,wg)


	wg.Wait()

	fmt.Println("time elapsed", time.Since(t).String())


}
func generateJobs(count int, jobs chan<- int, wg *sync.WaitGroup) {
	for i := 0; i < count; i++ {
		wg.Add(1)
		jobs <- i
	}
}


func saveUsersInfo(countWorker int, users chan User, wg *sync.WaitGroup) {
	for i := 0; i < countWorker; i++ {
		go func() {
			for user := range users {
				fmt.Printf("WRITING FILE FOR UID %d\n", user.id)

				// create file
				filename := fmt.Sprintf("logs/uid%d.txt", user.id)
				file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
				if err != nil {
					return
				}

				_, err = file.WriteString(user.getInfo())
				if err != nil {
					return
				}

				time.Sleep(time.Second)

				// wait for all users to be saved
				wg.Done()
			}
		}()
	}
}
func generateUser(countWorker int,jobs <-chan int, users chan<- User )  {
	for i := 0; i < countWorker; i++ {
		go func() {
			for j := range jobs {
				users <- User{
					id:    j,
					email: fmt.Sprintf("user%d@company.com", j),
					logs:  generateLogs(rand.Intn(1000)),
				}
				fmt.Printf("generated user %d\n", j)
				time.Sleep(time.Millisecond * 100)
			}
			close(users)
		}()
	}
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			timestamp: time.Now(),
			action:    actions[rand.Intn(len(actions)-1)],
		}
	}

	return logs
}
