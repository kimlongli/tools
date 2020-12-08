ssserver -p 11503 -k lilin33221 -m aes-256-cfb --user nobody -d start
cd tools && ./server_linux_amd64 -t 127.0.0.1:11503 -l :4000 -mode fast3 -crypt none -nocomp
