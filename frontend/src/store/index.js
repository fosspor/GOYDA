import { create } from 'zustand';

export const useUserStore = create((set) => ({
  user: null,
  loading: false,
  setUser: (user) => set({ user }),
  setLoading: (loading) => set({ loading }),
  logout: () => set({ user: null }),
}));

export const useRouteStore = create((set) => ({
  routes: [],
  selectedRoute: null,
  loading: false,
  setRoutes: (routes) => set({ routes }),
  setSelectedRoute: (route) => set({ selectedRoute: route }),
  setLoading: (loading) => set({ loading }),
}));

export const useLocationStore = create((set) => ({
  locations: [],
  selectedLocation: null,
  filters: {},
  setLocations: (locations) => set({ locations }),
  setSelectedLocation: (location) => set({ selectedLocation: location }),
  setFilters: (filters) => set({ filters }),
}));
