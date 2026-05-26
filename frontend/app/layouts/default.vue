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
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  SidebarProvider,
  SidebarRail,
  SidebarTrigger,
  useSidebar,
} from "~/components/ui/sidebar";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "~/components/ui/collapsible";

const { user, logout, isLoading } = useAuth();
const router = useRouter();
const route = useRoute();
const { isDark, toggle: toggleDark, init: initDark } = useDarkMode();
const sidebarOpen = ref(true);
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
  { title: "Даалгавар", icon: ClipboardList, to: "/tasks", accent: true },
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
  { title: "Ярилцлага", icon: Calendar, to: "/interviews", accent: true },
  {
    title: "Даалгавар",
    icon: ClipboardList,
    accent: true,
    children: [
      { title: "Даалгаврын сан", to: "/tasks", accent: true },
      { title: "Даалгаврын хариу", to: "/submissions", accent: true },
    ],
  },
];

const menu = computed(() => {
  if (user?.value?.role === "recruiter") return menuRecruiter;
  return menuUser;
});

const bottomMenu = computed(() => [
  {
    title: user?.value?.role === "recruiter" ? "Компани" : "Профайл",
    icon: User,
    to: "/profile",
    accent: true,
  },
  { title: "Тохиргоо", icon: Settings, to: "/settings", accent: true },
]);

const recruiterNeedsCompany = computed(
  () => user.value?.role === "recruiter" && !user.value?.company_id,
);

function isActive(path: string) {
  if (path === "/") return route.path === "/";
  return route.path.startsWith(path);
}

