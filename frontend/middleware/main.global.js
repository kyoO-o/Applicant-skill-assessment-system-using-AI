export default defineNuxtRouteMiddleware(async (to, from) => {
  // Authentication check
  let user = useCustomer()

  if (user.value === null) {
    await useFetchMe()
  }
  if (user.value === null && to.name !== "index" && to.name !== "terms") {
    return navigateTo("/login")
  }
  if (user.value !== null && to.name === "index" && user.value.total_storage_size > 0) {
    return navigateTo("/")
  }
})
