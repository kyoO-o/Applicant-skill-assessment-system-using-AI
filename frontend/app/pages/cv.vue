<script setup lang="ts">
import { toast } from "vue-sonner";
import {
  Loader2,
  Plus,
  Trash2,
  Upload,
  Printer,
  Save,
  ScrollText,
  Briefcase,
  GraduationCap,
  Brain,
  BookOpen,
  FlaskConical,
  Building2,
  Trophy,
} from "lucide-vue-next";
import type {
  CVProfile,
  WorkExperience,
  Education,
  Language,
  Training,
  Exam,
  Internship,
  Award,
} from "../composables/types";
import {
  emptyCV,
  emptyWorkExperience,
  emptyEducation,
  emptyLanguage,
  emptyTraining,
  emptyExam,
  emptyInternship,
  emptyAward,
} from "../composables/types";

definePageMeta({ middleware: "auth" });

const cvAPI = useCVAPI();
const { user } = useAuth();

const loading = ref(true);
const saving = ref(false);
const parsing = ref(false);

const cv = ref<CVProfile>(emptyCV());

const DRIVER_LICENSES = ["A", "B", "C", "D", "E", "M"];
const GENDERS = ["Эрэгтэй", "Эмэгтэй"];
const MARITAL_STATUSES = ["Гэрлэсэн", "Гэрлээгүй", "Салсан", "Бэлэвсэн"];
const LANGUAGE_LEVELS = [
  "Анхан шатны",
  "Дунд",
  "Дэвшилтэт",
  "Чөлөөт",
  "Эх хэл",
];

async function load() {
  try {
    const data = await cvAPI.get();
    cv.value = { ...emptyCV(), ...data };
  } catch {
    toast.error("CV мэдээлэл ачааллаж чадсангүй");
  } finally {
    loading.value = false;
  }
}

await load();

async function save() {
  saving.value = true;
  try {
    const saved = await cvAPI.save(cv.value);
    cv.value = { ...emptyCV(), ...saved };
    toast.success("CV хадгалагдлаа");
  } catch {
    toast.error("Хадгалахад алдаа гарлаа");
  } finally {
    saving.value = false;
  }
}

const fileInput = ref<HTMLInputElement | null>(null);

async function onImportPDF(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;

  parsing.value = true;
  try {
    const parsed = await cvAPI.parse(file);
    cv.value = { ...cv.value, ...parsed };
    toast.success("CV амжилттай уншлаа. Мэдээллийг шалгаад хадгалаарай.");
  } catch {
    toast.error("PDF уншихад алдаа гарлаа");
  } finally {
    parsing.value = false;
    input.value = "";
  }
}

function print() {
  window.print();
}

// Driver licenses
function toggleLicense(lic: string) {
  const idx = cv.value.driver_licenses.indexOf(lic);
  if (idx === -1) cv.value.driver_licenses.push(lic);
  else cv.value.driver_licenses.splice(idx, 1);
}

// Generic list item operations
function addItem<T>(list: T[], factory: () => T) {
  list.push(factory());
}
function removeItem<T>(list: T[], index: number) {
  list.splice(index, 1);
}

// Skill tag helpers
const skillInput = ref<Record<string, string>>({
  personal: "",
  professional: "",
  computer: "",
  art: "",
  sport: "",
});

function addSkill(field: keyof CVProfile, inputKey: string) {
  const val = (skillInput.value[inputKey] ?? "").trim();
  if (!val) return;
  const arr = cv.value[field] as string[];
  if (!arr.includes(val)) arr.push(val);
  skillInput.value[inputKey] = "";
}

function removeSkill(field: keyof CVProfile, index: number) {
  const arr = cv.value[field] as string[];
  arr.splice(index, 1);
}

function onSkillKeydown(
  e: KeyboardEvent,
  field: keyof CVProfile,
  inputKey: string,
) {
  if (e.key === "Enter") {
    e.preventDefault();
    addSkill(field, inputKey);
  }
}

// Sections for sidebar nav
const sections = [
  { id: "general", label: "Ерөнхий мэдээлэл", icon: ScrollText },
  { id: "about", label: "Миний тухай", icon: BookOpen },
  { id: "work", label: "Ажлын туршлага", icon: Briefcase },
  { id: "education", label: "Боловсрол", icon: GraduationCap },
  { id: "skills", label: "Ур чадвар", icon: Brain },
  { id: "trainings", label: "Сургалт, сертификат", icon: BookOpen },
  { id: "exams", label: "Шалгалт, тест", icon: FlaskConical },
  { id: "internships", label: "Дадлага", icon: Building2 },
  { id: "awards", label: "Шагнал урамшуулал", icon: Trophy },
];

