export const useUser = () => {
  return useState<User | null>("user", () => null);
};

export const useFetchMe = async () => {
  const userAPI = useUserAPI(0); // We don't have org id yet
  const user = useUser();
  try {
    const data = await userAPI.me();
    user.value = data;
  } catch (err: any) {
    // 401 is expected for guest users - silently fail
    if (err?.statusCode !== 401 && err?.status !== 401) {
      console.error(err);
    }
    // Keep user as null for guests
    user.value = null;
  }
};

export const useUserTempID = () => {
  return useState<string | null>("userTempID", () => null);
};

export const useFetchUserTempID = async () => {
  if (typeof localStorage === "undefined") return;

  const userTempID = useUserTempID();

  if (window.self !== window.top) {
    const tempID = localStorage.getItem("Agents-Temporary-User-ID");
    if (tempID) {
      userTempID.value = tempID;
    } else {
      userTempID.value = crypto.randomUUID();
      localStorage.setItem("Agents-Temporary-User-ID", userTempID.value);
    }
    return;
  }
};
