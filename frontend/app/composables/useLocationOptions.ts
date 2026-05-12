import locationMap from "../assets/location.json";

type LocationMap = Record<string, string[]>;

const typedLocationMap = locationMap as LocationMap;

export function useLocationOptions() {
  const cities = Object.keys(typedLocationMap);

  function districtsFor(city?: string | null) {
    if (!city) {
      return [];
    }

    return typedLocationMap[city] ?? [];
  }

  return {
    cities,
    districtsFor,
  };
}
