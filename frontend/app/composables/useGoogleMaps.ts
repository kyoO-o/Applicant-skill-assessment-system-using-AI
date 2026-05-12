export function useGoogleMaps() {
  function extractCoordinatesFromUrl(url: string) {
    const value = url.trim();
    if (!value) {
      return null;
    }

    const atMatch = value.match(/@(-?\d+(?:\.\d+)?),(-?\d+(?:\.\d+)?)/);
    if (atMatch) {
      return {
        lat: Number(atMatch[1]),
        lng: Number(atMatch[2]),
      };
    }

    const queryMatch = value.match(/[?&]q=(-?\d+(?:\.\d+)?),(-?\d+(?:\.\d+)?)/);
    if (queryMatch) {
      return {
        lat: Number(queryMatch[1]),
        lng: Number(queryMatch[2]),
      };
    }

    const searchMatch = value.match(/[?&]query=(-?\d+(?:\.\d+)?),(-?\d+(?:\.\d+)?)/);
    if (searchMatch) {
      return {
        lat: Number(searchMatch[1]),
        lng: Number(searchMatch[2]),
      };
    }

    return null;
  }

  function buildSearchUrl(query: string, lat?: number | null, lng?: number | null) {
    if (typeof lat === "number" && typeof lng === "number") {
      return `https://www.google.com/maps?q=${lat},${lng}`;
    }

    return `https://www.google.com/maps/search/${encodeURIComponent(query)}`;
  }

  return {
    extractCoordinatesFromUrl,
    buildSearchUrl,
  };
}
