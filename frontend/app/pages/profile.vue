<script setup lang="ts">
import { toast } from "vue-sonner";

import type { Company } from "../composables/types";
import type { SaveCompanyPayload } from "../composables/types/payload";

const { user, me } = useAuth();
const companyAPI = useCompanyAPI();
const router = useRouter();
const { cities, districtsFor } = useLocationOptions();
const { extractCoordinatesFromUrl, buildSearchUrl } = useGoogleMaps();

const company = ref<Company | null>(null);
const loadingCompany = ref(false);
const isSavingCompany = ref(false);
const companyError = ref("");

const companyForm = reactive({
  name: "",
  description: "",
  register_id: "",
  city: "",
  district: "",
  location_x: "",
  location_y: "",
  maps_url: "",
});

const isRecruiter = computed(() => user.value?.role === "recruiter");
const hasCompany = computed(() => Boolean(user.value?.company_id));

function displayName() {
  return user.value?.name || user.value?.full_name || "Your profile";
}

function fillCompanyForm(value: Company | null) {
  companyForm.name = value?.name || "";
  companyForm.description = value?.description || "";
  companyForm.register_id = value?.register_id || "";
  companyForm.city = value?.city || "";
  companyForm.district = value?.district || "";
  companyForm.location_x =
    typeof value?.location_x === "number" ? String(value.location_x) : "";
  companyForm.location_y =
    typeof value?.location_y === "number" ? String(value.location_y) : "";
  companyForm.maps_url = "";
}

function buildCompanyPayload(): SaveCompanyPayload {
  return {
    name: companyForm.name.trim(),
    description: companyForm.description.trim(),
    register_id: companyForm.register_id.trim(),
    city: companyForm.city.trim(),
    district: companyForm.district.trim(),
    location_x: companyForm.location_x.trim()
      ? Number(companyForm.location_x.trim())
      : undefined,
    location_y: companyForm.location_y.trim()
      ? Number(companyForm.location_y.trim())
      : undefined,
  };
}

const districtOptions = computed(() => districtsFor(companyForm.city));
const companyMapsUrl = computed(() =>
  buildSearchUrl(
    [companyForm.district, companyForm.city, companyForm.name]
      .filter(Boolean)
      .join(", "),
    companyForm.location_x ? Number(companyForm.location_x) : null,
    companyForm.location_y ? Number(companyForm.location_y) : null,
  ),
);

function applyCompanyMapsUrl() {
  const coordinates = extractCoordinatesFromUrl(companyForm.maps_url);
  if (!coordinates) {
    toast.error("Paste a valid Google Maps share link.");
    return;
  }

  companyForm.location_x = String(coordinates.lat);
  companyForm.location_y = String(coordinates.lng);
  toast.success("Company coordinates imported from Google Maps.");
}

async function loadCompany() {
  if (!isRecruiter.value || !hasCompany.value) {
    company.value = null;
    fillCompanyForm(null);
    return;
  }

  loadingCompany.value = true;
  try {
    company.value = await companyAPI.get();
    fillCompanyForm(company.value);
  } catch {
    company.value = null;
  } finally {
    loadingCompany.value = false;
  }
}

async function saveCompany() {
  companyError.value = "";
  isSavingCompany.value = true;
  const wasHasCompany = hasCompany.value;

  try {
    const savedCompany = await companyAPI.save(buildCompanyPayload());
    company.value = savedCompany;
    fillCompanyForm(savedCompany);
    await me();
    toast.success(
      wasHasCompany
        ? "Company profile updated."
        : "Company created. You can now manage recruiter jobs.",
    );
  } catch (error: any) {
    companyError.value =
      error?.data?.message || error?.message || "Unable to save the company.";
    toast.error(companyError.value);
  } finally {
    isSavingCompany.value = false;
  }
}

function openJobs() {
  if (isRecruiter.value && !hasCompany.value) {
    toast.warning("You need to add your company before opening recruiter jobs.");
    return;
  }

  router.push("/jobs");
}

await loadCompany();
</script>

