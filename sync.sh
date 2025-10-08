#/bin/bash

echo "[$(date '+%Y-%m-%d %H:%M:%S')]　sync" >> /var/www/homing/batch.log
curl -X POST http://homing_app_1/api/sync >> /var/www/homing/batch.log 2>&1