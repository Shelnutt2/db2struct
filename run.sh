
echo
cat readme_linhc.md
echo

#数据库host
dbhost=$IMPORT_DB_HOST

#库名
dbname="clerk"

#数据密码
pwd=$IMPORT_DB_PWD

#生成的model存放目录
path="/Users/linhaicheng/tmp/mis"


#path="/c/Users/admin/Desktop"
#path="/f/code/go/src/mis/internal/"run
#path="/f/code/go/src/oms/internal"
#dbname="oms"
#dbname="information_schema"

 if [ -z "$pwd" ];then
    echo "请设置环境变量（数据库密码） IMPORT_DB_PWD   举个票子: export IMPORT_DB_PWD=\"xxx\""
    exit 1
 fi

 if [ -z "$dbhost" ];then
    echo "请设置环境变量（数据库host） IMPORT_DB_HOST   举个票子: export IMPORT_DB_HOST=\"xxx\""
    exit 1
 fi


mkdir -p "$path/model"

function create() {
    table=$1
    st=$2
    ./db2struct --host $dbhost -d $dbname -t $table --package model --struct $st  -p $pwd --user root  --json  --gorm >$path/model/$table.go
    echo "file =  $path/model/$table.go"
    echo "create table $table to  struct $st  $pwd"
    cat $path/model/$table.go
}
echo "start..."

#create  表名 结构体名字
#create "mis_import_history" "MisImportHistory"

#mis------
#create "mis_import_history" "MisImportHistory"
#create "mis_flows" "MisFlows"
#create "mis_flow_template" "MisFlowTemplate"

#create "ats_detail" "AtsDetail"
#create "ats_monthly" "AtsMonthly"

#create "mis_accident_judge" "MisAccidentJudge"
#create "mis_holiday" "MisHoliday"
#create "ats_leave_detail" "MisAtsLeaveDetail"
#create "ats_duty" "AtsDuty"
create "ats_log" "AtsLog"




#oms===
#create "security_group" "OmsSecurityGroup"
#create "service_type" "OmsServiceType"
#

echo "success, please open $path/model"

