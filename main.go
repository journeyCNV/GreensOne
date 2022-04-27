package main

import (
	"GreensOne/gsweb"
	"GreensOne/gsweb/middleware"
	"GreensOne/test"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	gs := gsweb.NewGreensCore()
	gs.Use(middleware.Recovery())
	gs.Use(middleware.Cost())

	test.Register(gs)
	server := &http.Server{
		Handler: gs,
		Addr:    ":8080",
	}

	// 启动服务的goroutine
	go func() {
		server.ListenAndServe()
	}()

	// 当前的等待信号量
	quit := make(chan os.Signal)
	// 监控以下信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 阻塞当前等待信号
	<-quit

	// 控制优雅关闭最多等待5s
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 当监听到关闭进程的信号之后，就会执行下面的优雅关闭 graceful shuts down
	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("server shutdown ", err)
	}
}
