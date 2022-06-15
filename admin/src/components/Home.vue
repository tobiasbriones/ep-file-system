<template>
  <div class="px-5 py-4">
    <AdminUsers :users="users" />
  </div>
</template>

<script>
import AdminUsers from '@/components/Users';
import { Client } from '@/model/client';

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
    this.client = Client(this.handleUsers);
    this.client.init();
  },
  methods: {
    handleUsers(users) {
      this.users = users.map(u => JSON.parse(u));
    }
  }
};
</script>
