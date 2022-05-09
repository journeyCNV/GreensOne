package main

import (
	"github.com/journeycnv/greensone/gsweb"
	"github.com/journeycnv/greensone/gsweb/middleware"
	"github.com/journeycnv/greensone/test"
	"os"
)

func main() {
	gs := gsweb.NewGreensCore()
	gs.Use(middleware.Recovery(), middleware.LoggerDefault())
	gs.Bind(&test.DServiceProvider{})

	// 用户注册路由
	test.Register(gs)

	gs.Post("/upload_file", func(c *gsweb.Context) error {
		f, _ := c.FormFile("pic")
		dir, _ := os.Getwd()
		if f != nil {
			c.SaveUploadedFile(f, dir+"/pic")
			gsweb.LogInfo("upload ok", nil)
		}
		c.SetOkStatus()
		return nil
	})

	gs.Run()

	/**
	server := &http.Server{
		Handler: gs,
		Addr:    ":8080",
	}

	// 启动服务的goroutine
	go func() {
		server.ListenAndServe()
	}()

	// 服务关闭测试------------------------------------------------------------
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

	*/
}
