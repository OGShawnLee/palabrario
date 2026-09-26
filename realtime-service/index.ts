import Elysia from "elysia";

const app = new Elysia()
  .get("/", "Bienvenido a Palabrario Realtime Service (Bun + Elysia)")
  .get("/health", { status: "ok", service: "realtime-service" });

app.listen(process.env.PORT || 8001);
console.log(
  `⚡ [realtime-service] corriendo en http://${app.server?.hostname}:${app.server?.port}`,
);
