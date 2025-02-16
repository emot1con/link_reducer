@echo off
powershell -ExecutionPolicy Bypass -NoProfile -Command "Invoke-WebRequest -Uri 'http://localhost:8080/cron/delete-expired' -Method Post"
