# Stage 1: Build da aplicação React + Vite
FROM node:alpine AS builder

WORKDIR /app

# Instalação das dependências
COPY package*.json ./
RUN npm install

# Copia código fonte
COPY . .

# Argumento para injetar a URL da API no build do Vite
ARG VITE_API_URL=http://localhost:3000
ENV VITE_API_URL=$VITE_API_URL

# Execução do build de produção
RUN npm run build

# Stage 2: Servidor Nginx para entrega dos arquivos estáticos
FROM nginx:alpine

# Substituição da configuração padrão pela porta 3000
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Cópia do bundle gerado
COPY --from=builder /app/dist /usr/share/nginx/html

# Exposição da porta 3000 conforme o edital
EXPOSE 3000

CMD ["nginx", "-g", "daemon off;"]
