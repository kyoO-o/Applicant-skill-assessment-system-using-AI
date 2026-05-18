<script setup lang="ts">
import { toast } from "vue-sonner";
import {
  LayoutDashboard,
  BriefcaseBusiness,
  FileText,
  ScrollText,
  ClipboardList,
  Bot,
  User,
  Settings,
  Calendar,
  LogOut,
  Sun,
  Moon,
  Bell,
  ChevronDown,
  Users,
  ChevronLeft,
  PanelLeftClose,
  PanelLeftOpen,
} from "lucide-vue-next";

const { user, logout, isLoading } = useAuth();
const router = useRouter();
const route = useRoute();
const { isDark, toggle: toggleDark, init: initDark } = useDarkMode();
const {
  unreadCount,
  connect: connectWS,
  disconnect: disconnectWS,
  markAllRead,
} = useNotifications();

onMounted(() => {
  initDark();
  if (user.value) connectWS();
});

watch(
  () => user.value,
  (u) => {
    if (!u) disconnectWS();
  },
);

const sidebarCollapsed = ref(false);

const menuUser = [
  { title: "Нүүр хуудас", icon: LayoutDashboard, to: "/" },
  { title: "Ажлын байрууд", icon: BriefcaseBusiness, to: "/jobs" },
  { title: "Миний анкетууд", icon: FileText, to: "/applications" },
  { title: "CV Бүтээгч", icon: ScrollText, to: "/cv" },
  { title: "Ярилцлага", icon: Calendar, to: "/interviews", accent: true },
  { title: "Даалгаврууд", icon: ClipboardList, to: "/tasks" },
  { title: "AI Чатбот", icon: Bot, to: "/chat" },
];

const menuRecruiter = [
  { title: "Нүүр хуудас", icon: LayoutDashboard, to: "/" },
  { title: "Ажлын байрууд", icon: BriefcaseBusiness, to: "/jobs" },
  // { title: "Горилогчид", icon: Users, to: "/applications" },
  { title: "Ярилцлага", icon: Calendar, to: "/interviews", accent: true },
  { title: "Даалгаврууд", icon: ClipboardList, to: "/tasks" },
];

const bottomMenu = [
  {
    title: user?.value?.role === "recruiter" ? "Компани" : "Профайл",
    icon: User,
    to: "/profile",
  },
  { title: "Тохиргоо", icon: Settings, to: "/settings" },
];

const menu = computed(() => {
  if (user?.value?.role === "recruiter") return menuRecruiter;
  return menuUser;
});

const currentCrumb = computed(() => {
  const all = [...menu.value, ...bottomMenu];
  return all.find((item) => item.to === route.path)?.title || "Dashboard";
});

const recruiterNeedsCompany = computed(
  () => user.value?.role === "recruiter" && !user.value?.company_id,
);

function isActive(path: string) {
  if (path === "/") return route.path === "/";
  return route.path.startsWith(path);
}

function handleNav(path: string) {
  if (path === "/jobs" && recruiterNeedsCompany.value) {
    toast.warning("Ажлын байр нээхийн өмнө компанийн мэдээллээ нэмнэ үү.");
    router.push("/profile");
    return;
  }
  router.push(path);
}

async function handleLogout() {
  disconnectWS();
  await logout();
  await router.push("/login");
}

const userInitials = computed(() => {
  const name = user.value?.full_name || user.value?.name || "?";
  return name
    .split(" ")
    .map((s: string) => s[0])
    .slice(0, 2)
    .join("")
    .toUpperCase();
});

const config = useRuntimeConfig();
const userAvatarURL = computed(() => {
  const url = user.value?.profile_url;
  if (!url) return null;
  if (url.startsWith("http")) return url;
  return `${config.public.apiBase}${url}`;
});
</script>

