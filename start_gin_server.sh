sudo killall main
go build .
sudo setsid ./main > ginlog.txt 2>&1 &