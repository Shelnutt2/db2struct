#sh build.sh
echo


#path="/f/code/go/src/mis/internal/"
#path="/c/Users/admin/Desktop"
path="/f/code/go/src/oms/internal"


pwd=$IMPORT_DB_PWD
#dbname="clerk"
dbname="oms"

function create() {
    table=$1
    st=$2
      ./db2struct --host 192.168.67.9 -d $dbname -t $table --package model --struct $st  -p $pwd --user root  --json  --gorm >$path/model/$table.go
    #./db2struct --host 192.168.67.9 -d information_schema -t $table --package model --struct $st  -p $pwd --user root  --json  --gorm >$path/model/$table.go
    echo "create table $table to  struct $st  $pwd"
    cat $path/model/$table.go
}
echo "start..."


#create "mis_import_history" "MisImportHistory"
#create "mis_flows" "MisFlows"
#create "mis_flow_template" "MisFlowTemplate"

#create "ats_detail" "AtsDetail"
#create "ats_monthly" "AtsMonthly"

#create "mis_wage_tax" "MisWageTax"
#create "mis_user_info" "MisUserInfo"
#create "mis_lock" "MisLock"
#create "COLUMNS" "TColumn"

#oms===
create "brand" "OmsBrand"
create "cluster" "OmsCluster"
create "configure" "OmsConfigure"
create "configure_history" "OmsConfigureHistory"
create "ecs" "OmsEcs"
create "env" "OmsEnv"
create "env_history" "OmsEnvHistory"
create "image" "OmsImage"
create "instance" "OmsInstance"
create "instance_groups" "OmsInstanceGroups"
create "package" "OmsPackage"
create "script" "OmsScript"
create "script_history" "OmsScriptHistory"
create "security_group" "OmsSecurityGroup"
create "service_type" "OmsServiceType"


echo "success, please open $path"

