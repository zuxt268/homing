#/bin/bash

echo "[$(date '+%Y-%m-%d %H:%M:%S')] token" >> /var/www/homing/batch.log
curl -X POST http://homing_app_1/api/token/check >> /var/www/homing/batch.log 2>&1