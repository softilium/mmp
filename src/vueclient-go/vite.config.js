import { defineConfig } from "vite";
import plugin from "@vitejs/plugin-vue";
import mkcert from "vite-plugin-mkcert";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [plugin(), mkcert()],
  server: {
    port: 56056,
  },
});
