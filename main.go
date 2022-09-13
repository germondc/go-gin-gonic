package main

import (
    "log"
    "net"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

type info struct {
    IP       string `json:"ip"`
    Hostname string `json:"host"`
    Handler  string `json:"handler"`
    Fullpath string `json:"fullpath"`
    ClientIP string `json:"clientip"`
    Content  string `json:"content"`
    RemoteIP string `json:"remoteip"`
    Headers  map[string][]string `json:"headers"`
    Params   gin.Params `json:"params"`
}

func main() {
    router := gin.Default()
    router.GET("/info", getInfo)
    router.GET("/metrics", gin.WrapH(promhttp.Handler()))

    router.Run("0.0.0.0:8080")
}

func getInfo(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, genInfo(c))
}

func genInfo(c *gin.Context) info {
    host, err := os.Hostname()
    if err != nil {
        host = err.Error()
    }
    return info{
        IP: getOutboundIP().String(),
        Hostname: host,
        Handler: c.HandlerName(),
        Fullpath: c.FullPath(),
        ClientIP: c.ClientIP(),
        Content: c.ContentType(),
        RemoteIP: c.RemoteIP(),
        Headers: c.Request.Header,
        Params: c.Params,
    }
}

func getOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
}
