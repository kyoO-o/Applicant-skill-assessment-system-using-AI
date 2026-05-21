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
  ChevronRight,
  Users,
  ChevronLeft,
  PanelLeftClose,
  PanelLeftOpen,
} from "lucide-vue-next";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "~/components/ui/popover";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { icon } from "leaflet";

const { user, logout, isLoading } = useAuth();
const router = useRouter();
const route = useRoute();
const { isDark, toggle: toggleDark, init: initDark } = useDarkMode();
const {
  notifications,
  unreadCount,
  connect: connectWS,
  disconnect: disconnectWS,
  markAsRead,
  markAllRead,
  clear: clearNotifications,
} = useNotifications();

const notifOpen = ref(false);

function formatTimeAgo(date: Date): string {
  const diff = Math.floor((Date.now() - new Date(date).getTime()) / 1000);
  if (diff < 60) return `${diff}с өмнө`;
  if (diff < 3600) return `${Math.floor(diff / 60)}м өмнө`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}ц өмнө`;
  return `${Math.floor(diff / 86400)}х өмнө`;
}

function handleNotifClick(n: (typeof notifications.value)[0]) {
  markAsRead(n.id);
  notifOpen.value = false;
  router.push(n.route);
}

onMounted(() => {
  initDark();
  if (user.value) connectWS();
});

watch(
  () => user.value,
  (u) => {
    if (!u) {
      disconnectWS();
      clearNotifications();
    }
  },
);

const sidebarCollapsed = ref(false);

type MenuItem = {
  title: string;
  icon?: any;
  to?: string;
  accent?: boolean;
  children?: { title: string; to: string; accent?: boolean }[];
};

const menuUser: MenuItem[] = [
  { title: "Нүүр хуудас", icon: LayoutDashboard, to: "/", accent: true },
  {
    title: "Ажлын байрууд",
    icon: BriefcaseBusiness,
    to: "/jobs",
    accent: true,
  },
  {
    title: "Миний анкетууд",
    icon: FileText,
    to: "/applications",
    accent: true,
  },
  { title: "CV Бүтээгч", icon: ScrollText, to: "/cv", accent: true },
  { title: "Ярилцлага", icon: Calendar, to: "/interviews", accent: true },
  {
    title: "Даалгавар",
    icon: ClipboardList,
    to: "/tasks",
    accent: true,
  },
  { title: "AI Чатбот", icon: Bot, to: "/chat", accent: true },
];

const menuRecruiter: MenuItem[] = [
  { title: "Нүүр хуудас", icon: LayoutDashboard, to: "/", accent: true },
  {
    title: "Ажлын байрууд",
    icon: BriefcaseBusiness,
    to: "/jobs",
    accent: true,
  },
  // { title: "Горилогчид", icon: Users, to: "/applications" },
  { title: "Ярилцлага", icon: Calendar, to: "/interviews", accent: true },
  {
    title: "Даалгавар",
    icon: ClipboardList,
    children: [
      { title: "Даалгаврын сан", to: "/tasks", accent: true },
      { title: "Даалгаврын хариу", to: "/submissions", accent: true },
    ],
    accent: true,
  },
];

const openGroups = ref<string[]>([]);

function toggleGroup(title: string) {
  const idx = openGroups.value.indexOf(title);
  if (idx >= 0) openGroups.value.splice(idx, 1);
  else openGroups.value.push(title);
}

function isGroupOpen(title: string) {
  return openGroups.value.includes(title);
}

function isGroupActive(children: { to: string }[]) {
  return children.some((c) => route.path.startsWith(c.to));
}

const bottomMenu = [
  { title: "Мэдэгдэл", icon: Bell },
  {
    title: user?.value?.role === "recruiter" ? "Компани" : "Профайл",
    icon: User,
    to: "/profile",
    accent: true,
  },
  { title: "Тохиргоо", icon: Settings, to: "/settings", accent: true },
];

const menu = computed(() => {
  if (user?.value?.role === "recruiter") return menuRecruiter;
  return menuUser;
});

watch(
  () => route.path,
  () => {
    for (const item of menu.value) {
      if (item.children && isGroupActive(item.children)) {
        if (!openGroups.value.includes(item.title)) {
          openGroups.value.push(item.title);
        }
      }
    }
  },
  { immediate: true },
);

const currentCrumb = computed(() => {
  for (const item of [...menu.value, ...bottomMenu]) {
    if (item.to && item.to === route.path) return item.title;
    if (item.children) {
      const child = item.children.find((c) => route.path.startsWith(c.to));
      if (child) return child.title;
    }
  }
  return "Dashboard";
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
        <template v-for="item in menu" :key="item.to || item.title">
          <!-- Group item with children -->
          <template v-if="item.children">
            <button
              class="flex w-full h-[0px] items-center gap-3 whitespace-nowrap rounded-xl px-3 py-5 text-sm font-medium transition-colors duration-100"
              :class="[
                sidebarCollapsed ? 'justify-center' : '',
                'text-muted-foreground hover:bg-muted hover:text-foreground',
              ]"
              @click="
                sidebarCollapsed
                  ? (sidebarCollapsed = false)
                  : toggleGroup(item.title)
              "
            >
              <component :is="item.icon" class="h-4 w-4 flex-shrink-0" />
              <span v-if="!sidebarCollapsed">{{ item.title }}</span>
              <ChevronRight
                v-if="!sidebarCollapsed"
                class="ml-auto h-3 w-3 transition-transform duration-150"
                :class="isGroupOpen(item.title) ? 'rotate-90' : ''"
              />
            </button>
            <div
              v-if="isGroupOpen(item.title) && !sidebarCollapsed"
              class="ml-3 flex flex-col gap-0.5 border-l border-border pl-3"
            >
              <button
                v-for="child in item.children"
                :key="child.to"
                class="flex w-full items-center gap-3 whitespace-nowrap rounded-lg px-2 py-2 text-sm font-medium transition-colors duration-100"
                :class="
                  isActive(child.to)
                    ? 'bg-foreground/[0.07] text-foreground font-semibold'
                    : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                "
                @click="handleNav(child.to)"
              >
                {{ child.title }}
              </button>
            </div>
          </template>

          <!-- Regular item -->
          <button
            v-else
            class="flex w-full h-[0px] items-center gap-3 whitespace-nowrap rounded-xl px-3 py-5 text-sm font-medium transition-colors duration-100"
            :class="[
              sidebarCollapsed ? 'justify-center' : '',
              isActive(item.to!)
                ? item.accent
                  ? 'bg-primary/10 text-primary font-semibold'
                  : 'bg-foreground/[0.07] text-foreground font-semibold'
                : 'text-muted-foreground hover:bg-muted hover:text-foreground',
            ]"
            @click="handleNav(item.to!)"
          >
            <component
              :is="item.icon"
              class="h-4 w-4 flex-shrink-0"
              :class="isActive(item.to!) && item.accent ? 'text-primary' : ''"
            />
            <span v-if="!sidebarCollapsed">{{ item.title }}</span>
          </button>
        </template>
      </nav>

      <div class="flex-1" />

      <!-- Bottom Nav -->
      <nav class="flex flex-col gap-0.5 px-3.5">
        <template v-for="item in bottomMenu" :key="item.to || item.title">
          <!-- Notification item -->
          <Popover v-if="!item.to" v-model:open="notifOpen">
            <PopoverTrigger as-child>
              <button
                class="relative flex w-full items-center gap-3 whitespace-nowrap rounded-xl px-3 py-2.5 text-sm font-medium transition-colors duration-100 text-muted-foreground hover:bg-muted hover:text-foreground"
                :class="sidebarCollapsed ? 'justify-center' : ''"
              >
                <component :is="item.icon" class="h-4 w-4 flex-shrink-0" />
                <span v-if="!sidebarCollapsed">{{ item.title }}</span>
                <span
                  v-if="unreadCount > 0 && !sidebarCollapsed"
                  class="ml-auto flex h-5 w-5 items-center justify-center rounded-full bg-primary text-[10px] font-bold text-white"
                  >{{ unreadCount > 9 ? "9+" : unreadCount }}</span
                >
                <span
                  v-if="unreadCount > 0 && sidebarCollapsed"
                  class="absolute right-1.5 top-1.5 h-2 w-2 rounded-full bg-primary"
                />
              </button>
            </PopoverTrigger>
            <PopoverContent align="end" side="right" class="w-80 p-0 shadow-lg">
              <div
                class="flex items-center justify-between border-b border-border px-4 py-3"
              >
                <span class="text-sm font-semibold text-foreground"
                  >Мэдэгдлүүд</span
                >
                <button
                  v-if="unreadCount > 0"
                  class="text-[11px] text-primary hover:underline"
                  @click="markAllRead"
                >
                  Бүгдийг уншсан
                </button>
              </div>
              <div class="max-h-[360px] overflow-y-auto">
                <div
                  v-if="notifications.length === 0"
                  class="flex flex-col items-center justify-center gap-2 py-10 text-muted-foreground"
                >
                  <Bell class="h-6 w-6 opacity-40" />
                  <span class="text-xs">Мэдэгдэл байхгүй</span>
                </div>
                <button
                  v-for="n in notifications"
                  :key="n.id"
                  class="flex w-full items-start gap-3 border-b border-border px-4 py-3 text-left transition hover:bg-muted/50"
                  :class="n.read ? 'opacity-60' : ''"
                  @click="handleNotifClick(n)"
                >
                  <span
                    class="mt-1.5 h-2 w-2 flex-shrink-0 rounded-full"
                    :class="n.read ? 'bg-transparent' : 'bg-primary'"
                  />
                  <div class="min-w-0 flex-1">
                    <p class="truncate text-[13px] font-medium text-foreground">
                      {{ n.title }}
                    </p>
                    <p
                      class="mt-0.5 text-[12px] leading-snug text-muted-foreground line-clamp-2"
                    >
                      {{ n.body }}
                    </p>
                    <p class="mt-1 text-[11px] text-muted-foreground/60">
                      {{ formatTimeAgo(n.timestamp) }}
                    </p>
                  </div>
                </button>
              </div>
            </PopoverContent>
          </Popover>

          <!-- Regular nav item -->
          <button
            v-else
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
        </template>
      </nav>

      <!-- User footer card -->
      <div class="px-3.5 pb-4 mt-0.5">
        <DropdownMenu>
          <DropdownMenuTrigger as-child>
            <button
              class="flex w-full items-center gap-2.5 rounded-xl bg-muted/60 px-3 py-2.5 text-left transition hover:bg-muted cursor-pointer"
              :class="sidebarCollapsed ? 'justify-center' : ''"
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
              <div v-if="!sidebarCollapsed" class="min-w-0 flex-1">
                <p class="truncate text-[13px] font-medium text-foreground">
                  {{ user?.name }}
                </p>
                <p class="text-[11px] capitalize text-muted-foreground">
                  {{
                    user?.role === "recruiter" ? "Ажил олгогч" : "Ажил горилогч"
                  }}
                </p>
              </div>
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent side="right" align="start" class="w-52 mb-5">
            <DropdownMenuItem class="cursor-pointer" @click="toggleDark">
              <Sun v-if="isDark" class="mr-2 h-4 w-4" />
              <Moon v-else class="mr-2 h-4 w-4" />
              <span>{{ isDark ? "Гэрэлт горим" : "Харанхуй горим" }}</span>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              class="cursor-pointer text-destructive focus:text-destructive"
              :disabled="isLoading"
              @click="handleLogout"
            >
              <LogOut class="mr-2 h-4 w-4" />
              <span>Гарах</span>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </aside>

    <!-- ── Main column ──────────────────────────────────────── -->
    <div class="flex min-w-0 flex-1 flex-col">
      <!-- Page content -->
      <main class="flex-1 overflow-auto bg-background">
        <NuxtPage v-if="useRoute().meta.fullscreen" />
        <div v-else class="mx-auto max-w-7xl px-7 py-6">
          <NuxtPage />
        </div>
      </main>
    </div>
  </div>
</template>
