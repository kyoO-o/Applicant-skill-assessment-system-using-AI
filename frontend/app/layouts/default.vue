<script setup lang="ts">
import {
  ArrowRight,
  Calendar,
  ChartNoAxesColumnIncreasing,
  BriefcaseBusiness,
  User,
  LayoutGrid,
  Settings,
  LogOut,
  PanelLeft,
} from "lucide-vue-next";

const { user, logout, isLoading } = useAuth();
const router = useRouter();
const route = useRoute();

const menuUser = [
  { title: "Dashboard", icon: ChartNoAxesColumnIncreasing, to: "/" },
  { title: "Jobs", icon: BriefcaseBusiness, to: "/jobs" },
  { title: "My Zone", icon: LayoutGrid, to: "/my-zone" },
  { title: "Profile", icon: User, to: "/profile" },
  { title: "Settings", icon: Settings, to: "/settings" },
];

const menuRecruiter = [
  { title: "Dashboard", icon: ChartNoAxesColumnIncreasing, to: "/" },
  { title: "My Jobs", icon: BriefcaseBusiness, to: "/jobs" },
  { title: "Interviews", icon: Calendar, to: "/interviews" },
  { title: "Profile", icon: User, to: "/profile" },
  { title: "Settings", icon: Settings, to: "/settings" },
];

const menu = computed(() => {
  if (user?.value?.role === "recruiter") return menuRecruiter;
  return menuUser;
});

const currentPage = computed(() => {
  return menu.value.find((item) => item.to === route.path)?.title || "Workspace";
});

const currentSubtitle = computed(() => {
  if (user.value?.role === "recruiter") {
    return "Review your hiring pipeline, manage postings, and keep interviews moving.";
  }

  return "Track your applications, discover matches, and stay ready for the next step.";
});

async function handleLogout() {
  await logout();
  await router.push("/login");
}
</script>
<template>
  <div class="min-h-screen bg-background text-foreground">
    <SidebarProvider>
      <Sidebar class="border-r border-border bg-card">
        <SidebarHeader class="border-b border-border px-4 py-5">
          <div class="flex items-center gap-3">
            <div
              class="flex h-10 w-10 items-center justify-center rounded-2xl border border-border bg-foreground text-background shadow-sm"
            >
              <span class="text-sm font-semibold">MH</span>
            </div>

            <div>
              <h1 class="text-base font-semibold leading-none">MatchHire</h1>
              <p class="mt-1 text-xs text-muted-foreground">
                Applicant skill assessment
              </p>
            </div>
          </div>
        </SidebarHeader>

        <SidebarContent class="px-1 py-2">
          <SidebarGroup>
            <SidebarMenu class="gap-0">
              <SidebarMenuItem v-for="item in menu" :key="item.title">
                <SidebarMenuButton
                  @click="$router.push(item.to)"
                  class="group m-0 flex w-full items-center gap-3 rounded-2xl px-3 py-5 text-sm font-medium text-muted-foreground transition hover:bg-muted hover:text-foreground"
                  :data-active="$route.path === item.to"
                >
                  <component :is="item.icon" class="h-5 w-5 shrink-0" />

                  <span class="truncate">
                    {{ item.title }}
                  </span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroup>
        </SidebarContent>

        <SidebarFooter class="mt-auto border-t px-4 py-4">
          <div class="rounded-2xl border border-border bg-muted/40 px-3 py-3">
            <p class="text-xs font-medium text-muted-foreground">Logged in as</p>
            <p class="mt-1 text-sm font-semibold">{{ user?.name }}</p>
            <p class="mt-1 text-xs capitalize text-muted-foreground">
              {{ user?.role }}
            </p>
          </div>

          <Button
            variant="outline"
            class="mt-3 w-full justify-start gap-2"
            :disabled="isLoading"
            @click="handleLogout"
          >
            <LogOut class="h-4 w-4" />
            {{ isLoading ? "Logging out..." : "Log out" }}
          </Button>
        </SidebarFooter>
      </Sidebar>

      <div class="flex min-h-screen flex-1 flex-col bg-background">
        <header class="border-b border-border bg-card/80 px-5 py-4 backdrop-blur">
          <div
            class="mx-auto flex w-full max-w-7xl flex-col gap-4 md:flex-row md:items-center md:justify-between"
          >
            <div class="flex items-start gap-3">
              <div class="rounded-2xl border border-border bg-muted/40 p-2.5">
                <PanelLeft class="h-4 w-4" />
              </div>
              <div>
                <p class="text-xs font-medium uppercase tracking-[0.18em] text-muted-foreground">
                  {{ user?.role === "recruiter" ? "Recruiter Workspace" : "Applicant Workspace" }}
                </p>
                <h1 class="mt-1 text-2xl font-semibold tracking-tight">
                  {{ currentPage }}
                </h1>
                <p class="mt-1 text-sm text-muted-foreground">
                  {{ currentSubtitle }}
                </p>
              </div>
            </div>

            <NuxtLink
              to="/jobs"
              class="inline-flex items-center gap-2 self-start rounded-full border border-border bg-background px-4 py-2 text-sm font-medium transition hover:bg-muted md:self-auto"
            >
              {{ user?.role === "recruiter" ? "Open job manager" : "Explore jobs" }}
              <ArrowRight class="h-4 w-4" />
            </NuxtLink>
          </div>
        </header>

        <main class="mx-auto flex w-full max-w-7xl flex-1 flex-col px-5 py-6">
          <NuxtPage />
        </main>
      </div>
    </SidebarProvider>
  </div>
</template>
