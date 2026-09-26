const server = Bun.serve({
  port: Number(process.env.PORT) || 8001,
  hostname: "0.0.0.0",
  fetch(req, server) {
    const url = new URL(req.url);
    if (url.pathname === "/health") {
      return Response.json({ status: "ok", service: "realtime-service" });
    }
    return Response.json({ message: "Palabrario Realtime Service activo" });
  },
});

console.log(`⚡ [realtime-service] corriendo en http://${server.hostname}:${server.port}`);