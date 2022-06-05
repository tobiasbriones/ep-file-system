<template>
  <div class="px-5 py-4">
    <AdminUsers :users="users" />
  </div>
</template>

<script>
import AdminUsers from '@/components/Users';
import { newClient } from '@/model/client';

export default {
  name: 'AdminHome',
  components: { AdminUsers },
  data() {
    return {
      client: null,
      users: [
        {
          id: 1,
          channel: 'test'
        }
      ]
    };
  },
  created() {
    this.client = newClient(this.onConnected)
  },
  methods: {
    onConnected() {
      console.log("Connected!");
    }
  }
};
</script>