<template>
  <div class="space-y-6">
    <section
      class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm"
    >
      <p class="text-sm font-medium text-muted-foreground">Profile</p>
      <h1 class="mt-2 text-3xl font-semibold tracking-tight">
        {{ displayName() }}
      </h1>
      <p class="mt-2 text-sm leading-6 text-muted-foreground">
        This page keeps your current user model intact and presents it with the
        wireframe-inspired card layout.
      </p>
    </section>

    <section
      v-if="isRecruiter && !hasCompany"
      class="rounded-3xl border border-dashed border-border bg-muted/20 px-6 py-5"
    >
      <p class="text-sm font-semibold">Finish recruiter setup</p>
      <p class="mt-2 text-sm leading-6 text-muted-foreground">
        Your recruiter account is ready, but you need to create a company before posting jobs or opening the recruiter workspace.
      </p>
    </section>

    <section class="grid gap-4 lg:grid-cols-[1.3fr_0.9fr]">
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardTitle>Account details</CardTitle>
          <CardDescription>Synced from the authenticated backend user.</CardDescription>
        </CardHeader>
        <CardContent class="grid gap-4 sm:grid-cols-2">
          <div class="rounded-2xl border border-border bg-muted/30 p-4">
            <p class="text-xs uppercase tracking-[0.18em] text-muted-foreground">
              Email
            </p>
            <p class="mt-2 text-sm font-medium">{{ user?.email || "-" }}</p>
          </div>
          <div class="rounded-2xl border border-border bg-muted/30 p-4">
            <p class="text-xs uppercase tracking-[0.18em] text-muted-foreground">
              Role
            </p>
            <p class="mt-2 text-sm font-medium capitalize">{{ user?.role || "-" }}</p>
          </div>
          <div class="rounded-2xl border border-border bg-muted/30 p-4">
            <p class="text-xs uppercase tracking-[0.18em] text-muted-foreground">
              Company
            </p>
            <p class="mt-2 text-sm font-medium">{{ company?.name || user?.company_name || "-" }}</p>
          </div>
          <div class="rounded-2xl border border-border bg-muted/30 p-4">
            <p class="text-xs uppercase tracking-[0.18em] text-muted-foreground">
              Position
            </p>
            <p class="mt-2 text-sm font-medium">{{ user?.position || "-" }}</p>
          </div>
        </CardContent>
      </Card>

      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardTitle>Next</CardTitle>
          <CardDescription>Quick path based on your role.</CardDescription>
        </CardHeader>
        <CardContent class="space-y-3">
          <button
            type="button"
            class="block w-full rounded-2xl border border-border bg-muted/30 p-4 text-left text-sm font-medium transition hover:bg-muted/50"
            :class="isRecruiter && !hasCompany ? 'opacity-70' : ''"
            @click="openJobs"
          >
            {{ user?.role === "recruiter" ? "Manage job posts" : "Browse job matches" }}
          </button>
          <NuxtLink
            to="/settings"
            class="block rounded-2xl border border-border bg-muted/30 p-4 text-sm font-medium transition hover:bg-muted/50"
          >
            Review account settings
          </NuxtLink>
        </CardContent>
      </Card>
    </section>

    <Card v-if="isRecruiter" class="rounded-3xl border-border shadow-sm">
      <CardHeader>
        <CardTitle>{{ hasCompany ? "Company profile" : "Create your company" }}</CardTitle>
        <CardDescription>
          {{ hasCompany ? "Keep your recruiter company details up to date." : "This is required before you can create job posts." }}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div
          v-if="companyError"
          class="mb-4 rounded-2xl border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm text-destructive"
        >
          {{ companyError }}
        </div>

        <div v-if="loadingCompany" class="text-sm text-muted-foreground">
          Loading company details...
        </div>

        <form v-else class="grid gap-4 sm:grid-cols-2" @submit.prevent="saveCompany">
          <div class="space-y-2 sm:col-span-2">
            <Label for="company-name">Company name</Label>
            <Input id="company-name" v-model="companyForm.name" placeholder="Your company" />
          </div>

          <div class="space-y-2 sm:col-span-2">
            <Label for="company-description">Description</Label>
            <Textarea
              id="company-description"
              v-model="companyForm.description"
              rows="4"
              placeholder="Tell recruiters and applicants what your company does."
            />
          </div>

          <div class="space-y-2">
            <Label for="company-register-id">Register ID</Label>
            <Input id="company-register-id" v-model="companyForm.register_id" placeholder="A1234567" />
          </div>

          <div class="space-y-2">
            <Label for="company-city">Hot / Aimag</Label>
            <Select v-model="companyForm.city">
              <SelectTrigger id="company-city" class="w-full">
                <SelectValue placeholder="Select city or aimag" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="city in cities" :key="city" :value="city">
                  {{ city }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div class="space-y-2">
            <Label for="company-district">Duureg / Sum</Label>
            <Select v-model="companyForm.district" :disabled="!companyForm.city">
              <SelectTrigger id="company-district" class="w-full">
                <SelectValue placeholder="Select district or sum" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem
                  v-for="district in districtOptions"
                  :key="district"
                  :value="district"
                >
                  {{ district }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div class="sm:col-span-2 rounded-2xl border border-border bg-muted/20 p-4">
            <div class="flex items-start justify-between gap-3">
              <div>
                <p class="text-sm font-medium">Google Maps</p>
                <p class="mt-1 text-xs text-muted-foreground">
                  Open Google Maps, pick the company location, then paste the shared link to import coordinates.
                </p>
              </div>
              <Button type="button" variant="outline" as-child>
                <a :href="companyMapsUrl" target="_blank" rel="noreferrer">
                  Open Maps
                </a>
              </Button>
            </div>

            <div class="mt-4 grid gap-4 sm:grid-cols-[1fr_auto]">
              <Input
                v-model="companyForm.maps_url"
                placeholder="Paste Google Maps share link"
              />
              <Button type="button" variant="outline" @click="applyCompanyMapsUrl">
                Use Link
              </Button>
            </div>

            <div class="mt-4 grid gap-4 sm:grid-cols-2">
              <div class="space-y-2">
                <Label for="company-location-x">Latitude</Label>
                <Input
                  id="company-location-x"
                  v-model="companyForm.location_x"
                  placeholder="47.9184"
                />
              </div>
              <div class="space-y-2">
                <Label for="company-location-y">Longitude</Label>
                <Input
                  id="company-location-y"
                  v-model="companyForm.location_y"
                  placeholder="106.9177"
                />
              </div>
            </div>
          </div>

          <div class="sm:col-span-2 flex justify-end">
            <Button type="submit" :disabled="isSavingCompany">
              {{ isSavingCompany ? "Saving..." : hasCompany ? "Update company" : "Create company" }}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  </div>
</template>
