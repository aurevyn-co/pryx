import { Hono } from "hono"

const app = new Hono()

app.get("/health", (c) => c.json({ status: "ok" }))
app.get("*", (c) => c.text("Pryx web worker"))

export default app
