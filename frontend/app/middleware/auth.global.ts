const GUEST_ONLY_PATHS = new Set(["/login", "/register"]);

export default defineNuxtRouteMiddleware(async (to) => {
  const { user, me, initialized } = useAuth();

  if (!initialized.value) {
    await me();
  }

  if (!user.value && !GUEST_ONLY_PATHS.has(to.path)) {
    return navigateTo("/login");
  }
});
