<script setup lang="ts">
import { Menu } from "lucide-vue-next";

const navLinks = [
  { label: "Features", href: "#features" },
  { label: "For Applicants", href: "#roles" },
  { label: "For Recruiters", href: "#roles" },
];
</script>

<template>
  <header class="sticky top-0 z-40 border-b border-slate-200/80 bg-white/80 backdrop-blur-xl">
    <div class="mx-auto flex h-20 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
      <NuxtLink to="/" class="flex items-center gap-3">
        <span
          class="flex h-10 w-10 items-center justify-center rounded-2xl bg-gradient-to-br from-blue-600 to-violet-600 text-sm font-semibold text-white shadow-lg shadow-blue-200/60"
        >
          S
        </span>
        <span class="text-sm font-semibold tracking-wide text-slate-950">
          SkillMatch AI
        </span>
      </NuxtLink>

      <nav class="hidden items-center gap-8 md:flex">
        <a
          v-for="link in navLinks"
          :key="link.label"
          :href="link.href"
          class="text-sm font-medium text-slate-600 transition hover:text-slate-950"
        >
          {{ link.label }}
        </a>
      </nav>

      <div class="hidden items-center gap-3 md:flex">
        <Button as-child variant="ghost">
          <NuxtLink to="/login">Login</NuxtLink>
        </Button>
        <Button as-child class="shadow-sm shadow-blue-200/60">
          <NuxtLink to="/register">Get Started</NuxtLink>
        </Button>
      </div>

      <Sheet>
        <SheetTrigger as-child class="md:hidden">
          <Button variant="outline" size="icon" aria-label="Open navigation menu">
            <Menu class="h-5 w-5" />
          </Button>
        </SheetTrigger>
        <SheetContent side="right" class="w-[300px] border-slate-200">
          <SheetHeader class="text-left">
            <SheetTitle>SkillMatch AI</SheetTitle>
            <SheetDescription>
              Navigate the platform and access applicant or recruiter tools.
            </SheetDescription>
          </SheetHeader>

          <div class="mt-8 flex flex-col gap-4">
            <a
              v-for="link in navLinks"
              :key="link.label"
              :href="link.href"
              class="rounded-xl px-3 py-2 text-sm font-medium text-slate-700 transition hover:bg-slate-100 hover:text-slate-950"
            >
              {{ link.label }}
            </a>
          </div>

          <div class="mt-8 grid gap-3">
            <Button as-child variant="outline" class="w-full">
              <NuxtLink to="/login">Login</NuxtLink>
            </Button>
            <Button as-child class="w-full">
              <NuxtLink to="/register">Get Started</NuxtLink>
            </Button>
          </div>
        </SheetContent>
      </Sheet>
    </div>
  </header>
</template>
