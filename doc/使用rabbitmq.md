## docker内使用rabbitmq

1. 拉取镜像，

    坑：不加management下载的版本会不带web界面，无法通过浏览器直接看

    ```shell
    docker pull rabbitmq:management
    ```

2. 运行容器

    ```shell
     docker run -d --hostname rabbithost --name rabbit-node1 -p 5672:5672 -p 15672:15672 -v rabbitmq:/var/lib/rabbitmq rabbitmq:management
    ```

    rabbitmq:/var/lib/rabbitmq

    rabbitmq是相对目录，在宿主机上，相对于/var/lib/docker/volumes/，/var/lib/rabbitmq是容器内的绝对目录

    