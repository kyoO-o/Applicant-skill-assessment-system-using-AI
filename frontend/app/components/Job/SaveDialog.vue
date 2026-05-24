<script setup lang="ts">
import { Plus, X } from "lucide-vue-next";

import { JobStatus } from "../../composables/types";

const { cities, districtsFor } = useLocationOptions();

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
                  <Label for="job-title">Ажлын байрны нэр</Label>
                  <Input
                    id="job-title"
                    v-model="form.title"
                    placeholder="Frontend Developer"
                  />
                </div>
                <div class="space-y-2">
                  <Label for="job-contact">Холбоо барих мэдээлэл</Label>
                  <Input
                    id="job-contact"
                    v-model="form.contact_info"
                    placeholder="hr@company.mn or +976..."
                  />
                </div>
              </div>

              <div class="grid gap-4 sm:grid-cols-3">
                <div class="space-y-2">
                  <Label for="job-type">Хөдөлмөрийн гэрээний хэлбэр</Label>
                  <Input
                    id="job-type"
                    v-model="form.type"
                    placeholder="Full-time"
                  />
                </div>
                <div class="space-y-2">
                  <Label for="job-level">Мэргэжлийн түвшин</Label>
                  <Input
                    id="job-level"
                    v-model="form.level"
                    placeholder="Mid-level"
                  />
                </div>
                <div class="space-y-2">
                  <Label for="job-status">Төлөв</Label>
                  <Select v-model="form.status">
                    <SelectTrigger id="job-status" class="w-full">
                      <SelectValue placeholder="Төлөв сонгох" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem :value="JobStatus.Draft">Ноорог</SelectItem>
                      <SelectItem :value="JobStatus.Posted"
                        >Нийтлэгдсэн</SelectItem
                      >
                      <SelectItem :value="JobStatus.Closed"
                        >Хаагдсан</SelectItem
                      >
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div class="grid gap-4 sm:grid-cols-2">
                <div class="space-y-2">
                  <Label for="job-city">Хот / Аймаг</Label>
                  <Select v-model="form.city">
                    <SelectTrigger id="job-city" class="w-full">
                      <SelectValue placeholder="Хот эсвэл аймаг сонгох" />
                    </SelectTrigger>
                    <SelectContent>
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

              <div
                class="rounded-2xl border border-border bg-muted/20 p-4 space-y-2"
              >
                <Label>Газрын зураг дээр байршил сонгох</Label>
                <LocationSearch
                  v-model:model-x="form.location_x"
                  v-model:model-y="form.location_y"
                />
              </div>

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

              <div class="space-y-2">
                <Label for="job-description">Нэмэлт мэдээлэл</Label>
                <Textarea
                  id="job-description"
                  v-model="form.additional_info"
                  rows="5"
                  placeholder="Ажлын байрны тайлбар, хамрах хүрээ болон хүлээгдэж буй үр дүнг бичнэ үү."
                />
              </div>

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
                      <Plus class="mr-1 h-4 w-4" />
                      Нэмэх
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
                      <Plus class="mr-1 h-4 w-4" />
                      Нэмэх
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
                      <Plus class="mr-1 h-4 w-4" />
                      Нэмэх
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
                      <Plus class="mr-1 h-4 w-4" />
                      Нэмэх
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

              <div class="space-y-2">
                <Label>Байршлын тойм</Label>
                <p class="text-xs text-muted-foreground">
                  {{
                    [form.district, form.city].filter(Boolean).join(", ") ||
                    "Сонгогдоогүй"
                  }}
                </p>
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
