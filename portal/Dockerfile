# Portal Dockerfile - Nginx para servir landing page
FROM nginx:1.25-alpine

# Remover configuração padrão do nginx
RUN rm /etc/nginx/conf.d/default.conf

# Copiar configuração customizada
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copiar arquivos do portal
COPY index.html /usr/share/nginx/html/
COPY assets/ /usr/share/nginx/html/assets/
COPY app/ /usr/share/nginx/html/app/
COPY services/ /usr/share/nginx/html/services/

# Expor porta 80
EXPOSE 80

# Comando para iniciar nginx
CMD ["nginx", "-g", "daemon off;"]
