# Frontend Dockerfile - Nginx para servir arquivos estáticos
FROM nginx:1.25-alpine

# Instalar envsubst para substituição de variáveis
RUN apk add --no-cache gettext

# Remover configuração padrão do nginx
RUN rm /etc/nginx/conf.d/default.conf

# Copiar configuração customizada
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copiar arquivos do frontend
COPY dist/ /usr/share/nginx/html/

# Expor porta 80
EXPOSE 80

# Comando para iniciar nginx
CMD ["nginx", "-g", "daemon off;"]
