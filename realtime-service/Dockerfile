FROM oven/bun:1-alpine

WORKDIR /app

# Copiar solo dependencias primero para cachear capas de Docker
COPY package.json bun.lock* ./

RUN bun install

# Copiar el resto del código
COPY . .

EXPOSE 8001

CMD ["bun", "--watch", "run", "index.ts"]
