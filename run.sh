#sh build.sh
echo


path="/f/code/go/src/mis/internal/"

pwd=$IMPORT_DB_PWD

function create() {
    table=$1
    st=$2
    ./db2struct --host 192.168.67.9 -d clerk -t $table --package model --struct $st  -p $pwd --user root  --json  --gorm >$path/model/$table.go
    echo "create table $table to  struct $st "
    cat $path/model/$table.go
}
echo "start..."


#create "mis_import_history" "MisImportHistory"
create "mis_flows" "MisFlows"
create "mis_flow_Template" "MisFlowTemplate"

echo "success, please open $path"

