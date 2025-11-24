#/bin/bash

echo "[$(date '+%Y-%m-%d %H:%M:%S')]　sync" >> /var/www/homing/batch.log
curl -X POST http://localhost:8090/api/sync >> /var/www/homing/batch.log 2>&1