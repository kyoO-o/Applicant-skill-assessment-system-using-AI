<script setup lang="ts">
import { toast } from "vue-sonner";
import { Plus, X } from "lucide-vue-next";
import { JobStatus } from "../../composables/types";
import type { Company } from "../../composables/types";
import type { SaveJobPayload } from "../../composables/types/payload";

const { cities, districtsFor } = useLocationOptions();
const { extractCoordinatesFromUrl, buildSearchUrl } = useGoogleMaps();

const props = defineProps<{
  company?: Company | null;
  isSubmitting: boolean;
  errorMessage: string;
  initialTitle?: string;
  initialContactInfo?: string;
  initialType?: string;
  initialLevel?: string;
  initialStatus?: JobStatus;
  initialCity?: string;
  initialDistrict?: string;
  initialLocationX?: string;
  initialLocationY?: string;
  initialMinSalary?: string;
  initialMaxSalary?: string;
  initialAdditionalInfo?: string;
  initialDuties?: string[];
  initialRequirements?: string[];
  initialSkills?: string[];
  initialBonuses?: string[];
}>();

const emit = defineEmits<{
  submit: [payload: SaveJobPayload];
}>();

// When editing an existing job (initialTitle is set), start in override mode so fields are editable.
// For new jobs (no initialTitle), default to using company data (disabled).
const overrideContact = ref(!!props.initialTitle);
const overrideLocation = ref(!!props.initialTitle);

const form = reactive({
  title: props.initialTitle ?? "",
  contact_info: props.initialContactInfo ?? (props.company?.contact_info ?? ""),
  type: props.initialType ?? "Full-time",
  level: props.initialLevel ?? "Mid-level",
  status: props.initialStatus ?? JobStatus.Draft,
  city: props.initialCity ?? (props.company?.city ?? ""),
  district: props.initialDistrict ?? (props.company?.district ?? ""),
  location_x: props.initialLocationX ?? (props.company?.location_x ? String(props.company.location_x) : ""),
  location_y: props.initialLocationY ?? (props.company?.location_y ? String(props.company.location_y) : ""),
  maps_url: "",
  min_salary: props.initialMinSalary ?? "",
  max_salary: props.initialMaxSalary ?? "",
  additional_info: props.initialAdditionalInfo ?? "",
  duties: props.initialDuties?.length ? [...props.initialDuties] : [""],
  requirements: props.initialRequirements?.length ? [...props.initialRequirements] : [""],
  skills: props.initialSkills?.length ? [...props.initialSkills] : [""],
  bonuses: props.initialBonuses?.length ? [...props.initialBonuses] : [""],
});

watch(overrideContact, (val) => {
  if (!val && props.company) {
    form.contact_info = props.company.contact_info ?? "";
  }
});

watch(overrideLocation, (val) => {
  if (!val && props.company) {
    form.city = props.company.city ?? "";
    form.district = props.company.district ?? "";
    form.location_x = props.company.location_x ? String(props.company.location_x) : "";
    form.location_y = props.company.location_y ? String(props.company.location_y) : "";
  }
});

const districtOptions = computed(() => districtsFor(form.city));
const googleMapsUrl = computed(() =>
  buildSearchUrl(
    [form.district, form.city, form.title].filter(Boolean).join(", "),
    form.location_x ? Number(form.location_x) : null,
    form.location_y ? Number(form.location_y) : null,
  ),
);

function applyMapsUrl() {
  const coordinates = extractCoordinatesFromUrl(form.maps_url);
  if (!coordinates) {
    toast.error("Зөв Google Maps хуваалцах холбоосыг буулгана уу.");
    return;
  }
  form.location_x = String(coordinates.lat);
  form.location_y = String(coordinates.lng);
  toast.success("Координат Google Maps-аас импортлогдлоо.");
}

function addItem(list: string[]) { list.push(""); }

function removeItem(list: string[], index: number) {
  if (list.length === 1) { list[0] = ""; return; }
  list.splice(index, 1);
}

