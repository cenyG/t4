# Используем nginx для статического сайта
FROM nginx:alpine

# Копируем наш сайт в директорию для статики Nginx
COPY index.html /usr/share/nginx/html/

# Открываем порт 80 для веб-доступа
EXPOSE 80