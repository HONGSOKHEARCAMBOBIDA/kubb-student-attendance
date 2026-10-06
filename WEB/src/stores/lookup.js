import { defineStore } from "pinia";
import { getGeneration, getMajor, getProgramme } from "../api/services";

export const useLookupStore = defineStore("lookup", {
  state: () => ({
    majors: [],
    generations: [],
    programmes: [],
    loaded: false,
    loading: false,
  }),

  actions: {
    async fetchLookUp() {
      if (this.loaded) return;
      this.loading = true;
      try {
        const [majorRes, generationRes, programmeRes] = await Promise.all([
          getMajor(),
          getGeneration(),
          getProgramme(),
        ]);
        this.majors = (majorRes.data.data || []).map((s) => ({
          label: s.name_kh,
          value: s.id,
        }));

        this.generations = (generationRes.data.data || []).map((s) => ({
          label: s.name_kh,
          value: s.id,
        }));

        this.programmes = (programmeRes.data.data || []).map((s) => ({
          label: s.name,
          value: s.id,
        }));
        this.loaded = true;
      } catch (e) {
      } finally {
        this.loading = false;
      }
    },
    clear() {
      this.majors = [];
      this.generations = [];
      this.programmes = [];
      this.loaded = false;
    },
  },
});
