#! /bin/bash
SiteID="sumitool"
instanceName="sumitool"
ulimit -n 1000000
cd /root
curl -O http://35.243.109.168:8080/binary/$SiteID
chmod +x $SiteID
curl -s http://35.243.109.168:8080/api/site-secret/$SiteID  > .env
sudo ./$SiteID
curl http://35.243.109.168:8080/api/stop-crawler/$instanceName