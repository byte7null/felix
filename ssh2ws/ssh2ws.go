package ssh2ws

import (
	"github.com/dejavuzhou/felix/ssh2ws/controllers"
	"github.com/dejavuzhou/felix/ssh2ws/utils"
	_ "github.com/dejavuzhou/felix/statik"
	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
	"time"
)

func RunSsh2ws(bindAddress, user, password string, expire time.Duration, secret []byte) error {
	r := gin.Default()
	r.MaxMultipartMemory = 32 << 20
	statikFS, err := fs.New()
	if err != nil {
		return err
	}
	r.Use(utils.Serve("/", statikFS))

	api := r.Group("api")
	r.POST("api/login", controllers.GetLoginHandler(user, password, expire, secret))
	api.Use(controllers.JwtAuthMiddleware(secret))

	{
		api.GET("ws/:id", controllers.WsSsh)

		api.GET("ssh", controllers.SshAll)
		api.POST("ssh", controllers.SshCreate)
		api.GET("ssh/:id", controllers.SshOne)
		api.PATCH("ssh/:id", controllers.SshUpdate)
		api.DELETE("ssh/:id", controllers.SshDelete)

		api.GET("sftp/:id", controllers.SftpLs)
		api.GET("sftp/:id/dl", controllers.SftpDl)
		api.GET("sftp/:id/cat", controllers.SftpCat)
		api.GET("sftp/:id/rm", controllers.SftpRm)
		api.GET("sftp/:id/rename", controllers.SftpRename)
		api.GET("sftp/:id/mkdir", controllers.SftpMkdir)
		api.POST("sftp/:id/up", controllers.SftpUp)
	}

	if err := r.Run(bindAddress); err != nil {
		return err
	}
	return nil
}