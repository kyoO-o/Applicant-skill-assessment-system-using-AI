<template>
  <div
    class="flex min-h-screen w-screen bg-background font-sans text-foreground"
  >
    <!-- Left brand panel -->
    <div
      class="relative hidden flex-col overflow-hidden lg:flex lg:w-[46%]"
      style="
        background: linear-gradient(
          160deg,
          var(--primary) 0%,
          oklch(0.348 0.106 295) 100%
        );
        color: #fff;
      "
    >
      <div class="relative flex flex-1 flex-col p-12">
        <NuxtLink to="/" class="flex items-center gap-3 no-underline">
          <div
            class="flex h-9 w-9 items-center justify-center rounded-[10px] text-[18px] font-bold text-white"
            style="
              background: rgba(255, 255, 255, 0.2);
              backdrop-filter: blur(8px);
            "
          >
            m
          </div>
          <span class="text-[19px] font-semibold tracking-tight text-white"
            >Skillz</span
          >
        </NuxtLink>
        <div class="mt-auto max-w-sm">
          <h1
            class="text-[38px] font-semibold leading-[1.1] tracking-[-1px] text-white"
          >
            Тохирох ажлыг олоорой.<br />Сайн ажилтнаа олоорой.
          </h1>
          <p class="mt-5 text-[15px] leading-[1.55] text-white/80">
            AI-д суурилсан CV үнэлгээ, ярилцлагын хуваарь, даалгаврын менежмент
            — нэг дор.
          </p>
        </div>
        <div class="mt-10 flex items-center gap-4 text-[12px] text-white/60">
          <span>© 2026 Skillz</span>
        </div>
      </div>
    </div>

    <!-- Right form area -->
    <div class="flex flex-1 items-center justify-center px-8 py-10">
      <div class="w-full max-w-[420px]">
        <slot />
      </div>
    </div>
  </div>
</template>
