#!/bin/bash
port=""
tmpdir="/tmp"
localIp=$(hostname -I | awk '{print $1}')

wget_timeout=3
# todo:wangpz switch to real ip
# wget_server="10.128.5.116:8080"
wget_server="10.131.24.88:8080"
wget_root="software/n9e"
wget_cmd="wget --tries=2 --timeout=${wget_timeout} http://${wget_server}/${wget_root}/"

downloadFileType="wget"


if [[ ! -e "/etc/redhat-release" ]];then
    echo "unsupport system, continue anyway."
    # exit 2
fi

if [[ -e "/usr/bin/systemctl" ]];then
    osrelease=7
else
    osrelease=6
fi

echo "os release ${osrelease}"


function downloadFile(){
    local filename=$1
    local basefilename=`basename ${filename}`
    if [[ -e $basefilename && -s $basefilename ]];then return 0;fi
    if [[ ${downloadFileType} == "zbxcli" ]];then
        zbxcli download $filename
    else
        which wget &>>/dev/null
        if [[ $? -ne 0 ]];then 
            yum install -y wget
            if [[ $? -ne 0 ]];then echo "yum install wget failed.";exit 2;fi
        fi
        ${wget_cmd}${filename} -O ${basefilename}
    fi
    if [[ (! -e ${basefilename}) || (! -s ${basefilename}) ]];then
        echo "download file ${basefilename} failed."
        exit 4
    fi
}


function installNodeExporter(){
    local type_="$1"
    binDir="/usr/bin"
    case $type_ in
        "postgres")
        binDir="/usr/local/postgres_monitor_agent/bin"
        ;;
    esac


    if [[ ! -d "$binDir" ]];then mkdir -p "$binDir";fi
    cd $binDir
    if [[ ! -e "node_exporter" ]];then
        downloadFile "node_exporter"
    fi
    
    chmod +x ./node_exporter
    if [[ $osrelease -ge 7 ]];then
        echo "[Unit]
Description=node exporter service
Documentation=https://prometheus.io
After=network.target

[Service]
Type=simple
User=root
Group=root
ExecStart=${binDir}/node_exporter --collector.filesystem.ignored-fs-types=^(nfs4|fuse.mfs|tmpfs|rpc_pipefs)$
Restart=always

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/node_exporter.service
        systemctl daemon-reload
        systemctl restart node_exporter
    else
        nohup ${binDir}/node_exporter '--collector.filesystem.ignored-fs-types=^(nfs4|fuse.mfs|tmpfs|rpc_pipefs)$' &>/dev/null &
    fi
    
    ps -ef | grep -v grep | grep -wq node_exporter
    if [[ $? -eq 0 ]];then
        echo "start node_exporter successfully"'!'
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

function installPostgresExporter(){
    port=$1
    echo "install postgres_exporter at $port"
    binDir="/usr/local/postgres_monitor_agent/bin"
    if [[ ! -d "$binDir" ]];then mkdir -p "$binDir";fi
    
    if [ ! -e $binDir/"postgres_exporter" ]; then
        if [ -e "postgres_exporter" ]; then
            echo "use local file"
            cp "postgres_exporter" $binDir/
        fi
    fi    
    
    cd $binDir
    if [[ ! -e "postgres_exporter" ]];then
        echo "no postgres_exporter, download it."
        downloadFile "postgres_exporter"
    fi
    chmod +x ./postgres_exporter
    
    if [[ $osrelease -ge 7 ]];then
        echo "[Unit]
Description=postgres exporter service
Documentation=https://github.com/prometheus-community/postgres_exporter
After=network.target

[Service]
Type=simple
User=root
Group=root
ExecStart=${binDir}/postgres_exporter
Restart=always
Environment="DATA_SOURCE_URI=localhost:${port}/postgres?sslmode=disable" 
Environment="DATA_SOURCE_USER=postgres" 

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/postgres_exporter.service
        systemctl daemon-reload
        systemctl restart postgres_exporter.service
    else
        nohup ${binDir}/postgres_exporter &>/dev/null &
    fi
    
    ps -ef | grep -v grep | grep -wq postgres_exporter
    if [[ $? -eq 0 ]];then
        echo "start postgres_exporter successfully, metrics at http://${localIp}:9187/metrics . "'!'
    else
        echo "start postgres_exporter failed."
        exit 4
    fi

    # sudo journalctl -u postgres_exporter
}

