<script setup lang="ts">
import { toast } from "vue-sonner";
import { Plus, X } from "lucide-vue-next";

import { JobStatus } from "../../composables/types";

const { cities, districtsFor } = useLocationOptions();
const { extractCoordinatesFromUrl, buildSearchUrl } = useGoogleMaps();

const props = defineProps<{
  open: boolean;
  isEditing: boolean;
  isSubmitting: boolean;
  errorMessage: string;
  form: {
    title: string;
    contact_info: string;
    type: string;
    level: string;
    status: JobStatus;
    city: string;
    district: string;
    location_x: string;
    location_y: string;
    maps_url: string;
    min_salary: string;
    max_salary: string;
    additional_info: string;
    duties: string[];
    requirements: string[];
    skills: string[];
    bonuses: string[];
  };
}>();

const emit = defineEmits<{
  "update:open": [value: boolean];
  submit: [];
}>();

const dialogOpen = computed({
  get: () => props.open,
  set: (value: boolean) => emit("update:open", value),
});

function closeDialog() {
  emit("update:open", false);
}

function submitForm() {
  emit("submit");
}

const districtOptions = computed(() => districtsFor(props.form.city));
const googleMapsUrl = computed(() =>
  buildSearchUrl(
    [props.form.district, props.form.city, props.form.title].filter(Boolean).join(", "),
    props.form.location_x ? Number(props.form.location_x) : null,
    props.form.location_y ? Number(props.form.location_y) : null,
  ),
);

function applyMapsUrl() {
  const coordinates = extractCoordinatesFromUrl(props.form.maps_url);
  if (!coordinates) {
    toast.error("Paste a valid Google Maps share link.");
    return;
  }

  props.form.location_x = String(coordinates.lat);
  props.form.location_y = String(coordinates.lng);
  toast.success("Coordinates imported from Google Maps.");
}

function addItem(list: string[]) {
  list.push("");
}

function removeItem(list: string[], index: number) {
  if (list.length === 1) {
    list[0] = "";
    return;
  }

  list.splice(index, 1);
}
</script>

