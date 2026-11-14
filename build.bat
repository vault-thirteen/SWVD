@ECHO OFF

cd cmd\list
go build
move *.exe ..\..
cd ..\..

cd cmd\download
go build
move *.exe ..\..
cd ..\..

cd cmd\dwf
go build
move *.exe ..\..
cd ..\..
