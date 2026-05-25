<script setup lang="ts">
import { Plus, X } from "lucide-vue-next";
import {
  JobStatus,
  JOB_EMPLOYMENT_TYPES,
  JOB_LEVELS,
  DEPARTMENT,
} from "../../composables/types";
import type { SaveJobPayload } from "../../composables/types/payload";
import { useForm } from "vee-validate";
import { toTypedSchema } from "@vee-validate/zod";
import { jobSchema } from "~/utils/schemas";

const { cities, districtsFor } = useLocationOptions();

const props = defineProps<{
  open: boolean;
  isEditing: boolean;
  isSubmitting: boolean;
  errorMessage: string;
  initialTitle?: string;
  initialContactInfo?: string;
  initialType?: string;
  initialLevel?: string;
  initialDepartment?: string;
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
  "update:open": [value: boolean];
  submit: [payload: SaveJobPayload];
}>();

const dialogOpen = computed({
  get: () => props.open,
  set: (value: boolean) => emit("update:open", value),
});

function closeDialog() {
  emit("update:open", false);
}

const { validate, setFieldValue, errors } = useForm({
  validationSchema: toTypedSchema(jobSchema),
  initialValues: {
    title: props.initialTitle ?? "",
    contact_info: props.initialContactInfo ?? "",
    type: props.initialType ?? "Бүтэн цагийн",
    level: props.initialLevel ?? "Мэргэжилтэн",
    department: props.initialDepartment ?? "",
  },
});

const form = reactive({
  title: props.initialTitle ?? "",
  contact_info: props.initialContactInfo ?? "",
  type: props.initialType ?? "Бүтэн цагийн",
  level: props.initialLevel ?? "Мэргэжилтэн",
  department: props.initialDepartment ?? "",
  status: props.initialStatus ?? JobStatus.Draft,
  city: props.initialCity ?? "",
  district: props.initialDistrict ?? "",
  location_x: props.initialLocationX ?? "",
  location_y: props.initialLocationY ?? "",
  min_salary: props.initialMinSalary ?? "",
  max_salary: props.initialMaxSalary ?? "",
  additional_info: props.initialAdditionalInfo ?? "",
  duties: props.initialDuties?.length ? [...props.initialDuties] : [""],
  requirements: props.initialRequirements?.length
    ? [...props.initialRequirements]
    : [""],
  skills: props.initialSkills?.length ? [...props.initialSkills] : [""],
  bonuses: props.initialBonuses?.length ? [...props.initialBonuses] : [""],
});

watch(() => form.title, (v) => setFieldValue("title", v));
watch(() => form.contact_info, (v) => setFieldValue("contact_info", v));
watch(() => form.type, (v) => setFieldValue("type", v));
watch(() => form.level, (v) => setFieldValue("level", v));
watch(() => form.department, (v) => setFieldValue("department", v));

const districtOptions = computed(() => districtsFor(form.city));

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

async function submitForm() {
  console.log("[SaveDialog] submitForm called", { form: { ...form } });
  const { valid, errors: validationErrors } = await validate();
  console.log("[SaveDialog] validation result", { valid, errors: validationErrors });
  if (!valid) return;

  const city = form.city.trim();
  const district = form.district.trim();
  const location = [district, city].filter(Boolean).join(", ");
  const toList = (value: string[]) =>
    value.map((item) => item.trim()).filter(Boolean);

  const payload: SaveJobPayload = {
    title: form.title.trim(),
    location,
    additional_info: form.additional_info.trim(),
    contact_info: form.contact_info.trim(),
    type: form.type.trim(),
    level: form.level.trim(),
    department: form.department,
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
  };
  console.log("[SaveDialog] emitting submit with payload", payload);
  emit("submit", payload);
}
</script>