function scrollTo(id: string) {
  document
    .getElementById(id)
    ?.scrollIntoView({ behavior: "smooth", block: "start" });
}

// Print: computed full name
const fullName = computed(() =>
  [cv.value.last_name, cv.value.first_name].filter(Boolean).join(" "),
);
</script>

<template>
  <div>
    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-20">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <template v-else>
      <!-- ===== SCREEN VIEW ===== -->
      <div class="print:hidden space-y-6">
        <!-- Header -->
        <section class="rounded-3xl border border-border bg-card px-6 py-6">
          <div
            class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between"
          >
            <div class="space-y-1">
              <p class="text-sm font-medium text-muted-foreground">
                CV Бүтээгч
              </p>
              <h1 class="text-3xl font-semibold tracking-tight">Миний CV</h1>
              <p class="text-sm text-muted-foreground">
                Мэдээллээ бөглөөд PDF болгон татаж авах эсвэл PDF-ийг импортлон
                автоматаар бөглүүлэх боломжтой.
              </p>
            </div>
            <div class="flex flex-wrap gap-2">
              <input
                ref="fileInput"
                type="file"
                accept=".pdf"
                class="hidden"
                @change="onImportPDF"
              />
              <Button
                variant="outline"
                class="gap-2 rounded-full"
                :disabled="parsing"
                @click="fileInput?.click()"
              >
                <Loader2 v-if="parsing" class="h-4 w-4 animate-spin" />
                <Upload v-else class="h-4 w-4" />
                {{ parsing ? "Уншиж байна..." : "PDF импортлох" }}
              </Button>
              <Button
                variant="outline"
                class="gap-2 rounded-full"
                @click="print"
              >
                <Printer class="h-4 w-4" />
                PDF татах
              </Button>
              <Button
                class="gap-2 rounded-full"
                :disabled="saving"
                @click="save"
              >
                <Loader2 v-if="saving" class="h-4 w-4 animate-spin" />
                <Save v-else class="h-4 w-4" />
                {{ saving ? "Хадгалж байна..." : "Хадгалах" }}
              </Button>
            </div>
          </div>
        </section>

        <div class="flex gap-6">
          <!-- Sidebar nav -->
          <aside class="hidden w-52 shrink-0 lg:block">
            <div
              class="sticky top-6 space-y-1 rounded-3xl border border-border bg-card px-3 py-3"
            >
              <button
                v-for="s in sections"
                :key="s.id"
                type="button"
                class="flex w-full items-center gap-2.5 rounded-2xl px-3 py-2.5 text-left text-sm text-muted-foreground transition hover:bg-muted hover:text-foreground"
                @click="scrollTo(s.id)"
              >
                <component :is="s.icon" class="h-4 w-4 shrink-0" />
                {{ s.label }}
              </button>
            </div>
          </aside>

          <!-- Form -->
          <div class="min-w-0 flex-1 space-y-6">
            <!-- General Info -->
            <section
              id="general"
              class="rounded-3xl border border-border bg-card px-6 py-6"
            >
              <h2 class="mb-5 text-lg font-semibold flex items-center gap-2">
                <ScrollText class="h-5 w-5 text-muted-foreground" />
                Ерөнхий мэдээлэл
              </h2>
              <div class="grid gap-4 sm:grid-cols-2">
                <div class="space-y-1.5">
                  <label class="text-sm font-medium">Овог</label>
                  <Input v-model="cv.last_name" placeholder="Овог" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-sm font-medium">Нэр</label>
                  <Input v-model="cv.first_name" placeholder="Нэр" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-sm font-medium">Төрсөн огноо</label>
                  <Input v-model="cv.date_of_birth" type="date" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-sm font-medium">Хүйс</label>
                  <Select v-model="cv.gender">
                    <SelectTrigger class="rounded-xl">
                      <SelectValue placeholder="Сонгох" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem v-for="g in GENDERS" :key="g" :value="g">{{
                        g
                      }}</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div class="space-y-1.5">
                  <label class="text-sm font-medium">Регистрийн дугаар</label>
                  <Input v-model="cv.national_id" placeholder="АА00000000" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-sm font-medium">Гэрлэлтийн байдал</label>
                  <Select v-model="cv.marital_status">
                    <SelectTrigger class="rounded-xl">
                      <SelectValue placeholder="Сонгох" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem
                        v-for="m in MARITAL_STATUSES"
                        :key="m"
                        :value="m"
                        >{{ m }}</SelectItem
                      >
                    </SelectContent>
                  </Select>
                </div>
                <div class="space-y-1.5">
                  <label class="text-sm font-medium">Утасны дугаар</label>
                  <Input v-model="cv.phone" placeholder="+976 9900 0000" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-sm font-medium">И-мэйл</label>
                  <Input
                    v-model="cv.email"
                    type="email"
                    placeholder="mail@example.com"
                  />
                </div>
                <div class="space-y-1.5 sm:col-span-2">
                  <label class="text-sm font-medium">Хаяг</label>
                  <Input
                    v-model="cv.address"
                    placeholder="Улаанбаатар, Сүхбаатар дүүрэг..."
                  />
                </div>
                <div class="space-y-2 sm:col-span-2">
                  <label class="text-sm font-medium">Жолооны үнэмлэх</label>
                  <div class="flex flex-wrap gap-2">
                    <button
                      v-for="lic in DRIVER_LICENSES"
                      :key="lic"
                      type="button"
                      class="h-9 w-9 rounded-xl border text-sm font-semibold transition"
                      :class="
                        cv.driver_licenses.includes(lic)
                          ? 'border-primary bg-primary text-primary-foreground'
                          : 'border-border bg-muted text-muted-foreground hover:bg-muted/80'
                      "
                      @click="toggleLicense(lic)"
                    >
                      {{ lic }}
                    </button>
                  </div>
                </div>
              </div>
            </section>

            <!-- About -->
            <section
              id="about"
              class="rounded-3xl border border-border bg-card px-6 py-6"
            >
              <h2 class="mb-5 text-lg font-semibold flex items-center gap-2">
                <BookOpen class="h-5 w-5 text-muted-foreground" />
                Миний тухай
              </h2>
              <Textarea
                v-model="cv.about"
                placeholder="Өөрийнхөө тухай товч танилцуулга бичнэ үү..."
                class="min-h-28 resize-y rounded-xl"
              />
            </section>

            <!-- Work Experience -->
            <section
              id="work"
              class="rounded-3xl border border-border bg-card px-6 py-6"
            >
              <div class="mb-5 flex items-center justify-between">
                <h2 class="text-lg font-semibold flex items-center gap-2">
                  <Briefcase class="h-5 w-5 text-muted-foreground" />
                  Ажлын туршлага
                </h2>
                <Button
                  size="sm"
                  variant="outline"
                  class="gap-1.5 rounded-full"
                  @click="addItem(cv.work_experiences, emptyWorkExperience)"
                >
                  <Plus class="h-3.5 w-3.5" /> Нэмэх
                </Button>
              </div>
              <div class="space-y-4">
                <div
                  v-for="(exp, i) in cv.work_experiences"
                  :key="i"
                  class="relative rounded-2xl border border-border p-4"
                >
                  <button
                    type="button"
                    class="absolute right-3 top-3 rounded-lg p-1 text-muted-foreground hover:text-destructive"
                    @click="removeItem(cv.work_experiences, i)"
                  >
                    <Trash2 class="h-4 w-4" />
                  </button>
                  <div class="grid gap-3 sm:grid-cols-2">
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Байгууллага</label
                      >
                      <Input
                        v-model="exp.company"
                        placeholder="Компанийн нэр"
                      />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Албан тушаал</label
                      >
                      <Input v-model="exp.position" placeholder="Ажлын нэр" />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Эхлэх огноо</label
                      >
                      <Input v-model="exp.start_date" type="month" />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Дуусах огноо</label
                      >
                      <Input
                        v-model="exp.end_date"
                        type="month"
                        :disabled="exp.current"
                      />
                    </div>
                    <div class="flex items-center gap-2">
                      <Checkbox
                        :id="`work-current-${i}`"
                        v-model:checked="exp.current"
                      />
                      <label :for="`work-current-${i}`" class="text-sm"
                        >Одоогоор ажиллаж байна</label
                      >
                    </div>
                    <div class="space-y-1 sm:col-span-2">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Тайлбар</label
                      >
                      <Textarea
                        v-model="exp.description"
                        placeholder="Гүйцэтгэсэн үүрэг, амжилтууд..."
                        class="min-h-20 resize-y rounded-xl"
                      />
                    </div>
                  </div>
                </div>
                <p
                  v-if="!cv.work_experiences.length"
                  class="text-center text-sm text-muted-foreground py-4"
                >
                  Ажлын туршлага нэмэгдээгүй байна
                </p>
              </div>
            </section>

            <!-- Education -->
            <section
              id="education"
              class="rounded-3xl border border-border bg-card px-6 py-6"
            >
              <div class="mb-5 flex items-center justify-between">
                <h2 class="text-lg font-semibold flex items-center gap-2">
                  <GraduationCap class="h-5 w-5 text-muted-foreground" />
                  Боловсрол
                </h2>
                <Button
                  size="sm"
                  variant="outline"
                  class="gap-1.5 rounded-full"
                  @click="addItem(cv.education, emptyEducation)"
                >
                  <Plus class="h-3.5 w-3.5" /> Нэмэх
                </Button>
              </div>
              <div class="space-y-4">
                <div
                  v-for="(edu, i) in cv.education"
                  :key="i"
                  class="relative rounded-2xl border border-border p-4"
                >
                  <button
                    type="button"
                    class="absolute right-3 top-3 rounded-lg p-1 text-muted-foreground hover:text-destructive"
                    @click="removeItem(cv.education, i)"
                  >
                    <Trash2 class="h-4 w-4" />
                  </button>
                  <div class="grid gap-3 sm:grid-cols-2">
                    <div class="space-y-1 sm:col-span-2">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Сургуулийн нэр</label
                      >
                      <Input
                        v-model="edu.school"
                        placeholder="Монгол Улсын Их Сургууль"
                      />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Зэрэг</label
                      >
                      <Input
                        v-model="edu.degree"
                        placeholder="Бакалавр, Магистр..."
                      />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Чиглэл / Мэргэжил</label
                      >
                      <Input
                        v-model="edu.field"
                        placeholder="Компьютерийн шинжлэх ухаан"
                      />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Элссэн огноо</label
                      >
                      <Input v-model="edu.start_date" type="month" />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Төгссөн огноо</label
                      >
                      <Input
                        v-model="edu.end_date"
                        type="month"
                        :disabled="edu.current"
                      />
                    </div>
                    <div class="flex items-center gap-2">
                      <Checkbox
                        :id="`edu-current-${i}`"
                        v-model:checked="edu.current"
                      />
                      <label :for="`edu-current-${i}`" class="text-sm"
                        >Одоогоор суралцаж байна</label
                      >
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >GPA / Дундаж оноо</label
                      >
                      <Input v-model="edu.gpa" placeholder="3.8 / 4.0" />
                    </div>
                  </div>
                </div>
                <p
                  v-if="!cv.education.length"
                  class="text-center text-sm text-muted-foreground py-4"
                >
                  Боловсролын мэдээлэл нэмэгдээгүй байна
                </p>
              </div>
            </section>

            <!-- Skills -->
            <section
              id="skills"
              class="rounded-3xl border border-border bg-card px-6 py-6"
            >
              <h2 class="mb-5 text-lg font-semibold flex items-center gap-2">
                <Brain class="h-5 w-5 text-muted-foreground" />
                Мэдлэг ур чадвар
              </h2>
              <div class="space-y-6">
                <!-- Simple tag skill sections -->
                <div
                  v-for="(cfg, idx) in [
                    {
                      label: 'Хувийн ур чадвар',
                      field: 'personal_skills' as keyof CVProfile,
                      key: 'personal',
                      placeholder: 'Харилцааны ур чадвар...',
                    },
                    {
                      label: 'Мэргэжлийн ур чадвар',
                      field: 'professional_skills' as keyof CVProfile,
                      key: 'professional',
                      placeholder: 'JavaScript, Python...',
                    },
                    {
                      label: 'Компьютерын программ',
                      field: 'computer_skills' as keyof CVProfile,
                      key: 'computer',
                      placeholder: 'Microsoft Office, Photoshop...',
                    },
                    {
                      label: 'Урлагийн ур чадвар',
                      field: 'art_skills' as keyof CVProfile,
                      key: 'art',
                      placeholder: 'Дуу хөгжим, Зураг...',
                    },
                    {
                      label: 'Спортын ур чадвар',
                      field: 'sport_skills' as keyof CVProfile,
                      key: 'sport',
                      placeholder: 'Бөх, Сагсан бөмбөг...',
                    },
                  ]"
                  :key="cfg.key"
                >
                  <div class="space-y-2">
                    <label class="text-sm font-medium">{{ cfg.label }}</label>
                    <div class="flex flex-wrap gap-1.5 mb-2">
                      <span
                        v-for="(skill, si) in cv[cfg.field] as string[]"
                        :key="si"
                        class="inline-flex items-center gap-1 rounded-full border border-border bg-muted px-2.5 py-0.5 text-xs font-medium"
                      >
                        {{ skill }}
                        <button
                          type="button"
                          class="ml-0.5 text-muted-foreground hover:text-destructive"
                          @click="removeSkill(cfg.field, si)"
                        >
                          ×
                        </button>
                      </span>
                    </div>
                    <div class="flex gap-2">
                      <Input
                        v-model="skillInput[cfg.key]"
                        :placeholder="cfg.placeholder"
                        class="flex-1"
                        @keydown="
                          (e: KeyboardEvent) =>
                            onSkillKeydown(e, cfg.field, cfg.key)
                        "
                      />
                      <Button
                        type="button"
                        size="sm"
                        variant="outline"
                        class="shrink-0 rounded-xl"
                        @click="addSkill(cfg.field, cfg.key)"
                      >
                        <Plus class="h-3.5 w-3.5" />
                      </Button>
                    </div>
                  </div>
                </div>

                <!-- Languages -->
                <div class="space-y-2">
                  <div class="flex items-center justify-between">
                    <label class="text-sm font-medium"
                      >Гадаад хэлний мэдлэг</label
                    >
                    <Button
                      size="sm"
                      variant="outline"
                      class="gap-1 rounded-full"
                      @click="addItem(cv.languages, emptyLanguage)"
                    >
                      <Plus class="h-3.5 w-3.5" /> Нэмэх
                    </Button>
                  </div>
                  <div class="space-y-2">
                    <div
                      v-for="(lang, i) in cv.languages"
                      :key="i"
                      class="flex items-center gap-2"
                    >
                      <Input
                        v-model="lang.name"
                        placeholder="Англи, Орос..."
                        class="flex-1"
                      />
                      <Select v-model="lang.level">
                        <SelectTrigger class="w-40 rounded-xl">
                          <SelectValue placeholder="Түвшин" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem
                            v-for="lvl in LANGUAGE_LEVELS"
                            :key="lvl"
                            :value="lvl"
                            >{{ lvl }}</SelectItem
                          >
                        </SelectContent>
                      </Select>
                      <button
                        type="button"
                        class="p-1 text-muted-foreground hover:text-destructive"
                        @click="removeItem(cv.languages, i)"
                      >
                        <Trash2 class="h-4 w-4" />
                      </button>
                    </div>
                    <p
                      v-if="!cv.languages.length"
                      class="text-sm text-muted-foreground"
                    >
                      Хэлний мэдлэг нэмэгдээгүй байна
                    </p>
                  </div>
                </div>
              </div>
            </section>

            <!-- Trainings -->
            <section
              id="trainings"
              class="rounded-3xl border border-border bg-card px-6 py-6"
            >
              <div class="mb-5 flex items-center justify-between">
                <h2 class="text-lg font-semibold flex items-center gap-2">
                  <BookOpen class="h-5 w-5 text-muted-foreground" />
                  Сургалт, сертификат
                </h2>
                <Button
                  size="sm"
                  variant="outline"
                  class="gap-1.5 rounded-full"
                  @click="addItem(cv.trainings, emptyTraining)"
                >
                  <Plus class="h-3.5 w-3.5" /> Нэмэх
                </Button>
              </div>
              <div class="space-y-4">
                <div
                  v-for="(tr, i) in cv.trainings"
                  :key="i"
                  class="relative rounded-2xl border border-border p-4"
                >
                  <button
                    type="button"
                    class="absolute right-3 top-3 rounded-lg p-1 text-muted-foreground hover:text-destructive"
                    @click="removeItem(cv.trainings, i)"
                  >
                    <Trash2 class="h-4 w-4" />
                  </button>
                  <div class="grid gap-3 sm:grid-cols-2">
                    <div class="space-y-1 sm:col-span-2">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Сургалтын нэр</label
                      >
                      <Input v-model="tr.name" placeholder="Сургалтын нэр" />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Байгууллага</label
                      >
                      <Input
                        v-model="tr.organization"
                        placeholder="Зохион байгуулагч"
                      />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Огноо</label
                      >
                      <Input v-model="tr.date" type="month" />
                    </div>
                    <div class="space-y-1 sm:col-span-2">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Сертификат / Тэмдэглэл</label
                      >
                      <Input
                        v-model="tr.certificate"
                        placeholder="Сертификатын дугаар эсвэл нэр"
                      />
                    </div>
                  </div>
                </div>
                <p
                  v-if="!cv.trainings.length"
                  class="text-center text-sm text-muted-foreground py-4"
                >
                  Сургалт нэмэгдээгүй байна
                </p>
              </div>
            </section>

            <!-- Exams -->
            <section
              id="exams"
              class="rounded-3xl border border-border bg-card px-6 py-6"
            >
              <div class="mb-5 flex items-center justify-between">
                <h2 class="text-lg font-semibold flex items-center gap-2">
                  <FlaskConical class="h-5 w-5 text-muted-foreground" />
                  Шалгалт, тест
                </h2>
                <Button
                  size="sm"
                  variant="outline"
                  class="gap-1.5 rounded-full"
                  @click="addItem(cv.exams, emptyExam)"
                >
                  <Plus class="h-3.5 w-3.5" /> Нэмэх
                </Button>
              </div>
              <div class="space-y-4">
                <div
                  v-for="(ex, i) in cv.exams"
                  :key="i"
                  class="relative rounded-2xl border border-border p-4"
                >
                  <button
                    type="button"
                    class="absolute right-3 top-3 rounded-lg p-1 text-muted-foreground hover:text-destructive"
                    @click="removeItem(cv.exams, i)"
                  >
                    <Trash2 class="h-4 w-4" />
                  </button>
                  <div class="grid gap-3 sm:grid-cols-3">
                    <div class="space-y-1 sm:col-span-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Шалгалтын нэр</label
                      >
                      <Input
                        v-model="ex.name"
                        placeholder="IELTS, TOEFL, HSK..."
                      />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Оноо / Дүн</label
                      >
                      <Input v-model="ex.score" placeholder="7.5, B2..." />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Огноо</label
                      >
                      <Input v-model="ex.date" type="month" />
                    </div>
                  </div>
                </div>
                <p
                  v-if="!cv.exams.length"
                  class="text-center text-sm text-muted-foreground py-4"
                >
                  Шалгалт нэмэгдээгүй байна
                </p>
              </div>
            </section>

            <!-- Internships -->
            <section
              id="internships"
              class="rounded-3xl border border-border bg-card px-6 py-6"
            >
              <div class="mb-5 flex items-center justify-between">
                <h2 class="text-lg font-semibold flex items-center gap-2">
                  <Building2 class="h-5 w-5 text-muted-foreground" />
                  Дадлага
                </h2>
                <Button
                  size="sm"
                  variant="outline"
                  class="gap-1.5 rounded-full"
                  @click="addItem(cv.internships, emptyInternship)"
                >
                  <Plus class="h-3.5 w-3.5" /> Нэмэх
                </Button>
              </div>
              <div class="space-y-4">
                <div
                  v-for="(int, i) in cv.internships"
                  :key="i"
                  class="relative rounded-2xl border border-border p-4"
                >
                  <button
                    type="button"
                    class="absolute right-3 top-3 rounded-lg p-1 text-muted-foreground hover:text-destructive"
                    @click="removeItem(cv.internships, i)"
                  >
                    <Trash2 class="h-4 w-4" />
                  </button>
                  <div class="grid gap-3 sm:grid-cols-2">
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Байгууллага</label
                      >
                      <Input
                        v-model="int.company"
                        placeholder="Компанийн нэр"
                      />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Албан тушаал</label
                      >
                      <Input
                        v-model="int.position"
                        placeholder="Дадлагажигч, Туслах..."
                      />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Эхлэх огноо</label
                      >
                      <Input v-model="int.start_date" type="month" />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Дуусах огноо</label
                      >
                      <Input v-model="int.end_date" type="month" />
                    </div>
                    <div class="space-y-1 sm:col-span-2">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Тайлбар</label
                      >
                      <Textarea
                        v-model="int.description"
                        placeholder="Дадлагын үеийн туршлага, сурсан зүйлс..."
                        class="min-h-20 resize-y rounded-xl"
                      />
                    </div>
                  </div>
                </div>
                <p
                  v-if="!cv.internships.length"
                  class="text-center text-sm text-muted-foreground py-4"
                >
                  Дадлага нэмэгдээгүй байна
                </p>
              </div>
            </section>

            <!-- Awards -->
            <section
              id="awards"
              class="rounded-3xl border border-border bg-card px-6 py-6"
            >
              <div class="mb-5 flex items-center justify-between">
                <h2 class="text-lg font-semibold flex items-center gap-2">
                  <Trophy class="h-5 w-5 text-muted-foreground" />
                  Шагнал урамшуулал
                </h2>
                <Button
                  size="sm"
                  variant="outline"
                  class="gap-1.5 rounded-full"
                  @click="addItem(cv.awards, emptyAward)"
                >
                  <Plus class="h-3.5 w-3.5" /> Нэмэх
                </Button>
              </div>
              <div class="space-y-4">
                <div
                  v-for="(aw, i) in cv.awards"
                  :key="i"
                  class="relative rounded-2xl border border-border p-4"
                >
                  <button
                    type="button"
                    class="absolute right-3 top-3 rounded-lg p-1 text-muted-foreground hover:text-destructive"
                    @click="removeItem(cv.awards, i)"
                  >
                    <Trash2 class="h-4 w-4" />
                  </button>
                  <div class="grid gap-3 sm:grid-cols-2">
                    <div class="space-y-1 sm:col-span-2">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Шагналын нэр</label
                      >
                      <Input v-model="aw.name" placeholder="Шагналын нэр" />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Байгууллага</label
                      >
                      <Input
                        v-model="aw.organization"
                        placeholder="Шагнал олгосон байгууллага"
                      />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Огноо</label
                      >
                      <Input v-model="aw.date" type="month" />
                    </div>
                    <div class="space-y-1 sm:col-span-2">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Тайлбар</label
                      >
                      <Textarea
                        v-model="aw.description"
                        placeholder="Шагналын тухай дэлгэрэнгүй мэдээлэл..."
                        class="min-h-20 resize-y rounded-xl"
                      />
                    </div>
                  </div>
                </div>
                <p
                  v-if="!cv.awards.length"
                  class="text-center text-sm text-muted-foreground py-4"
                >
                  Шагнал нэмэгдээгүй байна
                </p>
              </div>
            </section>

            <!-- Bottom save -->
            <div class="flex justify-end gap-2 pb-6">
              <Button
                variant="outline"
                class="gap-2 rounded-full"
                @click="print"
              >
                <Printer class="h-4 w-4" /> PDF татах
              </Button>
              <Button
                class="gap-2 rounded-full"
                :disabled="saving"
                @click="save"
              >
                <Loader2 v-if="saving" class="h-4 w-4 animate-spin" />
                <Save v-else class="h-4 w-4" />
                {{ saving ? "Хадгалж байна..." : "Хадгалах" }}
              </Button>
            </div>
          </div>
        </div>
      </div>

      <!-- ===== PRINT VIEW ===== -->
      <div
        class="hidden print:block font-sans text-sm text-black leading-relaxed"
      >
        <!-- Header -->
        <div class="text-center border-b-2 border-gray-800 pb-4 mb-5">
          <h1 class="text-2xl font-bold tracking-wide">
            {{ fullName || "Нэр оруулаагүй" }}
          </h1>
          <div
            class="flex flex-wrap justify-center gap-3 mt-2 text-xs text-gray-600"
          >
            <span v-if="cv.email">{{ cv.email }}</span>
            <span v-if="cv.phone">{{ cv.phone }}</span>
            <span v-if="cv.address">{{ cv.address }}</span>
          </div>
          <div
            class="flex flex-wrap justify-center gap-3 mt-1 text-xs text-gray-600"
          >
            <span v-if="cv.date_of_birth">Төрсөн: {{ cv.date_of_birth }}</span>
            <span v-if="cv.gender">{{ cv.gender }}</span>
            <span v-if="cv.marital_status">{{ cv.marital_status }}</span>
            <span v-if="cv.national_id">Рег: {{ cv.national_id }}</span>
            <span v-if="cv.driver_licenses.length"
              >Жолооны үнэмлэх: {{ cv.driver_licenses.join(", ") }}</span
            >
          </div>
        </div>

        <!-- About -->
        <div v-if="cv.about" class="mb-5">
          <h2
            class="text-base font-bold border-b border-gray-300 pb-1 mb-2 uppercase tracking-wider"
          >
            Миний тухай
          </h2>
          <p class="text-xs">{{ cv.about }}</p>
        </div>

        <!-- Work Experience -->
        <div v-if="cv.work_experiences.length" class="mb-5">
          <h2
            class="text-base font-bold border-b border-gray-300 pb-1 mb-2 uppercase tracking-wider"
          >
            Ажлын туршлага
          </h2>
          <div v-for="(exp, i) in cv.work_experiences" :key="i" class="mb-3">
            <div class="flex justify-between items-start">
              <div>
                <p class="font-semibold">{{ exp.position }}</p>
                <p class="text-gray-600 text-xs">{{ exp.company }}</p>
              </div>
              <p class="text-xs text-gray-500 shrink-0 ml-4">
                {{ exp.start_date }} —
                {{ exp.current ? "Одоог хүртэл" : exp.end_date }}
              </p>
            </div>
            <p v-if="exp.description" class="mt-1 text-xs text-gray-700">
              {{ exp.description }}
            </p>
          </div>
        </div>

        <!-- Education -->
        <div v-if="cv.education.length" class="mb-5">
          <h2
            class="text-base font-bold border-b border-gray-300 pb-1 mb-2 uppercase tracking-wider"
          >
            Боловсрол
          </h2>
          <div v-for="(edu, i) in cv.education" :key="i" class="mb-3">
            <div class="flex justify-between items-start">
              <div>
                <p class="font-semibold">{{ edu.school }}</p>
                <p class="text-gray-600 text-xs">
                  {{ edu.degree
                  }}<span v-if="edu.field"> — {{ edu.field }}</span
                  ><span v-if="edu.gpa">, GPA: {{ edu.gpa }}</span>
                </p>
              </div>
              <p class="text-xs text-gray-500 shrink-0 ml-4">
                {{ edu.start_date }} —
                {{ edu.current ? "Одоог хүртэл" : edu.end_date }}
              </p>
            </div>
          </div>
        </div>

        <!-- Skills -->
        <div class="mb-5 grid grid-cols-2 gap-4">
          <div v-if="cv.personal_skills.length">
            <h2
              class="text-sm font-bold border-b border-gray-300 pb-1 mb-1 uppercase tracking-wider"
            >
              Хувийн ур чадвар
            </h2>
            <p class="text-xs">{{ cv.personal_skills.join(", ") }}</p>
          </div>
          <div v-if="cv.professional_skills.length">
            <h2
              class="text-sm font-bold border-b border-gray-300 pb-1 mb-1 uppercase tracking-wider"
            >
              Мэргэжлийн ур чадвар
            </h2>
            <p class="text-xs">{{ cv.professional_skills.join(", ") }}</p>
          </div>
          <div v-if="cv.computer_skills.length">
            <h2
              class="text-sm font-bold border-b border-gray-300 pb-1 mb-1 uppercase tracking-wider"
            >
              Компьютерын программ
            </h2>
            <p class="text-xs">{{ cv.computer_skills.join(", ") }}</p>
          </div>
          <div v-if="cv.languages.length">
            <h2
              class="text-sm font-bold border-b border-gray-300 pb-1 mb-1 uppercase tracking-wider"
            >
              Гадаад хэл
            </h2>
            <p class="text-xs">
              {{ cv.languages.map((l) => `${l.name} (${l.level})`).join(", ") }}
            </p>
          </div>
          <div v-if="cv.art_skills.length">
            <h2
              class="text-sm font-bold border-b border-gray-300 pb-1 mb-1 uppercase tracking-wider"
            >
              Урлагийн ур чадвар
            </h2>
            <p class="text-xs">{{ cv.art_skills.join(", ") }}</p>
          </div>
          <div v-if="cv.sport_skills.length">
            <h2
              class="text-sm font-bold border-b border-gray-300 pb-1 mb-1 uppercase tracking-wider"
            >
              Спортын ур чадвар
            </h2>
            <p class="text-xs">{{ cv.sport_skills.join(", ") }}</p>
          </div>
        </div>

        <!-- Trainings -->
        <div v-if="cv.trainings.length" class="mb-5">
          <h2
            class="text-base font-bold border-b border-gray-300 pb-1 mb-2 uppercase tracking-wider"
          >
            Сургалт, сертификат
          </h2>
          <div
            v-for="(tr, i) in cv.trainings"
            :key="i"
            class="mb-2 flex justify-between"
          >
            <div>
              <p class="font-semibold text-xs">{{ tr.name }}</p>
              <p class="text-xs text-gray-600">
                {{ tr.organization
                }}<span v-if="tr.certificate"> · {{ tr.certificate }}</span>
              </p>
            </div>
            <p class="text-xs text-gray-500 shrink-0 ml-4">{{ tr.date }}</p>
          </div>
        </div>

        <!-- Exams -->
        <div v-if="cv.exams.length" class="mb-5">
          <h2
            class="text-base font-bold border-b border-gray-300 pb-1 mb-2 uppercase tracking-wider"
          >
            Шалгалт, тест
          </h2>
          <div
            v-for="(ex, i) in cv.exams"
            :key="i"
            class="mb-1 flex justify-between"
          >
            <p class="text-xs font-semibold">{{ ex.name }}</p>
            <p class="text-xs text-gray-600">
              {{ ex.score }} <span class="text-gray-400">· {{ ex.date }}</span>
            </p>
          </div>
        </div>

        <!-- Internships -->
        <div v-if="cv.internships.length" class="mb-5">
          <h2
            class="text-base font-bold border-b border-gray-300 pb-1 mb-2 uppercase tracking-wider"
          >
            Дадлага
          </h2>
          <div v-for="(int, i) in cv.internships" :key="i" class="mb-3">
            <div class="flex justify-between items-start">
              <div>
                <p class="font-semibold text-xs">{{ int.position }}</p>
                <p class="text-gray-600 text-xs">{{ int.company }}</p>
              </div>
              <p class="text-xs text-gray-500 shrink-0 ml-4">
                {{ int.start_date }} — {{ int.end_date }}
              </p>
            </div>
            <p v-if="int.description" class="mt-1 text-xs text-gray-700">
              {{ int.description }}
            </p>
          </div>
        </div>

        <!-- Awards -->
        <div v-if="cv.awards.length" class="mb-5">
          <h2
            class="text-base font-bold border-b border-gray-300 pb-1 mb-2 uppercase tracking-wider"
          >
            Шагнал урамшуулал
          </h2>
          <div
            v-for="(aw, i) in cv.awards"
            :key="i"
            class="mb-2 flex justify-between"
          >
            <div>
              <p class="font-semibold text-xs">{{ aw.name }}</p>
              <p class="text-xs text-gray-600">{{ aw.organization }}</p>
            </div>
            <p class="text-xs text-gray-500 shrink-0 ml-4">{{ aw.date }}</p>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
