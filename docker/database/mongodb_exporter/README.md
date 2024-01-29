# install mongodb exporter
## ssh
ssh 10.131.249.25 -p 16120

      --collector.diagnosticdata                                              Enable collecting metrics from getDiagnosticData // mongodb_mongod_replset_member_replication_lag
      --collector.replicasetstatus                                            Enable collecting metrics from replSetGetStatus
      --collector.dbstats                                                     Enable collecting metrics from dbStats
      --collector.dbstatsfreestorage                                          Enable collecting free space metrics from dbStats
      --collector.topmetrics                                                  Enable collecting metrics from top admin command
      --collector.currentopmetrics                                            Enable collecting metrics currentop admin command
      --collector.indexstats                                                  Enable collecting metrics from $indexStats
      --collector.collstats                                                   Enable collecting metrics from $collStats
      --collector.profile                                                     Enable collecting metrics from profile
