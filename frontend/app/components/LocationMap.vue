<script setup lang="ts">
const props = defineProps<{
  lat: number;
  lng: number;
  label?: string;
}>();

const mapContainer = ref<HTMLDivElement | null>(null);
let leafletMap: any = null;

onMounted(async () => {
  if (import.meta.server || !mapContainer.value) return;
  const L = await import("leaflet");
  await import("leaflet/dist/leaflet.css");

  delete (L.Icon.Default.prototype as any)._getIconUrl;
  L.Icon.Default.mergeOptions({
    iconUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png",
    iconRetinaUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png",
    shadowUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png",
  });

  leafletMap = L.map(mapContainer.value, {
    zoomControl: true,
    dragging: true,
    scrollWheelZoom: false,
  }).setView([props.lat, props.lng], 15);

  L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
    attribution: "© OpenStreetMap",
    maxZoom: 19,
  }).addTo(leafletMap);

  const m = L.marker([props.lat, props.lng]).addTo(leafletMap);
  if (props.label) m.bindPopup(props.label).openPopup();
});

onUnmounted(() => {
  leafletMap?.remove();
  leafletMap = null;
});
</script>

<template>
  <div ref="mapContainer" class="h-52 w-full overflow-hidden rounded-2xl border border-border" />
</template>
