# docker内建立mysql容器，并进行主从配置

1.拉取mysql5.7镜像

```shell
docker pull mysql:5.7
```

2.运行mysql，-d表示后台运行，-p是将host主机的3307映射到mysql容器的3306端口，--name是将运行的mysql容器命名为mysql-master，密码设置为123456

```shell
docker run -p 3307:3306 --name mysql-master -e MYSQL_ROOT_PASSWORD=123456 -d mysql:5.7
```

3.进入mysql-master的命令行

```shell
docker exec -it mysql-master /bin/bash
```

4.登录mysql，用户为root用户，密码123456

```shell
mysql -uroot -p123456
```

