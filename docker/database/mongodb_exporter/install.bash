#!/bin/bash

localIp=$(hostname -I | awk '{print $1}')

# exporter的bin目录, ${exporter_bin_dir}/mongodb_exporter 文件存在时，不会再下载文件
exporter_bin_dir="/usr/local/bin"
# mongodb的bin目录
mongo_dir="/usr/local/mongodb4.2/bin"
# exporter下载服务器地址
wget_server="10.131.24.88:8080"
# exporter下载根目录
wget_root="software/n9e"
wget_cmd="wget --timeout=5 http://$wget_server/$wget_root/"

if [[ -e "/usr/bin/systemctl" ]];then
    osrelease=7
else
    osrelease=6
fi

function downloadExporter() {
    local file=$1

    # 本地目录
    if [ ! -d $exporter_bin_dir ]; then
        echo "Create exporter folder: $exporter_bin_dir"
        mkdir -p $exporter_bin_dir
    fi
    
    # 如果本地目录的exporter不存在
    if [ ! -e ${exporter_bin_dir}/$file ]; then
        if [ -e ${file} ]; then
            # 如果当前目录存在exporter文件，直接使用当前目录的文件，不下载了
            echo "Use local ${file}"
            cp ${file} ${exporter_bin_dir}/${file}
        else
            echo "Downloading $file"
            ${wget_cmd}${file} -O ${exporter_bin_dir}/${file}
        fi
    fi
    
    # 确认文件已存在
    cd $exporter_bin_dir
    if [[ (! -e ${file}) || (! -s ${file}) ]];then
        echo "Download ${file} failed."
        exit 4
    else
        echo "Load ${file} success."
    fi
}

function addExporterUser(){
    cd $mongo_dir
    echo 'var expect = db.getUser("mongodb_exporter")
if (!expect){
// return;
db.createUser({
user: "mongodb_exporter",
      pwd:  "mz.com",   // or cleartext passwordpasswordPrompt()
      roles: [ {
  "role":"clusterMonitor",
  "db":"admin"
  },{
         "role":"read",
         "db":"local"
      } ]
})
}' > /tmp/create_user.js
    echo "Create mongodb_exporter user in mongodb"
    ./mongo admin --port 28030 /tmp/create_user.js
    rm /tmp/create_user.js
}

function installService(){
    chmod +x ${exporter_bin_dir}/mongodb_exporter

    mongo_url=$(netstat -tlnp | grep "mongo[d|s] " |  awk '{split($4,arr,":"); print "mongodb://localIp:"arr[2]}' | tr '\n' ',')
    real_url=$(echo $mongo_url | sed "s/localIp/"${localIp}"/g")
    echo $real_url

    if [[ $osrelease -eq 7 ]];then
        echo "[Unit]
Description=mongodb exporter service
Documentation=https://prometheus.io
After=network.target

[Service]
Type=simple
User=root
Group=root
ExecStart=${exporter_bin_dir}/mongodb_exporter --no-mongodb.direct-connect --mongodb.uri=${real_url} --compatible-mode --collector.diagnosticdata --collector.replicasetstatus
Restart=on-failure

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/mongodb_exporter.service
        systemctl daemon-reload
        systemctl enable mongodb_exporter
        echo "Start mongodb_exporter as service, address http://${localIp}:9216/scrape?target=${real_url}"
        systemctl restart mongodb_exporter
    else
        echo "Start mongodb_exporter in nohup, address http://${localIp}:9216/scrape?target=${real_url}"
        nohup ${exporter_bin_dir}/mongodb_exporter --mongodb.uri=${real_url} --compatible-mode --collector.diagnosticdata --collector.replicasetstatus > /dev/null 2>&1 &
    fi
}

function installNodeExporter(){
    downloadExporter "node_exporter"
    
    chmod +x ${exporter_bin_dir}/node_exporter
    if [[ $osrelease -ge 7 ]];then
        echo "[Unit]
Description=node exporter service
Documentation=https://prometheus.io
After=network.target

[Service]
Type=simple
User=root
Group=root
ExecStart=${exporter_bin_dir}/node_exporter --collector.filesystem.ignored-fs-types=^(nfs4|fuse.mfs|tmpfs|rpc_pipefs)$
Restart=always

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/node_exporter.service
        systemctl daemon-reload
        systemctl restart node_exporter
    else
        nohup ${exporter_bin_dir}/node_exporter '--collector.filesystem.ignored-fs-types=^(nfs4|fuse.mfs|tmpfs|rpc_pipefs)$' &>/dev/null &
    fi
    
    ps -ef | grep -v grep | grep -wq node_exporter
    if [[ $? -eq 0 ]];then
        echo "start node_exporter successfully, address http://${localIp}:9100/metrics"
    else
        echo "start node_exporter failed."
        exit 4
    fi
    
}

function removeNodeExporter(){
    if [[ $osrelease -ge 7 ]];then
        systemctl stop node_exporter
    else
        ps -ef | awk '/node_exporter/{print $2}' | xargs kill -9
    fi
    ps -ef | grep -v grep | grep -wq node_exporter
    if [[ $? -ne 0 ]];then
        echo "finish to stop node_exporter"'!'
    else
        echo "stop node_exporter failed."
        exit 4
    fi
}

function uninstallService(){
    echo 'uninstall'
    if [[ $osrelease -ge 7 ]];then
        systemctl stop mongodb_exporter
        systemctl disable mongodb_exporter
    else
        ps -ef | awk '/mongodb_exporter/{print $2}' | xargs kill -9
    fi
}

function showHelp(){
    echo -e "Using $0 [install|uninstall] mongodb_exporter"
}

function main(){
    if [[ $# -lt 1 ]];then
        showHelp
        exit 2
    fi

    if [ $1 == 'install' ]; then
        installNodeExporter
        downloadExporter "mongodb_exporter"
        # addExporterUser
        installService
    elif [ $1 == 'uninstall' ]; then
        uninstallService
        removeNodeExporter
    else
        echo 'Invalid command'
        showHelp
        exit 2
    fi
}

function buildMongoUrl(){
    sudo netstat -tlnp | grep mongo |  awk '{split($4,arr,":"); print "mongodb://localhost:"arr[2]}' | tr '\n' ','
}

main $*