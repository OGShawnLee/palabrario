FROM oven/bun:1-alpine

WORKDIR /app

# Copiar solo dependencias primero para cachear capas de Docker
COPY package.json bun.lockb* bun.lock* ./

RUN bun install

# Copiar el resto del código
COPY . .

EXPOSE 8000

# Por defecto en desarrollo ejecuta con watch (se sobreescribe con compose si se desea)
CMD ["bun", "--watch", "run", "index.ts"]
