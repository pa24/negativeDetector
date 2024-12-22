#!/bin/bash

# Прекращаем выполнение скрипта при ошибке
set -e

# Pull the latest code from GitHub
echo "Pulling the latest code from GitHub..."
git pull origin main

# Build and deploy using Docker Compose
echo "Building and deploying services using Docker Compose..."
docker compose -f docker-compose.yml down          # Останавливаем старые контейнеры
docker compose -f docker-compose.yml up --build -d # Собираем и запускаем сервисы в фоне

# Выводим список работающих контейнеров
echo "Running containers:"
docker ps
