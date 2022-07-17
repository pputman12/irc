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

func readLines(path string) ([]string, error) {
        file, err := os.Open(path)
        if err != nil {
                return nil, err
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
                log.Fatalln(err)
        }

        words, err := readLines("words.txt")

        if err != nil {
                log.Fatalln(err)
        }

        rand.Seed(time.Now().UTC().UnixNano())
        x := rand.Intn(len(words) - 1)
        nick := words[x]

        log.Println("Using nick ", nick)

        channel := "#bot-test-chan"
        message := "FREE PREZPUSYGRAB!!!!"
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
                                                time.Sleep(time.Second)
                                                c.WriteMessage(&irc.Message{
                                                        Command: "PRIVMSG",
                                                        Params: []string{
                                                                channel,
                                                                message,
                                                        },
                                                })
                                        }
                                        cmd := exec.Command("sudo", "shutdown", "-h", "now")
                                        cmd.Run()
                                        err := cmd.Run()
                                        if err != nil {
                                                fmt.Println(err)
                                        }
                                }()

                        } else if m.Command == "PRIVMSG" && c.FromChannel(m) {
                                // Create a handler on all messages.
                                c.WriteMessage(&irc.Message{
                                        Command: "PRIVMSG",
                                        Params: []string{
                                                m.Params[0],
                                                m.Trailing(),
                                        },
                                })
                        }
                }),
        }
        // Create the client
        client := irc.NewClient(conn, config)
        err = client.Run()
        if err != nil {
                log.Fatalln(err)
        }
}
