export const useCustomer = () => {
  return useState("user", () => null)
}

export const useFetchMe = async () => {
  const user = useCustomer()
  const { data: c, error } = await useFetch("/api/me")
  if (error.value) {
    user.value = null
  }
  user.value = c.value
}
