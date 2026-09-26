import Elysia from "elysia";

const app = new Elysia()
  .get("/", "Bienvenido a Palabrario Core API (Bun + Elysia)")
  .get("/health", () => ({ status: "ok", service: "core-service" }));

app.listen(process.env.PORT || 8000);

console.log(
  `🚀 [core-service] corriendo en http://${app.server?.hostname}:${app.server?.port}`,
);
