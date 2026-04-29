import type { UserFilter } from "../../types/usermod/user_filter";
import type { User } from "../../types";
import type { FetchList, SaveUserPayload } from "~/composables/types/payload";

export class UserAPI {
  baseURL: string;
  orgID: number;

  constructor(baseURL: string, orgID: number) {
    this.baseURL = baseURL;
    this.orgID = orgID;
  }

  setOrgID(orgID: number) {
    this.orgID = orgID;
  }

  me() {
    return $fetch<User>(`${this.baseURL}/api/me`);
  }

  temporaryID() {
    return $fetch<string>(`${this.baseURL}/pub/temporary-id`);
  }

  saveMyProfile(item: Partial<User>, profile: File | undefined = undefined) {
    const data = JSON.stringify(item);

    const formData = new FormData();
    formData.set("data", data);
    if (profile) {
      formData.set("file", profile);
    }
    return $fetch<User>(`${this.baseURL}/api/me`, {
      method: "PUT",
      body: formData,
    });
  }

  async saveUserProfile(
    item: Partial<User>,
    profile: File | undefined = undefined,
  ) {
    const data = JSON.stringify(item);

    const formData = new FormData();
    formData.set("data", data);
    if (profile) {
      formData.set("file", profile);
    }

    return $fetch<User>(
      `${this.baseURL}/api/orgs/${this.orgID}/users/${item.id}`,
      {
        method: "PUT",
        body: formData,
      },
    );
  }

  getAll(
    filter: UserFilter = {},
    page: number | undefined = undefined,
    size: number | undefined = undefined,
  ) {
    return $fetch<FetchList<User>>(
      `${this.baseURL}/api/orgs/${this.orgID}/users`,
      {
        query: {
          ...filter,
          page,
          size,
        },
      },
    );
  }

  get(id: number) {
    return $fetch<FetchList<User>>(
      `${this.baseURL}/api/orgs/${this.orgID}/users/${id}`,
    );
  }

  save(item: SaveUserPayload) {
    return $fetch<User>(`${this.baseURL}/api/orgs/${this.orgID}/users`, {
      method: "PUT",
      body: item,
    });
  }

  remove(id: number) {
    return $fetch<string>(
      `${this.baseURL}/api/orgs/${this.orgID}/users/${id}`,
      {
        method: "DELETE",
      },
    );
  }
}
