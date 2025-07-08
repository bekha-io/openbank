#!/bin/bash

SERVER_USER=shuhrat_admin
SERVER_IP=172.164.248.87
KEY_PATH=~/.ssh/neobankone_key.pem
PROJECT_DIR=/home/$SERVER_USER/openbank

echo "===> Подключаемся к серверу и деплоим"

ssh -i $KEY_PATH $SERVER_USER@$SERVER_IP << EOF
  echo "Переключаемся в директорию проекта"
  cd $PROJECT_DIR
  
  echo "Пуллим изменения из ветки origin/fineract"
  git pull origin fineract
  
  echo "Запускаем Docker Compose"
  docker compose up -d --build
  
  echo "Деплой завершён"
EOF