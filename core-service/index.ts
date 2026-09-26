const server = Bun.serve({
  port: Number(process.env.PORT) || 8000,
  hostname: "0.0.0.0",
  fetch(req) {
    const url = new URL(req.url);
    if (url.pathname === "/health") {
      return Response.json({ status: "ok", service: "core-service" });
    }
    return Response.json({ message: "Bienvenido a Palabrario Core API (Bun)" });
  },
});

console.log(`🚀 [core-service] corriendo en http://${server.hostname}:${server.port}`);