function removePgExporter(){
    if [[ $osrelease -ge 7 ]];then
        systemctl stop postgres_exporter
        echo "service postgres_exporter stopped"
    else
        ps -ef | awk '/postgres_exporter/{print $2}' | xargs kill -9 &>/dev/null
    fi
}


function createPostgresExporterUser(){
    return
	echo "CREATE OR REPLACE FUNCTION __tmp_create_user() returns void as \$\$
BEGIN
  IF NOT EXISTS (
          SELECT                       -- SELECT list can stay empty for this
          FROM   pg_catalog.pg_user
          WHERE  usename = 'p8s_exporter') THEN
    CREATE USER p8s_exporter;
  END IF;
END;
\$\$ language plpgsql;

SELECT __tmp_create_user();
DROP FUNCTION __tmp_create_user();

ALTER USER p8s_exporter WITH PASSWORD 'mz.com';
ALTER USER p8s_exporter SET SEARCH_PATH TO p8s_exporter,pg_catalog;

GRANT CONNECT ON DATABASE postgres TO p8s_exporter;
GRANT pg_monitor to p8s_exporter;
" > /tmp/pg_monitor_role_setup.sql

    
    # sudo apt install postgresql-client-common
	# DB_HOST="10.10.10.1" DB_USER="postgres"
    # export PGPASSWORD="app_user"
    postgresOutput=$(psql -U postgres -c "select version();")
    postgresVersion=$(echo $postgresOutput  | grep -o "PostgreSQL 1[0-9]")
    if [[ $postgresVersion == "" ]];then
        echo $postgresOutput
        echo "unsupported postgres versoin, exit"
        exit 2
    else
        # localhost 不用用户密码
        #psql -h $DB_HOST -U $DB_USER -f /tmp/pg_monitor_role_setup.sql    
        #psql -h $DB_HOST -U $DB_USER -c "GRANT pg_monitor to p8s_exporter;"
        psql -U postgres -f /tmp/pg_monitor_role_setup.sql    
        psql -U postgres -c "GRANT pg_monitor to p8s_exporter;"
    fi
}

function checkPort(){
    port=$(echo $1 | grep -oE "[0-9]+" )
    if [[ "x$port" == "x" ]];then echo "Invalid [redis|mysql] port:$1";exit 2;fi
}

function showHelp(){
    echo -e "Using [install|remove] [pg_exporter]"
}

function main(){
    if [[ $# -lt 2 ]];then
        echo $*
        showHelp
        exit 2
    fi
    
    if [[ "$1" == "install" ]];then
        if [[ "$2" == "pg_exporter" ]];then
            # checkPort $3
            port=$(echo $3 | grep -oE "[0-9]+" )
            if [[ "x$port" == "x" ]];then port="5432";fi
            # createPostgresExporterUser
            installPostgresExporter $port
            downloadFile "node_exporter"
            installNodeExporter
            # journalctl -u postgres_exporter
        elif [[ "$2" == "node_exporter" ]];then
            downloadFile "node_exporter"
            installNodeExporter
        fi
    elif [[ "$1" == "remove" ]];then
        if [[ "$2" == "pg_exporter" ]];then
            # checkPort $3
            removePgExporter
            #unRegisterConsul "$2" "$localIp" $port
        fi
    fi
}
# journalctl -u postgres_exporter.service -f
# sudo service postgres_exporter stop
main $*