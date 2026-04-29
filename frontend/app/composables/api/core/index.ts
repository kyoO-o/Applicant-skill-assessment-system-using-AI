import { UserAPI } from "./users";

const baseURL = "/m/core";

export const useUserAPI = (orgID: number) => new UserAPI(baseURL, orgID);
