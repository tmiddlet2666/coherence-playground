#!/bin/bash
#

for file in alerts-dashboard.json coherence-dashboard-main.json federation-details-dashboard.json http-servers-summary-dashboard.json \
            member-details-dashboard.json proxy-server-detail-dashboard.json services-summary-dashboard.json cache-details-dashboard.json \
            federation-summary-dashboard.json  kubernetes-summary-dashboard.json members-summary-dashboard.json proxy-servers-summary-dashboard.json \
            caches-summary-dashboard.json elastic-data-summary-dashboard.json grpc-summary-dashboard.json machines-summary-dashboard.json \
            persistence-summary-dashboard.json service-details-dashboard.json coherence-executors-summary.json coherence-executor-detail.json
do
   echo "${file}"
   curl -s https://raw.githubusercontent.com/oracle/coherence-operator/master/dashboards/grafana/${file} -o dashboards/${file}
done
