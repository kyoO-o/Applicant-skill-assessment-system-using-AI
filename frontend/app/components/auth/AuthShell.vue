<script setup lang="ts">
interface Props {
  badge?: string;
  title?: string;
  description?: string;
  highlights?: Array<{ label: string; value: string }>;
}

withDefaults(defineProps<Props>(), {
  badge: "Skillz",
});
</script>

<template>
  <div
    class="flex min-h-screen w-screen overflow-auto bg-background font-sans text-foreground"
  >
    <!-- ── Left brand panel ──────────────────────────────── -->
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
      <!-- Subtle noise texture overlay -->
      <div
        class="pointer-events-none absolute inset-0 opacity-[0.04]"
        style="
          background-image: url('data:image/svg+xml,%3Csvg width=&quot;60&quot; height=&quot;60&quot; viewBox=&quot;0 0 60 60&quot; xmlns=&quot;http://www.w3.org/2000/svg&quot;%3E%3Cg fill=&quot;none&quot; fill-rule=&quot;evenodd&quot;%3E%3Cg fill=&quot;%23ffffff&quot; fill-opacity=&quot;1&quot;%3E%3Cpath d=&quot;M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z&quot;/%3E%3C/g%3E%3C/g%3E%3C/svg%3E');
        "
      />

      <div class="relative flex flex-1 flex-col p-12">
        <!-- Logo -->
        <NuxtLink to="/" class="flex items-center gap-3 no-underline">
          <!-- <div
            class="flex h-9 w-9 items-center justify-center rounded-[10px] text-[18px] font-bold text-white"
            style="
              background: rgba(255, 255, 255, 0.2);
              backdrop-filter: blur(8px);
            "
          >
            m
          </div> -->
          <img class="block h-8 w-8" src="/icons/logo.svg" alt="Skillz logo" />
          <span class="text-[19px] font-semibold tracking-tight text-white"
            >Skillz</span
          >
        </NuxtLink>

        <!-- Hero text -->
        <div class="mt-auto max-w-sm">
          <h1
            class="text-[40px] font-semibold leading-[1.1] tracking-[-1px] text-white"
            style="white-space: pre-line"
          >
            {{ title || "Тохирох ажлаа олоорой.\nСайн ажилтнаа олоорой." }}
          </h1>
          <p class="mt-5 text-[16px] leading-[1.55] text-white/80">
            {{
              description ||
              "AI-д суурилсан CV үнэлгээ, ярилцлагын хуваарь, даалгаврын менежмент — нэг дор."
            }}
          </p>

          <!-- AI preview card -->
          <div
            class="mt-9 rounded-[20px] p-5"
            style="
              background: rgba(255, 255, 255, 0.12);
              backdrop-filter: blur(8px);
              border: 1px solid rgba(255, 255, 255, 0.15);
            "
          >
            <div class="mb-2.5 flex items-center gap-2">
              <svg
                width="14"
                height="14"
                viewBox="0 0 20 20"
                fill="none"
                stroke="currentColor"
                stroke-width="1.6"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="text-white"
              >
                <path
                  d="M10 2l1.7 4.8L16.5 8.5 11.7 10.2 10 15l-1.7-4.8L3.5 8.5 8.3 6.8z"
                />
              </svg>
              <span
                class="text-[11px] font-semibold uppercase tracking-[1px] text-white"
                >AI үнэлгээ</span
              >
            </div>
            <p class="text-[18px] font-semibold text-white">
              "Маш сайн тохирол — 87/100"
            </p>
            <p class="mt-1 text-[13px] text-white/80">
              Хэдхэн секундэд ажлын байрны шаардлагатай тохирол хийнэ.
            </p>
          </div>
        </div>

        <!-- Footer -->
        <div class="mt-10 flex items-center gap-4 text-[12px] text-white/60">
          <span>© 2026 Skillz</span>
          <span>·</span>
          <a class="text-white/60 hover:text-white/90 cursor-pointer"
            >Нууцлал</a
          >
          <span>·</span>
          <a class="text-white/60 hover:text-white/90 cursor-pointer">Нөхцөл</a>
        </div>
      </div>
    </div>

    <!-- ── Right form area ───────────────────────────────── -->
    <div class="flex flex-1 items-center justify-center px-8 py-10">
      <!-- Mobile logo -->
      <div class="absolute left-6 top-6 lg:hidden">
        <NuxtLink to="/" class="flex items-center gap-2.5 no-underline">
          <!-- <div
            class="flex h-8 w-8 items-center justify-center rounded-[10px] text-[15px] font-bold text-white"
            style="
              background: linear-gradient(
                135deg,
                var(--primary),
                oklch(0.348 0.106 295)
              );
            "
          >
            m
          </div> -->
          <img class="block h-8 w-8" src="/icons/logo.svg" alt="Skillz logo" />
          <span class="text-[16px] font-semibold text-foreground">Skillz</span>
        </NuxtLink>
      </div>

      <div class="w-full max-w-[400px]">
        <slot />
      </div>
    </div>
  </div>
</template>
