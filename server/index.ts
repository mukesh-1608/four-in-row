import { spawn } from "child_process";

// This file is a lightweight wrapper to start the Go server
// because the Replit environment expects 'npm run dev' -> 'tsx server/index.ts'.
// The actual backend application logic is entirely in Go (main.go).

console.log("Bootstrapping Go server...");

const go = spawn("go", ["run", "main.go"], { 
  stdio: "inherit",
  env: process.env 
});

go.on("error", (err) => {
  console.error("Failed to start Go server:", err);
});

go.on("close", (code) => {
  console.log(`Go server exited with code ${code}`);
  process.exit(code ?? 0);
});