function addBenefit(desc: string) {
  const existing = form.bonuses.filter(Boolean);
  if (!existing.includes(desc)) form.bonuses = [...existing, desc];
}

function isBenefitAdded(desc: string) {
  return form.bonuses.filter(Boolean).includes(desc);
}

function submitForm() {
  const city = form.city.trim();
  const district = form.district.trim();
  const location = [district, city].filter(Boolean).join(", ");
  const toList = (value: string[]) => value.map((item) => item.trim()).filter(Boolean);

  emit("submit", {
    title: form.title.trim(),
    location,
    additional_info: form.additional_info.trim(),
    contact_info: form.contact_info.trim(),
    type: form.type.trim(),
    level: form.level.trim(),
    city: city || undefined,
    district: district || undefined,
    location_x: form.location_x.trim() ? Number(form.location_x.trim()) : undefined,
    location_y: form.location_y.trim() ? Number(form.location_y.trim()) : undefined,
    min_salary: form.min_salary ? Number(form.min_salary) : 0,
    max_salary: form.max_salary ? Number(form.max_salary) : 0,
    status: form.status,
    duties: toList(form.duties),
    requirements: toList(form.requirements),
    skills: toList(form.skills),
    bonuses: toList(form.bonuses),
  });
}
</script>

<template>
  <form class="space-y-6" @submit.prevent="submitForm">
    <div
      v-if="errorMessage"
      class="rounded-2xl border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm text-destructive"
    >
      {{ errorMessage }}
    </div>

    <!-- Title -->
    <div class="space-y-2">
      <Label for="job-title">Ажлын байрны нэр</Label>
      <Input id="job-title" v-model="form.title" placeholder="Frontend Developer" required />
    </div>

    <!-- Type + Level + Status -->
    <div class="grid gap-4 sm:grid-cols-3">
      <div class="space-y-2">
        <Label for="job-type">Хөдөлмөрийн гэрээний хэлбэр</Label>
        <Input id="job-type" v-model="form.type" placeholder="Full-time" />
      </div>
      <div class="space-y-2">
        <Label for="job-level">Мэргэжлийн түвшин</Label>
        <Input id="job-level" v-model="form.level" placeholder="Mid-level" />
      </div>
      <div class="space-y-2">
        <Label for="job-status">Төлөв</Label>
        <Select v-model="form.status">
          <SelectTrigger id="job-status" class="w-full">
            <SelectValue placeholder="Төлөв сонгох" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem :value="JobStatus.Draft">Ноорог</SelectItem>
            <SelectItem :value="JobStatus.Posted">Нийтлэгдсэн</SelectItem>
            <SelectItem :value="JobStatus.Closed">Хаагдсан</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>

    <!-- Contact info -->
    <div class="rounded-2xl border border-border bg-muted/20 p-4 space-y-3">
      <div class="flex items-center justify-between gap-3">
        <div>
          <p class="text-sm font-medium">Холбоо барих мэдээлэл</p>
          <p v-if="company && !overrideContact" class="mt-0.5 text-xs text-green-600">
            ✓ Компанийн мэдээлэл ашиглагдлаа
          </p>
        </div>
        <div v-if="company" class="flex items-center gap-2 text-xs text-muted-foreground">
          <Checkbox v-model:checked="overrideContact" />
          <span class="cursor-pointer select-none" @click="overrideContact = !overrideContact">Өөр мэдээлэл оруулах</span>
        </div>
      </div>
      <Input
        v-model="form.contact_info"
        placeholder="hr@company.mn | +976 99000000"
        :disabled="!overrideContact && !!company"
      />
    </div>

    <!-- Location -->
    <div class="rounded-2xl border border-border bg-muted/20 p-4 space-y-4">
      <div class="flex items-center justify-between gap-3">
        <div>
          <p class="text-sm font-medium">Байршил</p>
          <p v-if="company && !overrideLocation" class="mt-0.5 text-xs text-green-600">
            ✓ Компанийн байршил ашиглагдлаа
          </p>
        </div>
        <div v-if="company" class="flex items-center gap-2 text-xs text-muted-foreground">
          <Checkbox v-model:checked="overrideLocation" />
          <span class="cursor-pointer select-none" @click="overrideLocation = !overrideLocation">Өөр байршил оруулах</span>
        </div>
      </div>

      <div class="grid gap-4 sm:grid-cols-2">
        <div class="space-y-2">
          <Label for="job-city">Хот / Аймаг</Label>
          <Select v-model="form.city" :disabled="!overrideLocation && !!company">
            <SelectTrigger id="job-city" class="w-full">
              <SelectValue placeholder="Хот эсвэл аймаг сонгох" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem v-for="city in cities" :key="city" :value="city">{{ city }}</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div class="space-y-2">
          <Label for="job-district">Дүүрэг / Сум</Label>
          <Select v-model="form.district" :disabled="(!overrideLocation && !!company) || !form.city">
            <SelectTrigger id="job-district" class="w-full">
              <SelectValue placeholder="Дүүрэг эсвэл сум сонгох" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem v-for="d in districtOptions" :key="d" :value="d">{{ d }}</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <div v-if="overrideLocation || !company">
        <div class="flex items-start justify-between gap-3">
          <div>
            <p class="text-sm font-medium">Google Maps координат</p>
            <p class="mt-1 text-xs text-muted-foreground">Хуваалцах холбоосыг буулган координатыг импортлоно уу.</p>
          </div>
          <Button type="button" variant="outline" size="sm" as-child>
            <a :href="googleMapsUrl" target="_blank" rel="noreferrer">Нээх</a>
          </Button>
        </div>
        <div class="mt-3 grid gap-3 sm:grid-cols-[1fr_auto]">
          <Input v-model="form.maps_url" placeholder="Google Maps холбоос буулгах" />
          <Button type="button" variant="outline" @click="applyMapsUrl">Ашиглах</Button>
        </div>
        <div class="mt-3 grid gap-3 sm:grid-cols-2">
          <div class="space-y-1.5">
            <Label for="job-location-x">Өргөрөг</Label>
            <Input id="job-location-x" v-model="form.location_x" placeholder="47.9184" />
          </div>
          <div class="space-y-1.5">
            <Label for="job-location-y">Уртраг</Label>
            <Input id="job-location-y" v-model="form.location_y" placeholder="106.9177" />
          </div>
        </div>
      </div>
    </div>

    <!-- Salary -->
    <div class="grid gap-4 sm:grid-cols-2">
      <div class="space-y-2">
        <Label for="job-min-salary">Хамгийн бага цалин</Label>
        <Input id="job-min-salary" v-model="form.min_salary" type="number" min="0" placeholder="1800000" />
      </div>
      <div class="space-y-2">
        <Label for="job-max-salary">Хамгийн их цалин</Label>
        <Input id="job-max-salary" v-model="form.max_salary" type="number" min="0" placeholder="3200000" />
      </div>
    </div>

    <!-- Description -->
    <div class="space-y-2">
      <Label for="job-description">Нэмэлт мэдээлэл</Label>
      <Textarea
        id="job-description"
        v-model="form.additional_info"
        rows="5"
        placeholder="Ажлын байрны тайлбар, хамрах хүрээ болон хүлээгдэж буй үр дүнг бичнэ үү."
      />
    </div>

    <!-- Duties + Requirements -->
    <div class="grid gap-6 lg:grid-cols-2">
      <div class="space-y-3">
        <div class="flex items-center justify-between">
          <Label>Үүрэг хариуцлага</Label>
          <Button type="button" variant="outline" size="sm" @click="addItem(form.duties)">
            <Plus class="mr-1 h-4 w-4" />Нэмэх
          </Button>
        </div>
        <div class="space-y-3">
          <div v-for="(_, i) in form.duties" :key="`duty-${i}`" class="flex items-start gap-2">
            <Textarea
              v-model="form.duties[i]"
              placeholder="Нэг үүрэг бичнэ үү"
              rows="3"
              class="min-h-[80px] resize-y"
            />
            <Button type="button" variant="outline" size="icon" class="mt-1 shrink-0" @click="removeItem(form.duties, i)">
              <X class="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>

      <div class="space-y-3">
        <div class="flex items-center justify-between">
          <Label>Шаардлагууд</Label>
          <Button type="button" variant="outline" size="sm" @click="addItem(form.requirements)">
            <Plus class="mr-1 h-4 w-4" />Нэмэх
          </Button>
        </div>
        <div class="space-y-3">
          <div v-for="(_, i) in form.requirements" :key="`req-${i}`" class="flex items-start gap-2">
            <Textarea
              v-model="form.requirements[i]"
              placeholder="Нэг шаардлага бичнэ үү"
              rows="3"
              class="min-h-[80px] resize-y"
            />
            <Button type="button" variant="outline" size="icon" class="mt-1 shrink-0" @click="removeItem(form.requirements, i)">
              <X class="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>
    </div>

    <!-- Skills + Bonuses -->
    <div class="grid gap-6 lg:grid-cols-2">
      <div class="space-y-3">
        <div class="flex items-center justify-between">
          <Label>Ур чадвар</Label>
          <Button type="button" variant="outline" size="sm" @click="addItem(form.skills)">
            <Plus class="mr-1 h-4 w-4" />Нэмэх
          </Button>
        </div>
        <div class="space-y-3">
          <div v-for="(_, i) in form.skills" :key="`skill-${i}`" class="flex items-start gap-2">
            <Textarea
              v-model="form.skills[i]"
              placeholder="Нэг ур чадвар бичнэ үү"
              rows="2"
              class="min-h-[64px] resize-y"
            />
            <Button type="button" variant="outline" size="icon" class="mt-1 shrink-0" @click="removeItem(form.skills, i)">
              <X class="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>

      <div class="space-y-3">
        <div class="flex items-center justify-between">
          <Label>Нэмэлт давуу тал / Урамшуулал</Label>
          <Button type="button" variant="outline" size="sm" @click="addItem(form.bonuses)">
            <Plus class="mr-1 h-4 w-4" />Нэмэх
          </Button>
        </div>
        <!-- Company benefits quick-add -->
        <div
          v-if="company?.benefits?.length"
          class="flex flex-wrap gap-2 rounded-2xl border border-dashed border-border bg-muted/20 px-3 py-2"
        >
          <p class="w-full text-xs text-muted-foreground">Компанийн давуу талаас нэмэх:</p>
          <button
            v-for="benefit in company.benefits"
            :key="benefit.id"
            type="button"
            :class="[
              'rounded-full border px-3 py-1 text-xs transition',
              isBenefitAdded(benefit.description)
                ? 'border-primary bg-primary text-primary-foreground'
                : 'border-border bg-background hover:bg-muted',
            ]"
            @click="addBenefit(benefit.description)"
          >
            {{ benefit.description }}
          </button>
        </div>
        <div class="space-y-3">
          <div v-for="(_, i) in form.bonuses" :key="`bonus-${i}`" class="flex items-start gap-2">
            <Textarea
              v-model="form.bonuses[i]"
              placeholder="Нэг давуу тал бичнэ үү"
              rows="2"
              class="min-h-[64px] resize-y"
            />
            <Button type="button" variant="outline" size="icon" class="mt-1 shrink-0" @click="removeItem(form.bonuses, i)">
              <X class="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>
    </div>

    <div class="flex justify-end gap-3 border-t border-border pt-6">
      <slot name="actions">
        <Button type="submit" :disabled="isSubmitting">
          {{ isSubmitting ? "Хадгалж байна..." : "Хадгалах" }}
        </Button>
      </slot>
    </div>
  </form>
</template>