<template>
  <Dialog v-model:open="dialogOpen">
    <DialogContent
      class="max-h-[90vh] w-[min(96vw,1100px)] max-w-none overflow-hidden rounded-3xl p-0"
    >
      <div class="flex max-h-[90vh] flex-col">
        <DialogHeader class="border-b border-border px-6 py-5">
          <DialogTitle>
            {{ isEditing ? "Ажлын байр засах" : "Шинэ ажлын байр нэмэх" }}
          </DialogTitle>
          <DialogDescription>
            Ажлын байрны мэдээллийг бөглөж хадгална уу.
          </DialogDescription>
        </DialogHeader>

        <form
          class="flex min-h-0 flex-1 flex-col"
          @submit.prevent="submitForm"
        >
          <div class="min-h-0 flex-1 overflow-y-auto px-6 py-5">
            <div class="grid gap-4">
              <div
                v-if="errorMessage"
                class="rounded-2xl border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm text-destructive"
              >
                {{ errorMessage }}
              </div>

              <!-- Title + Department -->
              <div class="grid gap-4 sm:grid-cols-2">
                <div class="space-y-2">
                  <Label for="job-title">Ажлын байрны нэр</Label>
                  <Input
                    id="job-title"
                    v-model="form.title"
                    placeholder="Frontend Developer"
                  />
                  <p v-if="errors.title" class="text-sm text-destructive">
                    {{ errors.title }}
                  </p>
                </div>
                <div class="space-y-2">
                  <Label for="job-department">Салбар / Чиглэл</Label>
                  <Select v-model="form.department">
                    <SelectTrigger id="job-department" class="w-full">
                      <SelectValue placeholder="Салбар сонгох" />
                    </SelectTrigger>
                    <SelectContent class="h-80">
                      <SelectItem v-for="d in DEPARTMENT" :key="d" :value="d">
                        {{ d }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <p v-if="errors.department" class="text-sm text-destructive">
                    {{ errors.department }}
                  </p>
                </div>
              </div>

              <!-- Contact info -->
              <div class="space-y-2">
                <Label for="job-contact">Холбоо барих мэдээлэл</Label>
                <Input
                  id="job-contact"
                  v-model="form.contact_info"
                  placeholder="hr@company.mn or +976..."
                />
                <p v-if="errors.contact_info" class="text-sm text-destructive">
                  {{ errors.contact_info }}
                </p>
              </div>

              <!-- Type + Level + Status -->
              <div class="grid gap-4 sm:grid-cols-3">
                <div class="space-y-2">
                  <Label for="job-type">Ажиллах цагийн төрөл</Label>
                  <Select v-model="form.type">
                    <SelectTrigger id="job-type" class="w-full">
                      <SelectValue placeholder="Хэлбэр сонгох" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem
                        v-for="t in JOB_EMPLOYMENT_TYPES"
                        :key="t"
                        :value="t"
                      >
                        {{ t }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <p v-if="errors.type" class="text-sm text-destructive">
                    {{ errors.type }}
                  </p>
                </div>
                <div class="space-y-2">
                  <Label for="job-level">Мэргэжлийн түвшин</Label>
                  <Select v-model="form.level">
                    <SelectTrigger id="job-level" class="w-full">
                      <SelectValue placeholder="Түвшин сонгох" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem
                        v-for="l in JOB_LEVELS"
                        :key="l"
                        :value="l"
                      >
                        {{ l }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <p v-if="errors.level" class="text-sm text-destructive">
                    {{ errors.level }}
                  </p>
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

              <!-- City + District -->
              <div class="grid gap-4 sm:grid-cols-2">
                <div class="space-y-2">
                  <Label for="job-city">Хот / Аймаг</Label>
                  <Select v-model="form.city">
                    <SelectTrigger id="job-city" class="w-full">
                      <SelectValue placeholder="Хот эсвэл аймаг сонгох" />
                    </SelectTrigger>
                    <SelectContent class="h-80">
                      <SelectItem
                        v-for="city in cities"
                        :key="city"
                        :value="city"
                      >
                        {{ city }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div class="space-y-2">
                  <Label for="job-district">Дүүрэг / Сум</Label>
                  <Select v-model="form.district" :disabled="!form.city">
                    <SelectTrigger id="job-district" class="w-full">
                      <SelectValue placeholder="Дүүрэг эсвэл сум сонгох" />
                    </SelectTrigger>
                    <SelectContent class="h-80">
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

              <!-- Map -->
              <div class="rounded-2xl border border-border bg-muted/20 p-4 space-y-2">
                <Label>Газрын зураг дээр байршил сонгох</Label>
                <LocationSearch
                  v-model:model-x="form.location_x"
                  v-model:model-y="form.location_y"
                />
              </div>

              <!-- Salary -->
              <div class="grid gap-4 sm:grid-cols-2">
                <div class="space-y-2">
                  <Label for="job-min-salary">Хамгийн бага цалин</Label>
                  <Input
                    id="job-min-salary"
                    v-model="form.min_salary"
                    type="number"
                    min="0"
                    placeholder="1800000"
                  />
                </div>
                <div class="space-y-2">
                  <Label for="job-max-salary">Хамгийн их цалин</Label>
                  <Input
                    id="job-max-salary"
                    v-model="form.max_salary"
                    type="number"
                    min="0"
                    placeholder="3200000"
                  />
                </div>
              </div>

              <!-- Additional info -->
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
              <div class="grid gap-4 lg:grid-cols-2">
                <div class="space-y-2">
                  <div class="flex items-center justify-between gap-3">
                    <Label>Үүрэг хариуцлага</Label>
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      @click="addItem(form.duties)"
                    >
                      <Plus class="mr-1 h-4 w-4" />Нэмэх
                    </Button>
                  </div>
                  <div class="space-y-2">
                    <div
                      v-for="(_, index) in form.duties"
                      :key="`duty-${index}`"
                      class="flex items-center gap-2"
                    >
                      <Input
                        v-model="form.duties[index]"
                        placeholder="Нэг үүрэг бичнэ үү"
                      />
                      <Button
                        v-if="index > 0"
                        type="button"
                        variant="outline"
                        size="icon"
                        @click="removeItem(form.duties, index)"
                      >
                        <X class="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </div>
                <div class="space-y-2">
                  <div class="flex items-center justify-between gap-3">
                    <Label>Шаардлагууд</Label>
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      @click="addItem(form.requirements)"
                    >
                      <Plus class="mr-1 h-4 w-4" />Нэмэх
                    </Button>
                  </div>
                  <div class="space-y-2">
                    <div
                      v-for="(_, index) in form.requirements"
                      :key="`requirement-${index}`"
                      class="flex items-center gap-2"
                    >
                      <Input
                        v-model="form.requirements[index]"
                        placeholder="Нэг шаардлага бичнэ үү"
                      />
                      <Button
                        v-if="index > 0"
                        type="button"
                        variant="outline"
                        size="icon"
                        @click="removeItem(form.requirements, index)"
                      >
                        <X class="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Skills + Bonuses -->
              <div class="grid gap-4 lg:grid-cols-2">
                <div class="space-y-2">
                  <div class="flex items-center justify-between gap-3">
                    <Label>Ур чадвар</Label>
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      @click="addItem(form.skills)"
                    >
                      <Plus class="mr-1 h-4 w-4" />Нэмэх
                    </Button>
                  </div>
                  <div class="space-y-2">
                    <div
                      v-for="(_, index) in form.skills"
                      :key="`skill-${index}`"
                      class="flex items-center gap-2"
                    >
                      <Input
                        v-model="form.skills[index]"
                        placeholder="Нэг ур чадвар бичнэ үү"
                      />
                      <Button
                        v-if="index > 0"
                        type="button"
                        variant="outline"
                        size="icon"
                        @click="removeItem(form.skills, index)"
                      >
                        <X class="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </div>
                <div class="space-y-2">
                  <div class="flex items-center justify-between gap-3">
                    <Label>Хөнгөлөлт / Урамшуулал</Label>
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      @click="addItem(form.bonuses)"
                    >
                      <Plus class="mr-1 h-4 w-4" />Нэмэх
                    </Button>
                  </div>
                  <div class="space-y-2">
                    <div
                      v-for="(_, index) in form.bonuses"
                      :key="`bonus-${index}`"
                      class="flex items-center gap-2"
                    >
                      <Input
                        v-model="form.bonuses[index]"
                        placeholder="Нэг давуу тал бичнэ үү"
                      />
                      <Button
                        v-if="index > 0"
                        type="button"
                        variant="outline"
                        size="icon"
                        @click="removeItem(form.bonuses, index)"
                      >
                        <X class="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <DialogFooter class="border-t border-border px-6 py-4">
            <Button type="button" variant="outline" @click="closeDialog">
              Болих
            </Button>
            <Button type="submit" :disabled="isSubmitting">
              {{
                isSubmitting
                  ? "Хадгалж байна..."
                  : isEditing
                    ? "Шинэчлэх"
                    : "Үүсгэх"
              }}
            </Button>
          </DialogFooter>
        </form>
      </div>
    </DialogContent>
  </Dialog>
</template>
