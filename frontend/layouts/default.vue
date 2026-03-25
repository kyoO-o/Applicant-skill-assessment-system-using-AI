<script setup>
const user = useCustomer()
const route = useRoute()
const { $toast, $listen, $clear } = useNuxtApp()
const firstInitial = computed(() => {
  if (user.value?.name === "") {
    const initials = user.value.email
      .split(" ")
      .map((s) => s[0])
      .join("")
    return initials.slice(0, 2).toUpperCase()
  } else if (user.value?.name && typeof user.value.name === "string") {
    if (user.value.name.split(" ").length === 1) {
      return user.value.name.slice(0, 1)
    } else {
      const initials = user.value.name
        .split(" ")
        .map((s) => s[0])
        .join("")
      return initials.slice(0, 2).toUpperCase()
    }
  } else {
    return "😀"
  }
})
console.log("User:", user.value)
const dialog = ref(false)
</script>
<template>hel</template>
<style scoped>
.no-shadow.v-btn--hovered {
  box-shadow: none !important;
}
.footer {
  background-color: white;
  z-index: 1;
}
</style>
