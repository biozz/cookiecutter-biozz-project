import { defineConfig, loadEnv } from "vite";
import solidPlugin from "vite-plugin-solid";

const env = loadEnv("all", process.cwd());

export default defineConfig({
  plugins: [solidPlugin()],
  server: {
    port: parseInt(env.VITE_PORT),
  },

  build: {
    target: "esnext",
    minify: true,
  },
});
