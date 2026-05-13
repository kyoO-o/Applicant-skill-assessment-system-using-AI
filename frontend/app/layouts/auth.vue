<template>
  <div
    class="relative min-h-screen overflow-hidden bg-[linear-gradient(180deg,#f8fbff_0%,#ffffff_45%,#f8fafc_100%)]"
  >
    <div
      class="absolute inset-x-0 top-0 -z-10 h-80 bg-[radial-gradient(circle_at_top,_rgba(59,130,246,0.18),_transparent_45%),radial-gradient(circle_at_top_right,_rgba(124,58,237,0.16),_transparent_30%)]"
    />

    <header class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
      <NuxtLink
        to="/"
        class="inline-flex items-center gap-3 text-sm font-semibold text-slate-900 transition hover:text-blue-700"
      >
        <span
          class="flex h-10 w-10 items-center justify-center rounded-2xl bg-gradient-to-br from-blue-600 to-violet-600 text-white shadow-lg shadow-blue-200/60"
        >
          S
        </span>
        SkillMatch AI
      </NuxtLink>
    </header>

    <main class="mx-auto flex min-h-[calc(100vh-5.5rem)] max-w-7xl items-center px-4 pb-10 sm:px-6 lg:px-8">
      <div class="grid w-full gap-10 lg:grid-cols-[1.05fr_0.95fr] lg:items-center">
        <section class="hidden space-y-6 lg:block">
          <Badge class="rounded-full px-3 py-1">HR Технологийн Платформ</Badge>
          <div class="space-y-4">
            <h1 class="max-w-xl text-5xl font-semibold tracking-tight text-slate-950">
              AI-д суурилсан ур чадварын үнэлгээгээр илүү оновчтой ажилд авах
            </h1>
            <p class="max-w-2xl text-base leading-7 text-slate-600">
              Нэр дэвшигчийн үнэлгээ, ажлын байрны тохирол, ойлгомжтой шийдвэр дэмжлэгт
              зориулсан орчин үеийн ажлын орчин.
            </p>
          </div>

          <div class="grid max-w-2xl gap-4 sm:grid-cols-3">
            <div class="rounded-3xl border border-white/80 bg-white/80 p-5 shadow-sm backdrop-blur">
              <p class="text-sm text-slate-500">Үнэлгээ</p>
              <p class="mt-2 text-lg font-semibold text-slate-950">CV шинжилгээ</p>
            </div>
            <div class="rounded-3xl border border-white/80 bg-white/80 p-5 shadow-sm backdrop-blur">
              <p class="text-sm text-slate-500">Тохирол</p>
              <p class="mt-2 text-lg font-semibold text-slate-950">Үүрэгт тохирох шинжилгээ</p>
            </div>
            <div class="rounded-3xl border border-white/80 bg-white/80 p-5 shadow-sm backdrop-blur">
              <p class="text-sm text-slate-500">Ажлын урсгал</p>
              <p class="mt-2 text-lg font-semibold text-slate-950">Ажил олгогчид тохирсон</p>
            </div>
          </div>
        </section>

        <section class="w-full">
          <slot />
        </section>
      </div>
    </main>
  </div>
</template>
