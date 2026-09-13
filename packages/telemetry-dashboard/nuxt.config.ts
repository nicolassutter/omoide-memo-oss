import { heyApiPlugin } from "@hey-api/vite-plugin";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  devtools: { enabled: true },
  runtimeConfig: {
    public: {
      telemetryServerBaseUrl: "http://localhost:9999",
    },
  },
  vite: {
    plugins: [
      heyApiPlugin({
        vite: {
          apply: "serve",
        },
      }),
    ],
  },
});
