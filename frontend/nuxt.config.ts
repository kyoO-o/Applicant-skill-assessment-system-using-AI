import vuetify, { transformAssetUrls } from "vite-plugin-vuetify"

const isDev = process.env.NODE_ENV !== "production"

const domain = process.env.DOMAIN

console.log("Server type:", process.env.SERVER_TYPE)

const nitro =
  process.env.SERVER_TYPE != "local"
    ? {
        routeRules: {
          "pub/**": { proxy: "http://localhost:3000/api/**" },
        },
      }
    : {
        routeRules: {
          "/api/**": { proxy: "http://localhost:4000/api/**" },
          "/pub/**": { proxy: "http://localhost:4000/pub/**" },
        },
        devProxy: {
          "/api": {
            target: "http://localhost:4000/api/",
            changeOrigin: true,
            ws: true,
          },
          "/pub": {
            target: "http://localhost:4000/pub/",
            changeOrigin: true,
            ws: true,
          },
        },
      }

export default defineNuxtConfig({
  devtools: { enabled: true },

  build: {
    transpile: ["vuetify"],
  },

  plugins: [
    "~/plugins/toast.ts",
    "~/plugins/formatNumberWithKAndM.js",
    "~/plugins/vue-zoomer.js",
    "~/plugins/mitt.ts",
  ],

  modules: [
    (_options, nuxt) => {
      nuxt.hooks.hook("vite:extendConfig", (config) => {
        // @ts-expect-error
        config.plugins.push(
          vuetify({
            autoImport: true,
            styles: {
              configFile: "./assets/settings.scss",
            },
          }),
        )
      })
    },
  ],

  ssr: false,

  sourcemap: {
    server: false,
    client: false,
  },

  runtimeConfig: {
    public: {
      dev: isDev,
      domain: domain,
    },
  },

  vite: {
    vue: {
      template: {
        transformAssetUrls,
      },
    },
    css: {
      preprocessorOptions: {
        scss: {
          silenceDeprecations: ["legacy-js-api"],
        },
        sass: {
          silenceDeprecations: ["legacy-js-api"],
        },
      },
    },
  },

  nitro,

  app: {
    head: {
      title: "SkillRecruit",
      titleTemplate: "%s",
      link: [
        { rel: "stylesheet", href: "/feather/style.css" },
        { rel: "icon", type: "image/x-icon", href: "/icons/icon.svg" },
        { rel: "logo-brand", sizes: "180x180", href: "/icons/icon.svg" },
        {
          rel: "icon",
          type: "image/svg",
          sizes: "32x32",
          href: "/icons/icon.svg",
        },
        {
          rel: "icon",
          type: "image/svg",
          sizes: "16x16",
          href: "/icons/icon.svg",
        },
        {
          rel: "stylesheet",
          href: "https://fonts.googleapis.com/css2?family=Rubik:wght@400;500;700&display=swap",
        },
      ],
    },
  },

  compatibilityDate: "2024-12-18",
})
