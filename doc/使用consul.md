## 运行consul

```shell
docker run -d -p 8500:8500 -h node1 --name node1  consul agent -server -client 0.0.0.0 -ui
```