<template>
  <Dialog v-model:open="dialogOpen">
    <DialogContent class="max-h-[90vh] w-[min(96vw,1100px)] max-w-none overflow-hidden rounded-3xl p-0">
      <div class="flex max-h-[90vh] flex-col">
      <DialogHeader class="border-b border-border px-6 py-5">
        <DialogTitle>
          {{ isEditing ? "Edit job post" : "Create a new job post" }}
        </DialogTitle>
        <DialogDescription>
          Keep the structure from the wireframe and save directly to the
          backend.
        </DialogDescription>
      </DialogHeader>

      <form class="flex min-h-0 flex-1 flex-col" @submit.prevent="submitForm">
        <div class="min-h-0 flex-1 overflow-y-auto px-6 py-5">
        <div class="grid gap-4">
        <div
          v-if="errorMessage"
          class="rounded-2xl border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm text-destructive"
        >
          {{ errorMessage }}
        </div>

        <div class="grid gap-4 sm:grid-cols-2">
          <div class="space-y-2">
            <Label for="job-title">Job title</Label>
            <Input
              id="job-title"
              v-model="form.title"
              placeholder="Frontend Developer"
            />
          </div>
          <div class="space-y-2">
            <Label for="job-contact">Contact info</Label>
            <Input
              id="job-contact"
              v-model="form.contact_info"
              placeholder="hr@company.mn or +976..."
            />
          </div>
        </div>

        <div class="grid gap-4 sm:grid-cols-3">
          <div class="space-y-2">
            <Label for="job-type">Employment type</Label>
            <Input
              id="job-type"
              v-model="form.type"
              placeholder="Full-time"
            />
          </div>
          <div class="space-y-2">
            <Label for="job-level">Seniority</Label>
            <Input
              id="job-level"
              v-model="form.level"
              placeholder="Mid-level"
            />
          </div>
          <div class="space-y-2">
            <Label for="job-status">Status</Label>
            <Select v-model="form.status">
              <SelectTrigger id="job-status" class="w-full">
                <SelectValue placeholder="Select status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem :value="JobStatus.Draft">Draft</SelectItem>
                <SelectItem :value="JobStatus.Posted">Posted</SelectItem>
                <SelectItem :value="JobStatus.Closed">Closed</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        <div class="grid gap-4 sm:grid-cols-2">
          <div class="space-y-2">
            <Label for="job-city">Hot / Aimag</Label>
            <Select v-model="form.city">
              <SelectTrigger id="job-city" class="w-full">
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
            <Label for="job-district">Duureg / Sum</Label>
            <Select v-model="form.district" :disabled="!form.city">
              <SelectTrigger id="job-district" class="w-full">
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
        </div>

        <div class="rounded-2xl border border-border bg-muted/20 p-4">
          <div class="flex items-start justify-between gap-3">
            <div>
              <p class="text-sm font-medium">Google Maps</p>
              <p class="mt-1 text-xs text-muted-foreground">
                Open Google Maps, choose the location, then paste the shared link to import coordinates.
              </p>
            </div>
            <Button type="button" variant="outline" as-child>
              <a :href="googleMapsUrl" target="_blank" rel="noreferrer">
                Open Maps
              </a>
            </Button>
          </div>

          <div class="mt-4 grid gap-4 sm:grid-cols-[1fr_auto]">
            <Input
              v-model="form.maps_url"
              placeholder="Paste Google Maps share link"
            />
            <Button type="button" variant="outline" @click="applyMapsUrl">
              Use Link
            </Button>
          </div>

          <div class="mt-4 grid gap-4 sm:grid-cols-2">
            <div class="space-y-2">
              <Label for="job-location-x">Latitude</Label>
              <Input
                id="job-location-x"
                v-model="form.location_x"
                placeholder="47.9184"
              />
            </div>
            <div class="space-y-2">
              <Label for="job-location-y">Longitude</Label>
              <Input
                id="job-location-y"
                v-model="form.location_y"
                placeholder="106.9177"
              />
            </div>
          </div>
        </div>

        <div class="grid gap-4 sm:grid-cols-2">
          <div class="space-y-2">
            <Label for="job-min-salary">Min salary</Label>
            <Input
              id="job-min-salary"
              v-model="form.min_salary"
              type="number"
              min="0"
              placeholder="1800000"
            />
          </div>
          <div class="space-y-2">
            <Label for="job-max-salary">Max salary</Label>
            <Input
              id="job-max-salary"
              v-model="form.max_salary"
              type="number"
              min="0"
              placeholder="3200000"
            />
          </div>
        </div>

        <div class="space-y-2">
          <Label for="job-description">Additional info</Label>
          <Textarea
            id="job-description"
            v-model="form.additional_info"
            rows="5"
            placeholder="Describe the role, scope, and outcomes."
          />
        </div>

        <div class="grid gap-4 lg:grid-cols-2">
          <div class="space-y-2">
            <div class="flex items-center justify-between gap-3">
              <Label>Duties</Label>
              <Button type="button" variant="outline" size="sm" @click="addItem(form.duties)">
                <Plus class="mr-1 h-4 w-4" />
                Add
              </Button>
            </div>
            <div class="space-y-2">
              <div v-for="(_, index) in form.duties" :key="`duty-${index}`" class="flex items-center gap-2">
                <Input v-model="form.duties[index]" placeholder="Write one duty" />
                <Button type="button" variant="outline" size="icon" @click="removeItem(form.duties, index)">
                  <X class="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>
          <div class="space-y-2">
            <div class="flex items-center justify-between gap-3">
              <Label>Requirements</Label>
              <Button type="button" variant="outline" size="sm" @click="addItem(form.requirements)">
                <Plus class="mr-1 h-4 w-4" />
                Add
              </Button>
            </div>
            <div class="space-y-2">
              <div v-for="(_, index) in form.requirements" :key="`requirement-${index}`" class="flex items-center gap-2">
                <Input v-model="form.requirements[index]" placeholder="Write one requirement" />
                <Button type="button" variant="outline" size="icon" @click="removeItem(form.requirements, index)">
                  <X class="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>
        </div>

        <div class="grid gap-4 lg:grid-cols-2">
          <div class="space-y-2">
            <div class="flex items-center justify-between gap-3">
              <Label>Skills</Label>
              <Button type="button" variant="outline" size="sm" @click="addItem(form.skills)">
                <Plus class="mr-1 h-4 w-4" />
                Add
              </Button>
            </div>
            <div class="space-y-2">
              <div v-for="(_, index) in form.skills" :key="`skill-${index}`" class="flex items-center gap-2">
                <Input v-model="form.skills[index]" placeholder="Write one skill" />
                <Button type="button" variant="outline" size="icon" @click="removeItem(form.skills, index)">
                  <X class="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>
          <div class="space-y-2">
            <div class="flex items-center justify-between gap-3">
              <Label>Bonuses</Label>
              <Button type="button" variant="outline" size="sm" @click="addItem(form.bonuses)">
                <Plus class="mr-1 h-4 w-4" />
                Add
              </Button>
            </div>
            <div class="space-y-2">
              <div v-for="(_, index) in form.bonuses" :key="`bonus-${index}`" class="flex items-center gap-2">
                <Input v-model="form.bonuses[index]" placeholder="Write one bonus" />
                <Button type="button" variant="outline" size="icon" @click="removeItem(form.bonuses, index)">
                  <X class="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>
        </div>

        <div class="space-y-2">
          <Label>Location summary</Label>
          <p class="text-xs text-muted-foreground">
            Location summary: {{ [form.district, form.city].filter(Boolean).join(", ") || "Not selected" }}
          </p>
        </div>
        </div>
        </div>

        <DialogFooter class="border-t border-border px-6 py-4">
          <Button type="button" variant="outline" @click="closeDialog">
            Cancel
          </Button>
          <Button type="submit" :disabled="isSubmitting">
            {{
              isSubmitting
                ? "Saving..."
                : isEditing
                  ? "Update Job"
                  : "Create Job"
            }}
          </Button>
        </DialogFooter>
      </form>
      </div>
    </DialogContent>
  </Dialog>
</template>