<template>
  <div class="flex h-screen overflow-hidden bg-background text-foreground">
    <!-- ── Sidebar ─────────────────────────────────────────── -->
    <aside
      class="flex flex-shrink-0 flex-col border-r border-border bg-background transition-[width] duration-200 ease-out"
      :style="{ width: sidebarCollapsed ? '68px' : '232px' }"
    >
      <!-- Logo -->

      <div
        class="flex cursor-pointer items-center pb-5 pt-5 px-3.5"
        :class="!sidebarCollapsed ? 'justify-between' : 'justify-center'"
      >
        <div
          class="flex cursor-pointer items-center gap-2"
          @click="router.push('/')"
        >
          <div
            v-if="!sidebarCollapsed"
            class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-[10px] text-sm font-bold text-white"
            style="
              background: linear-gradient(
                135deg,
                var(--primary),
                oklch(0.348 0.106 295)
              );
            "
          >
            m
          </div>
          <div
            v-if="!sidebarCollapsed"
            class="text-md font-semibold tracking-tight text-foreground"
          >
            MatchHire
          </div>
        </div>
        <button
          class="flex h-8 w-8 items-center justify-center rounded-lg text-muted-foreground transition hover:bg-muted hover:text-foreground"
          @click="sidebarCollapsed = !sidebarCollapsed"
        >
          <PanelLeftClose v-if="!sidebarCollapsed" class="h-4 w-4" />
          <PanelLeftOpen v-else class="h-4 w-4" />
        </button>
      </div>

      <!-- Primary Nav -->
      <nav class="flex flex-col gap-0.5 px-3.5">
        <button
          v-for="item in menu"
          :key="item.to"
          class="flex w-full h-[0px] items-center gap-3 whitespace-nowrap rounded-xl px-3 py-5 text-sm font-medium transition-colors duration-100"
          :class="[
            sidebarCollapsed ? 'justify-center' : '',
            isActive(item.to)
              ? item.accent
                ? 'bg-primary/10 text-primary font-semibold'
                : 'bg-foreground/[0.07] text-foreground font-semibold'
              : 'text-muted-foreground hover:bg-muted hover:text-foreground',
          ]"
          @click="handleNav(item.to)"
        >
          <component
            :is="item.icon"
            class="h-4 w-4 flex-shrink-0"
            :class="isActive(item.to) && item.accent ? 'text-primary' : ''"
          />
          <span v-if="!sidebarCollapsed">{{ item.title }}</span>
        </button>
      </nav>

      <div class="flex-1" />

      <!-- Bottom Nav -->
      <nav class="flex flex-col gap-0.5 px-3.5">
        <button
          v-for="item in bottomMenu"
          :key="item.to"
          class="flex w-full items-center gap-3 whitespace-nowrap rounded-xl px-3 py-2.5 text-sm font-medium transition-colors duration-100"
          :class="[
            sidebarCollapsed ? 'justify-center' : '',
            isActive(item.to)
              ? 'bg-foreground/[0.07] text-foreground font-semibold'
              : 'text-muted-foreground hover:bg-muted hover:text-foreground',
          ]"
          @click="handleNav(item.to)"
        >
          <component :is="item.icon" class="h-4 w-4 flex-shrink-0" />
          <span v-if="!sidebarCollapsed">{{ item.title }}</span>
        </button>
      </nav>

      <!-- User footer card -->
      <div class="px-3.5 pb-4 mt-0.5">
        <div
          v-if="!sidebarCollapsed"
          class="flex items-center gap-2.5 rounded-xl bg-muted/60 px-3 py-2.5"
        >
          <div
            class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full overflow-hidden bg-muted-foreground/30 text-xs font-semibold text-foreground"
          >
            <img
              v-if="userAvatarURL"
              :src="userAvatarURL"
              alt="avatar"
              class="h-full w-full object-cover"
            />
            <span v-else>{{ userInitials }}</span>
          </div>
          <div class="min-w-0 flex-1">
            <p class="truncate text-[13px] font-medium text-foreground">
              {{ user?.name }}
            </p>
            <p class="text-[11px] capitalize text-muted-foreground">
              {{ user?.role === "recruiter" ? "Ажил олгогч" : "Ажил горилогч" }}
            </p>
          </div>
        </div>
        <div v-else class="flex justify-center">
          <div
            class="flex h-8 w-8 items-center justify-center rounded-full overflow-hidden bg-muted-foreground/30 text-xs font-semibold text-foreground"
          >
            <img
              v-if="userAvatarURL"
              :src="userAvatarURL"
              alt="avatar"
              class="h-full w-full object-cover"
            />
            <span v-else>{{ userInitials }}</span>
          </div>
        </div>
      </div>
    </aside>

    <!-- ── Main column ──────────────────────────────────────── -->
    <div class="flex min-w-0 flex-1 flex-col">
      <!-- TopBar -->
      <header
        class="flex h-16 flex-shrink-0 items-center justify-between border-b border-border bg-background px-7"
      >
        <!-- Left: collapse toggle + breadcrumb -->
        <div class="flex items-center gap-3">
          <div
            class="flex items-center gap-2 text-[14px] text-muted-foreground"
          >
            <span class="font-medium text-foreground">{{ currentCrumb }}</span>
          </div>
        </div>

        <!-- Right: actions -->
        <div class="flex items-center gap-2">
          <!-- Dark mode toggle -->
          <button
            class="flex h-9 w-9 items-center justify-center rounded-full border border-border bg-background text-muted-foreground transition hover:bg-muted hover:text-foreground"
            @click="toggleDark"
          >
            <Sun v-if="isDark" class="h-[15px] w-[15px]" />
            <Moon v-else class="h-[15px] w-[15px]" />
          </button>

          <!-- Notifications -->
          <button
            class="relative flex h-9 w-9 items-center justify-center rounded-full border border-border bg-background text-muted-foreground transition hover:bg-muted hover:text-foreground"
            @click="markAllRead"
          >
            <Bell class="h-[15px] w-[15px]" />
            <span
              v-if="unreadCount > 0"
              class="absolute right-1.5 top-1.5 flex h-4 w-4 items-center justify-center rounded-full bg-primary text-[9px] font-bold text-white"
              >{{ unreadCount > 9 ? "9+" : unreadCount }}</span
            >
          </button>

          <!-- User pill -->
          <button
            class="flex items-center gap-2 rounded-full border border-border bg-background px-3 py-1.5 text-sm font-medium text-foreground transition hover:bg-muted"
            @click="handleNav('/profile')"
          >
            <div
              class="flex h-6 w-6 items-center justify-center rounded-full overflow-hidden bg-muted-foreground/30 text-[10px] font-semibold"
            >
              <img
                v-if="userAvatarURL"
                :src="userAvatarURL"
                alt="avatar"
                class="h-full w-full object-cover"
              />
              <span v-else>{{ userInitials }}</span>
            </div>
            <span class="capitalize">{{
              user?.role === "recruiter" ? "Олгогч" : "Горилогч"
            }}</span>
            <ChevronDown class="h-3 w-3 text-muted-foreground" />
          </button>

          <!-- Logout -->
          <button
            class="flex h-9 w-9 items-center justify-center rounded-full border border-border bg-background text-muted-foreground transition hover:bg-muted hover:text-foreground"
            :disabled="isLoading"
            @click="handleLogout"
          >
            <LogOut class="h-4 w-4" />
          </button>
        </div>
      </header>

      <!-- Page content -->
      <main class="flex-1 overflow-auto bg-background">
        <div class="mx-auto max-w-7xl px-7 py-6">
          <NuxtPage />
        </div>
      </main>
    </div>
  </div>
</template>
