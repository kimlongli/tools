FROM ubuntu:16.04
WORKDIR /root
RUN apt update
RUN apt install python-pip -y
RUN apt install gcc -y
RUN apt install vim -y
RUN apt install git -y
#RUN apt-get install software-properties-common -y
#RUN apt-add-repository universe -y
#RUN apt-get update -y
RUN pip install shadowsocks

# 下载kcptun
RUN git clone https://github.com/kimlongli/tools.git

ADD start.sh /root

# 启动shadowsocks server
CMD chmod +x ./start.sh && ./start.sh