function isGroupActive(children: { to: string }[]) {
  return children.some((c) => route.path.startsWith(c.to));
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
  <SidebarProvider v-model:open="sidebarOpen">
    <Sidebar
      collapsible="icon"
      class="border-r border-border px-2 group-data-[collapsible=icon]:px-0"
    >
      <!-- ── Header ─────────────────────────────────────────── -->
      <SidebarHeader class="p-0 py-2">
        <div
          class="flex items-center justify-between px-3.5 pb-2 pt-4 group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-2"
        >
          <div
            class="flex cursor-pointer items-center gap-2 group-data-[collapsible=icon]:invisible group-data-[collapsible=icon]:w-0 group-data-[collapsible=icon]:overflow-hidden group-data-[collapsible=icon]:pointer-events-none"
            @click="router.push('/')"
          >
            <!-- <div
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
            </div> -->
            <img
              class="block h-8 w-8"
              src="/icons/logo.svg"
              alt="Skillz logo"
            />
            <div class="text-md font-semibold tracking-tight text-foreground">
              Skillz
            </div>
          </div>
          <SidebarTrigger />
        </div>
      </SidebarHeader>

      <!-- ── Content ────────────────────────────────────────── -->
      <SidebarContent>
        <!-- Primary menu group -->
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <template v-for="item in menu" :key="item.to || item.title">
                <!-- Group with children -->
                <Collapsible
                  v-if="item.children"
                  as-child
                  :default-open="isGroupActive(item.children)"
                  class="group/collapsible"
                >
                  <SidebarMenuItem>
                    <CollapsibleTrigger as-child>
                      <SidebarMenuButton
                        :tooltip="item.title"
                        @click="!sidebarOpen && (sidebarOpen = true)"
                      >
                        <component :is="item.icon" />
                        <span>{{ item.title }}</span>
                        <ChevronRight
                          class="ml-auto h-3 w-3 transition-transform duration-150 group-data-[state=open]/collapsible:rotate-90"
                        />
                      </SidebarMenuButton>
                    </CollapsibleTrigger>
                    <CollapsibleContent>
                      <SidebarMenuSub>
                        <SidebarMenuSubItem
                          v-for="child in item.children"
                          :key="child.to"
                        >
                          <SidebarMenuSubButton
                            :is-active="isActive(child.to)"
                            @click="handleNav(child.to)"
                            class="cursor-pointer"
                          >
                            <span>{{ child.title }}</span>
                          </SidebarMenuSubButton>
                        </SidebarMenuSubItem>
                      </SidebarMenuSub>
                    </CollapsibleContent>
                  </SidebarMenuItem>
                </Collapsible>

                <!-- Regular item -->
                <SidebarMenuItem v-else>
                  <SidebarMenuButton
                    :tooltip="item.title"
                    :is-active="isActive(item.to!)"
                    :class="[
                      isActive(item.to!) && item.accent
                        ? 'bg-primary/10 text-primary font-semibold hover:bg-primary/15 hover:text-primary data-[active=true]:bg-primary/10 data-[active=true]:text-primary'
                        : '',
                    ]"
                    class="cursor-pointer"
                    @click="handleNav(item.to!)"
                  >
                    <component
                      :is="item.icon"
                      :class="
                        isActive(item.to!) && item.accent ? 'text-primary' : ''
                      "
                    />
                    <span>{{ item.title }}</span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </template>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        <!-- Spacer pushes bottom group down -->
        <div class="flex-1" />

        <!-- Bottom menu group -->
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <!-- Notifications -->
              <!-- <SidebarMenuItem>
                <Popover v-model:open="notifOpen">
                  <PopoverTrigger as-child>
                    <SidebarMenuButton
                      tooltip="Мэдэгдэл"
                      class="cursor-pointer relative"
                    >
                      <Bell />
                      <span>Мэдэгдэл</span>
                      <span
                        v-if="unreadCount > 0"
                        class="ml-auto flex h-5 w-5 items-center justify-center rounded-full bg-primary text-[10px] font-bold text-white group-data-[collapsible=icon]:hidden"
                      >
                        {{ unreadCount > 9 ? "9+" : unreadCount }}
                      </span>
                      <span
                        v-if="unreadCount > 0"
                        class="absolute right-1.5 top-1.5 h-2 w-2 rounded-full bg-primary hidden group-data-[collapsible=icon]:block"
                      />
                    </SidebarMenuButton>
                  </PopoverTrigger>
                  <PopoverContent
                    align="end"
                    side="right"
                    class="w-80 p-0 shadow-lg"
                  >
                    <div
                      class="flex items-center justify-between border-b border-border px-4 py-3"
                    >
                      <span class="text-sm font-semibold text-foreground">
                        Мэдэгдлүүд
                      </span>
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
                          <p
                            class="truncate text-[13px] font-medium text-foreground"
                          >
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
              </SidebarMenuItem> -->

              <!-- Other bottom items -->
              <SidebarMenuItem v-for="item in bottomMenu" :key="item.to">
                <SidebarMenuButton
                  :tooltip="item.title"
                  :is-active="isActive(item.to)"
                  class="cursor-pointer"
                  @click="handleNav(item.to)"
                >
                  <component :is="item.icon" />
                  <span>{{ item.title }}</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <!-- ── Footer ─────────────────────────────────────────── -->
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <SidebarMenuButton
                  size="lg"
                  class="cursor-pointer data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
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
                  <div class="grid flex-1 text-left text-sm leading-tight">
                    <span class="truncate text-[13px] font-medium">
                      {{ user?.name }}
                    </span>
                    <span
                      class="truncate text-[11px] capitalize text-muted-foreground"
                    >
                      {{
                        user?.role === "recruiter"
                          ? "Ажил олгогч"
                          : "Ажил горилогч"
                      }}
                    </span>
                  </div>
                </SidebarMenuButton>
              </DropdownMenuTrigger>
              <DropdownMenuContent side="right" align="start" class="w-52 mb-3">
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
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>

      <SidebarRail />
    </Sidebar>

    <!-- ── Main column ──────────────────────────────────────── -->
    <SidebarInset>
      <main class="flex-1 bg-background">
        <NuxtPage v-if="useRoute().meta.fullscreen" />
        <div v-else class="mx-auto w-full px-7 py-6">
          <NuxtPage />
        </div>
      </main>
    </SidebarInset>
  </SidebarProvider>
</template>
