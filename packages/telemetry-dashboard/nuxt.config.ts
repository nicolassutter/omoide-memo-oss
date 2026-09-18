import tailwindcss from '@tailwindcss/vite'

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  app: {
    head: {
      title: 'Omoide Memo Telemetry',
      meta: [
        {
          name: 'description',
          content: 'Self-agnostic operational analytics platform',
        },
      ],
    },
  },
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  modules: ['shadcn-nuxt', '@nuxtjs/color-mode'],
  css: ['~/assets/css/tailwind.css'],
  colorMode: {
    classSuffix: '',
  },
  shadcn: {
    /**
     * Prefix for all the imported components
     */
    prefix: '',
    /**
     * Directory that the components live in
     */
    componentDir: './app/components/ui',
  },
  vite: {
    plugins: [tailwindcss()],
  },
})
