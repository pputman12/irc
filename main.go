package main

import (
        "os/exec"
        "bufio"
        "fmt"
        "gopkg.in/irc.v3"
        "log"
        "math/rand"
        "net"
        "os"
        "time"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[seededRand.Intn(len(charset)-1)]
    }
    return string(b)
}


func shutdown() {
    cmd := exec.Command("sudo", "shutdown", "-h", "now")
    cmd.Run()
    err := cmd.Run()
    if err != nil {
            fmt.Println(err)
    }
}


func readLines(path string) ([]string, error) {
        file, err := os.Open(path)
        if err != nil {
            shutdown()
        }
        defer file.Close()

        var lines []string
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                lines = append(lines, scanner.Text())
        }
        return lines, scanner.Err()
}

func main() {
        conn, err := net.Dial("tcp", "irc.undernet.org:6667")
        if err != nil {
            shutdown()
        }

        words, err := readLines("words.txt")

        if err != nil {
            shutdown()
        }

        rand.Seed(time.Now().UTC().UnixNano())
        x := rand.Intn(len(words) - 1)
        nick := words[x]

        log.Println("Using nick ", nick)

        channel := "#gulag"

        randLength := 30

        charSet := "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"
        joinMessage := fmt.Sprintf("JOIN %s", channel)
        config := irc.ClientConfig{
                Nick: nick,
                User: nick,
                Name: nick,
                Handler: irc.HandlerFunc(func(c *irc.Client, m *irc.Message) {
                        if m.Command == "001" {
                                // 001 is a welcome event, so we join channels there

                                c.Write(joinMessage)
                                go func() {
                                        for i := 0; i < 50; i++ {
                                                randString := StringWithCharset(randLength, charSet)
                                                message := "I tried to harass PrezPusyGrab but all I got was this lousy bot   " + randString
                                                time.Sleep(time.Second)
                                                err := c.WriteMessage(&irc.Message{
                                                        Command: "PRIVMSG",
                                                        Params: []string{
                                                                channel,
                                                                message,
                                                        },
                                                })
                                                if err != nil {
                                                    shutdown()                                                    
                                                }
                                        }
                                        shutdown()
                                }()

                        } 
                }),
        }
        // Create the client
        client := irc.NewClient(conn, config)
        err = client.Run()
        if err != nil {
                shutdown()
        }
}
