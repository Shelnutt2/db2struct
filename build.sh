
#!/bin/bash -e

name="db2struct"
start_seconds=$(date +%s);

#export CGO_ENABLED="0"
export GO111MODULE=on
go build -o $name  cmd/db2struct/main.go

if [ $? -ne 0 ];then
  end_seconds=$(date +%s);
  echo
  echo "---> build time: "$((end_seconds-start_seconds))"s, build error!!!"
  echo
  exit 1
else
  end_seconds=$(date +%s);
  echo
  ls -lt |grep $name
  echo
  echo "Build Success,took "$((end_seconds-start_seconds))"s"
  echo
  exit 0
fi




