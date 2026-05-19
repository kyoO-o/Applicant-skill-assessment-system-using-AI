<script setup lang="ts">
import { Plus, X } from "lucide-vue-next";
import { JobStatus } from "../../composables/types";
import type { Company } from "../../composables/types";
import type { SaveJobPayload } from "../../composables/types/payload";
import { useForm } from "vee-validate";
import { toTypedSchema } from "@vee-validate/zod";
import { jobSchema } from "~/utils/schemas";

const { cities, districtsFor } = useLocationOptions();

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

const overrideContact = ref(!!props.initialTitle);
const overrideLocation = ref(!!props.initialTitle);

const { validate, setValues, errors } = useForm({
  validationSchema: toTypedSchema(jobSchema),
});

// Split initial bonuses into company-benefit-selected vs manually typed
const companyBenefitSet = new Set(props.company?.benefits?.map(b => b.description) ?? []);
const selectedBenefits = ref<string[]>(
  (props.initialBonuses ?? []).filter(b => companyBenefitSet.has(b)),
);
const initialManualBonuses = (props.initialBonuses ?? []).filter(b => !companyBenefitSet.has(b));

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
  min_salary: props.initialMinSalary ?? "",
  max_salary: props.initialMaxSalary ?? "",
  additional_info: props.initialAdditionalInfo ?? "",
  duties: props.initialDuties?.length ? [...props.initialDuties] : [""],
  requirements: props.initialRequirements?.length ? [...props.initialRequirements] : [""],
  skills: props.initialSkills?.length ? [...props.initialSkills] : [""],
  bonuses: initialManualBonuses.length ? [...initialManualBonuses] : [""],
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

function addItem(list: string[]) { list.push(""); }

function removeItem(list: string[], index: number) {
  if (list.length === 1) { list[0] = ""; return; }
  list.splice(index, 1);
}

function toggleBenefit(desc: string) {
  const idx = selectedBenefits.value.indexOf(desc);
  if (idx === -1) selectedBenefits.value.push(desc);
  else selectedBenefits.value.splice(idx, 1);
}

function isBenefitAdded(desc: string) {
  return selectedBenefits.value.includes(desc);
}

async function submitForm() {
  setValues({
    title: form.title,
    contact_info: form.contact_info,
    type: form.type,
    level: form.level,
  });
  const { valid } = await validate();
  if (!valid) return;

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
    bonuses: [...selectedBenefits.value, ...toList(form.bonuses)],
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
      <Input id="job-title" v-model="form.title" placeholder="Frontend Developer" />
      <p v-if="errors.title" class="text-sm text-destructive">{{ errors.title }}</p>
    </div>

    <!-- Type + Level + Status -->
    <div class="grid gap-4 sm:grid-cols-3">
      <div class="space-y-2">
        <Label for="job-type">Хөдөлмөрийн гэрээний хэлбэр</Label>
        <Input id="job-type" v-model="form.type" placeholder="Full-time" />
        <p v-if="errors.type" class="text-sm text-destructive">{{ errors.type }}</p>
      </div>
      <div class="space-y-2">
        <Label for="job-level">Мэргэжлийн түвшин</Label>
        <Input id="job-level" v-model="form.level" placeholder="Mid-level" />
        <p v-if="errors.level" class="text-sm text-destructive">{{ errors.level }}</p>
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
          <Checkbox id="override-contact-cb" v-model:checked="overrideContact" />
          <label for="override-contact-cb" class="cursor-pointer select-none">Өөр мэдээлэл оруулах</label>
        </div>
      </div>
      <Input
        v-model="form.contact_info"
        placeholder="hr@company.mn | +976 99000000"
        :disabled="!overrideContact && !!company"
      />
      <p v-if="errors.contact_info" class="text-sm text-destructive">{{ errors.contact_info }}</p>
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
          <Checkbox id="override-location-cb" v-model:checked="overrideLocation" />
          <label for="override-location-cb" class="cursor-pointer select-none">Өөр байршил оруулах</label>
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

      <div>
        <Label class="mb-2 block">Газрын зураг дээр байршил сонгох</Label>
        <LocationSearch
          v-model:model-x="form.location_x"
          v-model:model-y="form.location_y"
          :disabled="!overrideLocation && !!company"
        />
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
            <Textarea v-model="form.duties[i]" placeholder="Нэг үүрэг бичнэ үү" rows="3" class="min-h-[80px] resize-y" />
            <Button v-if="i > 0" type="button" variant="outline" size="icon" class="mt-1 shrink-0" @click="removeItem(form.duties, i)">
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
            <Textarea v-model="form.requirements[i]" placeholder="Нэг шаардлага бичнэ үү" rows="3" class="min-h-[80px] resize-y" />
            <Button v-if="i > 0" type="button" variant="outline" size="icon" class="mt-1 shrink-0" @click="removeItem(form.requirements, i)">
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
            <Textarea v-model="form.skills[i]" placeholder="Нэг ур чадвар бичнэ үү" rows="2" class="min-h-[64px] resize-y" />
            <Button v-if="i > 0" type="button" variant="outline" size="icon" class="mt-1 shrink-0" @click="removeItem(form.skills, i)">
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

        <!-- Company benefits: click to toggle, no textarea -->
        <div
          v-if="company?.benefits?.length"
          class="rounded-2xl border border-dashed border-border bg-muted/20 px-3 py-2"
        >
          <p class="mb-2 text-xs text-muted-foreground">Компанийн давуу талаас сонгох:</p>
          <div class="flex flex-wrap gap-2">
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
              @click="toggleBenefit(benefit.description)"
            >
              {{ benefit.description }}
            </button>
          </div>
        </div>

        <!-- Manual input textareas only -->
        <div class="space-y-3">
          <div v-for="(_, i) in form.bonuses" :key="`bonus-${i}`" class="flex items-start gap-2">
            <Textarea v-model="form.bonuses[i]" placeholder="Нэмэлт давуу тал гараар бичнэ үү" rows="2" class="min-h-[64px] resize-y" />
            <Button v-if="i > 0" type="button" variant="outline" size="icon" class="mt-1 shrink-0" @click="removeItem(form.bonuses, i)">
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
