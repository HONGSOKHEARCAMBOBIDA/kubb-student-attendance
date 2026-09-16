import { defineStore } from 'pinia'
import { getuserdata } from '../api/services'

export const useUserDataStore = defineStore('userdata', {
  state: () => ({
    userdata: null,
  }),

  getters: {
    classid: (state) => state.userdata?.class_id || '',
    name: (state) => state.userdata?.name || '',
    permissions: (state) => state.userdata?.permissions || [],
  },

  actions: {
    async getuserdata() {
      const res = await getuserdata()
      this.userdata = res.data.data || ""
      
    },

    clear() {
      this.userdata = null
    }
  }
})