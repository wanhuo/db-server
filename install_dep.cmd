: Install dependence
cd ..\..\
echo %cd%
set GOPATH=%cd%
go get -v github.com/jingwanglong/cellnet
go get -v github.com/astaxie/beego
go get -v github.com/go-sql-driver/mysql
go get -v github.com/kardianos/service
go get -v github.com/golang/protobuf
go get -v github.com/davyxu/goobjfmt
pause

