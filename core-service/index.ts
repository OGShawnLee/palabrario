import Elysia from "elysia";
import DefinitionRouter from "@router/definition.router";

const app = new Elysia()
  .get("/", "Bienvenido a Palabrario Core API (Bun + Elysia)")
  .get("/health", () => ({ status: "ok", service: "core-service" }))
  .use(DefinitionRouter);

app.listen(process.env.PORT || 8000);

console.log(
  `🚀 [core-service] corriendo en http://${app.server?.hostname}:${app.server?.port}`,
);
