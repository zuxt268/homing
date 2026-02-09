#/bin/bash

echo "[$(date '+%Y-%m-%d %H:%M:%S')]　sync: instagram -> wordpress" >> /var/www/homing/batch.log
curl -X POST http://localhost:8090/api/sync/wordpress-instagram >> /var/www/homing/batch.log 2>&1

echo "[$(date '+%Y-%m-%d %H:%M:%S')]　sync: instagram -> GBP" >> /var/www/homing/batch.log
curl -X POST http://localhost:8090/api/sync/business-instagram >> /var/www/homing/batch.log 2>&1