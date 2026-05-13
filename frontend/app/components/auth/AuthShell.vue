<script setup lang="ts">
interface AuthHighlight {
  label: string;
  value: string;
}

interface Props {
  badge?: string;
  title: string;
  description: string;
  highlights?: AuthHighlight[];
}

withDefaults(defineProps<Props>(), {
  badge: "SkillMatch AI",
  highlights: () => [
    { label: "Үнэлгээ", value: "AI-д суурилсан ур чадварын шинжилгээ" },
    { label: "Ажлын урсгал", value: "Ажил олгогч болон горилогчдод" },
    { label: "Хандалт", value: "Нэг аюулгүй ажлын орчин" },
  ],
});
</script>

<template>
  <div
    class="relative min-h-screen overflow-hidden bg-[linear-gradient(180deg,#f7fbff_0%,#ffffff_45%,#eef6ff_100%)]"
  >
    <div
      class="absolute inset-x-0 top-0 -z-10 h-96 bg-[radial-gradient(circle_at_top_left,_rgba(14,165,233,0.18),_transparent_38%),radial-gradient(circle_at_top_right,_rgba(59,130,246,0.18),_transparent_34%),radial-gradient(circle_at_center,_rgba(14,116,144,0.08),_transparent_45%)]"
    />

    <div class="mx-auto flex min-h-screen max-w-7xl items-center px-4 py-10 sm:px-6 lg:px-8">
      <div class="grid w-full gap-10 lg:grid-cols-[1.05fr_0.95fr] lg:items-center">
        <section class="hidden space-y-8 lg:block">
          <NuxtLink
            to="/"
            class="inline-flex items-center gap-3 text-sm font-semibold tracking-[0.18em] text-cyan-900 uppercase transition hover:text-cyan-700"
          >
            <span
              class="flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br from-cyan-600 via-sky-600 to-blue-700 text-lg font-bold text-white shadow-lg shadow-sky-200/70"
            >
              S
            </span>
            {{ badge }}
          </NuxtLink>

          <div class="space-y-4">
            <p class="text-sm font-semibold uppercase tracking-[0.24em] text-cyan-700">
              Орчин үеийн ажилд авах платформ
            </p>
            <h1 class="max-w-xl text-5xl font-semibold tracking-tight text-slate-950">
              {{ title }}
            </h1>
            <p class="max-w-2xl text-base leading-7 text-slate-600">
              {{ description }}
            </p>
          </div>

          <div class="grid max-w-2xl gap-4 sm:grid-cols-3">
            <div
              v-for="item in highlights"
              :key="item.label"
              class="rounded-[1.75rem] border border-white/80 bg-white/75 p-5 shadow-sm backdrop-blur"
            >
              <p class="text-sm text-slate-500">
                {{ item.label }}
              </p>
              <p class="mt-2 text-lg font-semibold leading-6 text-slate-950">
                {{ item.value }}
              </p>
            </div>
          </div>
        </section>

        <section class="w-full">
          <div class="mb-6 lg:hidden">
            <NuxtLink
              to="/"
              class="inline-flex items-center gap-3 text-sm font-semibold text-cyan-950 transition hover:text-cyan-700"
            >
              <span
                class="flex h-11 w-11 items-center justify-center rounded-2xl bg-gradient-to-br from-cyan-600 via-sky-600 to-blue-700 text-lg font-bold text-white shadow-lg shadow-sky-200/70"
              >
                S
              </span>
              {{ badge }}
            </NuxtLink>
          </div>

          <slot />
        </section>
      </div>
    </div>
  </div>
</template>
