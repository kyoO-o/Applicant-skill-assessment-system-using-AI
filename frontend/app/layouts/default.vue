<script setup lang="ts">
import {
  Calendar,
  ChartNoAxesColumnIncreasing,
  BriefcaseBusiness,
  User,
  LayoutGrid,
  Settings,
} from "lucide-vue-next";

const { user } = useAuth();

const menuUser = [
  { title: "Тайлан", icon: ChartNoAxesColumnIncreasing, to: "/" },
  { title: "Ажлын зар", icon: BriefcaseBusiness, to: "/jobs" },
  { title: "My zone", icon: LayoutGrid, to: "/my-zone" },
  { title: "Profile", icon: User, to: "/profile" },
  { title: "Settings", icon: Settings, to: "/settings" },
];

const menuRecruiter = [
  { title: "Тайлан", icon: ChartNoAxesColumnIncreasing, to: "/" },
  { title: "Миний зар", icon: BriefcaseBusiness, to: "/jobs" },
  { title: "Уулзалт", icon: Calendar, to: "/applicants" },
  // { title: "Task", icon: "clipboard-list", to: "/task" },
  { title: "Profile", icon: User, to: "/profile" },
  { title: "Settings", icon: Settings, to: "/settings" },
];

const menu = computed(() => {
  if (user?.value?.role === "recruiter") return menuRecruiter;
  return menuUser;
});
</script>
<template>
  <div class="min-h-screen bg-background text-foreground">
    <SidebarProvider>
      <Sidebar class="border-r bg-background/95 backdrop-blur">
        <!-- Header -->
        <SidebarHeader class="px-4 py-5 border-b">
          <div class="flex items-center gap-3">
            <div
              class="flex h-10 w-10 items-center justify-center rounded-xl border bg-muted shadow-sm"
            >
              <span class="text-lg font-bold">S</span>
            </div>

            <div>
              <h1 class="text-base font-semibold leading-none">Skillz</h1>
              <p class="mt-1 text-xs text-muted-foreground">AI Recruitment</p>
            </div>
          </div>
        </SidebarHeader>

        <!-- Content -->
        <SidebarContent class="px-1 py-2">
          <SidebarGroup>
            <SidebarMenu class="gap-0">
              <SidebarMenuItem v-for="item in menu" :key="item.title">
                <SidebarMenuButton
                  @click="$router.push(item.to)"
                  class="group flex w-full items-center m-0 rounded-xl px-3 py-5 text-sm font-medium text-muted-foreground hover:bg-accent hover:text-accent-foreground"
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

        <!-- Footer -->
        <SidebarFooter class="mt-auto border-t px-4 py-4">
          <div class="rounded-xl bg-muted/50 px-3 py-3">
            <p class="text-xs font-medium">Logged in</p>
            <p class="mt-1 text-xs text-muted-foreground">
              {{ user?.role || "Guest" }}
            </p>
          </div>
        </SidebarFooter>
      </Sidebar>
    </SidebarProvider>

    <main class="min-h-screen">
      <slot />
    </main>
  </div>
</template>
