#!/bin/bash

local_ip=$(hostname -I | awk '{print $1}')
wget_server="10.131.24.88:8080"
wget_root="software/n9e"
bin_dir="/usr/local/bin"
tmp_dir="/tmp"

if [[ -e "/usr/bin/systemctl" ]];then
    osrelease=7
else
    osrelease=6
fi

echo "os release ${osrelease}"

# downloadFile to $tempdir
# args: $1:filename
function downloadFile(){
    local filename=$1
    local fileurl=http://${wget_server}/${wget_root}/${filename}
    local filepath=${tmp_dir}/${filename}
    rm -f ${filepath}
    echo "Download file ${filename} from ${fileurl}"
    wget ${fileurl} -O ${filepath}
    if [[ (! -e ${filepath}) || (! -s ${filepath}) ]];then
        echo "Download file ${filename} failed."
        exit 4
    fi
}

function removeSerivce(){
    local service_name=$1
    local service_file_path=/etc/systemd/system/${service_name}.service
    local bin_file_path=${bin_dir}/${service_name}
    if [[ $osrelease -ge 7 ]];then
        if [[ (-e ${service_file_path}) ]];then
            echo "stop ${service_name} service"
            systemctl stop ${service_name}
            rm ${service_file_path}
        fi
    else
        ps -ef | awk "/${service_name}/{print \$2}" | xargs kill -9
    fi
    if [[ (-e ${bin_file_path}) ]];then
        rm -f ${bin_file_path}
    fi
}

# download node_exporter and install as service
# systemctl stop node_exporter
function installNodeExporter(){
    removeSerivce "node_exporter"
    downloadFile "node_exporter"
    local tmp_file_path=${tmp_dir}/node_exporter
    local bin_file_path=${bin_dir}/node_exporter
    cp ${tmp_file_path} ${bin_file_path}
    if [ ! -e ${bin_file_path} ]; then
        echo "Copy ${tmp_file_path} to ${bin_file_path} failed."
        exit 4
    fi 
    chmod +x ${bin_file_path}
    if [[ $osrelease -ge 7 ]];then
        echo "[Unit]
Description=node exporter service
Documentation=https://prometheus.io
After=network.target

[Service]
Type=simple
User=root
Group=root
ExecStart=${bin_file_path} --collector.filesystem.ignored-fs-types=^(nfs4|fuse.mfs|tmpfs|rpc_pipefs)$
Restart=always

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/node_exporter.service
        systemctl daemon-reload
        systemctl restart node_exporter
    else
        nohup ${bin_file_path} '--collector.filesystem.ignored-fs-types=^(nfs4|fuse.mfs|tmpfs|rpc_pipefs)$' &>/dev/null &
    fi
    
    ps -ef | grep -v grep | grep -wq node_exporter
    if [[ $? -eq 0 ]];then
        echo "start node_exporter successfully, address: http://${local_ip}:9100/metrics"
    else
        echo "start node_exporter failed."
        exit 4
    fi
}

function installPostgresExporter(){
    removeSerivce "postgres_exporter"
    downloadFile "postgres_exporter"
    local tmp_file_path=${tmp_dir}/postgres_exporter
    local bin_file_path=${bin_dir}/postgres_exporter
    cp ${tmp_file_path} ${bin_file_path}
    if [ ! -e ${bin_file_path} ]; then
        echo "Copy ${tmp_file_path} to ${bin_file_path} failed."
        exit 4
    fi 
    chmod +x ${bin_file_path}
    
    port=$1
    echo "install postgres_exporter at $port"
    if [[ $osrelease -ge 7 ]];then
        echo "[Unit]
Description=postgres exporter service
Documentation=https://github.com/prometheus-community/postgres_exporter
After=network.target

[Service]
Type=simple
User=root
Group=root
ExecStart=${bin_file_path}
Restart=always
Environment="DATA_SOURCE_URI=localhost:${port}/postgres?sslmode=disable" 
Environment="DATA_SOURCE_USER=postgres" 

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/postgres_exporter.service
        systemctl daemon-reload
        systemctl restart postgres_exporter.service
    else
        nohup ${bin_file_path}/postgres_exporter &>/dev/null &
    fi
    
    ps -ef | grep -v grep | grep -wq postgres_exporter
    if [[ $? -eq 0 ]];then
        echo "start postgres_exporter successfully, metrics at http://${local_ip}:9187/metrics . "'!'
    else
        echo "start postgres_exporter failed."
        exit 4
    fi

    # sudo journalctl -u postgres_exporter
}


function installMongodbExporter(){
    removeSerivce "mongodb_exporter"
    downloadFile "mongodb_exporter"
    
    local tmp_file_path=${tmp_dir}/mongodb_exporter
    local bin_file_path=${bin_dir}/mongodb_exporter
    cp ${tmp_file_path} ${bin_file_path}
    if [ ! -e ${bin_file_path} ]; then
        echo "Copy ${tmp_file_path} to ${bin_file_path} failed."
        exit 4
    fi 
    chmod +x ${bin_file_path}

    mongo_url=$(netstat -tlnp | grep "mongo[d|s] " |  awk '{split($4,arr,":"); print "mongodb://localIp:"arr[2]}' | tr '\n' ',')
    real_url=$(echo $mongo_url | sed "s/localIp/"${local_ip}"/g")
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
ExecStart=${bin_file_path} --no-mongodb.direct-connect --mongodb.uri=${real_url} --compatible-mode --collector.diagnosticdata --collector.replicasetstatus
Restart=on-failure

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/mongodb_exporter.service
        systemctl daemon-reload
        systemctl enable mongodb_exporter
        echo "Start mongodb_exporter as service, address http://${local_ip}:9216/scrape?target=${real_url}"
        systemctl restart mongodb_exporter
    else
        echo "Start mongodb_exporter in nohup, address http://${local_ip}:9216/scrape?target=${real_url}"
        nohup ${bin_file_path} --mongodb.uri=${real_url} --compatible-mode --collector.diagnosticdata --collector.replicasetstatus > /dev/null 2>&1 &
    fi
}


function showHelp(){
    echo -e "Using [install|remove] [node_exporter|mongodb_exporter|postgres_exporter]"
}

function main(){
    if [[ $# -lt 2 ]];then
        echo "missing args"
        showHelp
        exit 2
    fi
    if [[ "$1" == "install" ]];then
        if [[ "$2" == "node_exporter" ]];then
            installNodeExporter
        elif [[ "$2" == "mongodb_exporter" ]];then
            echo "install mongodb_exporter"
            installNodeExporter
            installMongodbExporter
        elif [[ "$2" == "postgres_exporter" ]];then
            echo "install postgres_exporter"
            installNodeExporter
            port=$(echo $3 | grep -oE "[0-9]+" )
            if [[ "x$port" == "x" ]];then port="5432";fi
            installPostgresExporter $port
        fi
    elif [[ "$1" == "remove" ]];then
        if [[ "$2" == "node_exporter" ]];then
            echo "remove node_exporter"
            removeSerivce "node_exporter"
        elif [[ "$2" == "mongodb_exporter" ]];then
            echo "remove mongodb_exporter"
            removeSerivce "mongodb_exporter"
        elif [[ "$2" == "postgres_exporter" ]];then
            echo "remove postgres_exporter"
            removeSerivce "postgres_exporter"
        fi
    fi
}

main $*