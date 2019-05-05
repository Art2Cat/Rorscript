server {
  listen 80 default_server;
  listen [::]:80 default_server;

  root /usr/share/nginx/html;

  index index.html;

  server_name art2cat.com www.art2cat.com;

  location / {
    try_files $uri $uri/ =404;
  }
}
