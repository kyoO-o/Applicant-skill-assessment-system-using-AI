<script setup lang="ts">
import { Eye, EyeOff } from "lucide-vue-next";

defineOptions({ inheritAttrs: false });
const model = defineModel<string>();
const show = ref(false);
</script>

<template>
  <div class="relative">
    <Input
      v-bind="$attrs"
      v-model="model"
      :type="show ? 'text' : 'password'"
      class="pr-10"
    />
    <button
      type="button"
      tabindex="-1"
      class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground transition hover:text-foreground"
      @click="show = !show"
    >
      <EyeOff v-if="show" class="h-4 w-4" />
      <Eye v-else class="h-4 w-4" />
    </button>
  </div>
</template>
