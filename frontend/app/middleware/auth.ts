export default defineNuxtRouteMiddleware(async () => {
  const { user, me, initialized } = useAuth();

  if (!initialized.value) {
    await me();
  }

  if (!user.value) {
    return navigateTo("/login");
  }
});
