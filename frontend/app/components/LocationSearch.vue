<script setup lang="ts">
import { Search, MapPin, X } from "lucide-vue-next";

const props = defineProps<{
  modelX: string;
  modelY: string;
  disabled?: boolean;
}>();

const emit = defineEmits<{
  "update:modelX": [value: string];
  "update:modelY": [value: string];
}>();

type NominatimResult = {
  place_id: number;
  display_name: string;
  lat: string;
  lon: string;
};

const query = ref("");
const results = ref<NominatimResult[]>([]);
const showDropdown = ref(false);
let debounceTimer: ReturnType<typeof setTimeout>;

const mapContainer = ref<HTMLDivElement | null>(null);
let leafletMap: any = null;
let marker: any = null;

// Treat 0/0 as "no coordinates set" (DB default when never saved)
function isValid(x: string, y: string) {
  const lat = Number(x);
  const lng = Number(y);
  return !!x && !!y && (lat !== 0 || lng !== 0);
}

const hasCoordinates = computed(() => isValid(props.modelX, props.modelY));

const DEFAULT_LAT = 47.9184;
const DEFAULT_LNG = 106.9177;

onMounted(async () => {
  if (import.meta.server) return;
  const L = await import("leaflet");
  await import("leaflet/dist/leaflet.css");

  delete (L.Icon.Default.prototype as any)._getIconUrl;
  L.Icon.Default.mergeOptions({
    iconUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png",
    iconRetinaUrl:
      "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png",
    shadowUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png",
  });

  if (!mapContainer.value) return;

  const hasInit = isValid(props.modelX, props.modelY);
  const initLat = hasInit ? Number(props.modelX) : DEFAULT_LAT;
  const initLng = hasInit ? Number(props.modelY) : DEFAULT_LNG;

  leafletMap = L.map(mapContainer.value, { zoomControl: true }).setView(
    [initLat, initLng],
    hasInit ? 15 : 12,
  );

  L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
    attribution: "© OpenStreetMap",
    maxZoom: 19,
  }).addTo(leafletMap);

  if (hasInit) {
    marker = L.marker([initLat, initLng], { draggable: true }).addTo(
      leafletMap,
    );
    marker.on("dragend", () => {
      const pos = marker.getLatLng();
      emit("update:modelX", String(pos.lat));
      emit("update:modelY", String(pos.lng));
    });
  }

  leafletMap.on("click", (e: any) => {
    if (props.disabled) return;
    const { lat, lng } = e.latlng;
    emit("update:modelX", String(lat));
    emit("update:modelY", String(lng));
    if (marker) {
      marker.setLatLng([lat, lng]);
    } else {
      marker = L.marker([lat, lng], { draggable: true }).addTo(leafletMap);
      marker.on("dragend", () => {
        const pos = marker.getLatLng();
        emit("update:modelX", String(pos.lat));
        emit("update:modelY", String(pos.lng));
      });
    }
  });

  // Use browser geolocation as default center when no coordinates set
  if (!hasInit && navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      (pos) =>
        leafletMap?.setView([pos.coords.latitude, pos.coords.longitude], 14),
      () => {
        /* denied — stay on UB */
      },
      { timeout: 5000 },
    );
  }
});

watch([() => props.modelX, () => props.modelY], async ([x, y]) => {
  if (!leafletMap || !isValid(x, y)) return;
  const lat = Number(x);
  const lng = Number(y);
  leafletMap.setView([lat, lng], 15);
  if (marker) {
    marker.setLatLng([lat, lng]);
  } else {
    const L = await import("leaflet");
    marker = L.marker([lat, lng], { draggable: true }).addTo(leafletMap);
    marker.on("dragend", () => {
      const pos = marker.getLatLng();
      emit("update:modelX", String(pos.lat));
      emit("update:modelY", String(pos.lng));
    });
  }
});

onUnmounted(() => {
  leafletMap?.remove();
  leafletMap = null;
  marker = null;
});

async function search() {
  const q = query.value.trim();
  if (q.length < 2) {
    results.value = [];
    showDropdown.value = false;
    return;
  }
  try {
    const data = await $fetch<NominatimResult[]>(
      `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(q)}&format=json&limit=6&accept-language=mn,en`,
      { headers: { "User-Agent": "SkillAssessmentApp/1.0" } },
    );
    results.value = data;
    showDropdown.value = data.length > 0;
  } catch {
    /* ignore */
  }
}

function onInput() {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(search, 400);
}

function select(r: NominatimResult) {
  emit("update:modelX", r.lat);
  emit("update:modelY", r.lon);
  query.value = r.display_name.split(",").slice(0, 2).join(",").trim();
  showDropdown.value = false;
}

function hideDropdown() {
  setTimeout(() => {
    showDropdown.value = false;
  }, 150);
}

function clear() {
  emit("update:modelX", "");
  emit("update:modelY", "");
  query.value = "";
  results.value = [];
  if (marker) {
    marker.remove();
    marker = null;
  }
}
</script>

<template>
  <div class="space-y-3">
    <div class="relative">
      <Search
        class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground"
      />
      <Input
        v-model="query"
        :disabled="disabled"
        placeholder="Хаяг хайх... (жишээ: Сүхбаатар дүүрэг)"
        class="pl-9 pr-8"
        @input="onInput"
        @focus="showDropdown = results.length > 0"
        @blur="hideDropdown"
      />
      <button
        v-if="query && !disabled"
        type="button"
        tabindex="-1"
        class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground transition hover:text-foreground"
        @click="clear"
      >
        <X class="h-3.5 w-3.5" />
      </button>
      <div
        v-if="showDropdown"
        class="absolute z-50 mt-1 w-full overflow-hidden rounded-2xl border border-border bg-popover"
      >
        <button
          v-for="r in results"
          :key="r.place_id"
          type="button"
          class="flex w-full items-start gap-2 px-3 py-2.5 text-left transition-colors hover:bg-muted"
          @mousedown.prevent="select(r)"
        >
          <MapPin class="mt-0.5 h-3.5 w-3.5 shrink-0 text-muted-foreground" />
          <span class="line-clamp-2 text-xs">{{ r.display_name }}</span>
        </button>
      </div>
    </div>

    <div
      ref="mapContainer"
      class="h-64 w-full overflow-hidden rounded-2xl border border-border"
      :class="disabled ? 'pointer-events-none opacity-60' : 'cursor-crosshair'"
    />

    <div v-if="hasCoordinates" class="flex items-center gap-2 text-xs">
      <MapPin class="h-3.5 w-3.5 shrink-0 text-green-600" />
      <span class="font-medium text-green-700">
        {{ Number(modelX).toFixed(5) }}, {{ Number(modelY).toFixed(5) }}
      </span>
      <button
        v-if="!disabled"
        type="button"
        class="ml-auto text-muted-foreground transition hover:text-foreground"
        @click="clear"
      >
        <X class="h-3.5 w-3.5" />
      </button>
    </div>
  </div>
</template>
