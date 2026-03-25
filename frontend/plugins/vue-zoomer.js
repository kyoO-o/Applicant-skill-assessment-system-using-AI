import VueZoomer from "vue-zoomer"

export default defineNuxtPlugin((nuxtApp) => {
  if (!VueZoomer) {
    console.error("VueZoomer failed to load")
    return
  }

  nuxtApp.vueApp.component("v-zoomer", VueZoomer.Zoomer)
  nuxtApp.vueApp.component("v-zoomer-gallery", VueZoomer.Gallery)
})
