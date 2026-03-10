docker image build -t bs_build .
docker  run -it -d --restart=always -p 10903:4000/udp  bs_build