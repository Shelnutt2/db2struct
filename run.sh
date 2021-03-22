sh build.sh

echo
cat readme_linhc.md
echo

#数据库host
dbhost=$IMPORT_DB_HOST

#库名
dbname="clerk"

#项目名
projectname="mis"

#数据密码
pwd=$IMPORT_DB_PWD

#生成的model存放目录
path="/Users/linhaicheng/go/src/mis/internal"



 if [ -z "$pwd" ];then
    echo "请设置环境变量（数据库密码） IMPORT_DB_PWD   举个票子: export IMPORT_DB_PWD=\"xxx\""
    exit 1
 fi



mkdir -p "$path/model"

function create() {
    #表名
    table=$1
    #结体体名字
    st=$2
    #表-中文名
    note=$3
    ./db2struct --host $dbhost -d $dbname -t $table --package model --struct $st --project $projectname --note $note --target $path -p $pwd --user root  --json  --gorm 
    #echo "file =  $path/model/$table.go"
    #echo "create table $table to  struct $st  $pwd"
    #cat $path/model/$table.go
}
echo "start..."

#create  表名 结构体名字

#mis===
#create "service_type" "misServiceType"
create "mis_setting_ats" "MisSettingAts" "考勤设置"

echo "success, please open $path/model"

