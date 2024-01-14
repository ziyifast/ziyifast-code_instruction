pkgbuild的所有参数如下：

-root：指定要打包的文件夹的路径。

-identifier：指定包的标识符（ID）。

-install-location：指定安装包的安装位置。

-scripts：指定安装前后脚本文件的路径。

-version：指定包的版本号。

-sign：指定签名文件的路径，以进行安全签名。

-keychain：指定要用于签名的钥匙链。

-timestamp：指定时间戳服务器的URL，用于验证签名。

-verbose：显示详细信息。

-nopayload：仅创建空载荷。

-filter：指定要包含或排除的文件和文件夹。

-template：指定要使用的pkgproj模板文件的路径。

sudo pkgbuild --root /Users/ziyi2/GolandProjects/ziyi/build --identifier com.test.pkg.project --version 1.0.1 --install-location /tmp/ziyi_pkg   --scripts /Users/ziyi2/scripts  /Users/ziyi2/Desktop/project.pkg
# --root /Users/ziyi2/GolandProjects/ziyi/build：指定要打包的文件夹
# --install-location /tmp/ziyi_pkg 指定安装包的安装位置(tmp目录，默认如果文件超过3天没人访问则被删除)
# --scripts /Users/ziyi2/scripts 指定preinstall、postintall脚本所在文件夹（在安装前、安装后执行的脚本）
#### 注意：脚本文件名必须为：preinstall、postintall且有可执行权限
