#! /bin/sh
build_service(){
    echo "$1 compiling..."
    go build -o ./service/bin/$1.exe ./service/$1/main.go
    lsres=`ls ./service/bin | grep $1`
    if [ lsres = '' ]
    then
        echo "$1 compile failed"
    else
        echo "$1 compiled"
    fi
}

#linux重定向 > 文件不存在就新建，存在就覆盖， >> 文件不存在就新建，存在就追加
#2>&1标准错误输出到标准输出（也就是屏幕），1前面如果不加&会被当成普通文件，加了才是标准输出
run_service(){
    echo "start to run $1..."
    nohup ./service/bin/$1.exe >> $logpath/$1.log 2>&1 &
}

check_service(){
    res=`ps -a | grep "./service/bin/$1.exe"`
    if [ res=='' ]
    then
        echo "failed to run $1"
    fi
}


services="dbproxy
account
upload
apigateway
"
mkdir -p ./service/bin/ && rm -f ./service/bin/*
logpath=./data/log/cloud-store
mkdir -p $logpath
echo "start to compile:$1"

for service in $services
do
    build_service $service
done

for service in $services
do
    run_service $service
done

for service in $services
do
    sleep 5
    check_service $service
done

